package shopee

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"golang.org/x/net/proxy"
)

const (
	UserAgent = "go-marketplace-sdk/1.0.0"
)

type AppConfig struct {
	PartnerID    int
	PartnerKey   string
	RedirectURL  string
	Client       *ShopeeClient
	APIURL       string
	EnableSocks5 bool
	SockAddress  string
}

type ShopeeClient struct {
	Client    *http.Client
	appConfig AppConfig
	log       LeveledLoggerInterface
	baseURL   *url.URL

	// max number of retries, defaults to 0 for no retries see WithRetry option
	retries  int
	attempts int

	ShopID      uint64
	MerchantID  uint64
	AccessToken string
	Auth        AuthService
	Util        UtilService
	Chat        ChatService
	Product     ProductService
	Order       OrderService
}

// A general response error
type ResponseError struct {
	Status  int
	Message string
	Errors  []string
}

// NewClient returns a new Shopee API client with an already authenticated  and
// a.NewClient(shopName, token, opts) is equivalent to NewClient(a, shopName, token, opts)
func NewClient(app AppConfig, opts ...Option) *ShopeeClient {
	baseURL, err := url.Parse(app.APIURL)
	if err != nil {
		panic(err)
	}

	var proxyURL *url.URL
	var transport *http.Transport
	if app.EnableSocks5 {
		proxyURL, err = url.Parse(fmt.Sprintf("socks5://%s", app.SockAddress))
		if err != nil {
			panic(err)
		}

		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			panic(err)
		}

		transport = &http.Transport{Dial: dialer.Dial}
	}

	c := &ShopeeClient{
		Client:    &http.Client{},
		log:       &LeveledLogger{},
		appConfig: app,
		baseURL:   baseURL,
	}

	if app.EnableSocks5 {
		c.Client.Transport = transport
	}

	c.Auth = &AuthServiceOp{client: c}
	c.Util = &UtilServiceOp{client: c}
	c.Chat = &ChatServiceOp{client: c}
	c.Product = &ProductServiceOp{client: c}
	c.Order = &OrderServiceOp{client: c}

	// apply any options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// GetStatus returns http  response status
func (e ResponseError) GetStatus() int {
	return e.Status
}

// GetMessage returns response error message
func (e ResponseError) GetMessage() string {
	return e.Message
}

// GetErrors returns response errors list
func (e ResponseError) GetErrors() []string {
	return e.Errors
}

func (e ResponseError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	sort.Strings(e.Errors)
	s := strings.Join(e.Errors, ", ")

	if s != "" {
		return s
	}

	return "Unknown Error"
}

// ResponseDecodingError occurs when the response body from Shopee could
// not be parsed.
type ResponseDecodingError struct {
	Body    []byte
	Message string
	Status  int
}

func (e ResponseDecodingError) Error() string {
	return e.Message
}

// An error specific to a rate-limiting response. Embeds the ResponseError to
// allow consumers to handle it the same was a normal ResponseError.
type RateLimitError struct {
	ResponseError
	RetryAfter int
}

