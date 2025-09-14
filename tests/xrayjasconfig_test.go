//go:build itest

package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	xrayServices "github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func initXrayJasConfigTest(t *testing.T) (xrayServerPort int, xrayDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) {
	var err error
	initXrayTest(t)
	xrayServerPort = xray.StartXrayMockServer(t)
	xrayDetails = GetXrayDetails()
	client, err = jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	assert.NoError(t, err)
	return
}

func TestIsTokenValidationEnabled(t *testing.T) {
	xrayServerPort, xrayDetails, client := initXrayJasConfigTest(t)
	testsJasConfigService := xrayServices.NewJasConfigService(client)
	testsJasConfigService.XrayDetails = xrayDetails
	testsJasConfigService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")

	result, err := testsJasConfigService.GetJasConfigTokenValidation()
	assert.NoError(t, err)
	assert.Equal(t, result, true)
}
