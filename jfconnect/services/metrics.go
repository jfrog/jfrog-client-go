package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const LogMetricApiEndpoint = "api/v1/backoffice/metrics/log"

type VisibilityMetric struct {
	Value  int    `json:"value"`
	Name   string `json:"metrics_name"`
	Labels any    `json:"labels"`
}

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

func (jcs *JfConnectService) PostVisibilityMetric(metric VisibilityMetric) error {
	metricJson, err := json.Marshal(metric)
	if err != nil {
		return errorutils.CheckError(err)
	}
	details := jcs.GetJfConnectDetails()
	httpClientDetails := details.CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	url := clientutils.AddTrailingSlashIfNeeded(details.GetUrl())
	url += LogMetricApiEndpoint
	resp, body, err := jcs.client.SendPost(url, metricJson, &httpClientDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated, http.StatusOK)
}
