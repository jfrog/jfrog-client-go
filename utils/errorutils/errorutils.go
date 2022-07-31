package errorutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
// We use body instead of resp.Body beacuse resp.body disappears after resp.Close()
func CheckResponseStatus(resp *http.Response, body []byte, expectedStatusCodes ...int) error {
	for _, statusCode := range expectedStatusCodes {
		if statusCode == resp.StatusCode {
			return nil
		}
	}
	return CheckError(GenerateResponseError(resp.Status, generateErrorString(body)))
}

func GenerateResponseError(status, body string) error {
	return fmt.Errorf("server response: %s\n%s", status, body)
}

func generateErrorString(bodyArray []byte) string {
	var content bytes.Buffer
	if err := json.Indent(&content, bodyArray, "", "  "); err != nil {
		return string(bodyArray)
	}
	return content.String()
}
