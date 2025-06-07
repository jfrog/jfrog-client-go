package evidence

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	artifactoryAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	evidence "github.com/jfrog/jfrog-client-go/evidence/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
)

var dsseRaw = []byte("someData")

var evidenceData = evidence.EvidenceDetails{
	SubjectUri:  "someUri",
	DSSEFileRaw: dsseRaw,
	ProviderId:  "someProviderId",
}

func TestEvidenceServicesManager_UploadEvidence(t *testing.T) {
	handlerFunc, requestNum := createDefaultHandlerFunc(t)

	mockServer, evdService := createMockServer(t, handlerFunc)
	defer mockServer.Close()

	_, err := evdService.UploadEvidence(evidenceData)
	assert.NoError(t, err)
	assert.Equal(t, 0, *requestNum)
}

func createMockServer(t *testing.T, testHandler http.HandlerFunc) (*httptest.Server, *evidence.EvidenceService) {
	testServer := httptest.NewServer(testHandler)

	rtDetails := artifactoryAuth.NewArtifactoryDetails()
	rtDetails.SetUrl(testServer.URL + "/")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)
	return testServer, evidence.NewEvidenceService(rtDetails, client)
}

func createDefaultHandlerFunc(t *testing.T) (http.HandlerFunc, *int) {
	requestNum := 0
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/v1/subject" {
			w.WriteHeader(http.StatusOK)
			requestNum++
			writeMockStatusResponse(t, w, dsseRaw)
		}
	}, &requestNum
}

func writeMockStatusResponse(t *testing.T, w http.ResponseWriter, payload []byte) {
	content, err := json.Marshal(payload)
	assert.NoError(t, err)
	_, err = w.Write(content)
	assert.NoError(t, err)
}
