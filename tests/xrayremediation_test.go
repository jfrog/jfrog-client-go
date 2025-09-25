//go:build itest

package tests

import (
	"github.com/CycloneDX/cyclonedx-go"
	"github.com/stretchr/testify/require"
)

func createTestBOM() *cyclonedx.BOM {
	return &cyclonedx.BOM{
		Components: &[]cyclonedx.Component{
			{
				Type:    cyclonedx.ComponentTypeLibrary,
				Name:    "test-component-1",
				Version: "1.0.0",
			},
			{
				Type:    cyclonedx.ComponentTypeLibrary,
				Name:    "test-component-2",
				Version: "2.0.0",
			},
		},
	}
}


func TestCveRemediation(t *testing.T) {
	initXrayTest(t)
	bom := createTestBOM()
	cves := []string{"CVE-2023-1234", "CVE-2023-5678"}
	response, err := xrayManager.CveRemediation(bom, cves...)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.Remediations, 2)
	for _, remediation := range response.Remediations {
		require.Contains(t, cves, remediation.CVE)
		require.NotEmpty(t, remediation.FixVersions)
		require.NotEmpty(t, remediation.Component)
	}
}

func TestArtifactRemediation(t *testing.T) {
	initXrayTest(t)
	bom := createTestBOM()
	response, err := xrayManager.ArtifactRemediation(bom)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.Remediations, 2)
	for _, remediation := range response.Remediations {
		require.NotEmpty(t, remediation.FixVersions)
		require.NotEmpty(t, remediation.Component)
	}
}