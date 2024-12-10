package tests

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetConfigurationProfileByName(t *testing.T) {
	initXscTest(t, services.ConfigProfileMinXscVersion)

	mockServer, configProfileService := createXscMockServerForConfigProfile(t)
	defer mockServer.Close()

	configProfile, err := configProfileService.GetConfigurationProfileByName("default-test-profile")
	assert.NoError(t, err)

	profileFileContent, err := os.ReadFile("testdata/configprofile/configProfileExample.json")
	assert.NoError(t, err)
	var configProfileForComparison services.ConfigProfile
	err = json.Unmarshal(profileFileContent, &configProfileForComparison)
	assert.NoError(t, err)
	assert.Equal(t, &configProfileForComparison, configProfile)
}

func createXscMockServerForConfigProfile(t *testing.T) (mockServer *httptest.Server, configProfileService *services.ConfigurationProfileService) {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/xsc/api/v1/profile/default-test-profile" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			content, err := os.ReadFile("testdata/configprofile/configProfileExample.json")
			assert.NoError(t, err)
			_, err = w.Write(content)
			assert.NoError(t, err)
		} else {
			assert.Fail(t, "received an unexpected request")
		}
	}))

	xscDetails := GetXscDetails()
	xscDetails.SetUrl(mockServer.URL + "/xsc")
	xscDetails.SetAccessToken("")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	configProfileService = services.NewConfigurationProfileService(client)
	configProfileService.XscDetails = xscDetails
	return
}