// Creates an API request. A relative URL can be provided in urlStr, which will
// be resolved to the BaseURL of the Client. Relative URLS should always be
// specified without a preceding slash. If specified, the value pointed to by
// body is JSON encoded and included as the request body.
func (c *ShopeeClient) NewRequest(method, relPath string, body, options, headers interface{}) (*http.Request, error) {
	rel, err := url.Parse(relPath)
	if err != nil {
		return nil, err
	}

	// Make the full url based on the relative path
	u := c.baseURL.ResolveReference(rel)

	// Add custom options
	if options != nil {
		optionsQuery, err := query.Values(options)
		if err != nil {
			return nil, err
		}

		for k, values := range u.Query() {
			for _, v := range values {
				optionsQuery.Add(k, v)
			}
		}
		u.RawQuery = optionsQuery.Encode()
	}

	// A bit of JSON ceremony
	var js []byte = nil
	if body != nil {
		js, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	c.makeSignature(req)

	return req, nil
}

func (c *ShopeeClient) WithShop(sid uint64, tok string) *ShopeeClient {
	c.ShopID = sid
	c.AccessToken = tok
	return c
}

func (c *ShopeeClient) WithMerchant(mid uint64, tok string) *ShopeeClient {
	c.MerchantID = mid
	c.AccessToken = tok
	return c
}

func (c *ShopeeClient) WithToken(tok string) *ShopeeClient {
	c.AccessToken = tok
	return c
}

// https://open.shopee.com/documents?module=87&type=2&id=58&version=2
func (c *ShopeeClient) makeSignature(req *http.Request) (string, int64) {
	ts := time.Now().Unix()
	path := req.URL.Path

	var baseStr string

	u := req.URL

	query := u.Query()
	query.Add("partner_id", fmt.Sprintf("%v", c.appConfig.PartnerID))

	if c.ShopID != 0 {
		// Shop APIs: partner_id, api path, timestamp, access_token, shop_id
		baseStr = fmt.Sprintf("%d%s%d%s%d", c.appConfig.PartnerID, path, ts, c.AccessToken, c.ShopID)
		query.Add("shop_id", fmt.Sprintf("%v", c.ShopID))
		query.Add("access_token", c.AccessToken)
	} else if c.MerchantID != 0 {
		// Merchant APIs: partner_id, api path, timestamp, access_token, merchant_id
		baseStr = fmt.Sprintf("%d%s%d%s%d", c.appConfig.PartnerID, path, ts, c.AccessToken, c.MerchantID)
		query.Add("merchant_id", fmt.Sprintf("%v", c.MerchantID))
		query.Add("access_token", c.AccessToken)
	} else {
		// Public APIs: partner_id, api path, timestamp
		baseStr = fmt.Sprintf("%d%s%d", c.appConfig.PartnerID, path, ts)
	}
	h := hmac.New(sha256.New, []byte(c.appConfig.PartnerKey))
	h.Write([]byte(baseStr))
	result := hex.EncodeToString(h.Sum(nil))

	query.Add("timestamp", fmt.Sprintf("%v", ts))
	query.Add("sign", result)

	u.RawQuery = query.Encode()
	req.URL = u

	return result, ts
}

// doGetHeaders executes a request, decoding the response into `v` and also returns any response headers.
func (c *ShopeeClient) doGetHeaders(req *http.Request, v interface{}, skipBody bool) (http.Header, error) {
	var resp *http.Response
	var err error

	retries := c.retries
	c.attempts = 0
	c.logRequest(req, skipBody)

	for {
		c.attempts++

		resp, err = c.Client.Do(req)
		c.logResponse(resp)
		if err != nil {
			return nil, err //http client errors, not api responses
		}

		respErr := CheckResponseError(resp)
		if respErr == nil {
			break // no errors, break out of the retry loop
		}

		// retry scenario, close resp and any continue will retry
		resp.Body.Close()

		if retries <= 1 {
			return nil, respErr
		}

		if rateLimitErr, isRetryErr := respErr.(RateLimitError); isRetryErr {
			// back off and retry

			wait := time.Duration(rateLimitErr.RetryAfter) * time.Second
			c.log.Debugf("rate limited waiting %s", wait.String())
			time.Sleep(wait)
			retries--
			continue
		}

		var doRetry bool
		switch resp.StatusCode {
		case http.StatusServiceUnavailable:
			c.log.Debugf("service unavailable, retrying")
			doRetry = true
			retries--
		}

		if doRetry {
			continue
		}

		// no retry attempts, just return the err
		return nil, respErr
	}

	c.logResponse(resp)
	defer resp.Body.Close()

	if v != nil {
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&v)
		if err != nil {
			return nil, err
		}
	}

	return resp.Header, nil
}

// skipBody: if upload image, skip log its binary
func (c *ShopeeClient) logRequest(req *http.Request, skipBody bool) {
	if req == nil {
		return
	}
	if req.URL != nil {
		c.log.Debugf("%s: %s", req.Method, req.URL.String())
	}
	if !skipBody {
		c.logBody(&req.Body, "SENT: %s")
	}
}

func (c *ShopeeClient) logResponse(res *http.Response) {
	if res == nil {
		return
	}
	c.log.Debugf("RECV %d: %s", res.StatusCode, res.Status)
	c.logBody(&res.Body, "RESP: %s")
}

func (c *ShopeeClient) logBody(body *io.ReadCloser, format string) {
	if body == nil || *body == nil {
		return
	}
	b, _ := ioutil.ReadAll(*body)
	if len(b) > 0 {
		c.log.Debugf(format, string(b))
	}
	*body = ioutil.NopCloser(bytes.NewBuffer(b))
}

func wrapSpecificError(r *http.Response, err ResponseError) error {
	// TODO: check rate-limit error for shopee
	if err.Status == http.StatusTooManyRequests {
		f, _ := strconv.ParseFloat(r.Header.Get("Retry-After"), 64)
		return RateLimitError{
			ResponseError: err,
			RetryAfter:    int(f),
		}
	}

	// if err.Status == http.StatusSeeOther {
	// todo
	// The response to the request can be found under a different URL in the
	// Location header and can be retrieved using a GET method on that resource.
	// }

	if err.Status == http.StatusNotAcceptable {
		err.Message = http.StatusText(err.Status)
	}

	return err
}

