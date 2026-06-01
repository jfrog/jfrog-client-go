package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
)

const queryData = `{"query":"someGraphqlQuery"}`

func TestOnemodelService_Query(t *testing.T) {
	handlerFunc, requestNum := createOnemodelHandlerFunc(t)

	mockServer, onemodelService := createMockonemodelServer(t, handlerFunc)
	defer mockServer.Close()

	_, err := onemodelService.Query([]byte(queryData))
	assert.NoError(t, err)
	assert.Equal(t, 1, *requestNum)
}

func createMockonemodelServer(t *testing.T, testHandler http.HandlerFunc) (*httptest.Server, *onemodelService) {
	testServer := httptest.NewServer(testHandler)

	serviceDetails := auth.NewArtifactoryDetails()
	serviceDetails.SetUrl(testServer.URL + "/")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)
	return testServer, &onemodelService{serviceDetails: &serviceDetails, client: client}
}

func createOnemodelHandlerFunc(t *testing.T) (http.HandlerFunc, *int) {
	requestNum := 0
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/v1/graphql" {
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			w.WriteHeader(http.StatusOK)
			requestNum++
			writeMockonemodelResponse(t, w, []byte(queryData))
		}
	}, &requestNum
}

func writeMockonemodelResponse(t *testing.T, w http.ResponseWriter, payload []byte) {
	content, err := json.Marshal(payload)
	assert.NoError(t, err)
	_, err = w.Write(content)
	assert.NoError(t, err)
}
