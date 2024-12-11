package tests

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
	xscutils "github.com/jfrog/jfrog-client-go/xsc/services/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const configProfileWithoutRepo = "default-test-profile"

func TestGetConfigurationProfileByName(t *testing.T) {
	initXscTest(t, services.ConfigProfileMinXscVersion, "")

	xrayVersion, err := GetXrayDetails().GetVersion()
	require.NoError(t, err)

	mockServer, configProfileService := createXscMockServerForConfigProfile(t, xrayVersion)
	defer mockServer.Close()

	configProfile, err := configProfileService.GetConfigurationProfileByName(configProfileWithoutRepo)
	assert.NoError(t, err)

	profileFileContent, err := os.ReadFile("testdata/configprofile/configProfileExample.json")
	assert.NoError(t, err)
	var configProfileForComparison services.ConfigProfile
	err = json.Unmarshal(profileFileContent, &configProfileForComparison)
	assert.NoError(t, err)
	assert.Equal(t, &configProfileForComparison, configProfile)
}

func TestGetConfigurationProfileByUrl(t *testing.T) {
	initXscTest(t, "", services.ConfigProfileByUrlMinXrayVersion)

	xrayVersion, err := GetXrayDetails().GetVersion()
	require.NoError(t, err)

	mockServer, configProfileService := createXscMockServerForConfigProfile(t, xrayVersion)
	defer mockServer.Close()

	configProfile, err := configProfileService.GetConfigurationProfileByUrl(mockServer.URL)
	assert.NoError(t, err)

	profileFileContent, err := os.ReadFile("testdata/configprofile/configProfileWithRepoExample.json")
	assert.NoError(t, err)
	var configProfileForComparison services.ConfigProfile
	err = json.Unmarshal(profileFileContent, &configProfileForComparison)
	assert.NoError(t, err)
	assert.Equal(t, &configProfileForComparison, configProfile)

}

func createXscMockServerForConfigProfile(t *testing.T, xrayVersion string) (mockServer *httptest.Server, configProfileService *services.ConfigurationProfileService) {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiUrlPart := "api/v1/"
		var isXrayAfterXscMigration bool
		if isXrayAfterXscMigration = xscutils.IsXscXrayInnerService(xrayVersion); isXrayAfterXscMigration {
			apiUrlPart = ""
		}

		switch {
		case strings.Contains(r.RequestURI, "/xsc/"+apiUrlPart+"profile/"+configProfileWithoutRepo):
			w.WriteHeader(http.StatusOK)
			content, err := os.ReadFile("testdata/configprofile/configProfileExample.json")
			assert.NoError(t, err)
			_, err = w.Write(content)
			assert.NoError(t, err)

		case strings.Contains(r.RequestURI, "xray/api/v1/xsc/profile_repos") && isXrayAfterXscMigration:
			w.WriteHeader(http.StatusOK)
			content, err := os.ReadFile("testdata/configprofile/configProfileWithRepoExample.json")
			assert.NoError(t, err)
			_, err = w.Write(content)
			assert.NoError(t, err)
		default:
			assert.Fail(t, "received an unexpected request")
		}
	}))

	xscDetails := GetXscDetails()
	xscDetails.SetUrl(mockServer.URL + "/xsc")
	xscDetails.SetAccessToken("")

	xrayDetails := GetXrayDetails()
	xrayDetails.SetUrl(mockServer.URL + "/xray")
	xrayDetails.SetAccessToken("")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	configProfileService = services.NewConfigurationProfileService(client)
	configProfileService.XscDetails = xscDetails
	configProfileService.XrayDetails = xrayDetails
	return
}
