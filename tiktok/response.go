package tiktok

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type BaseResponse struct {
	RequestID string `json:"request_id"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	Code      string `json:"code"`
}

// A general response error
type ResponseError struct {
	Status  int
	Message string
	Errors  []string
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
	tiktokError := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer func() {
		// already read out, reload for next process
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}()

	if len(bodyBytes) > 0 {
		err := json.Unmarshal(bodyBytes, &tiktokError)
		if err != nil {
			return ResponseDecodingError{
				Body:    bodyBytes,
				Message: err.Error(),
				Status:  r.StatusCode,
			}
		}
	}

	if tiktokError.Code == 0 && http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	responseError := ResponseError{
		Status:  r.StatusCode,
		Message: fmt.Sprintf("tiktok-%d [%s]", tiktokError.Code, tiktokError.Message),
	}

	return wrapSpecificError(r, responseError)
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
