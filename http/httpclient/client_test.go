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

func TestDefaultTransportEnablesConnectionReuse(t *testing.T) {
	client, err := ClientBuilder().Build()
	assert.NoError(t, err)

	transport, ok := client.GetClient().Transport.(*http.Transport)
	if !ok {
		t.Skip("transport is not *http.Transport (e.g. custom certs); skipping connection reuse config check")
	}

	assert.True(t, transport.ForceAttemptHTTP2, "ForceAttemptHTTP2 should be true for HTTP/2 support")
	assert.Equal(t, 6, transport.MaxIdleConnsPerHost, "MaxIdleConnsPerHost should be 6 for connection reuse to the same host")
}
