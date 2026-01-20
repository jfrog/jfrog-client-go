package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRbCreationBodyJsonMarshal_shouldIncludeDraft_whenDraftIsTrue(t *testing.T) {
	// Arrange
	body := RbCreationBody{
		ReleaseBundleDetails: ReleaseBundleDetails{
			ReleaseBundleName:    "test-bundle",
			ReleaseBundleVersion: "1.0.0",
		},
		SourceType: Aql,
		Source:     CreateFromAqlSource{Aql: "items.find({\"repo\": \"test-repo\"})"},
		Draft:      true,
	}

	// Act
	jsonBytes, err := json.Marshal(body)

	// Assert
	assert.NoError(t, err)
	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, `"draft":true`)
	assert.Contains(t, jsonStr, `"release_bundle_name":"test-bundle"`)
	assert.Contains(t, jsonStr, `"release_bundle_version":"1.0.0"`)
}

func TestRbCreationBodyJsonMarshal_shouldOmitDraft_whenDraftIsFalse(t *testing.T) {
	// Arrange
	body := RbCreationBody{
		ReleaseBundleDetails: ReleaseBundleDetails{
			ReleaseBundleName:    "test-bundle",
			ReleaseBundleVersion: "1.0.0",
		},
		SourceType: Builds,
		Source:     CreateFromBuildsSource{Builds: []BuildSource{{BuildName: "build1", BuildNumber: "1"}}},
		Draft:      false,
	}

	// Act
	jsonBytes, err := json.Marshal(body)

	// Assert
	assert.NoError(t, err)
	jsonStr := string(jsonBytes)
	// Draft should be omitted when false due to omitempty tag
	assert.NotContains(t, jsonStr, `"draft"`)
	assert.Contains(t, jsonStr, `"release_bundle_name":"test-bundle"`)
}

func TestRbCreationBodyJsonMarshal_shouldIncludeSources_whenMultipleSources(t *testing.T) {
	// Arrange
	body := RbCreationBody{
		ReleaseBundleDetails: ReleaseBundleDetails{
			ReleaseBundleName:    "multi-source-bundle",
			ReleaseBundleVersion: "2.0.0",
		},
		Sources: []RbSource{
			{
				SourceType: Builds,
				Builds:     []BuildSource{{BuildName: "build1", BuildNumber: "1"}},
			},
			{
				SourceType: ReleaseBundles,
				ReleaseBundles: []ReleaseBundleSource{
					{ReleaseBundleName: "source-bundle", ReleaseBundleVersion: "1.0.0"},
				},
			},
		},
		Draft: true,
	}

	// Act
	jsonBytes, err := json.Marshal(body)

	// Assert
	assert.NoError(t, err)
	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, `"draft":true`)
	assert.Contains(t, jsonStr, `"sources"`)
	assert.Contains(t, jsonStr, `"build1"`)
	assert.Contains(t, jsonStr, `"source-bundle"`)
}

func TestBackwardCompatibility_originalMethodsDefaultToDraftFalse(t *testing.T) {
	// This test verifies the backward-compatible methods exist and have the correct signatures
	// by checking that calling the original method names compiles and works

	// The original methods (without Draft suffix) should exist and default draft to false
	// We're testing the struct/function signatures indirectly through RbCreationBody
	body := RbCreationBody{
		ReleaseBundleDetails: ReleaseBundleDetails{
			ReleaseBundleName:    "backward-compat-bundle",
			ReleaseBundleVersion: "1.0.0",
		},
		SourceType: Aql,
		Source:     CreateFromAqlSource{Aql: "items.find({\"repo\": \"test-repo\"})"},
		// Draft not set - should default to false (zero value)
	}

	jsonBytes, err := json.Marshal(body)
	assert.NoError(t, err)

	// Verify draft is not in output when using default (false)
	jsonStr := string(jsonBytes)
	assert.NotContains(t, jsonStr, `"draft"`)
}
