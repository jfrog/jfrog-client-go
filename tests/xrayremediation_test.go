//go:build itest

package tests

import (
	"strconv"
	"testing"
	"github.com/CycloneDX/cyclonedx-go"
	"github.com/stretchr/testify/require"

	"github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
)

func createTestBOM() *cyclonedx.BOM {
	return &cyclonedx.BOM{
		Components: &[]cyclonedx.Component{
			{
				BOMRef:  "npm://test-component-1@1.0.0",
				Type:    cyclonedx.ComponentTypeLibrary,
				Name:    "test-component-1",
				Version: "1.0.0",
				
			},
			{
				BOMRef:  "npm://test-component-2@2.0.0",
				Type:    cyclonedx.ComponentTypeLibrary,
				Name:    "test-component-2",
				Version: "2.0.0",
			},
		},
		Vulnerabilities: &[]cyclonedx.Vulnerability{
			{
				ID: "CVE-2023-1234",
				Ratings: &[]cyclonedx.VulnerabilityRating{[
					{
						Severity: cyclonedx.SeverityCritical,
					},
				]},
				Affects: &[]cyclonedx.Affects{
					{
						Ref: "npm://test-component-1@1.0.0",
					},
				},
			},
			{
				ID: "CVE-2023-5678",
				Ratings: &[]cyclonedx.VulnerabilityRating{[
					{
						Severity: cyclonedx.SeverityHigh,
					},
				]},
				Affects: &[]cyclonedx.Affects{
					{
						Ref: "npm://test-component-2@2.0.0",
					},
				},
			},
		},
	}
}


func TestCveRemediation(t *testing.T) {
	initXrayTest(t)
	
	xrayServerPort := xray.StartXrayMockServer(t)
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	if err != nil {
		t.Error(err)
	}
	remediationService := services.NewRemediationService(client)
	remediationService.XrayDetails = xrayDetails
	remediationService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/")

	cves := []string{"CVE-2023-1234", "CVE-2023-5678"}

	cvesToRemediationOptions, err := remediationService.RemediationByCve(createTestBOM())
	require.NoError(t, err)
	require.NotNil(t, cvesToRemediationOptions)

	require.Len(t, cvesToRemediationOptions, 2)
	for cve, remediationOptions := range cvesToRemediationOptions {
		require.Contains(t, cves, cve)		
		for _, remediation := range remediationOptions {
			require.NotEmpty(t, remediation.Steps)
		}
	}
}