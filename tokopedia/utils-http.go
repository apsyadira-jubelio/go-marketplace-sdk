package tokopedia

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// CreateAndDo performs a web request to Shopee with the given method (GET,
// POST, PUT, DELETE) and relative path
// The data, options and resource arguments are optional and only relevant in
// certain situations.
// If the data argument is non-nil, it will be used as the body of the request
// for POST and PUT requests.
func (c *TokopediaClient) CreateAndDo(method, relPath string, data, options, headers, resource interface{}) error {
	defer func() {
		// clear for next call
		c.ShopID = ""
		c.AccessToken = ""
		c.AuthToken = ""

	}()

	_, err := c.createAndDoGetHeaders(method, relPath, data, options, headers, resource)
	if err != nil {
		return err
	}

	return nil
}

// createAndDoGetHeaders creates an executes a request while returning the response headers.
func (c *TokopediaClient) createAndDoGetHeaders(method, relPath string, data, options, headers, resource interface{}) (http.Header, error) {
	if strings.HasPrefix(relPath, "/") {
		// make sure it's a relative path
		relPath = strings.TrimLeft(relPath, "/")
	}

	req, err := c.NewRequest(method, relPath, data, options, headers)
	if err != nil {
		return nil, err
	}

	// log.Printf("path in createAndDoGetHeaders:%s\n", relPath)
	return c.doGetHeaders(req, resource, false)
}

// Get performs a GET request for the given path and saves the result in the
// given resource.
func (c *TokopediaClient) Get(path string, resource, options interface{}) error {
	return c.CreateAndDo("GET", path, nil, options, nil, resource)
}

// Post performs a POST request for the given path and saves the result in the
// given resource.
func (c *TokopediaClient) Post(path string, data, resource interface{}) error {
	return c.CreateAndDo("POST", path, data, nil, nil, resource)
}

// Put performs a PUT request for the given path and saves the result in the
// given resource.
func (c *TokopediaClient) Put(path string, data, resource interface{}) error {
	return c.CreateAndDo("PUT", path, data, nil, nil, resource)
}

// Delete performs a DELETE request for the given path
func (c *TokopediaClient) Delete(path string) error {
	return c.CreateAndDo("DELETE", path, nil, nil, nil, nil)
}

// Creates an API request. A relative URL can be provided in urlStr, which will
// be resolved to the BaseURL of the Client. Relative URLS should always be
// specified without a preceding slash. If specified, the value pointed to by
// body is JSON encoded and included as the request body.
func (c *TokopediaClient) NewRequest(method, relPath string, body, options, headers interface{}) (*http.Request, error) {
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
			c.log.Errorf("[NewRequest] error in marshall:%+v\n", err)
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(js))
	if err != nil {
		c.log.Errorf("[NewRequest] error in create new request:%+v\n", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if c.AccessToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.AccessToken)
	}

	if c.AuthToken != "" {
		req.Header.Add("Authorization", "Basic "+c.AuthToken)
	}

	return req, nil
}

// Upload performs a Upload request for the given path and saves the result in the
// given resource.
func (c *TokopediaClient) Upload(relPath, fieldname, filename string, resource interface{}) error {
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
func (c *TokopediaClient) NewfileUploadRequest(relPath, paramName, filename string) (*http.Request, error) {
	if strings.HasPrefix(relPath, "/") {
		// make sure it's a relative path
		relPath = strings.TrimLeft(relPath, "/")
	}

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

	if c.AccessToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.AccessToken)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", UserAgent)

	return req, nil
}

// doGetHeaders executes a request, decoding the response into `v` and also returns any response headers.
func (c *TokopediaClient) doGetHeaders(req *http.Request, v interface{}, skipBody bool) (http.Header, error) {
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
		defer resp.Body.Close()

		if retries <= 1 {
			return nil, respErr
		}

		if _, ok := respErr.(RateLimitError); ok {
			rateLimitErr := respErr.(RateLimitError)
			// back off and retry
			wait := time.Duration(rateLimitErr.RetryAfter) * time.Second
			time.Sleep(wait)
			retries--
			continue
		}

		var doRetry bool
		switch resp.StatusCode {
		case http.StatusServiceUnavailable:
			c.log.Errorf("service unavailable, retrying")
			doRetry = true
			retries--
		}

		if doRetry {
			continue
		}

		// no retry attempts, just return the err
		return nil, respErr
	}

	defer resp.Body.Close()

	if v != nil {
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&v)
		if err != nil {
			return nil, err
		}
		// for chat/product only
		// addHeadersToResponse(v, resp)
	}

	return resp.Header, nil
}

// addHeadersToResponse adds common headers to different response structs
func addHeadersToResponse(v interface{}, resp *http.Response) {
	headers := map[string]string{
		"X-Ratelimit-Full-Reset-After": resp.Header.Get("X-Ratelimit-Full-Reset-After"),
		"X-Ratelimit-Limit":            resp.Header.Get("X-Ratelimit-Limit"),
		"X-Ratelimit-Remaining":        resp.Header.Get("X-Ratelimit-Remaining"),
		"X-Ratelimit-Reset-After":      resp.Header.Get("X-Ratelimit-Reset-After"),
	}

	switch response := v.(type) {
	case *ReplyListResponse:
		response.Header.StatusCode = resp.StatusCode
		response.Header.HTTPHeader = headers
	case *MessageResponse:
		response.Header.StatusCode = resp.StatusCode
		response.Header.HTTPHeader = headers
	case *SendMessageResponse:
		response.Header.StatusCode = resp.StatusCode
		response.Header.HTTPHeader = headers
	case *ProductInfoResponse:
		response.Header.StatusCode = resp.StatusCode
		response.Header.HTTPHeader = headers
	default:
		if response == nil {
			log.Println("response header is nil")
		} else {
			log.Println("response header is not mandatory")
		}
	}
}
