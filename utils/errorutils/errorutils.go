package errorutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	return CheckError(GenerateResponseError(resp.Status, generateErrorString(resp.Body)))
}

func GenerateResponseError(status, body string) error {
	return fmt.Errorf("server response: %s\n%s", status, body)
}

func generateErrorString(body io.ReadCloser) string {
	bodyArray, err := ioutil.ReadAll(body)
	if err != nil {
		return ""
	}
	var content bytes.Buffer
	if err = json.Indent(&content, bodyArray, "", "  "); err != nil {
		return string(bodyArray)
	} else {
		return content.String()
	}
}
