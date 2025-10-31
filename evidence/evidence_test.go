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
