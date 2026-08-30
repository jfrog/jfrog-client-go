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

func initXrayZeroTouchRemediationTest(t *testing.T) (xrayServerPort int, xrayDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) {
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

func TestZeroTouchRemediationService_Remediate_NpmBuildTool_NoChanges(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayZeroTouchRemediationTest(t)
	input := `{"lockfileVersion":3}`
	svc := xrayServices.NewZeroTouchRemediationService(client)
	svc.XrayDetails = xrayDetails
	svc.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")
	resp, disabled, err := svc.Remediate(xrayServices.ComponentResolutionRequest{
		BuildTool: "npm",
		Repo:      "npm-virtual",
		Lockfile:  input,
	})
	require.NoError(t, err)
	assert.False(t, disabled)
	assert.Equal(t, input, resp.Lockfile)
	assert.Empty(t, resp.Changes)
}

func TestZeroTouchRemediationService_Remediate_Disabled_NoChanges(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayZeroTouchRemediationTest(t)
	input := `{"lockfileVersion":3}`
	svc := xrayServices.NewZeroTouchRemediationService(client)
	svc.XrayDetails = xrayDetails
	svc.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")
	resp, disabled, err := svc.Remediate(xrayServices.ComponentResolutionRequest{
		BuildTool: "ztr-disabled",
		Repo:      "npm-virtual",
		Lockfile:  input,
	})
	require.NoError(t, err)
	assert.True(t, disabled)
	assert.Equal(t, input, resp.Lockfile)
	assert.Empty(t, resp.Changes)
}

func TestZeroTouchRemediationService_Remediate_NpmBuildTool_Changes(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayZeroTouchRemediationTest(t)
	input := `{"name":"xray-simple-npm-app","lockfileVersion":3,"packages":{}}`
	svc := xrayServices.NewZeroTouchRemediationService(client)
	svc.XrayDetails = xrayDetails
	svc.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")
	resp, disabled, err := svc.Remediate(xrayServices.ComponentResolutionRequest{
		BuildTool: "npm",
		Repo:      "npm-virtual",
		Lockfile:  input,
	})
	require.NoError(t, err)
	assert.False(t, disabled)
	assert.NotEqual(t, input, resp.Lockfile)
	assert.NotEmpty(t, resp.Changes)
	assert.Contains(t, resp.Lockfile, "lodash")
	assert.Equal(t, "lodash", resp.Changes[0].Package)
}
