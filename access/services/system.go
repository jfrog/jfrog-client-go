package services

import (
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

const (
	rootCertificateApi = "api/v1/cert/root"
	readinessApi       = "api/v1/system/readiness"
)

type SystemService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

func NewSystemService(client *jfroghttpclient.JfrogHttpClient) *SystemService {
	return &SystemService{client: client}
}

func (ss *SystemService) GetRootCertificate() (string, error) {
	anonymous := httputils.HttpClientDetails{}
	resp, body, _, err := ss.client.SendGet(ss.ServiceDetails.GetUrl()+rootCertificateApi, true, &anonymous)
	if err != nil {
		return "", err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return "", errorutils.CheckErrorf("GET %s: %s", rootCertificateApi, err.Error())
	}
	return strings.TrimSpace(string(body)), nil
}

func (ss *SystemService) CheckReadiness() (int, error) {
	httpDetails := ss.ServiceDetails.CreateHttpClientDetails()
	resp, _, _, err := ss.client.SendGet(ss.ServiceDetails.GetUrl()+readinessApi, true, &httpDetails)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}
