package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

const (
	sourceRepo       = "source-repo"
	targetRepo       = "target-repo"
	dockerRepo       = "docker-repo"
	targetDockerRepo = "target-docker-repo"
	tag              = "tag"
	targetTag        = "target-tag"
	copy             = true
)

func TestDockerPromote(t *testing.T) {
	// Create mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/docker/"+sourceRepo+"/v2/promote", r.URL.Path)

		// Check body
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		response := services.DockerPromoteBody{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, targetRepo, response.TargetRepo)
		assert.Equal(t, dockerRepo, response.DockerRepository)
		assert.Equal(t, targetDockerRepo, response.TargetDockerRepository)
		assert.Equal(t, tag, response.Tag)
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
	client, err := httpclient.ArtifactoryClientBuilder().
		SetInsecureTls(true).
		SetServiceDetails(&rtDetails).
		Build()
	assert.NoError(t, err, "Failed to create Artifactory client: %v\n")

	// Create docker promote service
	dockerPromoteService := services.NewDockerPromoteService(client)
	dockerPromoteService.ArtDetails = rtDetails
	return dockerPromoteService
}

func createTestParams() services.DockerPromoteParams {
	params := services.NewDockerPromoteParams(sourceRepo, targetRepo, dockerRepo)
	params.TargetDockerRepository = targetDockerRepo
	params.Tag = tag
	params.TargetTag = targetTag
	params.Copy = copy
	return params
}
