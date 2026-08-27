package redirect

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newZeroRetryClient(t *testing.T) *jfroghttpclient.JfrogHttpClient {
	client, err := jfroghttpclient.JfrogClientBuilder().SetRetries(0).Build()
	require.NoError(t, err)
	return client
}

func TestNewEndpointBoundary(t *testing.T) {
	t.Run("bare host covers root", func(t *testing.T) {
		boundary, err := NewEndpointBoundary("https://example.com")
		require.NoError(t, err)
		assert.Equal(t, "/", boundary.path)
		assert.NoError(t, boundary.Validate("https://example.com"))
		assert.NoError(t, boundary.Validate("https://example.com/api/pypi/repo/file.whl"))
	})

	t.Run("trailing slash is normalized", func(t *testing.T) {
		withoutSlash, err := NewEndpointBoundary("https://example.com/api/pypi/repo")
		require.NoError(t, err)
		withSlash, err := NewEndpointBoundary("https://example.com/api/pypi/repo/")
		require.NoError(t, err)
		assert.Equal(t, withoutSlash, withSlash)
	})

	t.Run("IPv6 and default port are normalized", func(t *testing.T) {
		boundary, err := NewEndpointBoundary("https://[2001:db8::1]/api/pypi/repo/")
		require.NoError(t, err)
		assert.NoError(t, boundary.Validate("https://[2001:db8::1]:443/api/pypi/repo/pkg.whl"))
		assert.Error(t, boundary.Validate("https://[2001:db8::2]/api/pypi/repo/pkg.whl"))
	})

	for _, rawURL := range []string{
		"https://user@example.com/api/pypi/repo/",
		"https://example.com/api/pypi/repo/?token=value",
		"https://example.com/api/pypi/repo/#fragment",
		"ftp://example.com/api/pypi/repo/",
	} {
		t.Run("rejects invalid boundary "+rawURL, func(t *testing.T) {
			_, err := NewEndpointBoundary(rawURL)
			require.Error(t, err)
		})
	}
}

func TestEndpointBoundaryValidateRejectsEscapes(t *testing.T) {
	boundary, err := NewEndpointBoundary("https://example.com/api/pypi/repo/")
	require.NoError(t, err)

	for _, rawURL := range []string{
		"https://example.com/api/pypi/repo/%2Fsecret",
		"https://example.com/api/pypi/repo/%252Fsecret",
		"https://example.com/api/pypi/repo/%25252e%25252e%25252fsecret",
		"https://example.com/api/pypi/repo/%255csecret",
	} {
		t.Run(rawURL, func(t *testing.T) {
			err := boundary.Validate(rawURL)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "ambiguous escaped request path")
		})
	}
}

func TestEndpointBoundaryValidateRejectsOutsideEndpoint(t *testing.T) {
	boundary, err := NewEndpointBoundary("https://example.com/api/pypi/repo/")
	require.NoError(t, err)

	for _, rawURL := range []string{
		"http://example.com/api/pypi/repo/pkg.whl",
		"https://other.example.com/api/pypi/repo/pkg.whl",
		"https://example.com/api/pypi/other/pkg.whl",
		"https://user@example.com/api/pypi/repo/pkg.whl",
	} {
		t.Run(rawURL, func(t *testing.T) {
			assert.Error(t, boundary.Validate(rawURL))
		})
	}
}

func TestSendWithBoundedRedirectsLimit(t *testing.T) {
	tests := []struct {
		name         string
		redirects    int32
		wantErr      bool
		wantRequests int32
	}{
		{name: "three redirects succeed", redirects: MaxAuthenticatedRedirects, wantRequests: 4},
		{name: "fourth redirect is rejected", redirects: MaxAuthenticatedRedirects + 1, wantErr: true, wantRequests: 4},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestNumber := requests.Add(1)
				if requestNumber <= test.redirects {
					http.Redirect(w, r, fmt.Sprintf("/api/pypi/repo/%d", requestNumber), http.StatusFound)
					return
				}
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte("ok")); err != nil {
					t.Errorf("failed writing response: %v", err)
				}
			}))
			defer server.Close()
			// Zero retries required; see SendWithBoundedRedirects.
			client := newZeroRetryClient(t)
			boundary, err := NewEndpointBoundary(server.URL + "/api/pypi/repo/")
			require.NoError(t, err)
			details := httputils.HttpClientDetails{}

			resp, body, err := SendWithBoundedRedirects(client, http.MethodGet,
				server.URL+"/api/pypi/repo/0", &details, boundary, MaxAuthenticatedRedirects)
			if test.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "redirect limit")
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "ok", string(body))
			}
			assert.Equal(t, test.wantRequests, requests.Load())
		})
	}
}

func TestSendWithBoundedRedirectsRejectsLateHopEscape(t *testing.T) {
	// The first redirect is legitimate; only the second hop escapes the boundary.
	// Every hop must be validated, not just the initial request.
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch requests.Add(1) {
		case 1:
			http.Redirect(w, r, "/api/pypi/repo/inner", http.StatusFound)
		case 2:
			http.Redirect(w, r, "/api/pypi/other-repo/secret", http.StatusFound)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()
	client := newZeroRetryClient(t)
	boundary, err := NewEndpointBoundary(server.URL + "/api/pypi/repo/")
	require.NoError(t, err)
	details := httputils.HttpClientDetails{}

	_, _, err = SendWithBoundedRedirects(client, http.MethodGet,
		server.URL+"/api/pypi/repo/start", &details, boundary, MaxAuthenticatedRedirects)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsafe redirect")
	assert.Equal(t, int32(2), requests.Load(), "must stop at the escaping hop, not follow it")
}
