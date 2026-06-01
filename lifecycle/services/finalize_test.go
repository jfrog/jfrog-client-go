package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinalizeOperation_getOperationRestApi_shouldReturnCorrectPath(t *testing.T) {
	// Arrange
	operation := finalizeOperation{
		rbDetails: ReleaseBundleDetails{
			ReleaseBundleName:    "my-bundle",
			ReleaseBundleVersion: "1.0.0",
		},
	}

	// Act
	result := operation.getOperationRestApi()

	// Assert
	assert.Equal(t, "api/v2/release_bundle/my-bundle/1.0.0/finalize", result)
}

func TestFinalizeOperation_getRequestBody_shouldReturnNil(t *testing.T) {
	// Arrange
	operation := finalizeOperation{
		rbDetails: ReleaseBundleDetails{
			ReleaseBundleName:    "my-bundle",
			ReleaseBundleVersion: "1.0.0",
		},
	}

	// Act
	result := operation.getRequestBody()

	// Assert
	assert.Nil(t, result)
}

func TestFinalizeOperation_getOperationSuccessfulMsg_shouldReturnCorrectMessage(t *testing.T) {
	// Arrange
	operation := finalizeOperation{}

	// Act
	result := operation.getOperationSuccessfulMsg()

	// Assert
	assert.Equal(t, "Release Bundle successfully finalized", result)
}

func TestFinalizeOperation_getOperationParams_shouldReturnParams(t *testing.T) {
	// Arrange
	params := CommonOptionalQueryParams{
		ProjectKey: "my-project",
		Async:      true,
	}
	operation := finalizeOperation{
		params: params,
	}

	// Act
	result := operation.getOperationParams()

	// Assert
	assert.Equal(t, "my-project", result.ProjectKey)
	assert.True(t, result.Async)
}

func TestFinalizeOperation_getSigningKeyName_shouldReturnSigningKey(t *testing.T) {
	// Arrange
	operation := finalizeOperation{
		signingKeyName: "my-gpg-key",
	}

	// Act
	result := operation.getSigningKeyName()

	// Assert
	assert.Equal(t, "my-gpg-key", result)
}

func TestBuildFinalizeQueryParams_shouldIncludeProjectKey_whenProvided(t *testing.T) {
	// Arrange
	params := CommonOptionalQueryParams{
		ProjectKey: "my-project",
		Async:      false,
	}

	// Act
	result := buildFinalizeQueryParams(params)

	// Assert
	assert.Equal(t, "my-project", result["project"])
	assert.Equal(t, "false", result["async"])
}

func TestBuildFinalizeQueryParams_shouldSetAsyncTrue_whenAsyncEnabled(t *testing.T) {
	// Arrange
	params := CommonOptionalQueryParams{
		Async: true,
	}

	// Act
	result := buildFinalizeQueryParams(params)

	// Assert
	assert.Equal(t, "true", result["async"])
}
