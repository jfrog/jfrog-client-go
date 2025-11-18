package evidence

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	artifactoryAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	evidence "github.com/jfrog/jfrog-client-go/evidence/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
)

var dsseRaw = []byte("someData")

func createMockServer(t *testing.T, testHandler http.HandlerFunc) (*httptest.Server, *evidence.EvidenceService) {
	testServer := httptest.NewServer(testHandler)
	rtDetails := artifactoryAuth.NewArtifactoryDetails()
	rtDetails.SetUrl(testServer.URL + "/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)
	return testServer, evidence.NewEvidenceService(rtDetails, client)
}

func writeMockStatusResponse(t *testing.T, w http.ResponseWriter, payload []byte) {
	_, err := w.Write(payload)
	assert.NoError(t, err)
}

type UploadEvidenceMockHandlerConfig struct {
	VersionEndpointStatusCode   int
	SubjectEndpointStatusCode   int
	SubjectEndpointResponseBody []byte

	CapturedProviderId *string // Pointer to capture the providerId from /api/v1/subject request
}

func createUploadEvidenceMockHandler(t *testing.T, config *UploadEvidenceMockHandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/system/version":
			w.WriteHeader(config.VersionEndpointStatusCode)

		case strings.HasPrefix(r.URL.Path, "/api/v1/subject/") || r.URL.Path == "/api/v1/subject":
			parsedRequestURI, err := url.Parse(r.RequestURI)
			assert.NoError(t, err, "Failed to parse RequestURI: %s", r.RequestURI)

			queryParams, err := url.ParseQuery(parsedRequestURI.RawQuery)
			assert.NoError(t, err, "Failed to parse query from parsed RequestURI.RawQuery: %s", parsedRequestURI.RawQuery)

			providerId := queryParams.Get("providerId")

			if config.CapturedProviderId != nil {
				*config.CapturedProviderId = providerId
			}

			w.WriteHeader(config.SubjectEndpointStatusCode)
			if config.SubjectEndpointResponseBody != nil {
				writeMockStatusResponse(t, w, config.SubjectEndpointResponseBody)
			}

		default:
			t.Errorf("Unexpected request URI: %s", r.RequestURI)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func TestUploadEvidence_ProviderIdLogic(t *testing.T) {
	tests := []struct {
		name                        string
		initialEvidenceData         evidence.EvidenceDetails
		versionStatusCode           int
		expectedProviderIdInRequest string
		expectedUploadError         bool
	}{
		{
			name:                        "ProviderId supported: Should be passed",
			initialEvidenceData:         evidence.EvidenceDetails{SubjectUri: "sub1", DSSEFileRaw: dsseRaw, ProviderId: "testProvider123"},
			versionStatusCode:           http.StatusOK,
			expectedProviderIdInRequest: "testProvider123",
			expectedUploadError:         false,
		},
		{
			name:                        "ProviderId not supported (404 version): Should be empty",
			initialEvidenceData:         evidence.EvidenceDetails{SubjectUri: "sub2", DSSEFileRaw: dsseRaw, ProviderId: "testProvider123"},
			versionStatusCode:           http.StatusNotFound,
			expectedProviderIdInRequest: "",
			expectedUploadError:         false,
		},
		{
			name:                        "ProviderId not supported (error version): Should be empty",
			initialEvidenceData:         evidence.EvidenceDetails{SubjectUri: "sub3", DSSEFileRaw: dsseRaw, ProviderId: "testProvider123"},
			versionStatusCode:           http.StatusInternalServerError,
			expectedProviderIdInRequest: "",
			expectedUploadError:         false,
		},
		{
			name:                        "ProviderId supported but initially empty: Should remain empty",
			initialEvidenceData:         evidence.EvidenceDetails{SubjectUri: "sub4", DSSEFileRaw: dsseRaw, ProviderId: ""},
			versionStatusCode:           http.StatusOK,
			expectedProviderIdInRequest: "",
			expectedUploadError:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedProviderId string

			mockHandlerConfig := &UploadEvidenceMockHandlerConfig{
				VersionEndpointStatusCode:   tt.versionStatusCode,
				SubjectEndpointStatusCode:   http.StatusOK,
				SubjectEndpointResponseBody: []byte("mocked_do_operation_response"),
				CapturedProviderId:          &capturedProviderId,
			}

			handlerFunc := createUploadEvidenceMockHandler(t, mockHandlerConfig)
			mockServer, evdService := createMockServer(t, handlerFunc)
			defer mockServer.Close()

			_, err := evdService.UploadEvidence(tt.initialEvidenceData)

			if tt.expectedUploadError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedProviderIdInRequest, capturedProviderId)
		})
	}
}
func TestIsEvidenceSupportsProviderId(t *testing.T) {
	tests := []struct {
		name           string
		handlerFunc    func(*testing.T) (http.HandlerFunc, *int)
		expectedResult bool
	}{
		{
			name:           "Version supports providerId",
			handlerFunc:    createDefaultHandlerFuncVersion,
			expectedResult: true,
		},
		{
			name:           "Version does not support providerId",
			handlerFunc:    createErrorHandlerFuncVersion,
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc, _ := tt.handlerFunc(t)
			mockServer, evdService := createMockServer(t, handlerFunc)
			defer mockServer.Close()
			result := evdService.IsEvidenceSupportsProviderId()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func createDefaultHandlerFuncVersion(t *testing.T) (http.HandlerFunc, *int) {
	requestNum := 0
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/v1/version" {
			w.WriteHeader(http.StatusOK)
			requestNum++
			versionPayload, err := json.Marshal(map[string]string{"version": "1.0.0"})
			assert.NoError(t, err)
			writeMockStatusResponse(t, w, versionPayload)
		}
	}, &requestNum
}

func createErrorHandlerFuncVersion(t *testing.T) (http.HandlerFunc, *int) {
	requestNum := 0
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/v1/version" {
			w.WriteHeader(http.StatusNotFound)
			requestNum++
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}, &requestNum
}

