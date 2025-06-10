package tests

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	srcRepo   = "maven-dev-virtual"
	trgRepo   = "maven-release-local"
	buildName = "lib-systest"
	buildNum  = "feature/JIRA_ID-1"
	status    = "Released"
	copyVal   = false
)

func TestBuildPromote(t *testing.T) {
	initArtifactoryTest(t)
	// Create mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/build/promote/"+buildName+"/"+buildNum, r.URL.Path)

		// Check body
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		response := services.BuildPromotionBody{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, srcRepo, response.SourceRepo)
		assert.Equal(t, trgRepo, response.TargetRepo)
		assert.Equal(t, copyVal, *response.Copy)

		// Send response 200 OK
		w.WriteHeader(http.StatusOK)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Promote build
	service := createPromoteService(t, ts.URL)

	err := service.BuildPromote(createTestPromotionParams())
	assert.NoError(t, err)
}

func createPromoteService(t *testing.T, url string) *services.PromoteService {
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

	// Create build promote service
	buildPromoteService := services.NewPromotionService(client)
	buildPromoteService.ArtDetails = rtDetails
	return buildPromoteService
}

func createTestPromotionParams() services.PromotionParams {
	params := services.NewPromotionParams()
	params.BuildName = buildName
	params.BuildNumber = buildNum
	params.TargetRepo = trgRepo
	params.Status = status
	params.SourceRepo = srcRepo

	return params
}
