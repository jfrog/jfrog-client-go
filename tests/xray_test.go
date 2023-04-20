package tests

import (
	"strconv"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xray/services"
)

var testsXrayEntitlementsService *services.EntitlementsService

func TestXrayVersion(t *testing.T) {
	initXrayTest(t)
	version, err := GetXrayDetails().GetVersion()
	if err != nil {
		t.Error(err)
	}
	if version == "" {
		t.Error("Expected a version, got empty string")
	}
}

func TestXrayEntitlementsService(t *testing.T) {
	initXrayTest(t)
	xrayServerPort := xray.StartXrayMockServer()
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	if err != nil {
		t.Error(err)
	}
	testsXrayEntitlementsService = services.NewEntitlementsService(client)
	testsXrayEntitlementsService.XrayDetails = xrayDetails
	testsXrayEntitlementsService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/")

	// Run tests
	tests := []struct {
		name      string
		featureId string
		expected  bool
	}{
		{name: "userEntitled", featureId: xray.ContextualAnalysisFeatureId, expected: true},
		{name: "userNotEntitled", featureId: xray.BadFeatureId, expected: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testEntitlements(t, test.featureId, test.expected)
		})
	}
}

func testEntitlements(t *testing.T, featureId string, expected bool) {
	result, err := testsXrayEntitlementsService.IsEntitled(featureId)
	if err != nil {
		t.Error(err)
	}
	if result != expected {
		t.Error("Expected:", expected, "Got: ", result)
	}
}

func initXrayTest(t *testing.T) {
	if !*TestXray {
		t.Skip("Skipping xray test. To run xray test add the '-test.xray=true' option.")
	}
}
