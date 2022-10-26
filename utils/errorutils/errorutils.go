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

// Check expected status codes and return error if needed
func CheckResponseStatus(resp *http.Response, expectedStatusCodes ...int) error {
	for _, statusCode := range expectedStatusCodes {
		if statusCode == resp.StatusCode {
			return nil
		}
	}
	// Add resp.Body to error response if exists
	errorBody, _ := io.ReadAll(resp.Body)
	return CheckError(GenerateResponseError(resp.Status, string(errorBody)))
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
	return CheckError(GenerateResponseError(resp.Status, generateErrorString(body)))
}

func GenerateResponseError(status, body string) error {
	responseErrString := "server response: " + status
	if body != "" {
		responseErrString = responseErrString + "\n" + body
	}
	return fmt.Errorf(responseErrString)
}

func generateErrorString(bodyArray []byte) string {
	var content bytes.Buffer
	if err := json.Indent(&content, bodyArray, "", "  "); err != nil {
		return string(bodyArray)
	}
	return content.String()
}
