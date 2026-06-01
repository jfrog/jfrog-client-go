package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRbUpdateBodyJsonMarshal_shouldIncludeAddSources_whenSourcesProvided(t *testing.T) {
	// Arrange
	rbDetails := ReleaseBundleDetails{
		ReleaseBundleName:    "test-bundle",
		ReleaseBundleVersion: "1.0.0",
	}
	rbUpdateBody := RbUpdateBody{
		ReleaseBundleDetails: rbDetails,
		AddSources: []RbSource{
			{
				SourceType: Artifacts,
				Artifacts: []ArtifactSource{
					{Path: "repo/path/artifact.jar"},
				},
			},
		},
	}

	// Act
	marshaled, err := json.Marshal(rbUpdateBody)

	// Assert
	assert.NoError(t, err)
	expected := `{
		"release_bundle_name": "test-bundle",
		"release_bundle_version": "1.0.0",
		"add_sources": [
			{
				"source_type": "artifacts",
				"artifacts": [
					{"path": "repo/path/artifact.jar"}
				]
			}
		]
	}`
	assert.JSONEq(t, expected, string(marshaled))
}

func TestRbUpdateBodyJsonMarshal_shouldOmitAddSources_whenEmpty(t *testing.T) {
	// Arrange
	rbUpdateBody := RbUpdateBody{
		ReleaseBundleDetails: ReleaseBundleDetails{
			ReleaseBundleName:    "test-bundle",
			ReleaseBundleVersion: "1.0.0",
		},
		AddSources: []RbSource{},
	}

	// Act
	marshaled, err := json.Marshal(rbUpdateBody)

	// Assert
	assert.NoError(t, err)
	// Empty slice should be omitted due to omitempty
	expected := `{"release_bundle_name": "test-bundle", "release_bundle_version": "1.0.0"}`
	assert.JSONEq(t, expected, string(marshaled))
}

func TestRbUpdateBodyJsonMarshal_shouldIncludeMultipleSources_whenMixedSourceTypes(t *testing.T) {
	// Arrange
	rbUpdateBody := RbUpdateBody{
		ReleaseBundleDetails: ReleaseBundleDetails{
			ReleaseBundleName:    "multi-source-bundle",
			ReleaseBundleVersion: "2.0.0",
		},
		AddSources: []RbSource{
			{
				SourceType: Artifacts,
				Artifacts: []ArtifactSource{
					{Path: "generic-local/file.txt"},
				},
			},
			{
				SourceType: Builds,
				Builds: []BuildSource{
					{
						BuildName:           "my-build",
						BuildNumber:         "123",
						BuildRepository:     "build-info",
						IncludeDependencies: true,
					},
				},
			},
		},
	}

	// Act
	marshaled, err := json.Marshal(rbUpdateBody)

	// Assert
	assert.NoError(t, err)
	expected := `{
		"release_bundle_name": "multi-source-bundle",
		"release_bundle_version": "2.0.0",
		"add_sources": [
			{
				"source_type": "artifacts",
				"artifacts": [
					{"path": "generic-local/file.txt"}
				]
			},
			{
				"source_type": "builds",
				"builds": [
					{
						"build_name": "my-build",
						"build_number": "123",
						"build_repository": "build-info",
						"include_dependencies": true
					}
				]
			}
		]
	}`
	assert.JSONEq(t, expected, string(marshaled))
}
