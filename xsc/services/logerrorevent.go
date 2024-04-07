package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const (
	postLogErrorAPI       = "api/v1/event/logMessage"
	LogErrorMinXscVersion = AnalyticsMetricsMinXscVersion
)

type LogErrorEventService struct {
	client     *jfroghttpclient.JfrogHttpClient
	XscDetails auth.ServiceDetails
}

type ExternalErrorLog struct {
	Log_level string `json:"log_level"`
	Source    string `json:"source"`
	Message   string `json:"message"`
}

func NewLogErrorEventService(client *jfroghttpclient.JfrogHttpClient) *LogErrorEventService {
	return &LogErrorEventService{client: client}
}

func (les *LogErrorEventService) SendLogErrorEvent(errorLog *ExternalErrorLog) error {
	httpClientDetails := les.XscDetails.CreateHttpClientDetails()
	requestContent, err := json.Marshal(errorLog)
	if err != nil {
		return fmt.Errorf("failed to convert POST request body's struct into JSON: %q", err)
	}
	url := utils.AddTrailingSlashIfNeeded(les.XscDetails.GetUrl()) + postLogErrorAPI
	response, body, err := les.client.SendPost(url, requestContent, &httpClientDetails)
	if err != nil {
		return fmt.Errorf("failed to send POST query to '%s': %s", url, err.Error())
	}
	return errorutils.CheckResponseStatusWithBody(response, body, http.StatusCreated)
}
