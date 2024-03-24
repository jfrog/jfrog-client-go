package tests

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var testsEventService *services.EventService

func TestXscPostEvent(t *testing.T) {
	xscDetails, client := initXscEventTest(t)
	testsEventService = services.NewEventService(client)
	testsEventService.XscDetails = xscDetails

	event := services.XscGeneralEvent{
		EventType:              1, // ?
		EventStatus:            "started",
		Product:                "cli",
		ProductVersion:         "2.53.1", // add cli version call
		IsDefaultConfig:        false,    // what is this?
		JfrogUser:              "gail",   // add cli user
		OsPlatform:             "mac",    // add
		OsArchitecture:         "arm",    // add
		MachineId:              "",       //?
		AnalyzerManagerVersion: "1.1.1",  //add
		JpdVersion:             "1.5",    //?,
	}
	msi, err := testsEventService.PostEvent(event)
	assert.NoError(t, err)
	assert.True(t, isValidUUID(msi))
}

func isValidUUID(str string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-1[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(str)
}

func initXscEventTest(t *testing.T) (xscDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) {
	var err error
	initXscTest(t)
	xscDetails = GetXscDetails()
	client, err = jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xscDetails.GetClientCertPath()).
		SetClientCertKeyPath(xscDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xscDetails.RunPreRequestFunctions).
		Build()
	assert.NoError(t, err)
	return
}
