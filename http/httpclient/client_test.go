package httpclient

import (
	"net/http"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/stretchr/testify/assert"
)

var shouldRetryCases = []struct {
	name                string
	status              int
	expectedRetry       bool
	preRetryInterceptor httputils.PreRetryInterceptor
}{
	// Status 200
	{"200", http.StatusOK, false, nil},
	{"200 with interceptor returning false", http.StatusOK, false, func() bool { return false }},
	{"200 with interceptor returning true", http.StatusOK, false, func() bool { return true }},

	// Status 502
	{"502", http.StatusBadGateway, true, nil},
	{"429", http.StatusTooManyRequests, true, nil},
	{"502 with interceptor returning false", http.StatusBadGateway, false, func() bool { return false }},
	{"502 with interceptor returning true", http.StatusBadGateway, true, func() bool { return true }},
}

func TestShouldRetry(t *testing.T) {
	httpClient, err := ClientBuilder().Build()
	assert.NoError(t, err)

	for _, testCase := range shouldRetryCases {
		t.Run(testCase.name, func(t *testing.T) {
			httpClientsDetails := &httputils.HttpClientDetails{}
			if testCase.preRetryInterceptor != nil {
				httpClientsDetails.AddPreRetryInterceptor(testCase.preRetryInterceptor)
			}
			shouldRetry := httpClient.shouldRetry(&http.Response{StatusCode: testCase.status}, httpClientsDetails)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedRetry, shouldRetry)
		})
	}
}
