package errorutils

import (
	"errors"
	"github.com/ztrue/tracerr"
	"io/ioutil"
	"net/http"
)

// Max stacktrace including the top frame (This file)
const maxStackTraceSize int = 4

// Use this function to allow showing the stacktrace after an error.
func WrapError(err error) error {
	if err == nil {
		return nil
	}
	return wrapError(tracerr.Wrap(err))
}

// Use this function to allow showing the stacktrace after an error.
func NewError(message string) error {
	return wrapError(tracerr.New(message))
}

func wrapError(err tracerr.Error) error {
	stackTrace := err.StackTrace()
	stackSize := len(stackTrace)
	if maxStackTraceSize < stackSize {
		stackSize = maxStackTraceSize // Limit the stacktrace size to 3.
	}
	// Remove the first frame in order to not trace this file in the stacktrace.
	return tracerr.CustomError(err, stackTrace[1:stackSize])
}

// Check expected status codes and return error if needed
func CheckResponseStatus(resp *http.Response, expectedStatusCodes ...int) error {
	for _, statusCode := range expectedStatusCodes {
		if statusCode == resp.StatusCode {
			return nil
		}
	}

	errorBody, _ := ioutil.ReadAll(resp.Body)
	return errors.New(resp.Status + " " + string(errorBody))
}
