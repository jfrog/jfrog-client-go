//go:build itest

package tests

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xray/services"
)

func TestArtifactStatus(t *testing.T) {
	initXrayTest(t)
	xrayServerPort := xray.StartXrayMockServer(t)
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	require.NoError(t, err)
	testsArtifactService := services.NewArtifactService(client)
	testsArtifactService.XrayDetails = xrayDetails
	testsArtifactService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/xray/")

	tests := []struct {
		name           string
		repo           string
		path           string
		expectedStatus services.ArtifactStatus
		expectedTime   string
	}{
		{
			name:           "completed-scan",
			repo:           "test-repo",
			path:           "path/to/artifact",
			expectedStatus: services.ArtifactStatusDone,
			expectedTime:   "2023-12-01T10:00:00Z",
		},
		{
			name:           "pending-scan",
			repo:           "test-repo",
			path:           "path/to/pending-artifact",
			expectedStatus: services.ArtifactStatusPending,
			expectedTime:   "2023-12-01T09:30:00Z",
		},
		{
			name:           "unsupported-artifact",
			repo:           "test-repo",
			path:           "path/to/unsupported-artifact",
			expectedStatus: services.ArtifactStatusNotSupported,
			expectedTime:   "2023-12-01T11:00:00Z",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := testsArtifactService.GetStatus(test.repo, test.path)
			require.NoError(t, err)
			require.NotNil(t, response)
			require.NotNil(t, response.Overall)
			require.NotNil(t, response.Details)
			require.NotNil(t, response.Details.Sca)
			require.NotNil(t, response.Details.ContextualAnalysis)
			require.NotNil(t, response.Details.Exposures)
			require.NotNil(t, response.Details.Violations)

			// Verify the overall status and timestamp
			assert.Equal(t, test.expectedStatus, response.Overall.Status)
			assert.Equal(t, test.expectedTime, response.Overall.Timestamp)

			// Verify that all details have timestamps
			assert.NotEmpty(t, response.Details.Sca.Timestamp)
			assert.NotEmpty(t, response.Details.ContextualAnalysis.Timestamp)
			assert.NotEmpty(t, response.Details.Exposures.Timestamp)
			assert.NotEmpty(t, response.Details.Violations.Timestamp)
		})
	}

	// Test specific scenario details for the completed scan
	t.Run("completed-scan-details", func(t *testing.T) {
		response, err := testsArtifactService.GetStatus("test-repo", "path/to/artifact")
		require.NoError(t, err)
		require.NotNil(t, response)
		require.NotNil(t, response.Details)
		require.NotNil(t, response.Details.Sca)
		require.NotNil(t, response.Details.ContextualAnalysis)
		require.NotNil(t, response.Details.Exposures)
		require.NotNil(t, response.Details.Violations)

		// Verify specific statuses for the completed scan
		assert.Equal(t, services.ArtifactStatusDone, response.Details.Sca.Status)
		assert.Equal(t, services.ArtifactStatusDone, response.Details.ContextualAnalysis.Status)
		assert.Equal(t, services.ArtifactStatusNotSupported, response.Details.Exposures.Status)
		assert.Equal(t, services.ArtifactStatusFailed, response.Details.Violations.Status)
	})

	// Test specific scenario details for the pending scan
	t.Run("pending-scan-details", func(t *testing.T) {
		response, err := testsArtifactService.GetStatus("test-repo", "path/to/pending-artifact")
		require.NoError(t, err)
		require.NotNil(t, response)
		require.NotNil(t, response.Details)
		require.NotNil(t, response.Details.Sca)
		require.NotNil(t, response.Details.ContextualAnalysis)
		require.NotNil(t, response.Details.Exposures)
		require.NotNil(t, response.Details.Violations)

		// Verify specific statuses for the pending scan
		assert.Equal(t, services.ArtifactStatusPending, response.Details.Sca.Status)
		assert.Equal(t, services.ArtifactStatusNotScanned, response.Details.ContextualAnalysis.Status)
		assert.Equal(t, services.ArtifactStatusNotSupported, response.Details.Exposures.Status)
		assert.Equal(t, services.ArtifactStatusNotScanned, response.Details.Violations.Status)
	})
}
