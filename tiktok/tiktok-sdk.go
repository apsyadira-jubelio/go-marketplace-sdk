package tiktok

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"net/http"
	"net/url"
	"time"
)

const (
	UserAgent = "go-marketplace-sdk/1.0.0"
)

type AppConfig struct {
	AppKey      string
	AppSecret   string
	RedirectURL string
	Client      *TiktokClient
	APIURL      string
}

type TiktokClient struct {
	Client    *http.Client
	log       LeveledLoggerInterface
	appConfig AppConfig
	baseURL   *url.URL

	// max number of retries, defaults to 0 for no retries see WithRetry option
	retries  int
	attempts int

	ShopChiper string
	// AccessToken string

	Auth AuthService
	Util UtilService
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

	// apply any options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *TiktokClient) WithShopChiper(chiper string) *TiktokClient {
	c.ShopChiper = chiper
	return c
}

func (c *TiktokClient) makeSignature(req *http.Request) (string, int64) {
	ts := time.Now().Unix()
	path := req.URL.Path

	var baseStr string

	u := req.URL

	query := u.Query()
	if c.ShopChiper != "" {
		baseStr = fmt.Sprintf("%s%d%s", path, ts, c.ShopChiper)
		query.Add("shop_chiper", fmt.Sprintf("%v", c.ShopChiper))
	} else {
		baseStr = fmt.Sprintf("%s%s%d", c.appConfig.AppKey, path, ts)
	}

	h := hmac.New(sha256.New, []byte(c.appConfig.AppKey))
	h.Write([]byte(baseStr))
	result := hex.EncodeToString(h.Sum(nil))

	query.Add("timestamp", fmt.Sprintf("%v", ts))
	query.Add("sign", result)

	u.RawQuery = query.Encode()
	req.URL = u

	return result, ts
}
