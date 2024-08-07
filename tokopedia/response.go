package tokopedia

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type BaseResponse struct {
	Header HeaderResponse `json:"header"`
}

type HeaderResponse struct {
	ProcessTime int               `json:"process_time"`
	Messages    string            `json:"messages"`
	Message     string            `json:"message"`
	Reason      string            `json:"reason"`
	StatusCode  int               `json:"http_code"`
	HTTPHeader  map[string]string `json:"http_header"`
}

// A general response error
type ResponseError struct {
	Header  HeaderResponse `json:"header"`
	Data    interface{}    `json:"data"`
	Message string         `json:"message"`
	ReqID   string         `json:"req_id"`
}

// GetMessage returns response error message
func (e ResponseError) GetMessage() string {
	return e.Header.Messages
}

// GetErrors returns response errors list
func (e ResponseError) GetErrors() string {
	return e.Header.Reason
}

func (e ResponseError) Error() string {
	if e.Header.Messages != "" {
		return e.Header.Messages + " " + e.Header.Reason
	} else if e.Message != "" {
		return e.Message
	}

	return "Unknown Error"
}

// ResponseDecodingError occurs when the response body from Tiktok could
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

func CheckResponseError(r *http.Response) error {
	tokopediaError := ResponseError{}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer func() {
		// already read out, reload for next process
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}()

	if len(bodyBytes) > 0 {
		err := json.Unmarshal(bodyBytes, &tokopediaError)
		if err != nil {
			return ResponseDecodingError{
				Body:    bodyBytes,
				Message: err.Error(),
				Status:  r.StatusCode,
			}
		}
	}

	if tokopediaError.Header.Reason == "" && http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	// headers := map[string]string{
	// 	"X-Ratelimit-Full-Reset-After": r.Header.Get("X-Ratelimit-Full-Reset-After"),
	// 	"X-Ratelimit-Limit":            r.Header.Get("X-Ratelimit-Limit"),
	// 	"X-Ratelimit-Remaining":        r.Header.Get("X-Ratelimit-Remaining"),
	// 	"X-Ratelimit-Reset-After":      r.Header.Get("X-Ratelimit-Reset-After"),
	// }
	// tokopediaError.Header.HTTPHeader = headers
	// tokopediaError.Header.StatusCode = r.StatusCode
	responseError := ResponseError{
		Header:  tokopediaError.Header,
		Data:    tokopediaError.Data,
		Message: tokopediaError.Message,
	}
	// log.Println(responseError.Header)

	return wrapSpecificError(r, responseError)
}

func wrapSpecificError(r *http.Response, err ResponseError) error {
	// TODO: check rate-limit error for tokopedia
	if r.StatusCode == http.StatusTooManyRequests {
		f, _ := strconv.ParseFloat(r.Header.Get("X-Ratelimit-Full-Reset-After"), 64)
		errRateLimit := RateLimitError{
			ResponseError: err,
			RetryAfter:    int(f),
		}
		return errRateLimit
	}

	if r.StatusCode == http.StatusNotAcceptable {
		err.Header.Messages = http.StatusText(r.StatusCode)
	}

	return err
}
