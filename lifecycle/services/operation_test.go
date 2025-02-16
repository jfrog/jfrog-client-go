package services

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockReleaseBundleOperation struct {
	restApi         string
	requestBody     any
	successfulMsg   string
	operationParams CommonOptionalQueryParams
	signingKeyName  string
}

func (m *MockReleaseBundleOperation) getOperationRestApi() string {
	return m.restApi
}

func (m *MockReleaseBundleOperation) getRequestBody() any {
	return m.requestBody
}

func (m *MockReleaseBundleOperation) getOperationSuccessfulMsg() string {
	return m.successfulMsg
}

func (m *MockReleaseBundleOperation) getOperationParams() CommonOptionalQueryParams {
	return m.operationParams
}

func (m *MockReleaseBundleOperation) getSigningKeyName() string {
	return m.signingKeyName
}

func TestPrepareRequest(t *testing.T) {
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err, "Failed to create JFrog HTTP client")

	lcDetails := auth.NewArtifactoryDetails()
	lcDetails.SetUrl("http://localhost:8081/artifactory")
	rbs := NewReleaseBundlesService(lcDetails, client)

	scenarios := []struct {
		name                string
		operationParams     CommonOptionalQueryParams
		signingKeyName      string
		expectedURL         string
		expectedContentType string
		expectedErr         bool
	}{
		{
			name:                "Valid Operation with All Parameters",
			operationParams:     CommonOptionalQueryParams{ProjectKey: "testProject", Async: true, PromotionType: "move"},
			signingKeyName:      "testKey",
			expectedURL:         "http://localhost:8081/artifactory/api/test?async=true&operation=move&project=testProject",
			expectedContentType: "application/json",
			expectedErr:         false,
		},
		{
			name:                "Empty Operation Parameters",
			operationParams:     CommonOptionalQueryParams{},
			signingKeyName:      "testKey",
			expectedURL:         "http://localhost:8081/artifactory/api/test?async=false",
			expectedContentType: "application/json",
			expectedErr:         false,
		},
		{
			name:                "No Signing Key Name",
			operationParams:     CommonOptionalQueryParams{ProjectKey: "testProject", Async: true, PromotionType: "move"},
			signingKeyName:      "", // No signing key
			expectedURL:         "http://localhost:8081/artifactory/api/test?async=true&operation=move&project=testProject",
			expectedContentType: "application/json",
			expectedErr:         false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			operation := &MockReleaseBundleOperation{
				restApi:         "api/test",
				requestBody:     nil,
				successfulMsg:   "Operation successful",
				operationParams: scenario.operationParams,
				signingKeyName:  scenario.signingKeyName,
			}

			requestFullUrl, httpClientDetails, err := prepareRequest(operation, rbs)

			if scenario.expectedErr {
				assert.Error(t, err, "Expected an error for scenario: %s", scenario.name)
			} else {
				assert.NoError(t, err, "Expected no error for scenario: %s", scenario.name)
				assert.Equal(t, scenario.expectedURL, requestFullUrl, "Unexpected request URL for scenario: %s", scenario.name)
				assert.Equal(t, scenario.expectedContentType, httpClientDetails.Headers["Content-Type"], "Unexpected Content-Type for scenario: %s", scenario.name)
				assert.Equal(t, scenario.signingKeyName, httpClientDetails.Headers["X-JFrog-Signing-Key-Name"], "Unexpected Signing Key Name for scenario: %s", scenario.name)
			}
		})
	}
}
