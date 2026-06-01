package errorutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Error modes (how should the application behave when the CheckError function is invoked):
type OnErrorHandler func(error) error

var CheckError = func(err error) error {
	return err
}

func CheckErrorf(format string, a ...interface{}) error {
	if len(a) > 0 {
		return CheckError(fmt.Errorf(format, a...))
	}
	return CheckError(errors.New(format))
}

// HttpResponseError is returned by CheckResponseStatus and CheckResponseStatusWithBody
// when the server returns an unexpected status code. It preserves the raw response so
// callers can produce structured error output (e.g. JSON to stderr) via errors.As,
// while Error() still produces the legacy "server response: ..." string for
// backward compatibility with existing string-matching callers and tests.
type HttpResponseError struct {
	StatusCode int
	Status     string
	Body       []byte
}

func (e *HttpResponseError) Error() string {
	msg := "server response: " + e.Status
	if len(e.Body) > 0 {
		msg += "\n" + GenerateErrorString(e.Body)
	}
	return msg
}

// Check expected status codes and return error if needed
func CheckResponseStatus(resp *http.Response, expectedStatusCodes ...int) error {
	for _, statusCode := range expectedStatusCodes {
		if statusCode == resp.StatusCode {
			return nil
		}
	}
	// Add resp.Body to error response if exists
	errorBody, _ := io.ReadAll(resp.Body)
	return CheckError(&HttpResponseError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       errorBody,
	})
}

// Check expected status codes and return error with body if needed
// We use body variable that was saved outside the resp.body object,
// Instead of resp.Body because resp.body disappears after resp.body.Close()
func CheckResponseStatusWithBody(resp *http.Response, body []byte, expectedStatusCodes ...int) error {
	for _, statusCode := range expectedStatusCodes {
		if statusCode == resp.StatusCode {
			return nil
		}
	}
	return CheckError(&HttpResponseError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       body,
	})
}

func GenerateResponseError(status, body string) error {
	responseErrString := "server response: " + status
	if body != "" {
		responseErrString = responseErrString + "\n" + body
	}
	return errors.New(responseErrString)
}

func GenerateErrorString(bodyArray []byte) string {
	var content bytes.Buffer
	if err := json.Indent(&content, bodyArray, "", "  "); err != nil {
		return string(bodyArray)
	}
	return content.String()
}
