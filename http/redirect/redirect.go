package redirect

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

const MaxAuthenticatedRedirects = 3

// EndpointBoundary represents an allowed scheme+host+path prefix that a request (and any
// redirect it follows) must stay within. It has no notion of what the boundary represents -
// callers decide what URL to construct it from.
type EndpointBoundary struct {
	scheme string
	host   string
	path   string
}

func NewEndpointBoundary(rawBaseURL string) (EndpointBoundary, error) {
	u, err := url.Parse(rawBaseURL)
	if err != nil || u.Host == "" || u.User != nil || u.RawPath != "" || u.RawQuery != "" || u.Fragment != "" ||
		(!strings.EqualFold(u.Scheme, "http") && !strings.EqualFold(u.Scheme, "https")) {
		return EndpointBoundary{}, fmt.Errorf("invalid endpoint boundary")
	}
	cleanPath, valid := normalizedBoundaryPath(u.Path)
	if !valid {
		return EndpointBoundary{}, fmt.Errorf("invalid endpoint boundary")
	}
	return EndpointBoundary{
		scheme: strings.ToLower(u.Scheme),
		host:   normalizedHost(u),
		path:   cleanPath,
	}, nil
}

func (b EndpointBoundary) Validate(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" || u.User != nil {
		return fmt.Errorf("invalid request URL")
	}
	if hasAmbiguousPathEscape(u.EscapedPath()) {
		return fmt.Errorf("ambiguous escaped request path")
	}
	requestPath, valid := normalizedBoundaryPath(u.Path)
	if !valid {
		return fmt.Errorf("request URL escapes the configured endpoint")
	}
	if strings.ToLower(u.Scheme) != b.scheme || normalizedHost(u) != b.host ||
		!strings.HasPrefix(requestPath, b.path) {
		return fmt.Errorf("request URL escapes the configured endpoint")
	}
	return nil
}

func normalizedBoundaryPath(rawPath string) (string, bool) {
	if rawPath == "" {
		rawPath = "/"
	}
	cleanPath := strings.TrimSuffix(path.Clean(rawPath), "/") + "/"
	return cleanPath, cleanPath == strings.TrimSuffix(rawPath, "/")+"/"
}

func hasAmbiguousPathEscape(escapedPath string) bool {
	for {
		lowerPath := strings.ToLower(escapedPath)
		if strings.Contains(lowerPath, "%2f") || strings.Contains(lowerPath, "%2e") ||
			strings.Contains(lowerPath, "%5c") {
			return true
		}
		decoded, err := url.PathUnescape(escapedPath)
		if err != nil {
			return true
		}
		if decoded == escapedPath {
			return false
		}
		escapedPath = decoded
	}
}

func normalizedHost(u *url.URL) string {
	port := u.Port()
	if port == "" {
		if strings.EqualFold(u.Scheme, "http") {
			port = "80"
		} else if strings.EqualFold(u.Scheme, "https") {
			port = "443"
		}
	}
	return strings.ToLower(u.Hostname()) + ":" + port
}

// SendWithBoundedRedirects sends requestURL and follows any redirect the server returns, but only
// if each hop (including the initial URL) validates against boundary. Credentials are only ever
// re-attached to a hop once it has been validated. Callers decide what boundary means to them -
// this function only compares URLs against it.
func SendWithBoundedRedirects(client *jfroghttpclient.JfrogHttpClient, method, requestURL string,
	details *httputils.HttpClientDetails, boundary EndpointBoundary, maxRedirects int,
) (*http.Response, []byte, error) {
	// JfrogHttpClient has no per-request redirect hook, so validate each hop before forwarding auth.
	if maxRedirects < 0 {
		return nil, nil, fmt.Errorf("redirect limit must be non-negative")
	}
	// HttpClient.Send retries on CheckRedirect errors, which would desync this hop counter.
	if retries := client.GetHttpClient().GetRetries(); retries != 0 {
		return nil, nil, fmt.Errorf("bounded redirects require a zero-retry client, got %d retries configured", retries)
	}
	currentURL := requestURL
	for redirects := 0; ; redirects++ {
		if err := boundary.Validate(currentURL); err != nil {
			return nil, nil, err
		}
		resp, body, redirectURL, err := client.Send(method, currentURL, nil, false, true, details.Clone(), "")
		if redirectURL == "" {
			return resp, body, err
		}
		if redirects == maxRedirects {
			return resp, body, fmt.Errorf("redirect limit of %d exceeded", maxRedirects)
		}
		if err := boundary.Validate(redirectURL); err != nil {
			return resp, body, fmt.Errorf("unsafe redirect: %w", err)
		}
		currentURL = redirectURL
	}
}
