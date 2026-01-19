//go:build itest

package tests

import (
	"strconv"
	"testing"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/catalog"
	"github.com/jfrog/jfrog-client-go/catalog/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

func initCatalogTest(t *testing.T) {
	if !*TestCatalog {
		t.Skip("Skipping catalog test. To run catalog test add the '-test.catalog=true' option.")
	}
}

func initCatalogEnrichTest(t *testing.T, params catalog.MockServerParams) (catalogServerPort int, catalogDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) {
	var err error
	initCatalogTest(t)
	catalogServerPort = catalog.StartCatalogMockServerWithParams(t, params)
	catalogDetails = GetXrayDetails()
	client, err = jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(catalogDetails.GetClientCertPath()).
		SetClientCertKeyPath(catalogDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(catalogDetails.RunPreRequestFunctions).
		Build()
	require.NoError(t, err)
	return
}

func TestGetVersionSucceeded(t *testing.T) {
	testCase := []struct {
		name    string
		success bool
	}{
		{
			name:    "Get Version Succeeded",
			success: true,
		},
		{
			name: "Get Version Not Available",
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			catalogServerPort, catalogDetails, client := initCatalogEnrichTest(t, catalog.MockServerParams{Alive: tc.success})
			testsVersionService := services.NewVersionService(client)
			testsVersionService.CatalogDetails = catalogDetails
			testsVersionService.CatalogDetails.SetUrl("http://localhost:" + strconv.Itoa(catalogServerPort) + "/catalog/")

			version, err := testsVersionService.GetVersion()
			if tc.success {
				assert.NoError(t, err)
				assert.NotEmpty(t, version)
			} else {
				assert.Error(t, err)
				assert.Equal(t, services.CatalogMinVersionForEnrichApi, version)
			}
		})
	}
}

func TestEnrichSucceeded(t *testing.T) {
	vulnerabilities := []cyclonedx.Vulnerability{
		{BOMRef: "CVE-2021-1234", Source: &cyclonedx.Source{Name: "NVD"}},
		{BOMRef: "CVE-2021-5678", Source: &cyclonedx.Source{Name: "NVD"}},
	}
	catalogServerPort, catalogDetails, client := initCatalogEnrichTest(t, catalog.MockServerParams{Alive: true, EnrichedVuln: vulnerabilities})
	testsEnrichService := services.NewEnrichService(client)
	testsEnrichService.CatalogDetails = catalogDetails
	testsEnrichService.CatalogDetails.SetUrl("http://localhost:" + strconv.Itoa(catalogServerPort) + "/catalog/")

	bom := cyclonedx.NewBOM()
	result, err := testsEnrichService.Enrich(bom)
	require.NoError(t, err)
	require.NotNil(t, result)
	if assert.NotNil(t, result.Vulnerabilities) {
		assert.Len(t, *result.Vulnerabilities, len(vulnerabilities))
	}
}
