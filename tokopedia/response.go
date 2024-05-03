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
	ProcessTime int    `json:"process_time"`
	Messages    string `json:"messages"`
	Reason      string `json:"reason"`
}

// A general response error
type ResponseError struct {
	Header HeaderResponse `json:"header"`
	Data   interface{}    `json:"data"`
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

	responseError := ResponseError{
		Header: tokopediaError.Header,
		Data:   tokopediaError.Data,
	}

	return wrapSpecificError(r, responseError)
}

func wrapSpecificError(r *http.Response, err ResponseError) error {
	// TODO: check rate-limit error for shopee
	if r.StatusCode == http.StatusTooManyRequests {
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

	if r.StatusCode == http.StatusNotAcceptable {
		err.Header.Messages = http.StatusText(r.StatusCode)
	}

	return err
}
