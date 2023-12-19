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
	PreRetryInterceptors  []PreRetryInterceptor
}

type PreRetryInterceptor func() (shouldRetry bool)

func (httpClientDetails HttpClientDetails) Clone() *HttpClientDetails {
	headers := make(map[string]string)
	utils.MergeMaps(httpClientDetails.Headers, headers)
	var transport *http.Transport
	if httpClientDetails.Transport != nil {
		transport = httpClientDetails.Transport.Clone()
	}
	return &HttpClientDetails{
		User:                  httpClientDetails.User,
		Password:              httpClientDetails.Password,
		ApiKey:                httpClientDetails.ApiKey,
		AccessToken:           httpClientDetails.AccessToken,
		Headers:               headers,
		Transport:             transport,
		DialTimeout:           httpClientDetails.DialTimeout,
		OverallRequestTimeout: httpClientDetails.OverallRequestTimeout,
		PreRetryInterceptors:  httpClientDetails.PreRetryInterceptors,
	}
}

func (httpClientDetails *HttpClientDetails) AddPreRetryInterceptor(preRetryInterceptors PreRetryInterceptor) {
	httpClientDetails.PreRetryInterceptors = append(httpClientDetails.PreRetryInterceptors, preRetryInterceptors)
}
