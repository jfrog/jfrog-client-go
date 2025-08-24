package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	DefaultSonarURL = "https://api.sonarcloud.io"
)

// Response models for SonarQube API endpoints
type QualityGatesAnalysis struct {
	ProjectStatus struct {
		Status     string `json:"status"`
		Conditions []struct {
			Status         string `json:"status"`
			MetricKey      string `json:"metricKey"`
			Comparator     string `json:"comparator"`
			PeriodIndex    int    `json:"periodIndex"`
			ErrorThreshold string `json:"errorThreshold"`
			ActualValue    string `json:"actualValue"`
		} `json:"conditions"`
		Periods []struct {
			Index int    `json:"index"`
			Mode  string `json:"mode"`
			Date  string `json:"date"`
		} `json:"periods"`
		IgnoredConditions bool `json:"ignoredConditions"`
	} `json:"projectStatus"`
}

type TaskDetails struct {
	Task struct {
		ID                 string      `json:"id"`
		Type               string      `json:"type"`
		ComponentID        string      `json:"componentId"`
		ComponentKey       string      `json:"componentKey"`
		ComponentName      string      `json:"componentName"`
		ComponentQualifier string      `json:"componentQualifier"`
		AnalysisID         string      `json:"analysisId"`
		Status             string      `json:"status"`
		SubmittedAt        string      `json:"submittedAt"`
		StartedAt          string      `json:"startedAt"`
		ExecutedAt         string      `json:"executedAt"`
		ExecutionTimeMs    int         `json:"executionTimeMs"`
		Logs               interface{} `json:"logs"`
		HasScannerContext  bool        `json:"hasScannerContext"`
		Organization       string      `json:"organization"`
	} `json:"task"`
}

type Service interface {
	GetQualityGateAnalysis(analysisID string) (*QualityGatesAnalysis, error)
	GetTaskDetails(ceTaskID string) (*TaskDetails, error)
	GetSonarIntotoStatementRaw(ceTaskID string) ([]byte, error)
}

type sonarService struct {
	client         *jfroghttpclient.JfrogHttpClient
	serviceDetails auth.ServiceDetails
}

func NewSonarService(serviceDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) Service {
	return &sonarService{serviceDetails: serviceDetails, client: client}
}

func (s *sonarService) GetSonarDetails() auth.ServiceDetails {
	return s.serviceDetails
}

func (s *sonarService) buildAuthHeader() string {
	// 1) Access token
	token := s.serviceDetails.GetAccessToken()
	if token != "" {
		return "Bearer " + token
	}
	// 2) API key
	if apiKey := s.serviceDetails.GetApiKey(); apiKey != "" {
		return "Bearer " + apiKey
	}
	// 3) Basic auth
	user := s.serviceDetails.GetUser()
	pass := s.serviceDetails.GetPassword()
	if user != "" && pass != "" {
		creds := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
		return "Basic " + creds
	}
	// 4) Env token fallback
	if val := os.Getenv("SONAR_TOKEN"); val != "" {
		return "Bearer " + val
	}
	if val := os.Getenv("SONARQUBE_TOKEN"); val != "" {
		return "Bearer " + val
	}
	return ""
}

func (s *sonarService) httpGetJSON(urlStr string) ([]byte, int, error) {
	httpClientDetails := s.GetSonarDetails().CreateHttpClientDetails()
	httpClientDetails.Headers["Authorization"] = s.buildAuthHeader()

	log.Debug("HTTP GET", urlStr)
	resp, body, _, err := s.client.SendGet(urlStr, true, &httpClientDetails)
	if err != nil {
		log.Debug("HTTP GET error for", urlStr, "error:", err.Error())
		return nil, 0, err
	}
	log.Debug("HTTP GET response for", urlStr, "status:", resp.StatusCode, "body:", string(body))
	return body, resp.StatusCode, nil
}

// GetSonarIntotoStatementRaw returns the raw JSON bytes of the in-toto statement.
// We return []byte instead of a typed object to avoid extra marshal/unmarshal cycles and
// to allow callers to augment the statement (e.g., add subject/stage) and sign it as-is.
func (s *sonarService) GetSonarIntotoStatementRaw(ceTaskID string) ([]byte, error) {
	if ceTaskID == "" {
		return nil, errorutils.CheckError(fmt.Errorf("missing ce task id for enterprise endpoint"))
	}
	sonarBaseURL := s.GetSonarDetails().GetUrl()
	if sonarBaseURL == "" {
		sonarBaseURL = DefaultSonarURL
	}
	baseURL := strings.TrimRight(sonarBaseURL, "/")
	u, _ := url.Parse(baseURL)
	hostname := u.Hostname()
	if hostname != "localhost" && net.ParseIP(hostname) == nil && !strings.HasPrefix(hostname, "api.") {
		baseURL = strings.Replace(baseURL, "://", "://api.", 1)
	}
	enterpriseURL := fmt.Sprintf("%s/dop-translation/jfrog-evidence/%s", baseURL, url.QueryEscape(ceTaskID))
	body, statusCode, err := s.httpGetJSON(enterpriseURL)
	if err != nil {
		return nil, errorutils.CheckErrorf("enterprise endpoint failed with status %d and response %s %v", statusCode, string(body), err)
	}
	if statusCode != 200 {
		return nil, errorutils.CheckErrorf("enterprise endpoint returned status %d: %s", statusCode, string(body))
	}
	return body, nil
}

func (s *sonarService) GetQualityGateAnalysis(analysisID string) (*QualityGatesAnalysis, error) {
	if analysisID == "" {
		return nil, errorutils.CheckError(fmt.Errorf("missing analysis id for quality gates endpoint"))
	}

	sonarBaseURL := s.GetSonarDetails().GetUrl()
	if sonarBaseURL == "" {
		sonarBaseURL = DefaultSonarURL
	}

	qualityGatesURL := fmt.Sprintf("%s/api/qualitygates/project_status?analysisId=%s", strings.TrimRight(sonarBaseURL, "/"), url.QueryEscape(analysisID))

	body, statusCode, err := s.httpGetJSON(qualityGatesURL)
	if err != nil {
		return nil, errorutils.CheckErrorf("quality gates endpoint failed with status %d: %v", statusCode, err)
	}
	if statusCode != 200 {
		return nil, errorutils.CheckErrorf("quality gates endpoint returned status %d: %s", statusCode, string(body))
	}

	var response QualityGatesAnalysis
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errorutils.CheckErrorf("failed to parse quality gates response: %v", err)
	}

	return &response, nil
}

func (s *sonarService) GetTaskDetails(ceTaskID string) (*TaskDetails, error) {
	if ceTaskID == "" {
		return nil, nil
	}

	sonarBaseURL := s.GetSonarDetails().GetUrl()
	if sonarBaseURL == "" {
		sonarBaseURL = DefaultSonarURL
	}

	taskURL := fmt.Sprintf("%s/api/ce/task?id=%s", strings.TrimRight(sonarBaseURL, "/"), url.QueryEscape(ceTaskID))
	body, statusCode, err := s.httpGetJSON(taskURL)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, errorutils.CheckErrorf("task endpoint returned status %d: %s", statusCode, string(body))
	}

	var response TaskDetails
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errorutils.CheckErrorf("failed to parse task response: %v", err)
	}

	return &response, nil
}
