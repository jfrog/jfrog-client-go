package tests

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const errorMessageContentForTest = "THIS IS NOT A REAL ERROR! This Error is posted as part of TestXscSendLogErrorEvent test"

func TestXscSendLogErrorEvent(t *testing.T) {
	initXscTest(t, services.LogErrorMinXscVersion, "")
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
		if r.RequestURI == "/xsc/api/v1/event/logMessage" && r.Method == http.MethodPost {
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

	xscDetails := GetXscDetails()
	xscDetails.SetUrl(mockServer.URL + "/xsc")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	logErrorService = services.NewLogErrorEventService(client)
	logErrorService.XscDetails = xscDetails
	return
}
