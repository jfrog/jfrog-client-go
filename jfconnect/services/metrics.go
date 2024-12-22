package services

import (
	"encoding/json"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const LogMetricApiEndpoint = "api/v1/backoffice/metrics/log"

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

type Labels struct {
	ProductID                            string `json:"product_id"`
	FeatureID                            string `json:"feature_id"`
	OIDCUsed                             string `json:"oidc_used"`
	JobID                                string `json:"job_id"`
	RunID                                string `json:"run_id"`
	GitRepo                              string `json:"git_repo"`
	GhTokenForCodeScanningAlertsProvided string `json:"gh_token_for_code_scanning_alerts_provided"`
}

type VisibilityMetric struct {
	Value       int    `json:"value"`
	MetricsName string `json:"metrics_name"`
	Labels      Labels `json:"labels"`
}
