//go:build itest

package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
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
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/build/promote/"+buildName+"/"+buildNum, r.URL.Path)
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		response := services.BuildPromotionBody{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, srcRepo, response.SourceRepo)
		assert.Equal(t, trgRepo, response.TargetRepo)
		assert.Equal(t, copyVal, *response.Copy)

		w.WriteHeader(http.StatusOK)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	service := createPromoteService(t, ts.URL)

	err := service.BuildPromote(createTestPromotionParams())
	assert.NoError(t, err)
}

func TestBuildPromote_WithSlashInBuildName(t *testing.T) {
	initArtifactoryTest(t)

	buildNameWithSlash := "lib-systest/JIRA-1"
	buildNumberWithSlash := "JIRA_ID-1"

	expectedEscapedPath := "/api/build/promote/" +
		url.QueryEscape(buildNameWithSlash) + "/" +
		url.QueryEscape(buildNumberWithSlash)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, expectedEscapedPath, r.URL.EscapedPath())

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		response := services.BuildPromotionBody{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, srcRepo, response.SourceRepo)
		assert.Equal(t, trgRepo, response.TargetRepo)
		assert.Equal(t, copyVal, *response.Copy)

		w.WriteHeader(http.StatusOK)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	service := createPromoteService(t, ts.URL)

	params := createTestPromotionParams()
	params.BuildName = buildNameWithSlash
	params.BuildNumber = buildNumberWithSlash

	err := service.BuildPromote(params)
	assert.NoError(t, err)
}

func createPromoteService(t *testing.T, url string) *services.PromoteService {
	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl(url + "/")

	client, err := jfroghttpclient.JfrogClientBuilder().
		SetInsecureTls(true).
		SetClientCertPath(rtDetails.GetClientCertPath()).
		SetClientCertKeyPath(rtDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(rtDetails.RunPreRequestFunctions).
		Build()
	assert.NoError(t, err, "Failed to create Artifactory client: %v\n")

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
