package httputils

import (
	"github.com/jfrog/jfrog-client-go/utils"
	"net/http"
	"time"
)

type HttpClientDetails struct {
	User        string
	Password    string
	AccessToken string
	Headers     map[string]string
	Transport   *http.Transport
	HttpTimeout time.Duration
}

func (httpClientDetails HttpClientDetails) Clone() *HttpClientDetails {
	headers := make(map[string]string)
	utils.MergeMaps(httpClientDetails.Headers, headers)
	return &HttpClientDetails{
		User:        httpClientDetails.User,
		Password:    httpClientDetails.Password,
		AccessToken: httpClientDetails.AccessToken,
		Headers:     headers}
}
