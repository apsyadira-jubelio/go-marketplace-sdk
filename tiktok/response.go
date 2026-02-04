package tiktok

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type BaseResponse struct {
	RequestID string `json:"request_id"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
}

// A general response error
type ResponseError struct {
	Message  string   `json:"message"`
	Status   int      `json:"status"`
	RequstID string   `json:"request_id"`
	Code     int      `json:"code"`
	Errors   []string `json:"errors,omitempty"`
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
		errorBytes, err := json.Marshal(e)
		if err != nil {
			return fmt.Sprintf("unknown error: %s", err)
		}

		return string(errorBytes)
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
	var tiktokError ResponseError

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while io.ReadAll:", err)
		return err
	}

	// Restore the body for further reading downstream
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	defer func() {
		// Defensive: makes sure Body is always available again
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}()

	if len(bodyBytes) == 0 {
		// If body is empty, consider only StatusCode
		if http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices {
			return nil
		}
		return fmt.Errorf("empty response body, status: %d", r.StatusCode)
	}

	contentType := r.Header.Get("Content-Type")
	isJSON := strings.Contains(contentType, "application/json")
	if isJSON {
		err := json.Unmarshal(bodyBytes, &tiktokError)
		if err != nil {
			log.Printf("[CheckResponseError] Failed to decode JSON. Body:\n%s\n", string(bodyBytes))
			return ResponseDecodingError{
				Body:    bodyBytes,
				Message: err.Error(),
				Status:  r.StatusCode,
			}
		}
	} else {
		// Non-JSON: maybe HTML error, log it
		if bodyBytes[0] == '<' {
			log.Printf("[CheckResponseError] Non-JSON (possible HTML) body received with status %d:\n%s\n", r.StatusCode, string(bodyBytes))
			return fmt.Errorf("unexpected non-JSON response received (status %d)", r.StatusCode)
		} else {
			log.Printf("[CheckResponseError] Unexpected non-JSON body received: %s", string(bodyBytes))
			return fmt.Errorf("unexpected non-JSON response body (status %d)", r.StatusCode)
		}
	}

	// Success: code==0 and status is 2xx
	if tiktokError.Code == 0 && http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	// Consolidated ResponseError
	responseError := ResponseError{
		Status:   r.StatusCode,
		Message:  tiktokError.Message,
		RequstID: tiktokError.RequstID,
		Code:     tiktokError.Code,
	}
	return wrapSpecificError(r, responseError)
}

// func CheckResponseError(r *http.Response) error {
// 	tiktokError := ResponseError{}
//
// 	bodyBytes, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		log.Println("error while io.ReadAll: ", err)
// 		return err
// 	}
//
// 	defer func() {
// 		// already read out, reload for next process
// 		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
// 	}()
//
// 	if len(bodyBytes) > 0 {
// 		err := json.Unmarshal(bodyBytes, &tiktokError)
// 		if err != nil {
// 			return ResponseDecodingError{
// 				Body:    bodyBytes,
// 				Message: err.Error(),
// 				Status:  r.StatusCode,
// 			}
// 		}
// 	}
//
// 	if tiktokError.Code == 0 && http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices {
// 		return nil
// 	}
//
// 	responseError := ResponseError{
// 		Status:   r.StatusCode,
// 		Message:  tiktokError.Message,
// 		RequstID: tiktokError.RequstID,
// 		Code:     tiktokError.Code,
// 	}
//
// 	return wrapSpecificError(r, responseError)
// }

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
