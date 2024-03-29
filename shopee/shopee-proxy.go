package shopee

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/net/proxy"
)

type ProxyClient struct {
	Client      *resty.Client
	AccessToken string
	appConfig   ProxyAppConfig

	ShopID     uint64
	MerchantID uint64

	baseURL *url.URL
}

type RequestOptions struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Auth    map[string]string `json:"auth,omitempty"`
	Gzip    bool              `json:"gzip,omitempty"`
	Body    interface{}       `json:"body,omitempty"`
	Qs      string            `json:"qs,omitempty"`
	Headers map[string]string `json:"headers"`
	JSON    bool              `json:"json,omitempty"`
	Timeout int               `json:"timeout,omitempty"`
}

type ProxyAppConfig struct {
	ProxyURL    string
	ProxyHost   string
	PartnerID   uint64
	PartnerKey  string
	ShopID      uint64
	MerchantID  uint64
	AccessToken string
	APIURL      *url.URL
	EnableLog   bool
	EnableRetry bool
	RetryCount  int
	MaxTimeout  time.Duration
	UseSocks5   bool
}

// New method creates a new HTTPProxy client.
func NewProxyClient(app ProxyAppConfig) *ProxyClient {

	if app.MaxTimeout == 0 {
		app.MaxTimeout = 30 * time.Second
	}

	var proxyURL *url.URL
	var transport *http.Transport

	// Create socks5 proxy agent if app.UseSocks5 is true
	if app.UseSocks5 {
		proxyURL, _ = url.Parse(fmt.Sprintf("socks5://%s", app.ProxyHost))
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			panic(err)
		}

		transport = &http.Transport{Dial: dialer.Dial}
	}

	// Create resty client
	client := resty.New().
		SetTimeout(app.MaxTimeout).
		OnBeforeRequest(
			func(c *resty.Client, r *resty.Request) error {
				if app.EnableLog {
					log.Printf("Request: %s %s", r.Method, r.URL)
				}

				return nil
			}).
		OnAfterResponse(
			func(c *resty.Client, r *resty.Response) error {
				if app.EnableLog {
					log.Println("Response Status", r.Status())
				}
				return nil
			})

	var baseURL *url.URL
	if app.UseSocks5 {
		baseURL, _ = url.Parse(app.APIURL.String())
		client.SetBaseURL(app.APIURL.String())
		client.SetTransport(transport)
	} else {
		baseURL, _ = url.Parse(app.ProxyURL)
		client.SetBaseURL(app.ProxyURL)
	}

	// Enable retry if enabled in config
	if app.EnableRetry {
		client.SetRetryCount(app.RetryCount).
			SetRetryWaitTime(5 * time.Second).
			SetRetryMaxWaitTime(1 * time.Minute).
			AddRetryCondition(
				func(r *resty.Response, err error) bool {
					if err != nil {
						log.Println("retry on Error:", err.Error())
						return true
					}
					if (r.StatusCode() == http.StatusBadGateway) || (r.StatusCode() == http.StatusGatewayTimeout) {
						log.Println("retry on Http Error:", r.StatusCode())
						return true
					}
					return false
				})

	}

	return &ProxyClient{
		Client:      client,
		AccessToken: app.AccessToken,
		ShopID:      app.ShopID,
		MerchantID:  app.MerchantID,
		appConfig:   app,
		baseURL:     baseURL,
	}
}

func (c *ProxyClient) WithShopID(sid uint64, tok string) *ProxyClient {
	c.ShopID = sid
	c.AccessToken = tok

	return c
}

func (c *ProxyClient) MakeProxySignature(path string) (string, int64) {
	ts := time.Now().Unix()
	var baseStr string

	if c.ShopID != 0 {
		baseStr = fmt.Sprintf("%d%s%d%s%d", c.appConfig.PartnerID, path, ts, c.AccessToken, c.ShopID)
	} else if c.MerchantID != 0 {
		baseStr = fmt.Sprintf("%d%s%d%s%d", c.appConfig.PartnerID, path, ts, c.AccessToken, c.MerchantID)
	} else {
		// Public APIs: partner_id, api path, timestamp
		baseStr = fmt.Sprintf("%d%s%d", c.appConfig.PartnerID, path, ts)
	}

	h := hmac.New(sha256.New, []byte(c.appConfig.PartnerKey))
	h.Write([]byte(baseStr))
	result := hex.EncodeToString(h.Sum(nil))

	return result, ts
}

func (c *ProxyClient) generateFullURL(relPath string) string {
	if strings.HasPrefix(relPath, "/") {
		// make sure it's a relative path
		relPath = strings.TrimLeft(relPath, "/")
	}

	// Combine the relative path with the Shopee API URL
	relPath = path.Join("api/v2", relPath)

	rel, err := url.Parse(relPath)
	if err != nil {
		return ""
	}

	// Make the full url based on the relative path
	u := c.appConfig.APIURL.ResolveReference(rel)

	// Generate the signature and timestamp for the request
	signature, timestamp := c.MakeProxySignature(u.Path)

	// Add query parameters to the URL with the signature and timestamp
	query := u.Query()
	query.Add("partner_id", fmt.Sprintf("%d", c.appConfig.PartnerID))

	if c.ShopID != 0 {
		query.Add("shop_id", fmt.Sprintf("%d", c.ShopID))
		query.Add("access_token", c.AccessToken)
	}

	query.Add("timestamp", fmt.Sprintf("%d", timestamp))
	query.Add("sign", signature)

	// Encode the query parameters and set them in the URL
	// This is necessary because the URL is used in the request body
	u.RawQuery = query.Encode()
	uri := u.String()

	return uri
}

func (h *ProxyClient) GetFullPath(relPath string) (res string) {
	// Make the full url based on the relative path
	// and generate the signature and timestamp for the request
	uri := h.generateFullURL(relPath)
	return uri
}

// SendFormDataRequest sends a multipart/form-data request with the given file and URI
// this is used to upload an image to our proxy server
func (c *ProxyClient) SendUploadRequest(file, relPath string) (res *resty.Response, err error) {

	// Download the image from the URL
	// Instead of using the io.Copy, we use http.Get to download the image directly
	respImage, err := http.Get(file)
	if err != nil {
		return nil, err
	}

	// Read the image data from the response body and close it
	imgData, err := io.ReadAll(respImage.Body)
	if err != nil {
		log.Fatalf("Error reading image data: %v", err)
	}

	// Check if image data is not empty
	if len(imgData) == 0 {
		log.Fatalf("Downloaded image data is empty")
	}

	// Close the response body when the function returns
	defer respImage.Body.Close()

	// Make the full url based on the relative path
	uri := c.generateFullURL(relPath)

	// Set up and execute the request with the image data as the body
	// and the multipart/form-data content type header
	// This is used to upload an image to our proxy server using the Shopee API
	resp, err := c.Client.R().
		SetResult(UploadImageResponse{}).
		SetHeader("Content-Type", "multipart/form-data").
		SetFileReader("file", "filename.jpeg", bytes.NewReader(imgData)).
		SetMultipartFormData(
			map[string]string{
				"url": uri,
			},
		).
		SetHeader("Connection", "keep-alive").
		Post("/api/proxy/upload-image")

	return resp, nil
}
