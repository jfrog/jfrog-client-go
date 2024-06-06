package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/stretchr/testify/assert"
)

func TestNewIgnoreRulesService(t *testing.T) {
	httpClient := &jfroghttpclient.JfrogHttpClient{}
	service := NewIgnoreRulesService(httpClient)
	assert.NotNil(t, service)
	assert.Equal(t, httpClient, service.client)
}

func TestCheckMinimumVersion(t *testing.T) {
	// Mock XrayDetails
	xrDetails := &mockServiceDetails{}

	// Mock IgnoreRulesService with mocked XrayDetails
	service := &IgnoreRulesService{
		XrayDetails: xrDetails,
	}

	// Test case: Xray version above minimum required
	xrDetails.version = "3.11.1"
	err := service.CheckMinimumVersion()
	assert.NoError(t, err)

	// Test case: Xray version below minimum required
	xrDetails.version = "3.10.9"
	err = service.CheckMinimumVersion()
	assert.Error(t, err)

	// Test case: Missing XrayDetails
	service.XrayDetails = nil
	err = service.CheckMinimumVersion()
	assert.Error(t, err)
}

func TestGetIgnoreRulesURL(t *testing.T) {
	// Mock IgnoreRulesService with mocked XrayDetails
	service := &IgnoreRulesService{
		XrayDetails: &mockServiceDetails{},
	}

	expectedURL := "http://example.com/api/v1/ignore_rules"
	assert.Equal(t, expectedURL, service.getIgnoreRulesURL())
}

func TestGetRuleIdUrl(t *testing.T) {
	// Mock IgnoreRulesService with mocked XrayDetails
	service := &IgnoreRulesService{
		XrayDetails: &mockServiceDetails{},
	}

	expectedURL := "http://example.com/api/v1/ignore_rules/ruleId"
	assert.Equal(t, expectedURL, service.getRuleIdUrl("ruleId"))
}

func TestGetParamMap(t *testing.T) {
	now := time.Now()
	params := &IgnoreRulesGetAllParams{
		Vulnerability:        "vuln",
		License:              "lic",
		Policy:               "pol",
		Watch:                "watch",
		ComponentName:        "compName",
		ComponentVersion:     "compVer",
		ArtifactName:         "artName",
		ArtifactVersion:      "artVer",
		BuildName:            "buildName",
		BuildVersion:         "buildVer",
		ReleaseBundleName:    "rbName",
		ReleaseBundleVersion: "rbVer",
		DockerLayer:          "dockerLayer",
		OrderBy:              "order",
		Direction:            "dir",
		PageNum:              1,
		NumOfRows:            10,
		ExpiresBefore:        now,
		ExpiresAfter:         now,
		ProjectKey:           "projKey",
	}
	paramMap := params.getParamMap()
	assert.NotNil(t, paramMap)
	assert.Equal(t, paramMap["vulnerability"], params.Vulnerability)

	assert.Equal(t, paramMap["license"], params.License)

	assert.Equal(t, paramMap["policy"], params.Policy)

	assert.Equal(t, paramMap["watch"], params.Watch)

	assert.Equal(t, paramMap["component_name"], params.ComponentName)

	assert.Equal(t, paramMap["component_version"], params.ComponentVersion)

	assert.Equal(t, paramMap["artifact_name"], params.ArtifactName)

	assert.Equal(t, paramMap["artifact_version"], params.ArtifactVersion)

	assert.Equal(t, paramMap["build_name"], params.BuildName)

	assert.Equal(t, paramMap["build_version"], params.BuildVersion)

	assert.Equal(t, paramMap["release_bundle_name"], params.ReleaseBundleName)

	assert.Equal(t, paramMap["release_bundle_version"], params.ReleaseBundleVersion)

	assert.Equal(t, paramMap["docker_layer"], params.DockerLayer)

	assert.Equal(t, paramMap["order_by"], params.OrderBy)

	assert.Equal(t, paramMap["direction"], params.Direction)

	assert.Equal(t, paramMap["page_num"], fmt.Sprintf("%d", params.PageNum))

	assert.Equal(t, paramMap["num_of_rows"], fmt.Sprintf("%d", params.NumOfRows))

	assert.Equal(t, paramMap["expires_before"], params.ExpiresBefore.UTC().Format(time.RFC3339))

	assert.Equal(t, paramMap["expires_after"], params.ExpiresAfter.UTC().Format(time.RFC3339))

	assert.Equal(t, paramMap["project_key"], params.ProjectKey)

}

// Mock implementation of auth.ServiceDetails
type mockServiceDetails struct {
	version string
}

func (m *mockServiceDetails) GetUser() string {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetPassword() string {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetApiKey() string {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetAccessToken() string {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetPreRequestFunctions() []auth.ServiceDetailsPreRequestFunc {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetClientCertPath() string {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetClientCertKeyPath() string {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetSshUrl() string {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetSshKeyPath() string {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetSshPassphrase() string {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetSshAuthHeaders() map[string]string {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetClient() *jfroghttpclient.JfrogHttpClient {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetUrl(url string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetUser(user string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetPassword(password string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetApiKey(apiKey string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetAccessToken(accessToken string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) AppendPreRequestFunction(requestFunc auth.ServiceDetailsPreRequestFunc) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetClientCertPath(certificatePath string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetClientCertKeyPath(certificatePath string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetSshUrl(url string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetSshKeyPath(sshKeyPath string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetSshPassphrase(sshPassphrase string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetSshAuthHeaders(sshAuthHeaders map[string]string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetClient(client *jfroghttpclient.JfrogHttpClient) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetDialTimeout(dialTimeout time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) SetOverallRequestTimeout(overallRequestTimeout time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) IsSshAuthHeaderSet() bool {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) IsSshAuthentication() bool {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) AuthenticateSsh(sshKey, sshPassphrase string) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) InitSsh() error {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) RunPreRequestFunctions(httpClientDetails *httputils.HttpClientDetails) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) CreateHttpClientDetails() httputils.HttpClientDetails {
	//TODO implement me
	panic("implement me")
}

func (m *mockServiceDetails) GetVersion() (string, error) {
	return m.version, nil
}

func (m *mockServiceDetails) GetUrl() string {
	return "http://example.com/"
}
