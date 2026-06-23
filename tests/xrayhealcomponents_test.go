//go:build itest

package tests

import (
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

func TestHealComponentsService_Heal_NpmBuildTool_NoChanges(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayHealComponentsTest(t)
	input := `{"lockfileVersion":3}`
	svc := xrayServices.NewComponentsHealService(client)
	svc.XrayDetails = xrayDetails
	svc.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")
	resp, err := svc.Heal(xrayServices.ComponentResolutionRequest{
		BuildTool: "npm",
		Repo:      "npm-virtual",
		Lockfile:  input,
	})
	require.NoError(t, err)
	assert.Equal(t, input, resp.Lockfile)
	assert.Empty(t, resp.Changes)
}

func TestHealComponentsService_Heal_SelfHealDisabled_NoChanges(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayHealComponentsTest(t)
	input := `{"lockfileVersion":3}`
	svc := xrayServices.NewComponentsHealService(client)
	svc.XrayDetails = xrayDetails
	svc.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")
	resp, err := svc.Heal(xrayServices.ComponentResolutionRequest{
		BuildTool: "self-heal-disabled",
		Repo:      "npm-virtual",
		Lockfile:  input,
	})
	require.NoError(t, err)
	assert.Equal(t, input, resp.Lockfile)
	assert.Empty(t, resp.Changes)
}

func TestHealComponentsService_Heal_MavenBuildTool_Changes(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayHealComponentsTest(t)
	inputPom := `<?xml version="1.0"?><project><artifactId>app</artifactId></project>`
	svc := xrayServices.NewComponentsHealService(client)
	svc.XrayDetails = xrayDetails
	svc.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")
	resp, err := svc.Heal(xrayServices.ComponentResolutionRequest{
		BuildTool: "maven",
		Repo:      "maven-virtual",
		Lockfile:  inputPom,
	})
	require.NoError(t, err)
	assert.NotEqual(t, inputPom, resp.Lockfile)
	assert.NotEmpty(t, resp.Changes)
	assert.Contains(t, resp.Lockfile, "<?xml")
	assert.Contains(t, resp.Lockfile, "spring-core")
	assert.Contains(t, resp.Lockfile, "5.3.39-0.cgr.4")
}
