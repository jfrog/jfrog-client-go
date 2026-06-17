//go:build itest

package tests

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	xrayServices "github.com/jfrog/jfrog-client-go/xray/services"
)

func initXrayHealComponentsTest(t *testing.T) (xrayServerPort int, xrayDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) {
	var err error
	initXrayTest(t)
	xrayServerPort = xray.StartXrayMockServer(t)
	xrayDetails = GetXrayDetails()
	client, err = jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	require.NoError(t, err)
	return
}


func TestHealComponentsService_Heal_NpmBuildTool(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayHealComponentsTest(t)
	input := json.RawMessage(`{"lockfileVersion":3}`)
	svc := xrayServices.NewComponentsHealService(client)
	svc.XrayDetails = xrayDetails
	svc.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")
	resp, err := svc.Heal(xrayServices.ComponentResolutionRequest{
		BuildTool: "npm",
		Repo:      "npm-virtual",
		Lockfile:  input,
	})
	require.NoError(t, err)
	assert.Equal(t, input, resp.Content)
	assert.Empty(t, resp.Changes)
}

func TestHealComponentsService_Heal_MavenBuildTool(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayHealComponentsTest(t)
	input := json.RawMessage(`{"lockFileVersion":1,"groupId":"demo","artifactId":"app"}`)
	svc := xrayServices.NewComponentsHealService(client)
	svc.XrayDetails = xrayDetails
	svc.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")
	resp, err := svc.Heal(xrayServices.ComponentResolutionRequest{
		BuildTool: "maven",
		Repo:      "maven-virtual",
		Lockfile:  input,
	})
	require.NoError(t, err)
	assert.NotEqual(t, input, resp.Content)
	assert.NotEmpty(t, resp.Changes)
	var healed map[string]any
	require.NoError(t, json.Unmarshal(resp.Content, &healed))
	assert.Equal(t, float64(1), healed["lockFileVersion"])
}
