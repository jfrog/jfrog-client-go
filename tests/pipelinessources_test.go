package tests

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/pipelines/auth"
	"github.com/jfrog/jfrog-client-go/pipelines/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testsDummyRepo                 = "some-user/project-examples"
	testsDummyBranch               = "master"
	testsDummyProjectIntegrationId = 123
)

func TestSources(t *testing.T) {
	t.Run("addPipelineSource", testAddPipelineSource)
}

func testAddPipelineSource(t *testing.T) {
	expectedSourceId := 123
	tls := createPipelinesTLSServer(t, http.MethodPost, expectedSourceId, http.StatusOK)
	defer tls.Close()

	sourcesService, err := createDummySourcesService(tls.URL)
	if err != nil {
		assert.NoError(t, err)
		return
	}

	sourceId, err := sourcesService.AddPipelineSource(testsDummyProjectIntegrationId, testsDummyRepo, testsDummyBranch, services.DefaultPipelinesFileFilter)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.Equal(t, expectedSourceId, sourceId)
}

func createPipelinesTLSServer(t *testing.T, expectedRequest string, expectedSourceId, expectedStatusCode int) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedRequest, r.Method)
		assert.Equal(t, "/"+services.SourcesRestApi, r.URL.Path)
		w.WriteHeader(expectedStatusCode)
		_, err := w.Write([]byte(fmt.Sprintf(`{"id": %d}`, expectedSourceId)))
		assert.NoError(t, err)
	})
	return httptest.NewTLSServer(handler)
}

func createDummySourcesService(tlsUrl string) (*services.SourcesService, error) {
	details := auth.NewPipelinesDetails()
	details.SetUrl(tlsUrl + "/")
	details.SetAccessToken("fake-token")

	client, err := jfroghttpclient.JfrogClientBuilder().
		SetInsecureTls(true).
		SetServiceDetails(&details).
		Build()
	if err != nil {
		return nil, err
	}

	sourcesService := services.NewSourcesService(client)
	sourcesService.ServiceDetails = details
	return sourcesService, nil
}
