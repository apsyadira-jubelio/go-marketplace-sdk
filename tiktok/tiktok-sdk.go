package tiktok

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"sort"

	"net/http"
	"net/url"
	"time"
)

const (
	UserAgent     = "go-marketplace-sdk/1.0.0"
	OpenAPIURL    = "https://open-api.tiktokglobalshop.com"
	AuthBaseURL   = "https://services.tiktokshop.com"
	LegacyAuthURL = "https://auth.tiktok-shops.com"
)

type AppConfig struct {
	AppKey      string
	AppSecret   string
	RedirectURL string
	Client      *TiktokClient
	APIURL      string
	Version     string
}

type TiktokClient struct {
	Client    *http.Client
	log       LeveledLoggerInterface
	appConfig AppConfig
	baseURL   *url.URL

	// max number of retries, defaults to 0 for no retries see WithRetry option
	retries  int
	attempts int

	ShopCipher  string
	AccessToken string
	ShopID      string

	Auth AuthService
	Util UtilService
	Chat ChatService
}

type CommonParamRequest struct {
	AccessToken    string
	ShopID         string
	ShopCipher     string
	ConversationID string
}

func NewClient(app AppConfig, opts ...Option) *TiktokClient {
	baseURL, err := url.Parse(app.APIURL)
	if err != nil {
		panic(err)
	}

	c := &TiktokClient{
		Client:    &http.Client{},
		log:       &LeveledLogger{},
		appConfig: app,
		baseURL:   baseURL,
	}

	c.Auth = &AuthServiceOp{client: c}
	c.Util = &UtilServiceOp{client: c}
	c.Chat = &ChatServiceOp{client: c}

	// apply any options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *TiktokClient) WithCommonParamRequest(param CommonParamRequest) *TiktokClient {
	if param.AccessToken != "" {
		c.AccessToken = param.AccessToken
	}
	if param.ShopID != "" {
		c.ShopID = param.ShopID
	}
	if param.ShopCipher != "" {
		c.ShopCipher = param.ShopCipher
	}

	return c
}

func (c *TiktokClient) WithShopCipher(cipher string) *TiktokClient {
	c.ShopCipher = cipher
	return c
}

func (c *TiktokClient) WithAccessToken(accessToken string) *TiktokClient {
	if accessToken == "" {
		return c
	}

	c.AccessToken = accessToken
	return c
}

func (c *TiktokClient) WithShopID(shopID string) *TiktokClient {
	if shopID == "" {
		return c
	}

	c.ShopID = shopID
	return c
}

func (c *TiktokClient) makeSignature(req *http.Request) string {
	ts := time.Now().Unix()
	u := req.URL

	query := u.Query()
	if c.ShopCipher != "" {
		query.Add("shop_cipher", fmt.Sprintf("%v", c.ShopCipher))
	}

	if c.ShopID != "" {
		query.Add("shop_id", fmt.Sprintf("%v", c.ShopID))
	}

	if c.AccessToken != "" {
		query.Add("access_token", fmt.Sprintf("%v", c.AccessToken))
	}

	query.Add("app_key", c.appConfig.AppKey)
	query.Add("timestamp", fmt.Sprintf("%v", ts))

	u.RawQuery = query.Encode()
	req.URL = u
	// only for not auth API
	signResult := c.CalSignAndGenerateSignature(req, c.appConfig.AppSecret)

	query.Add("sign", signResult)

	u.RawQuery = query.Encode()
	req.URL = u

	return signResult
}

func (c *TiktokClient) CalSignAndGenerateSignature(req *http.Request, secret string) string {
	queries := req.URL.Query()

	// extract all query parameters excluding sign and access_token
	keys := make([]string, len(queries))
	idx := 0
	for key := range queries {
		// params except 'sign' & 'access_token'
		if key != "sign" && key != "access_token" {
			keys[idx] = key
			idx++
		}
	}

	// reorder the parameters' key in alphabetical order
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// Concatenate all the parameters in the format of {key}{value}
	input := ""
	for _, key := range keys {
		input = input + key + queries.Get(key)
	}

	// append the request path
	input = req.URL.Path + input

	// if the request header Content-type is not multipart/form-data, append body to the end
	mediaType, _, _ := mime.ParseMediaType(req.Header.Get("Content-type"))
	if mediaType != "multipart/form-data" {
		body, _ := io.ReadAll(req.Body)
		input = input + string(body)

		req.Body.Close()
		// reset body after reading from the original
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	// wrap the string generated in step 5 with the App secret
	input = secret + input + secret

	return c.generateSHA256(input, secret)
}

func (c *TiktokClient) generateSHA256(input, secret string) string {
	// encode the digest byte stream in hexadecimal and use sha256 to generate sign with salt(secret)
	h := hmac.New(sha256.New, []byte(secret))

	if _, err := h.Write([]byte(input)); err != nil {
		// TODO error log
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}
