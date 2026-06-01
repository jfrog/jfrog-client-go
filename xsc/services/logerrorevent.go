package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	xscutils "github.com/jfrog/jfrog-client-go/xsc/services/utils"
	"net/http"
)

const (
	xscLogErrorApiSuffix           = "event/logMessage"
	xscDeprecatedLogErrorApiSuffix = "api/v1/" + xscLogErrorApiSuffix
	LogErrorMinXscVersion          = AnalyticsMetricsMinXscVersion
)

type LogErrorEventService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XscDetails      auth.ServiceDetails
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

type ExternalErrorLog struct {
	Log_level string `json:"log_level"`
	Source    string `json:"source"`
	Message   string `json:"message"`
}

func NewLogErrorEventService(client *jfroghttpclient.JfrogHttpClient) *LogErrorEventService {
	return &LogErrorEventService{client: client}
}

func (les *LogErrorEventService) sendLogErrorRequest(requestContent []byte) (url string, resp *http.Response, body []byte, err error) {
	if les.XrayDetails != nil {
		httpClientDetails := les.XrayDetails.CreateHttpClientDetails()
		url = utils.AddTrailingSlashIfNeeded(les.XrayDetails.GetUrl()) + xscutils.XscInXraySuffix + xscLogErrorApiSuffix
		resp, body, err = les.client.SendPost(utils.AppendScopedProjectKeyParam(url, les.ScopeProjectKey), requestContent, &httpClientDetails)
		return
	}
	// Backward compatibility
	httpClientDetails := les.XscDetails.CreateHttpClientDetails()
	url = utils.AddTrailingSlashIfNeeded(les.XscDetails.GetUrl()) + xscDeprecatedLogErrorApiSuffix
	resp, body, err = les.client.SendPost(utils.AppendScopedProjectKeyParam(url, les.ScopeProjectKey), requestContent, &httpClientDetails)
	return
}

func (les *LogErrorEventService) SendLogErrorEvent(errorLog *ExternalErrorLog) error {
	requestContent, err := json.Marshal(errorLog)
	if err != nil {
		return fmt.Errorf("failed to convert POST request body's struct into JSON: %q", err)
	}
	url, response, body, err := les.sendLogErrorRequest(requestContent)
	if err != nil {
		return fmt.Errorf("failed to send POST query to '%s': %s", url, err.Error())
	}
	return errorutils.CheckResponseStatusWithBody(response, body, http.StatusCreated)
}
