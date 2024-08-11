package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const queryData = `{"query":"someGraphqlQuery"}`

func TestMetadataService_Query(t *testing.T) {
	handlerFunc, requestNum := createMetadataHandlerFunc(t)

	mockServer, metadataService := createMockMetadataServer(t, handlerFunc)
	defer mockServer.Close()

	_, err := metadataService.Query([]byte(queryData))
	assert.NoError(t, err)
	assert.Equal(t, 1, *requestNum)
}

func createMockMetadataServer(t *testing.T, testHandler http.HandlerFunc) (*httptest.Server, *metadataService) {
	testServer := httptest.NewServer(testHandler)

	serviceDetails := auth.NewArtifactoryDetails()
	serviceDetails.SetUrl(testServer.URL + "/")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)
	return testServer, &metadataService{serviceDetails: &serviceDetails, client: client}
}

func createMetadataHandlerFunc(t *testing.T) (http.HandlerFunc, *int) {
	requestNum := 0
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/v1/query" {
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			w.WriteHeader(http.StatusOK)
			requestNum++
			writeMockMetadataResponse(t, w, []byte(queryData))
		}
	}, &requestNum
}

func writeMockMetadataResponse(t *testing.T, w http.ResponseWriter, payload []byte) {
	content, err := json.Marshal(payload)
	assert.NoError(t, err)
	_, err = w.Write(content)
	assert.NoError(t, err)
}
