package tokopedia

import (
	"net/http"
	"net/url"
)

const (
	UserAgent = "go-marketplace-sdk/1.0.0"
	AuthURL   = "https://accounts.tokopedia.com"
	APIURL    = "https://fs.tokopedia.net"
)

type AppConfig struct {
	ClientID     string
	ClientSecret string
	FsID         int
	Client       *TokopediaClient
	APIURL       string
	Version      string
}

type TokopediaClient struct {
	Client    *http.Client
	log       LeveledLoggerInterface
	appConfig AppConfig
	baseURL   *url.URL

	// max number of retries, defaults to 0 for no retries see WithRetry option
	retries  int
	attempts int

	AccessToken string
	AuthToken   string
	ShopID      string

	Auth    AuthService
	Chat    ChatService
	Product ProductService
	Shop    ShopService
}

type CommonParamRequest struct {
	AccessToken    string
	ShopID         string
	ShopCipher     string
	ConversationID string
}

func NewClient(app AppConfig, opts ...Option) *TokopediaClient {
	baseURL, err := url.Parse(app.APIURL)
	if err != nil {
		panic(err)
	}

	c := &TokopediaClient{
		Client:    &http.Client{},
		log:       &LeveledLogger{},
		appConfig: app,
		baseURL:   baseURL,
	}

	c.Auth = &AuthServiceOp{client: c}
	c.Chat = &ChatServiceOp{client: c}
	c.Product = &ProductServiceOp{client: c}
	c.Shop = &ShopServiceOp{client: c}

	// apply any options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *TokopediaClient) WithCommonParamRequest(param CommonParamRequest) *TokopediaClient {
	if param.AccessToken != "" {
		c.AccessToken = param.AccessToken
	}
	if param.ShopID != "" {
		c.ShopID = param.ShopID
	}

	return c
}

func (c *TokopediaClient) WithAccessToken(accessToken string) *TokopediaClient {
	if accessToken == "" {
		return c
	}

	c.AccessToken = accessToken
	return c
}

func (c *TokopediaClient) WithBasicAuth(token string) *TokopediaClient {
	if token == "" {
		return c
	}

	c.AuthToken = token
	return c
}

func (c *TokopediaClient) WithShopID(shopID string) *TokopediaClient {
	if shopID == "" {
		return c
	}

	c.ShopID = shopID
	return c
}