func TestUploadEvidence_URLEncodingFix(t *testing.T) {
	// test for double URL encoding bug
	// This test verifies that build names with spaces are encoded only once, not twice
	testCases := []struct {
		name              string
		subjectUri        string
		expectedPathInURL string
		description       string
	}{
		{
			name:              "Build name with spaces (original bug case)",
			subjectUri:        "ssfpoc-build-info/SSF Demo Electron/2-1762810453125.json",
			expectedPathInURL: "/api/v1/subject/ssfpoc-build-info/SSF%20Demo%20Electron/2-1762810453125.json",
			description:       "Spaces should be encoded as %20, not double-encoded as %2520",
		},
		{
			name:              "Build name with multiple spaces",
			subjectUri:        "test-build-info/My Build Name With Spaces/123.json",
			expectedPathInURL: "/api/v1/subject/test-build-info/My%20Build%20Name%20With%20Spaces/123.json",
			description:       "Multiple spaces should be properly encoded",
		},
		{
			name:              "Build name with special characters",
			subjectUri:        "proj-build-info/Build & Test + Deploy/456.json",
			expectedPathInURL: "/api/v1/subject/proj-build-info/Build%20&%20Test%20+%20Deploy/456.json",
			description:       "Special characters should be properly encoded",
		},
		{
			name:              "Simple build name without special chars",
			subjectUri:        "simple-build-info/SimpleBuildName/789.json",
			expectedPathInURL: "/api/v1/subject/simple-build-info/SimpleBuildName/789.json",
			description:       "Simple names should remain unchanged",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var capturedRequestPath string
			var capturedRawURL string

			// Create a mock handler that captures the actual request path
			mockHandler := func(w http.ResponseWriter, r *http.Request) {
				switch {
				case r.URL.Path == "/api/v1/system/version":
					w.WriteHeader(http.StatusOK)

				case strings.HasPrefix(r.URL.Path, "/api/v1/subject/"):
					// Capture both the decoded path and the raw URL for verification
					capturedRequestPath = r.URL.Path
					capturedRawURL = r.URL.String()
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(`{"success": true}`))
					assert.NoError(t, err)

				default:
					t.Errorf("Unexpected request path: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}

			mockServer, evdService := createMockServer(t, mockHandler)
			defer mockServer.Close()

			// Create evidence details with the test subject URI
			evidenceDetails := evidence.EvidenceDetails{
				SubjectUri:  tc.subjectUri,
				DSSEFileRaw: dsseRaw,
				ProviderId:  "test-provider",
			}

			// Execute the upload
			_, err := evdService.UploadEvidence(evidenceDetails)
			assert.NoError(t, err, "Upload should succeed for: %s", tc.description)

			// The main test: verify that the URL contains properly encoded paths
			// We check the raw URL string to see the actual encoding
			assert.Contains(t, capturedRawURL, tc.expectedPathInURL,
				"URL should contain properly encoded path for: %s. Expected path in: %s, Got URL: %s",
				tc.description, tc.expectedPathInURL, capturedRawURL)

			// Specific regression test for the original double encoding bug
			if strings.Contains(tc.subjectUri, "SSF Demo Electron") {
				// Verify that spaces are encoded as %20 and NOT as %2520 (double encoded)
				assert.Contains(t, capturedRawURL, "SSF%20Demo%20Electron",
					"URL should contain single-encoded spaces (%%20)")
				assert.NotContains(t, capturedRawURL, "SSF%2520Demo%2520Electron",
					"URL should NOT contain double-encoded spaces (%%2520)")

				t.Logf("✓ Regression test passed: %s", tc.description)
			}

			t.Logf("Test case: %s", tc.description)
			t.Logf("Subject URI: %s", tc.subjectUri)
			t.Logf("Captured decoded path: %s", capturedRequestPath)
			t.Logf("Captured raw URL: %s", capturedRawURL)
		})
	}
}
