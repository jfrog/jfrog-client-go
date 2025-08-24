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
	"github.com/jfrog/jfrog-client-go/evidence/services"
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
		if r.URL.Path == "/api/v1/system/version" {
			w.WriteHeader(http.StatusNotFound)
			requestNum++
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}, &requestNum
}

func TestIsEvidenceVersionSupportsFeature(t *testing.T) {
	tests := []struct {
		name                string
		supportedVersion    string
		handlerFunc         func(*testing.T) (http.HandlerFunc, *int)
		expectedResult      bool
		expectedError       bool
	}{
		{
			name:             "Version is higher than supported",
			supportedVersion: "1.0.0",
			handlerFunc:      createVersionHandlerFunc("2.0.0"),
			expectedResult:   true,
			expectedError:    false,
		},
		{
			name:             "Version is equal to supported",
			supportedVersion: "1.0.0",
			handlerFunc:      createVersionHandlerFunc("1.0.0"),
			expectedResult:   true,
			expectedError:    false,
		},
		{
			name:             "Version is lower than supported",
			supportedVersion: "2.0.0",
			handlerFunc:      createVersionHandlerFunc("1.0.0"),
			expectedResult:   false,
			expectedError:    false,
		},
		{
			name:             "Server returns 404",
			supportedVersion: "1.0.0",
			handlerFunc:      createErrorHandlerFuncVersion,
			expectedResult:   false,
			expectedError:    true,
		},
		{
			name:             "Empty supported version",
			supportedVersion: "",
			handlerFunc:      createVersionHandlerFunc("1.0.0"),
			expectedResult:   false,
			expectedError:    true,
		},
		{
			name:             "Invalid version format",
			supportedVersion: "invalid",
			handlerFunc:      createVersionHandlerFunc("1.0.0"),
			expectedResult:   false,
			expectedError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc, _ := tt.handlerFunc(t)
			mockServer, evdService := createMockServer(t, handlerFunc)
			defer mockServer.Close()

			result, err := evdService.IsEvidenceVersionSupportsFeature(tt.supportedVersion)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func createVersionHandlerFunc(version string) func(*testing.T) (http.HandlerFunc, *int) {
	return func(t *testing.T) (http.HandlerFunc, *int) {
		requestNum := 0
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v1/system/version" {
				w.WriteHeader(http.StatusOK)
				requestNum++
				versionPayload, err := json.Marshal(map[string]string{"version": version})
				assert.NoError(t, err)
				writeMockStatusResponse(t, w, versionPayload)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}, &requestNum
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		minVersion     string
		expectedResult bool
		expectedError  bool
	}{
		{
			name:           "Equal versions",
			version:        "1.0.0",
			minVersion:     "1.0.0",
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:           "Higher major version",
			version:        "2.0.0",
			minVersion:     "1.0.0",
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:           "Higher minor version",
			version:        "1.2.0",
			minVersion:     "1.1.0",
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:           "Higher patch version",
			version:        "1.0.2",
			minVersion:     "1.0.1",
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:           "Lower major version",
			version:        "1.0.0",
			minVersion:     "2.0.0",
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:           "Lower minor version",
			version:        "1.1.0",
			minVersion:     "1.2.0",
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:           "Lower patch version",
			version:        "1.0.1",
			minVersion:     "1.0.2",
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:           "Invalid version format - too few parts",
			version:        "1.0",
			minVersion:     "1.0.0",
			expectedResult: false,
			expectedError:  true,
		},
		{
			name:           "Invalid version format - too many parts",
			version:        "1.0.0.1",
			minVersion:     "1.0.0",
			expectedResult: false,
			expectedError:  true,
		},
		{
			name:           "Invalid version format - non-numeric",
			version:        "1.a.0",
			minVersion:     "1.0.0",
			expectedResult: false,
			expectedError:  true,
		},
		{
			name:           "Empty version",
			version:        "",
			minVersion:     "1.0.0",
			expectedResult: false,
			expectedError:  true,
		},
		{
			name:           "Empty minimum version",
			version:        "1.0.0",
			minVersion:     "",
			expectedResult: false,
			expectedError:  true,
		},
		{
			name:           "Whitespace handling",
			version:        " 1.0.0 ",
			minVersion:     " 1.0.0 ",
			expectedResult: true,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := services.CompareVersions(tt.version, tt.minVersion)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
