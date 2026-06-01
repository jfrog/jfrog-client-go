package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jfrog/jfrog-client-go/apptrust/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
)

const mockApplicationKey = "test-app-key"

var mockApplication = Application{
	ApplicationName: "Test Application",
	ApplicationKey:  "test-app-key",
	ProjectName:     "Test Project",
	ProjectKey:      "test-proj",
	Criticality:     "high",
	MaturityLevel:   "production",
}

func TestApplicationService_GetApplicationDetails_Success(t *testing.T) {
	handlerFunc, requestNum := createApplicationHandlerFunc(t, http.StatusOK, mockApplication)

	mockServer, applicationService := createMockApplicationServer(t, handlerFunc)
	defer mockServer.Close()

	application, err := applicationService.GetApplicationDetails(mockApplicationKey)
	assert.NoError(t, err)
	assert.NotNil(t, application)
	assert.Equal(t, "Test Application", application.ApplicationName)
	assert.Equal(t, "test-app-key", application.ApplicationKey)
	assert.Equal(t, "Test Project", application.ProjectName)
	assert.Equal(t, "test-proj", application.ProjectKey)
	assert.Equal(t, "high", application.Criticality)
	assert.Equal(t, "production", application.MaturityLevel)
	assert.Equal(t, 1, *requestNum)
}

func TestApplicationService_GetApplicationDetails_NotFound(t *testing.T) {
	handlerFunc, requestNum := createApplicationHandlerFunc(t, http.StatusNotFound, Application{})

	mockServer, applicationService := createMockApplicationServer(t, handlerFunc)
	defer mockServer.Close()

	application, err := applicationService.GetApplicationDetails("non-existent-key")
	assert.Error(t, err)
	assert.Nil(t, application)
	assert.Equal(t, 1, *requestNum)
}

func TestApplicationService_GetApplicationDetails_BadRequest(t *testing.T) {
	handlerFunc, requestNum := createApplicationHandlerFunc(t, http.StatusBadRequest, Application{})

	mockServer, applicationService := createMockApplicationServer(t, handlerFunc)
	defer mockServer.Close()

	application, err := applicationService.GetApplicationDetails("invalid-key")
	assert.Error(t, err)
	assert.Nil(t, application)
	assert.Equal(t, 1, *requestNum)
}

func TestApplicationService_GetApplicationDetails_Unauthorized(t *testing.T) {
	handlerFunc, requestNum := createApplicationHandlerFunc(t, http.StatusUnauthorized, Application{})

	mockServer, applicationService := createMockApplicationServer(t, handlerFunc)
	defer mockServer.Close()

	application, err := applicationService.GetApplicationDetails(mockApplicationKey)
	assert.Error(t, err)
	assert.Nil(t, application)
	assert.Equal(t, 1, *requestNum)
}

func TestApplicationService_GetApplicationDetails_Forbidden(t *testing.T) {
	handlerFunc, requestNum := createApplicationHandlerFunc(t, http.StatusForbidden, Application{})

	mockServer, applicationService := createMockApplicationServer(t, handlerFunc)
	defer mockServer.Close()

	application, err := applicationService.GetApplicationDetails(mockApplicationKey)
	assert.Error(t, err)
	assert.Nil(t, application)
	assert.Equal(t, 1, *requestNum)
}

func TestApplicationService_GetApplicationDetails_InvalidJSON(t *testing.T) {
	requestNum := 0
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		expectedURI := "/api/v1/applications/" + mockApplicationKey
		if r.RequestURI == expectedURI {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			w.WriteHeader(http.StatusOK)
			requestNum++
			// Write invalid JSON
			_, err := w.Write([]byte("invalid json"))
			assert.NoError(t, err)
		}
	}

	mockServer, applicationService := createMockApplicationServer(t, handlerFunc)
	defer mockServer.Close()

	application, err := applicationService.GetApplicationDetails(mockApplicationKey)
	assert.Error(t, err)
	assert.Nil(t, application)
	assert.Equal(t, 1, requestNum)
}

func TestApplicationService_NewApplicationService(t *testing.T) {
	apptrustDetails := auth.NewApptrustDetails()
	apptrustDetails.SetUrl("http://localhost:8081/")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	service := NewApplicationService(apptrustDetails, client)
	assert.NotNil(t, service)
	assert.Equal(t, apptrustDetails, service.GetApptrustDetails())
}

func TestApplicationService_GetApptrustDetails(t *testing.T) {
	apptrustDetails := auth.NewApptrustDetails()
	apptrustDetails.SetUrl("http://localhost:8081/")
	apptrustDetails.SetUser("testuser")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	service := NewApplicationService(apptrustDetails, client)
	retrievedDetails := service.GetApptrustDetails()

	assert.Equal(t, "http://localhost:8081/", retrievedDetails.GetUrl())
	assert.Equal(t, "testuser", retrievedDetails.GetUser())
}

func createMockApplicationServer(t *testing.T, testHandler http.HandlerFunc) (*httptest.Server, *ApplicationService) {
	testServer := httptest.NewServer(testHandler)

	apptrustDetails := auth.NewApptrustDetails()
	apptrustDetails.SetUrl(testServer.URL + "/")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	return testServer, NewApplicationService(apptrustDetails, client)
}

func createApplicationHandlerFunc(t *testing.T, statusCode int, response Application) (http.HandlerFunc, *int) {
	requestNum := 0
	return func(w http.ResponseWriter, r *http.Request) {
		expectedURI := "/api/v1/applications/" + mockApplicationKey
		if r.RequestURI == expectedURI || r.RequestURI == "/api/v1/applications/non-existent-key" || r.RequestURI == "/api/v1/applications/invalid-key" {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			w.WriteHeader(statusCode)
			requestNum++

			if statusCode == http.StatusOK {
				writeMockApplicationResponse(t, w, response)
			}
		}
	}, &requestNum
}

func writeMockApplicationResponse(t *testing.T, w http.ResponseWriter, response Application) {
	content, err := json.Marshal(response)
	assert.NoError(t, err)
	_, err = w.Write(content)
	assert.NoError(t, err)
}
