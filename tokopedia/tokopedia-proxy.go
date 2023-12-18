package tokopedia

import (
	"log"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty/v2"
)

type ProxyClient struct {
	client    *resty.Client
	appConfig ProxyAppConfig

	AccessToken string
	IsAuth      bool

	apiURL *url.URL
}

type RequestOptions struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Auth    interface{}       `json:"auth,omitempty"`
	Gzip    bool              `json:"gzip,omitempty"`
	Body    interface{}       `json:"body,omitempty"`
	Qs      interface{}       `json:"qs,omitempty"`
	Headers map[string]string `json:"headers"`
	JSON    bool              `json:"json,omitempty"`
	Timeout int               `json:"timeout,omitempty"`
}

type ProxyAppConfig struct {
	ProxyURL string
	AppID    int64

	ShopID      uint64
	AccessToken string

	APIURL *url.URL

	EnableLog   bool
	EnableRetry bool
	RetryCount  int
}

func NewHTTPProxy(app ProxyAppConfig) *ProxyClient {
	client := resty.New().
		SetBaseURL(app.ProxyURL).
		OnBeforeRequest(
			func(c *resty.Client, r *resty.Request) error {
				if app.EnableLog {
					log.Printf("=================================")
					log.Printf("Request: %s %s", r.Method, r.URL)
				}

				return nil
			}).
		OnAfterResponse(
			func(c *resty.Client, r *resty.Response) error {
				if app.EnableLog {
					log.Println("Response Status", r.Status())
					log.Printf("=================================")
				}
				return nil
			})

	if app.EnableRetry {
		client.SetRetryCount(app.RetryCount).
			SetRetryWaitTime(5 * time.Second).
			SetRetryMaxWaitTime(1 * time.Minute).
			AddRetryAfterErrorCondition()

	}

	return &ProxyClient{
		client:    client,
		appConfig: app,
		apiURL:    app.APIURL,
	}
}

func (c *ProxyClient) InitAuth() *ProxyClient {
	baseURL, err := url.Parse("https://accounts.tokopedia.com")
	if err != nil {
		log.Fatal(err)
	}

	c.IsAuth = true
	c.apiURL = baseURL
	c.appConfig.APIURL = baseURL

	return c
}

func (c *ProxyClient) SetAccessToken(token string) *ProxyClient {
	c.AccessToken = token

	return c
}

func (c *ProxyClient) generateFullURL(relPath string) string {
	if strings.HasPrefix(relPath, "/") {

		// make sure it's a relative path
		relPath = strings.TrimLeft(relPath, "/")
	}

	// Combine the relative path with the Tokopedia API if not init auth
	if !c.IsAuth {
		relPath = path.Join("v1", relPath)
	}

	rel, err := url.Parse(relPath)
	if err != nil {
		return ""
	}

	// Make the full url based on the relative path
	u := c.apiURL.ResolveReference(rel)

	// Encode the query parameters and set them in the URL
	// This is necessary because the URL is used in the request body
	uri := u.String()

	return uri
}

func (c *ProxyClient) CreateParams(relPath string, method string, data interface{}) (res RequestOptions) {
	uri := c.generateFullURL(relPath)

	options := RequestOptions{
		Method:  method,
		JSON:    true,
		Timeout: 60000,
		URL:     uri,
		Gzip:    true,
		Headers: map[string]string{
			"Connection": "keep-alive",
		},
	}

	if c.AccessToken != "" {
		options.Headers["Authorization"] = "Bearer " + c.AccessToken
	}

	if c.IsAuth {
		options.Auth = data
	}

	if method == "POST" && !c.IsAuth {
		options.Body = data
	}

	spew.Dump(options)
	return options
}

func (c *ProxyClient) SendRequest(relPath string, options RequestOptions) (*resty.Response, error) {
	// Set up and execute the request
	body := map[string]interface{}{
		"requestOption": options,
	}

	resp, err := c.client.R().
		SetBody(body).
		SetHeader("Content-Type", "application/json").
		Post(relPath)

	return resp, err
}
