package evidence

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	artifactoryAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	evidence "github.com/jfrog/jfrog-client-go/evidence/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			expectedPathInURL: "/evidence/api/v1/subject/ssfpoc-build-info/SSF%20Demo%20Electron/2-1762810453125.json",
			description:       "Spaces should be encoded as %20, not double-encoded as %2520",
		},
		{
			name:              "Build name with multiple spaces",
			subjectUri:        "test-build-info/My Build Name With Spaces/123.json",
			expectedPathInURL: "/evidence/api/v1/subject/test-build-info/My%20Build%20Name%20With%20Spaces/123.json",
			description:       "Multiple spaces should be properly encoded",
		},
		{
			name:              "Build name with special characters",
			subjectUri:        "proj-build-info/Build & Test + Deploy/456.json",
			expectedPathInURL: "/evidence/api/v1/subject/proj-build-info/Build%20&%20Test%20+%20Deploy/456.json",
			description:       "Special characters should be properly encoded",
		},
		{
			name:              "Simple build name without special chars",
			subjectUri:        "simple-build-info/SimpleBuildName/789.json",
			expectedPathInURL: "/evidence/api/v1/subject/simple-build-info/SimpleBuildName/789.json",
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
				case r.URL.Path == "/evidence/api/v1/system/version":
					w.WriteHeader(http.StatusOK)

				case strings.HasPrefix(r.URL.Path, "/evidence/api/v1/subject/"):
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

			// Create mock server with evidence URL prefix
			mockServer := httptest.NewServer(http.HandlerFunc(mockHandler))
			defer mockServer.Close()

			// Create evidence service with the correct evidence URL
			serviceDetails := artifactoryAuth.NewArtifactoryDetails()
			serviceDetails.SetUrl(mockServer.URL + "/evidence/")
			client, err := jfroghttpclient.JfrogClientBuilder().Build()
			assert.NoError(t, err)
			evdService := evidence.NewEvidenceService(serviceDetails, client)

			// Create evidence details with the test subject URI
			evidenceDetails := evidence.EvidenceDetails{
				SubjectUri:  tc.subjectUri,
				DSSEFileRaw: dsseRaw,
				ProviderId:  "test-provider",
			}

			// Execute the upload
			_, err = evdService.UploadEvidence(evidenceDetails)
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

func TestUploadEvidence_WithAttachments(t *testing.T) {
	var capturedBody []byte
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/system/version":
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("7.646.1"))
			assert.NoError(t, err)
		case strings.HasPrefix(r.URL.Path, "/api/v1/subject/") || r.URL.Path == "/api/v1/subject":
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			capturedBody = body
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write([]byte(`{"success": true}`))
			assert.NoError(t, err)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer testServer.Close()

	serviceDetails := artifactoryAuth.NewArtifactoryDetails()
	serviceDetails.SetUrl(testServer.URL + "/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)
	evdService := evidence.NewEvidenceService(serviceDetails, client)

	details := evidence.EvidenceDetails{
		SubjectUri:  "repo/path/file.txt",
		DSSEFileRaw: []byte(`{"payload":"abc","payloadType":"application/vnd.in-toto+json","signatures":[],"attachments":[{"repository":"example-repo-local","path":"tmp/a.txt","sha256":"abc123"}]}`),
		Attachments: []evidence.AttachmentDetails{{
			Repository: "example-repo-local",
			Path:       "tmp/a.txt",
			Sha256:     "abc123",
		}},
	}
	_, err = evdService.UploadEvidence(details)
	require.NoError(t, err)
	assert.JSONEq(t,
		`{"payload":"abc","payloadType":"application/vnd.in-toto+json","signatures":[],"attachments":[{"repository":"example-repo-local","path":"tmp/a.txt","sha256":"abc123"}]}`,
		string(capturedBody))
}

func TestUploadEvidence_WithAttachments_FailsOnUnsupportedVersion(t *testing.T) {
	subjectEndpointCalled := false
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/system/version":
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("7.645.0"))
			assert.NoError(t, err)
		case strings.HasPrefix(r.URL.Path, "/api/v1/subject/") || r.URL.Path == "/api/v1/subject":
			subjectEndpointCalled = true
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte(`{"success": true}`))
			assert.NoError(t, err)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer testServer.Close()

	serviceDetails := artifactoryAuth.NewArtifactoryDetails()
	serviceDetails.SetUrl(testServer.URL + "/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)
	evdService := evidence.NewEvidenceService(serviceDetails, client)

	_, err = evdService.UploadEvidence(evidence.EvidenceDetails{
		SubjectUri:  "repo/path/file.txt",
		DSSEFileRaw: []byte(`{"payload":"abc"}`),
		Attachments: []evidence.AttachmentDetails{{Repository: "r", Path: "p", Sha256: "s"}},
	})
	require.EqualError(t, err, "You are using JFrog Evidence version 7.645.0, while this operation requires version 7.646.1 or higher.")
	assert.False(t, subjectEndpointCalled)
}

