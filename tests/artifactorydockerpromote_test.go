package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
)

const (
	sourceRepo        = "source-repo"
	targetRepo        = "target-repo"
	sourceDockerImage = "source-docker-image"
	targetDockerImage = "target-docker-image"
	sourceTag         = "source-tag"
	targetTag         = "target-tag"
	copy              = true
)

func TestDockerPromote(t *testing.T) {
	initArtifactoryTest(t)
	// Create mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/docker/"+sourceRepo+"/v2/promote", r.URL.Path)

		// Check body
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		response := services.DockerPromoteBody{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, targetRepo, response.TargetRepo)
		assert.Equal(t, sourceDockerImage, response.DockerRepository)
		assert.Equal(t, targetDockerImage, response.TargetDockerRepository)
		assert.Equal(t, sourceTag, response.Tag)
		assert.Equal(t, targetTag, response.TargetTag)
		assert.Equal(t, copy, response.Copy)

		// Send response 200 OK
		w.WriteHeader(http.StatusOK)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Promote docker image
	service := createDockerPromoteService(t, ts.URL)
	err := service.PromoteDocker(createTestParams())
	assert.NoError(t, err)
}

func createDockerPromoteService(t *testing.T, url string) *services.DockerPromoteService {
	// Create artifactory details
	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl(url + "/")

	// Create http client
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetInsecureTls(true).
		SetClientCertPath(rtDetails.GetClientCertPath()).
		SetClientCertKeyPath(rtDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(rtDetails.RunPreRequestFunctions).
		Build()
	assert.NoError(t, err, "Failed to create Artifactory client: %v\n")

	// Create docker promote service
	dockerPromoteService := services.NewDockerPromoteService(rtDetails, client)
	return dockerPromoteService
}

func createTestParams() services.DockerPromoteParams {
	params := services.NewDockerPromoteParams(sourceDockerImage, sourceRepo, targetRepo)
	params.TargetDockerImage = targetDockerImage
	params.SourceTag = sourceTag
	params.TargetTag = targetTag
	params.Copy = copy
	return params
}
