//go:build itest

package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
	"github.com/jfrog/jfrog-client-go/xsc/services/utils"
)

const TestMultiScanId = "3472b4e2-bddc-11ee-a9c9-acde48001122"
const TestMultiScanIdResponse = "{\"multi_scan_id\": \"" + TestMultiScanId + "\"}"

var initialEvent = services.XscAnalyticsGeneralEvent{XscAnalyticsBasicGeneralEvent: services.XscAnalyticsBasicGeneralEvent{
	EventType:              services.CliEventType,
	EventStatus:            services.Started,
	Product:                services.CliProduct,
	ProductVersion:         "2.53.1",
	IsDefaultConfig:        false,
	JfrogUser:              "gail",
	OsPlatform:             "mac",
	OsArchitecture:         "arm64",
	MachineId:              "id",
	AnalyzerManagerVersion: "1.1.1",
	ProjectPath:            "/path/to/project",
}}

var finalEvent = services.XscAnalyticsGeneralEvent{XscAnalyticsBasicGeneralEvent: services.XscAnalyticsBasicGeneralEvent{
	EventType:              services.CliEventType,
	EventStatus:            services.Completed,
	Product:                services.CliProduct,
	ProductVersion:         "2.53.1",
	IsDefaultConfig:        false,
	JfrogUser:              "gail",
	OsPlatform:             "mac",
	OsArchitecture:         "arm64",
	MachineId:              "id",
	AnalyzerManagerVersion: "1.1.1",
	TotalFindings:          10,
	TotalIgnoredFindings:   5,
	TotalScanDuration:      "15s",
	ProjectPath:            "/path/to/project",
}}

func TestXscAddAndUpdateGeneralEvent(t *testing.T) {
	initXscTest(t, services.LogErrorMinXscVersion, utils.MinXrayVersionXscTransitionToXray)

	testCases := []struct {
		name        string
		xrayVersion string
	}{
		{
			name:        "Xray version with deprecated AddGeneralEvent",
			xrayVersion: "3.115.0",
		},
		{
			name:        "Xray version with new AddGeneralEvent",
			xrayVersion: "3.116.0",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer, analyticsService := createXscMockServerForGeneralEvent(t)
			defer mockServer.Close()

			msi, err := analyticsService.AddGeneralEvent(initialEvent, tc.xrayVersion)
			assert.NoError(t, err)

			// Validate that the event sent and saved properly in XSC.
			resp, err := analyticsService.GetGeneralEvent(msi)
			require.NoError(t, err)
			assert.Equal(t, initialEvent, *resp)

			finalizeEvent := services.XscAnalyticsGeneralEventFinalize{
				MultiScanId: msi,
				XscAnalyticsBasicGeneralEvent: services.XscAnalyticsBasicGeneralEvent{
					EventStatus:          services.Completed,
					TotalFindings:        10,
					TotalIgnoredFindings: 5,
					TotalScanDuration:    "15s",
				},
			}

			err = analyticsService.UpdateGeneralEvent(finalizeEvent)
			assert.NoError(t, err)

			// Validate that the event's update sent and saved properly in XSC.
			// We add suffix to the msi to enable the mock server to differentiate between the initial response to the final response
			resp, err = analyticsService.GetGeneralEvent(msi + "-final")
			assert.NoError(t, err)
			assert.Equal(t, finalEvent, *resp)
		})
	}
}

func createXscMockServerForGeneralEvent(t *testing.T) (mockServer *httptest.Server, analyticsService *services.AnalyticsEventService) {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.RequestURI, "/xray/api/v1/xsc/event") && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte(TestMultiScanIdResponse))
			assert.NoError(t, err)
		case r.RequestURI == "/xray/api/v1/xsc/event/"+TestMultiScanId && r.Method == http.MethodGet:
			// This is the first GET for the even before update
			w.WriteHeader(http.StatusOK)
			eventJson, err := json.Marshal(initialEvent)
			assert.NoError(t, err)
			_, err = w.Write(eventJson)
			assert.NoError(t, err)
		case strings.Contains(r.RequestURI, "/xray/api/v1/xsc/event") && r.Method == http.MethodPut:
			w.WriteHeader(http.StatusOK)
		case r.RequestURI == "/xray/api/v1/xsc/event/"+TestMultiScanId+"-final" && r.Method == http.MethodGet:
			// This is the second GET after Updating the event
			w.WriteHeader(http.StatusOK)
			eventJson, err := json.Marshal(finalEvent)
			assert.NoError(t, err)
			_, err = w.Write(eventJson)
			assert.NoError(t, err)
		default:
			assert.Fail(t, "received an unexpected request")
		}
	}))

	xrayDetails := GetXrayDetails()
	xrayDetails.SetUrl(mockServer.URL + "/xray")
	xrayDetails.SetAccessToken("")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	analyticsService = services.NewAnalyticsEventService(client)
	analyticsService.XrayDetails = xrayDetails
	return
}

func TestXscSendGitIntegrationEvent(t *testing.T) {
	initXscTest(t, services.LogErrorMinXscVersion, utils.MinXrayVersionGitIntegrationEvent)

	testCases := []struct {
		name        string
		xrayVersion string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Xray version below minimum",
			xrayVersion: "3.134.0",
			expectError: true,
			errorMsg:    "git integration event version error",
		},
		{
			name:        "Xray version at minimum",
			xrayVersion: "3.135.0",
			expectError: false,
		},
		{
			name:        "Xray version above minimum",
			xrayVersion: "3.136.0",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer, analyticsService := createXscMockServerForGitIntegrationEvent(t)
			defer mockServer.Close()

			event := services.GitIntegrationEvent{
				EventType:     "Source Code SBOM Results Upload",
				GitProvider:   "github",
				GitOwner:      "jfrog",
				GitRepository: "jfrog-cli-security",
				GitBranch:     "main",
				EventStatus:   "completed",
			}

			err := analyticsService.SendGitIntegrationEvent(event, tc.xrayVersion)

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func createXscMockServerForGitIntegrationEvent(t *testing.T) (mockServer *httptest.Server, analyticsService *services.AnalyticsEventService) {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.RequestURI, "/xray/api/v1/xsc/git_integration_event") && r.Method == http.MethodPost:
			// Validate request body
			var receivedEvent services.GitIntegrationEvent
			err := json.NewDecoder(r.Body).Decode(&receivedEvent)
			assert.NoError(t, err)

			// Validate required fields
			assert.NotEmpty(t, receivedEvent.EventType)
			assert.NotEmpty(t, receivedEvent.GitProvider)
			assert.NotEmpty(t, receivedEvent.GitOwner)
			assert.NotEmpty(t, receivedEvent.GitRepository)
			assert.NotEmpty(t, receivedEvent.GitBranch)
			assert.NotEmpty(t, receivedEvent.EventStatus)

			w.WriteHeader(http.StatusCreated)
		default:
			assert.Fail(t, "received an unexpected request: "+r.Method+" "+r.RequestURI)
		}
	}))

	xrayDetails := GetXrayDetails()
	xrayDetails.SetUrl(mockServer.URL + "/xray")
	xrayDetails.SetAccessToken("")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	analyticsService = services.NewAnalyticsEventService(client)
	analyticsService.XrayDetails = xrayDetails
	return
}