// shopee error maybe return status=200
// eg. {"error":"error_incalid_category.","message":"Invalid category ID","request_id":"2069449bd255af166cb52b0e15189d6d"}
// {"error":"error_category_is_block.","message":"Category is restricted","request_id":"97994a47af37a22da79cb910bfd9841a"}
func CheckResponseError(r *http.Response) error {
	shopeeError := struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}{}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer func() {
		// already read out, reload for next process
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}()

	if len(bodyBytes) > 0 {
		err := json.Unmarshal(bodyBytes, &shopeeError)
		if err != nil {
			return ResponseDecodingError{
				Body:    bodyBytes,
				Message: err.Error(),
				Status:  r.StatusCode,
			}
		}
	}

	if shopeeError.Error == "" && http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	responseError := ResponseError{
		Status:  r.StatusCode,
		Message: fmt.Sprintf("shopee-%s [%s]", shopeeError.Error, shopeeError.Message),
	}

	return wrapSpecificError(r, responseError)
}

// CreateAndDo performs a web request to Shopee with the given method (GET,
// POST, PUT, DELETE) and relative path
// The data, options and resource arguments are optional and only relevant in
// certain situations.
// If the data argument is non-nil, it will be used as the body of the request
// for POST and PUT requests.
func (c *ShopeeClient) CreateAndDo(method, relPath string, data, options, headers, resource interface{}) error {
	defer func() {
		// clear for next call
		c.ShopID = 0
		c.MerchantID = 0
		c.AccessToken = ""
	}()

	_, err := c.createAndDoGetHeaders(method, relPath, data, options, headers, resource)
	if err != nil {
		return err
	}
	return nil
}

// createAndDoGetHeaders creates an executes a request while returning the response headers.
func (c *ShopeeClient) createAndDoGetHeaders(method, relPath string, data, options, headers, resource interface{}) (http.Header, error) {
	if strings.HasPrefix(relPath, "/") {
		// make sure it's a relative path
		relPath = strings.TrimLeft(relPath, "/")
	}

	relPath = path.Join("api/v2", relPath)

	if data != nil {
		params := data.(map[string]interface{})
		params["partner_id"] = c.appConfig.PartnerID
		data = params
	}

	req, err := c.NewRequest(method, relPath, data, options, headers)
	if err != nil {
		return nil, err
	}

	return c.doGetHeaders(req, resource, false)
}

// Get performs a GET request for the given path and saves the result in the
// given resource.
func (c *ShopeeClient) Get(path string, resource, options interface{}) error {
	return c.CreateAndDo("GET", path, nil, options, nil, resource)
}

// Post performs a POST request for the given path and saves the result in the
// given resource.
func (c *ShopeeClient) Post(path string, data, resource interface{}) error {
	return c.CreateAndDo("POST", path, data, nil, nil, resource)
}

// Put performs a PUT request for the given path and saves the result in the
// given resource.
func (c *ShopeeClient) Put(path string, data, resource interface{}) error {
	return c.CreateAndDo("PUT", path, data, nil, nil, resource)
}

// Delete performs a DELETE request for the given path
func (c *ShopeeClient) Delete(path string) error {
	return c.CreateAndDo("DELETE", path, nil, nil, nil, nil)
}

// Upload performs a Upload request for the given path and saves the result in the
// given resource.
func (c *ShopeeClient) Upload(relPath, fieldname, filename string, resource interface{}) error {
	req, err := c.NewfileUploadRequest(relPath, fieldname, filename)
	if err != nil {
		return err
	}

	if _, err := c.doGetHeaders(req, resource, true); err != nil {
		return err
	}

	return nil
}

// Creates a new file upload http request with optional extra params
func (c *ShopeeClient) NewfileUploadRequest(relPath, paramName, filename string) (*http.Request, error) {
	if strings.HasPrefix(relPath, "/") {
		// make sure it's a relative path
		relPath = strings.TrimLeft(relPath, "/")
	}

	relPath = path.Join("api/v2", relPath)

	rel, err := url.Parse(relPath)
	if err != nil {
		return nil, err
	}

	// Make the full url based on the relative path
	u := c.baseURL.ResolveReference(rel)
	uri := u.String()

	// Replace os.Open with http.Get to fetch data from the URL
	resp, err := http.Get(filename)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, filepath.Base(uri))
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, resp.Body); err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", UserAgent)

	c.makeSignature(req)

	return req, nil
}
