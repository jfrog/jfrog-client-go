package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	sonarauth "github.com/jfrog/jfrog-client-go/sonar/auth"
	"github.com/stretchr/testify/assert"
)

// Local copy for unmarshalling expectations
type localStatement struct {
	Type          string                 `json:"_type"`
	PredicateType string                 `json:"predicateType"`
	Predicate     map[string]interface{} `json:"predicate"`
	CreatedAt     string                 `json:"createdAt"`
	CreatedBy     string                 `json:"createdBy"`
	Markdown      string                 `json:"markdown"`
}

func TestNewSonarService(t *testing.T) {
	serviceDetails := sonarauth.NewSonarDetails()
	serviceDetails.SetUrl("https://sonarcloud.io")
	serviceDetails.SetAccessToken("test-token")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	service := NewSonarService(serviceDetails, client)
	assert.NotNil(t, service)
	assert.IsType(t, &sonarService{}, service)
}

func TestSonarService_BuildAuthHeader(t *testing.T) {
	os.Unsetenv("SONAR_TOKEN")
	tests := []struct {
		name           string
		accessToken    string
		apiKey         string
		user           string
		password       string
		envToken       string
		expectedPrefix string
	}{
		{
			name:           "Access token takes precedence",
			accessToken:    "access-token",
			apiKey:         "api-key",
			user:           "user",
			password:       "password",
			envToken:       "",
			expectedPrefix: "Bearer access-token",
		},
		{
			name:           "API key used when no access token",
			accessToken:    "",
			apiKey:         "api-key",
			user:           "user",
			password:       "password",
			envToken:       "",
			expectedPrefix: "Bearer api-key",
		},
		{
			name:           "Basic auth used when no token",
			accessToken:    "",
			apiKey:         "",
			user:           "user",
			password:       "password",
			envToken:       "",
			expectedPrefix: "Basic dXNlcjpwYXNzd29yZA==",
		},
		{
			name:           "Environment token used when no service details",
			accessToken:    "",
			apiKey:         "",
			user:           "",
			password:       "",
			envToken:       "env-token",
			expectedPrefix: "Bearer env-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceDetails := sonarauth.NewSonarDetails()
			serviceDetails.SetAccessToken(tt.accessToken)
			serviceDetails.SetApiKey(tt.apiKey)
			serviceDetails.SetUser(tt.user)
			serviceDetails.SetPassword(tt.password)

			client, err := jfroghttpclient.JfrogClientBuilder().Build()
			assert.NoError(t, err)

			service := &sonarService{
				client:         client,
				serviceDetails: serviceDetails,
			}

			if tt.envToken != "" {
				os.Setenv("SONAR_TOKEN", tt.envToken)
				defer os.Unsetenv("SONAR_TOKEN")
			} else {
				os.Unsetenv("SONAR_TOKEN")
			}

			authHeader := service.buildAuthHeader()
			assert.Contains(t, authHeader, tt.expectedPrefix)
		})
	}
}

func TestSonarService_GetEnterpriseResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/dop-translation/jfrog-evidence/")

		response := localStatement{
			Type:          "test-type",
			PredicateType: "test-type",
			Predicate:     map[string]interface{}{"test": "data"},
			Markdown:      "test markdown",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	serviceDetails := sonarauth.NewSonarDetails()
	serviceDetails.SetUrl(server.URL)
	serviceDetails.SetAccessToken("test-token")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	service := NewSonarService(serviceDetails, client)

	body, err := service.GetSonarIntotoStatementRaw("test-task-id")
	assert.NoError(t, err)
	assert.NotNil(t, body)

	var result localStatement
	assert.NoError(t, json.Unmarshal(body, &result))
	assert.Equal(t, "test-type", result.PredicateType)
	assert.Equal(t, "test markdown", result.Markdown)
	assert.Equal(t, "data", result.Predicate["test"])
}

func TestSonarService_GetEnterpriseResponse_EmptyTaskID(t *testing.T) {
	serviceDetails := sonarauth.NewSonarDetails()
	serviceDetails.SetUrl("https://sonarcloud.io")
	serviceDetails.SetAccessToken("test-token")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	service := NewSonarService(serviceDetails, client)

	body, err := service.GetSonarIntotoStatementRaw("")
	assert.Error(t, err)
	assert.Nil(t, body)
	assert.Contains(t, err.Error(), "missing ce task id")
}

func TestSonarService_GetQualityGatesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/api/qualitygates/project_status")
		assert.Contains(t, r.URL.RawQuery, "analysisId=test-analysis-id")

		response := QualityGatesAnalysis{
			ProjectStatus: struct {
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
			}{
				Status: "OK",
				Conditions: []struct {
					Status         string `json:"status"`
					MetricKey      string `json:"metricKey"`
					Comparator     string `json:"comparator"`
					PeriodIndex    int    `json:"periodIndex"`
					ErrorThreshold string `json:"errorThreshold"`
					ActualValue    string `json:"actualValue"`
				}{
					{
						Status:         "OK",
						MetricKey:      "new_reliability_rating",
						Comparator:     "GT",
						PeriodIndex:    1,
						ErrorThreshold: "1",
						ActualValue:    "1",
					},
				},
				Periods: []struct {
					Index int    `json:"index"`
					Mode  string `json:"mode"`
					Date  string `json:"date"`
				}{
					{
						Index: 1,
						Mode:  "previous_version",
						Date:  "2025-01-15T10:30:00+0000",
					},
				},
				IgnoredConditions: false,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	serviceDetails := sonarauth.NewSonarDetails()
	serviceDetails.SetUrl(server.URL)
	serviceDetails.SetAccessToken("test-token")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	service := NewSonarService(serviceDetails, client)

	result, err := service.GetQualityGateAnalysis("test-analysis-id")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "OK", result.ProjectStatus.Status)
	assert.Len(t, result.ProjectStatus.Conditions, 1)
	assert.Equal(t, "new_reliability_rating", result.ProjectStatus.Conditions[0].MetricKey)
}

func TestSonarService_GetQualityGatesResponse_EmptyAnalysisID(t *testing.T) {
	serviceDetails := sonarauth.NewSonarDetails()
	serviceDetails.SetUrl("https://sonarcloud.io")
	serviceDetails.SetAccessToken("test-token")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	service := NewSonarService(serviceDetails, client)

	result, err := service.GetQualityGateAnalysis("")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "missing analysis id")
}

func TestSonarService_GetTaskResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/api/ce/task")
		assert.Contains(t, r.URL.RawQuery, "id=test-task-id")

		response := TaskDetails{
			Task: struct {
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
			}{
				ID:           "test-task-id",
				ComponentKey: "test-project",
				AnalysisID:   "test-analysis-id",
				Status:       "SUCCESS",
				Logs:         nil,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	serviceDetails := sonarauth.NewSonarDetails()
	serviceDetails.SetUrl(server.URL)
	serviceDetails.SetAccessToken("test-token")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	service := NewSonarService(serviceDetails, client)

	result, err := service.GetTaskDetails("test-task-id")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-task-id", result.Task.ID)
	assert.Equal(t, "test-project", result.Task.ComponentKey)
	assert.Equal(t, "test-analysis-id", result.Task.AnalysisID)
}

func TestSonarService_GetTaskResponse_EmptyTaskID(t *testing.T) {
	serviceDetails := sonarauth.NewSonarDetails()
	serviceDetails.SetUrl("https://sonarcloud.io")
	serviceDetails.SetAccessToken("test-token")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	service := NewSonarService(serviceDetails, client)

	result, err := service.GetTaskDetails("")
	assert.NoError(t, err)
	assert.Nil(t, result)
}
