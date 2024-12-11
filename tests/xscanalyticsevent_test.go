package tests

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var testsEventService *services.AnalyticsEventService

func TestXscAddAndUpdateGeneralEvent(t *testing.T) {
	xscDetails, client := initXscEventTest(t)
	testsEventService = services.NewAnalyticsEventService(client)
	testsEventService.XscDetails = xscDetails

	event := services.XscAnalyticsGeneralEvent{XscAnalyticsBasicGeneralEvent: services.XscAnalyticsBasicGeneralEvent{
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
	msi, err := testsEventService.AddGeneralEvent(event)
	assert.NoError(t, err)
	assert.True(t, isValidUUID(msi))

	// Validate that the event sent and saved properly in XSC.
	resp, err := testsEventService.GetGeneralEvent(msi)
	assert.NoError(t, err)
	assert.Equal(t, event, *resp)

	finalizeEvent := services.XscAnalyticsGeneralEventFinalize{
		MultiScanId: msi,
		XscAnalyticsBasicGeneralEvent: services.XscAnalyticsBasicGeneralEvent{
			EventStatus:          services.Completed,
			TotalFindings:        10,
			TotalIgnoredFindings: 5,
			TotalScanDuration:    "15s",
		},
	}

	err = testsEventService.UpdateGeneralEvent(finalizeEvent)
	assert.NoError(t, err)

	// Validate that the event's update sent and saved properly in XSC.
	resp, err = testsEventService.GetGeneralEvent(msi)
	assert.NoError(t, err)
	event.EventStatus = services.Completed
	event.TotalFindings = 10
	event.TotalIgnoredFindings = 5
	event.TotalScanDuration = "15s"
	assert.Equal(t, event, *resp)
}

func isValidUUID(str string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-1[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(str)
}

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
