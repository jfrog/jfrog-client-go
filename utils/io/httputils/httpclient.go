package httputils

import (
	"net/http"
	"time"

	"github.com/jfrog/jfrog-client-go/utils"
)

type HttpClientDetails struct {
	User                  string
	Password              string
	ApiKey                string
	AccessToken           string
	Headers               map[string]string
	Transport             *http.Transport
	DialTimeout           time.Duration
	OverallRequestTimeout time.Duration
	// Prior to each retry attempt, the list of PreRetryInterceptors is invoked sequentially. If any of these interceptors yields a 'false' response, the retry process stops instantly.
	PreRetryInterceptors []PreRetryInterceptor
}

type PreRetryInterceptor func() (shouldRetry bool)

func (hcd HttpClientDetails) Clone() *HttpClientDetails {
	headers := make(map[string]string)
	utils.MergeMaps(hcd.Headers, headers)
	var transport *http.Transport
	if hcd.Transport != nil {
		transport = hcd.Transport.Clone()
	}
	return &HttpClientDetails{
		User:                  hcd.User,
		Password:              hcd.Password,
		ApiKey:                hcd.ApiKey,
		AccessToken:           hcd.AccessToken,
		Headers:               headers,
		Transport:             transport,
		DialTimeout:           hcd.DialTimeout,
		OverallRequestTimeout: hcd.OverallRequestTimeout,
		PreRetryInterceptors:  hcd.PreRetryInterceptors,
	}
}

func (hcd *HttpClientDetails) AddPreRetryInterceptor(preRetryInterceptors PreRetryInterceptor) {
	hcd.PreRetryInterceptors = append(hcd.PreRetryInterceptors, preRetryInterceptors)
}

func (hcd *HttpClientDetails) SetContentTypeApplicationJson() {
	hcd.AddHeader("Content-Type", "application/json")
}

func (hcd *HttpClientDetails) AddHeader(headerName, headerValue string) {
	if hcd.Headers == nil {
		hcd.Headers = make(map[string]string)
	}
	hcd.Headers[headerName] = headerValue
}
