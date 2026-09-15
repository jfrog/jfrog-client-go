package httpclient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestSendPostFromReader(t *testing.T) {
	payload := "streamed-body"

	tests := []struct {
		name           string
		body           io.Reader
		details        httputils.HttpClientDetails
		handler        http.HandlerFunc
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			body: io.NopCloser(strings.NewReader(payload)),
			details: httputils.HttpClientDetails{
				AccessToken: "token",
				Headers:     map[string]string{"Content-Type": "application/octet-stream"},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "application/octet-stream", r.Header.Get("Content-Type"))
				assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.Equal(t, payload, string(body))
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write([]byte(`{"ok":true}`))
			},
			expectedStatus: http.StatusAccepted,
			expectedBody:   `{"ok":true}`,
		},
		{
			name: "does not fail on error status",
			body: strings.NewReader("x"),
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("bad request"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			t.Cleanup(server.Close)

			httpClient, err := ClientBuilder().Build()
			require.NoError(t, err)

			resp, body, err := httpClient.SendPostFromReader(server.URL, tt.body, tt.details)
			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, string(body))
		})
	}
}
