//go:build itest

package tests

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
	"github.com/jfrog/jfrog-client-go/xsc/services/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const errorMessageContentForTest = "THIS IS NOT A REAL ERROR! This Error is posted as part of TestXscSendLogErrorEvent test"

func TestXscSendLogErrorEvent(t *testing.T) {
	initXscTest(t, services.LogErrorMinXscVersion, utils.MinXrayVersionXscTransitionToXray)
	mockServer, logErrorService := createXscMockServerForLogEvent(t)
	defer mockServer.Close()

	event := &services.ExternalErrorLog{
		Log_level: "error",
		Source:    "cli",
		Message:   errorMessageContentForTest,
	}

	assert.NoError(t, logErrorService.SendLogErrorEvent(event))
}

func createXscMockServerForLogEvent(t *testing.T) (mockServer *httptest.Server, logErrorService *services.LogErrorEventService) {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/xray/api/v1/xsc/event/logMessage" && r.Method == http.MethodPost {
			var reqBody services.ExternalErrorLog
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&reqBody)
			assert.NoError(t, err, "Invalid JSON request body")
			if err != nil {
				return
			}

			assert.Equal(t, "error", reqBody.Log_level)
			assert.Equal(t, "cli", reqBody.Source)
			assert.Equal(t, errorMessageContentForTest, reqBody.Message)
			w.WriteHeader(http.StatusCreated)
			return
		} else {
			assert.Fail(t, "received an unexpected request")
		}
	}))

	xrayDetails := GetXrayDetails()
	xrayDetails.SetUrl(mockServer.URL + "/xray")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	logErrorService = services.NewLogErrorEventService(client)
	logErrorService.XrayDetails = xrayDetails
	return
}