func TestUploadEvidence_WithAttachments_FailsOnVersionApiError(t *testing.T) {
	subjectEndpointCalled := false
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/system/version":
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(`boom`))
			assert.NoError(t, err)
		case strings.HasPrefix(r.URL.Path, "/api/v1/subject/") || r.URL.Path == "/api/v1/subject":
			subjectEndpointCalled = true
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte(`{"success": true}`))
			assert.NoError(t, err)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer testServer.Close()

	serviceDetails := artifactoryAuth.NewArtifactoryDetails()
	serviceDetails.SetUrl(testServer.URL + "/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)
	evdService := evidence.NewEvidenceService(serviceDetails, client)

	_, err = evdService.UploadEvidence(evidence.EvidenceDetails{
		SubjectUri:  "repo/path/file.txt",
		DSSEFileRaw: []byte(`{"payload":"abc"}`),
		Attachments: []evidence.AttachmentDetails{{Repository: "r", Path: "p", Sha256: "s"}},
	})
	require.Error(t, err)
	assert.False(t, subjectEndpointCalled)
}

type DeleteEvidenceMockHandlerConfig struct {
	StatusCode              int
	CapturedSubjectRepoPath *string
	CapturedEvidenceName    *string
}

func createDeleteEvidenceMockHandler(t *testing.T, config *DeleteEvidenceMockHandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/v1/evidence/"):
			trimmed := strings.TrimPrefix(r.URL.Path, "/api/v1/evidence/")
			lastSlash := strings.LastIndex(trimmed, "/")
			if lastSlash < 0 {
				t.Errorf("Unexpected delete path format: %s", r.URL.Path)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			subjectRepoPathEscaped := trimmed[:lastSlash]
			evidenceNameEscaped := trimmed[lastSlash+1:]

			subjectRepoPath, err := url.PathUnescape(subjectRepoPathEscaped)
			assert.NoError(t, err)
			evidenceName, err := url.PathUnescape(evidenceNameEscaped)
			assert.NoError(t, err)

			if config.CapturedSubjectRepoPath != nil {
				*config.CapturedSubjectRepoPath = subjectRepoPath
			}
			if config.CapturedEvidenceName != nil {
				*config.CapturedEvidenceName = evidenceName
			}

			if config.StatusCode != 0 {
				w.WriteHeader(config.StatusCode)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		default:
			t.Errorf("Unexpected request URI: %s", r.RequestURI)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func TestDeleteEvidence_Success(t *testing.T) {
	subjectRepoPath := "repo/path with space/a/b"
	evidenceName := "evidence name.json"

	var capturedSubject string
	var capturedName string

	handler := createDeleteEvidenceMockHandler(t, &DeleteEvidenceMockHandlerConfig{
		StatusCode:              http.StatusNoContent,
		CapturedSubjectRepoPath: &capturedSubject,
		CapturedEvidenceName:    &capturedName,
	})
	mockServer, evdService := createMockServer(t, handler)
	defer mockServer.Close()

	err := evdService.DeleteEvidence(subjectRepoPath, evidenceName)
	assert.NoError(t, err)
	assert.Equal(t, subjectRepoPath, capturedSubject)
	assert.Equal(t, evidenceName, capturedName)
}

func TestDeleteEvidence_NotFound(t *testing.T) {
	handler := createDeleteEvidenceMockHandler(t, &DeleteEvidenceMockHandlerConfig{StatusCode: http.StatusNotFound})
	mockServer, evdService := createMockServer(t, handler)
	defer mockServer.Close()

	err := evdService.DeleteEvidence("repo/a", "missing.json")
	assert.Error(t, err)
}

func TestDeleteEvidence_ServerError(t *testing.T) {
	handler := createDeleteEvidenceMockHandler(t, &DeleteEvidenceMockHandlerConfig{StatusCode: http.StatusInternalServerError})
	mockServer, evdService := createMockServer(t, handler)
	defer mockServer.Close()

	err := evdService.DeleteEvidence("repo/a", "boom.json")
	assert.Error(t, err)
}
