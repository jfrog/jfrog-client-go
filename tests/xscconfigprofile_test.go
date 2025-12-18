//go:build itest

package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
	xscutils "github.com/jfrog/jfrog-client-go/xsc/services/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const configProfileWithoutRepo = "default-test-profile"

func TestGetConfigurationProfileByName(t *testing.T) {
	initXscTest(t, services.ConfigProfileMinXscVersion, xscutils.MinXrayVersionXscTransitionToXray)

	xrayVersion, err := GetXrayDetails().GetVersion()
	require.NoError(t, err)

	mockServer, configProfileService := createXscMockServerForConfigProfile(t, xrayVersion)
	defer mockServer.Close()

	configProfile, err := configProfileService.GetConfigurationProfileByName(configProfileWithoutRepo)
	assert.NoError(t, err)
	assert.Equal(t, getComparisonConfigProfile(), configProfile)
}

func TestGetConfigurationProfileByUrl(t *testing.T) {
	initXscTest(t, "", services.ConfigProfileByUrlMinXrayVersion)

	xrayVersion, err := GetXrayDetails().GetVersion()
	require.NoError(t, err)

	mockServer, configProfileService := createXscMockServerForConfigProfile(t, xrayVersion)
	defer mockServer.Close()

	configProfile, err := configProfileService.GetConfigurationProfileByUrl(mockServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, getComparisonConfigProfile(), configProfile)

}

func getComparisonConfigProfile() *services.ConfigProfile {
	return &services.ConfigProfile{
		ProfileName: "default-profile",
		GeneralConfig: services.GeneralConfig{
			ScannersDownloadPath:    "https://repo.example.com/releases",
			GeneralExcludePatterns:  []string{"*.log*", "*.tmp*"},
			FailUponAnyScannerError: true,
		},
		FrogbotConfig: services.FrogbotConfig{
			AggregateFixes:                      true,
			HideSuccessBannerForNoIssues:        false,
			BranchNameTemplate:                  "frogbot-${IMPACTED_PACKAGE}-${BRANCH_NAME_HASH}",
			PrTitleTemplate:                     "[🐸 Frogbot] Upgrade {IMPACTED_PACKAGE} to {FIX_VERSION}",
			CommitMessageTemplate:               "Upgrade {IMPACTED_PACKAGE} to {FIX_VERSION}",
			ShowSecretsAsPrComment:              false,
			CreateAutoFixPr:                     true,
			IncludeVulnerabilitiesAndViolations: false,
		},
		Modules: []services.Module{
			{
				ModuleName:   "default-module",
				PathFromRoot: ".",
				ScanConfig: services.ScanConfig{
					ScaScannerConfig: services.ScaScannerConfig{
						EnableScaScan:   true,
						ExcludePatterns: []string{"**/build/**"},
					},
					ContextualAnalysisScannerConfig: services.CaScannerConfig{
						EnableCaScan:    true,
						ExcludePatterns: []string{"**/docs/**"},
					},
					SastScannerConfig: services.SastScannerConfig{
						EnableSastScan:  true,
						ExcludePatterns: []string{"**/_test.go/**"},
						ExcludeRules:    []string{"xss-injection"},
					},
					SecretsScannerConfig: services.SecretsScannerConfig{
						EnableSecretsScan:   true,
						ValidateSecrets:     true,
						ExcludePatterns:     []string{"**/_test.go/**"},
						EnableCustomSecrets: true,
					},
					IacScannerConfig: services.IacScannerConfig{
						EnableIacScan:   true,
						ExcludePatterns: []string{"*.tfstate"},
					},
				},
			},
		},
	}
}

// TODO backwards compatability can be removed once old Xsc service is removed from all servers
func createXscMockServerForConfigProfile(t *testing.T, xrayVersion string) (mockServer *httptest.Server, configProfileService *services.ConfigurationProfileService) {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiUrlPart := "api/v1/"
		var isXrayAfterXscMigration bool
		if isXrayAfterXscMigration = xscutils.IsXscXrayInnerService(xrayVersion); isXrayAfterXscMigration {
			apiUrlPart = ""
		}

		switch {
		case (strings.Contains(r.RequestURI, "/xsc/"+apiUrlPart+"profile/"+configProfileWithoutRepo) && r.Method == http.MethodGet) ||
			strings.Contains(r.RequestURI, "xray/api/v1/xsc/profile_repos") && r.Method == http.MethodPost && isXrayAfterXscMigration:
			w.WriteHeader(http.StatusOK)
			content, err := os.ReadFile("testdata/configprofile/configProfileExample.json")
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
