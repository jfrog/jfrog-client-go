package tests

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
	"github.com/jfrog/jfrog-client-go/xsc/services/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
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
}}

func TestXscAddAndUpdateGeneralEvent(t *testing.T) { // TODO eran this is skipped - V
	// xscDetails, client := initXscEventTest(t)
	//testsEventService := services.NewAnalyticsEventService(client)
	//testsEventService.XscDetails = xscDetails

	initXscTest(t, services.LogErrorMinXscVersion, utils.MinXrayVersionXscTransitionToXray)
	mockServer, analyticsService := createXscMockServerForGeneralEvent(t)
	defer mockServer.Close()

	//event := initialEvent
	//msi, err := testsEventService.AddGeneralEvent(event)
	msi, err := analyticsService.AddGeneralEvent(initialEvent)
	assert.NoError(t, err)
	//assert.True(t, isValidUUID(msi))

	// Validate that the event sent and saved properly in XSC.
	//resp, err := testsEventService.GetGeneralEvent(msi)
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

	//err = testsEventService.UpdateGeneralEvent(finalizeEvent)
	err = analyticsService.UpdateGeneralEvent(finalizeEvent)
	assert.NoError(t, err)

	// Validate that the event's update sent and saved properly in XSC.
	// We add suffix to the msi to enable the mock server to differentiate between the initial response to the final response
	resp, err = analyticsService.GetGeneralEvent(msi + "-final")
	assert.NoError(t, err)
	assert.Equal(t, finalEvent, *resp)
	/*
		resp, err = testsEventService.GetGeneralEvent(msi)
		assert.NoError(t, err)
		event.EventStatus = services.Completed
		event.TotalFindings = 10
		event.TotalIgnoredFindings = 5
		event.TotalScanDuration = "15s"
		assert.Equal(t, event, *resp)

	*/
}

func isValidUUID(str string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-1[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(str)
}

/*
func initXscEventTest(t *testing.T) (xscDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) {
	var err error
	initXscTest(t, services.AnalyticsMetricsMinXscVersion, "")
	xscDetails = GetXscDetails()
	client, err = jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xscDetails.GetClientCertPath()).
		SetClientCertKeyPath(xscDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xscDetails.RunPreRequestFunctions).
		Build()
	assert.NoError(t, err)
	return
}

*/

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
