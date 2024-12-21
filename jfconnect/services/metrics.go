package services

import (
	"net/http"

	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

type JfConnectService struct {
	client         *jfroghttpclient.JfrogHttpClient
	serviceDetails *auth.ServiceDetails
}

func NewJfConnectService(serviceDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *JfConnectService {
	return &JfConnectService{serviceDetails: &serviceDetails, client: client}
}

func (js *JfConnectService) GetJfConnectDetails() auth.ServiceDetails {
	return *js.serviceDetails
}

func (js *JfConnectService) PostMetric(metric []byte) error {
	details := js.GetJfConnectDetails()
	httpClientDetails := details.CreateHttpClientDetails()
	rtUtils.SetContentType("application/json", &httpClientDetails.Headers)

	url := clientutils.AddTrailingSlashIfNeeded(details.GetUrl())
	url += "jfconnect/api/v1/backoffice/metrics/log"
	resp, body, err := js.client.SendPost(url, metric, &httpClientDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated, http.StatusOK)
}
