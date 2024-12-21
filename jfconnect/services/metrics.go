package services

import (
	"net/http"

	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const LogMetricApiEndpoint = "jfconnect/api/v1/backoffice/metrics/log"

type JfConnectService struct {
	client         *jfroghttpclient.JfrogHttpClient
	serviceDetails *auth.ServiceDetails
}

func NewJfConnectService(serviceDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *JfConnectService {
	return &JfConnectService{serviceDetails: &serviceDetails, client: client}
}

func (jcs *JfConnectService) GetJfConnectDetails() auth.ServiceDetails {
	return *jcs.serviceDetails
}

func (jcs *JfConnectService) PostMetric(metric []byte) error {
	details := jcs.GetJfConnectDetails()
	httpClientDetails := details.CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	url := clientutils.AddTrailingSlashIfNeeded(details.GetUrl())
	url += LogMetricApiEndpoint
	resp, body, err := jcs.client.SendPost(url, metric, &httpClientDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated, http.StatusOK)
}
