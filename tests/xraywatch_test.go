package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	artifactoryServices "github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"github.com/jfrog/jfrog-client-go/xray/services"

	"github.com/stretchr/testify/assert"
)

func TestXrayWatchAll(t *testing.T) {
	policy1Name := fmt.Sprintf("%s-%d", "jfrog-policy1", time.Now().Unix())
	err := createPolicy(policy1Name)
	assert.NoError(t, err)
	defer deletePolicy(policy1Name)

	policy2Name := fmt.Sprintf("%s-%d", "jfrog-policy2", time.Now().Unix())
	err = createPolicy(policy2Name)
	assert.NoError(t, err)
	defer deletePolicy(policy2Name)

	AllWatchName := fmt.Sprintf("%s-%d", "jfrog-client-go-tests-watch-all-repos", time.Now().Unix())
	paramsAllRepos := services.NewXrayWatchParams()
	paramsAllRepos.Name = AllWatchName
	paramsAllRepos.Description = "All Repos"
	paramsAllRepos.Active = true

	paramsAllRepos.Repositories.Type = services.WatchRepositoriesAll
	paramsAllRepos.Repositories.All.Filters.PackageTypes = []string{"NpM", "maven"}
	paramsAllRepos.Repositories.All.Filters.Names = []string{"example-name-1"}
	paramsAllRepos.Repositories.All.Filters.Paths = []string{"example-path-1"}
	paramsAllRepos.Repositories.All.Filters.MimeTypes = []string{"example-mime-type-1"}
	paramsAllRepos.Repositories.All.Filters.Properties = map[string]string{"some-key-1": "some-value-1"}

	paramsAllRepos.Repositories.ExcludePatterns = []string{"excludePath1", "excludePath2"}
	paramsAllRepos.Repositories.IncludePatterns = []string{"includePath1", "includePath2"}

	paramsAllRepos.Builds.Type = services.WatchBuildAll
	paramsAllRepos.Builds.All.BinMgrID = "default"
	paramsAllRepos.Policies = []services.XrayWatchPolicy{
		{
			Name: policy1Name,
			Type: "security",
		},
		{
			Name: policy2Name,
			Type: "security",
		},
	}

	err = testsXrayWatchService.Create(paramsAllRepos)
	assert.NoError(t, err)
	defer testsXrayWatchService.Delete(paramsAllRepos.Name)

	validateWatchGeneralSettings(t, paramsAllRepos)
	targetConfig, err := testsXrayWatchService.Get(paramsAllRepos.Name)
	assert.NoError(t, err)

	assert.Equal(t, []string{"excludePath1", "excludePath2"}, targetConfig.Repositories.ExcludePatterns)
	assert.Equal(t, []string{"includePath1", "includePath2"}, targetConfig.Repositories.IncludePatterns)
	assert.Equal(t, []string{"Maven", "Npm"}, targetConfig.Repositories.All.Filters.PackageTypes)
	assert.Equal(t, []string{"example-name-1"}, targetConfig.Repositories.All.Filters.Names)
	assert.Equal(t, []string{"example-path-1"}, targetConfig.Repositories.All.Filters.Paths)
	assert.Equal(t, []string{"example-mime-type-1"}, targetConfig.Repositories.All.Filters.MimeTypes)
	assert.Equal(t, map[string]string{"some-key-1": "some-value-1"}, targetConfig.Repositories.All.Filters.Properties)
	assert.Equal(t, services.WatchRepositoriesAll, targetConfig.Repositories.Type)

	assert.Equal(t, services.WatchBuildAll, targetConfig.Builds.Type)
	assert.Equal(t, "default", targetConfig.Builds.All.BinMgrID)

	targetConfig.Description = "Updated Description"
	targetConfig.Repositories.All.Filters.PackageTypes = []string{"generic", "pypi"}
	targetConfig.Repositories.All.Filters.Names = []string{"example-name-2"}
	targetConfig.Repositories.All.Filters.Paths = []string{"example-path-2"}
	targetConfig.Repositories.All.Filters.MimeTypes = []string{"example-mime-type-2"}
	targetConfig.Repositories.All.Filters.Properties = map[string]string{"some-key-2": "some-value-2"}

	targetConfig.Repositories.ExcludePatterns = []string{"excludePath3", "excludePath4"}
	targetConfig.Repositories.IncludePatterns = []string{"includePath3", "includePath4"}

	targetConfig.Builds.Type = services.WatchBuildAll
	targetConfig.Builds.All.BinMgrID = "default"
	targetConfig.Policies = []services.XrayWatchPolicy{
		{
			Name: policy2Name,
			Type: "security",
		},
	}
	err = testsXrayWatchService.Update(*targetConfig)
	assert.NoError(t, err)

	validateWatchGeneralSettings(t, *targetConfig)
	updatedTargetConfig, err := testsXrayWatchService.Get(paramsAllRepos.Name)
	assert.NoError(t, err)

	assert.Equal(t, []string{"excludePath3", "excludePath4"}, updatedTargetConfig.Repositories.ExcludePatterns)
	assert.Equal(t, []string{"includePath3", "includePath4"}, updatedTargetConfig.Repositories.IncludePatterns)
	assert.Equal(t, []string{"Generic", "Pypi"}, updatedTargetConfig.Repositories.All.Filters.PackageTypes)
	assert.Equal(t, []string{"example-name-2"}, updatedTargetConfig.Repositories.All.Filters.Names)
	assert.Equal(t, []string{"example-path-2"}, updatedTargetConfig.Repositories.All.Filters.Paths)
	assert.Equal(t, []string{"example-mime-type-2"}, updatedTargetConfig.Repositories.All.Filters.MimeTypes)
	assert.Equal(t, map[string]string{"some-key-2": "some-value-2"}, updatedTargetConfig.Repositories.All.Filters.Properties)
}

func TestXrayWatchSelectedRepos(t *testing.T) {
	policy1Name := fmt.Sprintf("%s-%d", "jfrog-policy1-pattern", time.Now().Unix())
	err := createPolicy(policy1Name)
	assert.NoError(t, err)
	defer deletePolicy(policy1Name)

	repo1Name := fmt.Sprintf("%s-%d", "jfrog-repo1", time.Now().Unix())
	createRepo(t, repo1Name)
	defer deleteRepo(t, repo1Name)
	repo2Name := fmt.Sprintf("%s-%d", "jfrog-repo2", time.Now().Unix())
	createRepo(t, repo2Name)
	defer deleteRepo(t, repo2Name)

	build1Name := fmt.Sprintf("%s-%d", "jfrog-build1", time.Now().Unix())
	err = createBuild(build1Name)
	assert.NoError(t, err)
	defer deleteBuild(build1Name)

	build2Name := fmt.Sprintf("%s-%d", "jfrog-build2", time.Now().Unix())
	err = createBuild(build2Name)
	assert.NoError(t, err)
	defer deleteBuild(build2Name)

	paramsSelectedRepos := services.NewXrayWatchParams()
	paramsSelectedRepos.Name = fmt.Sprintf("%s-%d", "jfrog-client-go-tests-watch-selected-repos", time.Now().Unix())
	paramsSelectedRepos.Description = "Selected Repos"
	paramsSelectedRepos.Active = true
	paramsSelectedRepos.Policies = []services.XrayWatchPolicy{
		{
			Name: policy1Name,
			Type: "security",
		},
	}

	var repos = map[string]services.XrayWatchRepository{}
	repo := services.NewXrayWatchRepository(repo1Name, "default")
	repo.Filters.PackageTypes = []string{"npm", "maven"}
	repo.Filters.Names = []string{"example-name"}
	repo.Filters.Paths = []string{"example-path"}
	repo.Filters.MimeTypes = []string{"example-mime-type"}
	repo.Filters.Properties = map[string]string{"some-key": "some-value"}

	repos[repo1Name] = repo

	anotherRepo := services.NewXrayWatchRepository(repo2Name, "default")
	anotherRepo.Filters.PackageTypes = []string{"nuget"}
	anotherRepo.Filters.Names = []string{"another-example-name"}
	anotherRepo.Filters.Paths = []string{"another-example-path"}
	anotherRepo.Filters.MimeTypes = []string{"another-example-mime-type"}
	anotherRepo.Filters.Properties = map[string]string{"another-key": "some-value"}

	repos[repo2Name] = anotherRepo

	paramsSelectedRepos.Repositories.Type = services.WatchRepositoriesByName
	paramsSelectedRepos.Repositories.Repositories = repos
	paramsSelectedRepos.Repositories.ExcludePatterns = []string{"selectedExcludePath1", "selectedExcludePath2"}
	paramsSelectedRepos.Repositories.IncludePatterns = []string{"selectedIncludePath1", "selectedIncludePath2"}

	paramsSelectedRepos.Builds.Type = services.WatchBuildByName
	paramsSelectedRepos.Builds.ByNames = map[string]services.XrayWatchBuildsByNameParams{}
	paramsSelectedRepos.Builds.ByNames[build1Name] = services.XrayWatchBuildsByNameParams{
		Name:     build1Name,
		BinMgrID: "default",
	}
	err = testsXrayWatchService.Create(paramsSelectedRepos)
	assert.NoError(t, err)
	defer testsXrayWatchService.Delete(paramsSelectedRepos.Name)
	validateWatchGeneralSettings(t, paramsSelectedRepos)

	targetConfig, err := testsXrayWatchService.Get(paramsSelectedRepos.Name)
	assert.NoError(t, err)
	assert.Equal(t, []string{"selectedExcludePath1", "selectedExcludePath2"}, targetConfig.Repositories.ExcludePatterns)
	assert.Equal(t, []string{"selectedIncludePath1", "selectedIncludePath2"}, targetConfig.Repositories.IncludePatterns)
	assert.Equal(t, services.WatchRepositoriesByName, targetConfig.Repositories.Type)

	assert.Equal(t, repo1Name, targetConfig.Repositories.Repositories[repo1Name].Name)
	assert.Equal(t, "default", targetConfig.Repositories.Repositories[repo1Name].BinMgrID)
	assert.Equal(t, []string{"Maven", "Npm"}, targetConfig.Repositories.Repositories[repo1Name].Filters.PackageTypes)
	assert.Equal(t, []string{"example-name"}, targetConfig.Repositories.Repositories[repo1Name].Filters.Names)
	assert.Equal(t, []string{"example-path"}, targetConfig.Repositories.Repositories[repo1Name].Filters.Paths)
	assert.Equal(t, []string{"example-mime-type"}, targetConfig.Repositories.Repositories[repo1Name].Filters.MimeTypes)
	assert.Equal(t, map[string]string{"some-key": "some-value"}, targetConfig.Repositories.Repositories[repo1Name].Filters.Properties)

	assert.Equal(t, repo2Name, targetConfig.Repositories.Repositories[repo2Name].Name)
	assert.Equal(t, "default", targetConfig.Repositories.Repositories[repo2Name].BinMgrID)
	assert.Equal(t, []string{"NuGet"}, targetConfig.Repositories.Repositories[repo2Name].Filters.PackageTypes)
	assert.Equal(t, []string{"another-example-name"}, targetConfig.Repositories.Repositories[repo2Name].Filters.Names)
	assert.Equal(t, []string{"another-example-path"}, targetConfig.Repositories.Repositories[repo2Name].Filters.Paths)
	assert.Equal(t, []string{"another-example-mime-type"}, targetConfig.Repositories.Repositories[repo2Name].Filters.MimeTypes)
	assert.Equal(t, map[string]string{"another-key": "some-value"}, targetConfig.Repositories.Repositories[repo2Name].Filters.Properties)

	assert.Equal(t, services.WatchBuildByName, targetConfig.Builds.Type)
	assert.Empty(t, targetConfig.Builds.All.ExcludePatterns)
	assert.Empty(t, targetConfig.Builds.All.IncludePatterns)

	assert.Equal(t, build1Name, targetConfig.Builds.ByNames[build1Name].Name)
	assert.Equal(t, "default", targetConfig.Builds.ByNames[build1Name].BinMgrID)

	targetConfig.Repositories.ExcludePatterns = []string{"excludePath-2"}
	targetConfig.Repositories.IncludePatterns = []string{"includePath-2", "fake-2"}
	targetConfig.Builds.ByNames[build2Name] = services.XrayWatchBuildsByNameParams{
		Name:     build2Name,
		BinMgrID: "default",
	}

	delete(targetConfig.Repositories.Repositories, repo2Name)

	updatedRepo1 := targetConfig.Repositories.Repositories[repo1Name]

	updatedRepo1.Filters.PackageTypes = []string{"Generic"}
	updatedRepo1.Filters.Names = []string{"example-name-2"}
	updatedRepo1.Filters.Paths = []string{"example-path-2"}
	updatedRepo1.Filters.MimeTypes = []string{"example-mime-type-2"}
	updatedRepo1.Filters.Properties = map[string]string{"some-key": "some-value-2"}

	targetConfig.Repositories.Repositories[repo1Name] = updatedRepo1

	err = testsXrayWatchService.Update(*targetConfig)
	assert.NoError(t, err)

	validateWatchGeneralSettings(t, *targetConfig)
	updatedTargetConfig, err := testsXrayWatchService.Get(paramsSelectedRepos.Name)
	assert.NoError(t, err)

	assert.Equal(t, []string{"excludePath-2"}, updatedTargetConfig.Repositories.ExcludePatterns)
	assert.Equal(t, []string{"fake-2", "includePath-2"}, updatedTargetConfig.Repositories.IncludePatterns)
	assert.Empty(t, updatedTargetConfig.Repositories.Repositories[repo2Name])

	assert.Equal(t, repo1Name, updatedTargetConfig.Repositories.Repositories[repo1Name].Name)
	assert.Equal(t, []string{"Generic"}, updatedTargetConfig.Repositories.Repositories[repo1Name].Filters.PackageTypes)
	assert.Equal(t, []string{"example-name-2"}, updatedTargetConfig.Repositories.Repositories[repo1Name].Filters.Names)
	assert.Equal(t, []string{"example-path-2"}, updatedTargetConfig.Repositories.Repositories[repo1Name].Filters.Paths)
	assert.Equal(t, []string{"example-mime-type-2"}, updatedTargetConfig.Repositories.Repositories[repo1Name].Filters.MimeTypes)
	assert.Equal(t, map[string]string{"some-key": "some-value-2"}, updatedTargetConfig.Repositories.Repositories[repo1Name].Filters.Properties)

}

func TestXrayWatchBuildsByPattern(t *testing.T) {
	policy1Name := fmt.Sprintf("%s-%d", "jfrog-policy1-pattern", time.Now().Unix())
	err := createPolicy(policy1Name)
	assert.NoError(t, err)
	defer deletePolicy(policy1Name)

	paramsBuildsByPattern := services.NewXrayWatchParams()
	paramsBuildsByPattern.Name = fmt.Sprintf("%s-%d", "jfrog-client-go-tests-watch-builds-by-pattern", time.Now().Unix())
	paramsBuildsByPattern.Description = "Builds By Pattern"
	paramsBuildsByPattern.Builds.Type = services.WatchBuildAll
	paramsBuildsByPattern.Builds.All.ExcludePatterns = []string{"excludePath"}
	paramsBuildsByPattern.Builds.All.IncludePatterns = []string{"includePath", "fake"}
	paramsBuildsByPattern.Builds.All.BinMgrID = "default"
	paramsBuildsByPattern.Policies = []services.XrayWatchPolicy{
		{
			Name: policy1Name,
			Type: "security",
		},
	}

	err = testsXrayWatchService.Create(paramsBuildsByPattern)
	assert.NoError(t, err)
	defer testsXrayWatchService.Delete(paramsBuildsByPattern.Name)
	validateWatchGeneralSettings(t, paramsBuildsByPattern)

	targetConfig, err := testsXrayWatchService.Get(paramsBuildsByPattern.Name)
	assert.NoError(t, err)
	assert.Equal(t, services.WatchBuildAll, targetConfig.Builds.Type)
	assert.Equal(t, []string{"excludePath"}, targetConfig.Builds.All.ExcludePatterns)
	assert.Equal(t, []string{"includePath", "fake"}, targetConfig.Builds.All.IncludePatterns)

	targetConfig.Builds.All.ExcludePatterns = []string{"excludePath-2"}
	targetConfig.Builds.All.IncludePatterns = []string{"includePath-2", "fake-2"}

	err = testsXrayWatchService.Update(*targetConfig)
	assert.NoError(t, err)

	validateWatchGeneralSettings(t, *targetConfig)
	updatedTargetConfig, err := testsXrayWatchService.Get(paramsBuildsByPattern.Name)
	assert.NoError(t, err)

	assert.Equal(t, []string{"excludePath-2"}, updatedTargetConfig.Builds.All.ExcludePatterns)
	assert.Equal(t, []string{"includePath-2", "fake-2"}, updatedTargetConfig.Builds.All.IncludePatterns)
}

func validateWatchGeneralSettings(t *testing.T, params services.XrayWatchParams) {
	targetConfig, err := testsXrayWatchService.Get(params.Name)
	assert.NoError(t, err)
	assert.Equal(t, params.Name, targetConfig.Name)
	assert.Equal(t, params.Description, targetConfig.Description)
	assert.Equal(t, params.Active, targetConfig.Active)
	assert.Equal(t, params.Policies, targetConfig.Policies)
	return
}

func createRepo(t *testing.T, repoKey string) {
	glp := artifactoryServices.NewGenericLocalRepositoryParams()
	glp.Key = repoKey
	glp.XrayIndex = &trueValue

	err := testsCreateLocalRepositoryService.Generic(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
}

func createBuild(buildName string) error {
	artDetails := GetRtDetails()
	artHTTPDetails := artDetails.CreateHttpClientDetails()

	utils.SetContentType("application/json", &artHTTPDetails.Headers)
	artClient, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}

	xrayDetails := GetXrayDetails()
	xrayHTTPDetails := xrayDetails.CreateHttpClientDetails()

	utils.SetContentType("application/json", &xrayHTTPDetails.Headers)
	xrayClient, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}

	dataArtifactoryBuild := ArtifactoryBuild{
		Name:       buildName,
		Version:    "1.0.0",
		Number:     "2",
		Started:    "2014-09-30T12:00:19.893+0300",
		Properties: map[string]interface{}{},
		Modules: []ArtifactoryModule{
			{
				ID: "example-mdule",
				Artifacts: []ArtifactoryArtifact{
					{
						Type: "gz",
						Sha1: "9d4336ff7bc2d2348aee4e27ad55e42110df4a80",
						Md5:  "b4918187cc9b3bf1b0772546d9398d7d",
						Name: "c.tar.gz",
					},
				},
			},
		},
	}
	requestContentArtifactoryBuild, err := json.Marshal(dataArtifactoryBuild)
	if err != nil {
		return errors.New("failed marshalling build " + buildName)
	}

	resp, _, err := artClient.SendPut(artDetails.GetUrl()+"api/build", requestContentArtifactoryBuild, artHTTPDetails)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errors.New("Status is not OK or NoContent - " + strconv.Itoa(resp.StatusCode))
	}

	// the build needs to be indexed before a watch can be associated with it.
	dataIndexBuild := struct {
		Names []string `json:"names"`
	}{
		Names: []string{buildName},
	}

	requestContentIndexBuild, err := json.Marshal(dataIndexBuild)

	resp, _, err = xrayClient.SendPost(xrayDetails.GetUrl()+"api/v1/binMgr/builds", requestContentIndexBuild, artHTTPDetails)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func deleteBuildIndex(buildName string) error {
	xrayDetails := GetXrayDetails()
	artHTTPDetails := xrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &artHTTPDetails.Headers)
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return nil
	}

	dataIndexBuild := struct {
		Names []string `json:"indexed_builds"`
	}{
		Names: []string{},
	}

	requestContentIndexBuild, err := json.Marshal(dataIndexBuild)

	resp, _, err := client.SendPut(xrayDetails.GetUrl()+"api/v1/binMgr/default/builds", requestContentIndexBuild, artHTTPDetails)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func deleteBuild(buildName string) error {
	err := deleteBuildIndex(buildName)
	if err != nil {
		return nil
	}

	artDetails := GetRtDetails()
	artHTTPDetails := artDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return nil
	}

	resp, _, err := client.SendDelete(artDetails.GetUrl()+"api/build/"+buildName+"?deleteAll=1", nil, artHTTPDetails)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("failed unmarshalling build " + resp.Status)
	}

	return nil
}

type ArtifactoryBuild struct {
	Version    string              `json:"version"`
	Name       string              `json:"name"`
	Number     string              `json:"number"`
	Started    string              `json:"started"`
	Properties interface{}         `json:"properties"`
	Modules    []ArtifactoryModule `json:"modules"`
}

type ArtifactoryModule struct {
	ID        string                `json:"id"`
	Artifacts []ArtifactoryArtifact `json:"artifacts"`
}

type ArtifactoryArtifact struct {
	Type string `json:"type"`
	Sha1 string `json:"sha1"`
	Md5  string `json:"md5"`
	Name string `json:"name"`
}

func createPolicy(policyName string) error {
	xrayDetails := GetXrayDetails()
	xrayHTTPDetails := xrayDetails.CreateHttpClientDetails()

	utils.SetContentType("application/json", &xrayHTTPDetails.Headers)
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	data := ArtifactoryPolicy{
		Name:        policyName,
		Description: "example policy",
		Type:        "security",
		Rules: []ArtifactoryPolicyRules{{
			Name:     "sec_rule",
			Priority: 1,
			Criteria: map[string]string{
				"min_severity": "medium",
			},
			Actions: ArtifactoryPolicyActions{
				Webhooks: []string{},
				BlockDownload: ArtifactoryPolicyActionsBlockDownload{
					Active:    true,
					Unscanned: false,
				},
				BlockReleaseBundleDistribution: true,
				FailBuild:                      true,
				NotifyDeployer:                 true,
				NotifyWatchRecipients:          true,
			},
		}},
	}

	requestContent, err := json.Marshal(data)
	if err != nil {
		return errors.New("failed marshalling policy " + policyName)
	}

	resp, _, err := client.SendPost(xrayDetails.GetUrl()+"api/v2/policies", requestContent, xrayHTTPDetails)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("Status is not Created - " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}

func deletePolicy(policyName string) error {
	xrayDetails := GetXrayDetails()
	xrayHTTPDetails := xrayDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &xrayHTTPDetails.Headers)
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return nil
	}

	resp, _, err := client.SendDelete(xrayDetails.GetUrl()+"api/v2/policies/"+policyName, nil, xrayHTTPDetails)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("failed to delete policy " + resp.Status)
	}
	return nil
}

type ArtifactoryPolicy struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Type        string                   `json:"type"`
	Rules       []ArtifactoryPolicyRules `json:"rules"`
}

type ArtifactoryPolicyRules struct {
	Name     string                   `json:"name"`
	Priority int                      `json:"priority"`
	Criteria map[string]string        `json:"criteria"`
	Actions  ArtifactoryPolicyActions `json:"actions"`
}

type ArtifactoryPolicyActions struct {
	Webhooks                       []string                              `json:"webhooks"`
	BlockDownload                  ArtifactoryPolicyActionsBlockDownload `json:"block_download"`
	BlockReleaseBundleDistribution bool                                  `json:"block_release_bundle_distribution"`
	FailBuild                      bool                                  `json:"fail_build"`
	NotifyDeployer                 bool                                  `json:"notify_deployer"`
	NotifyWatchRecipients          bool                                  `json:"notify_watch_recipients"`
}

type ArtifactoryPolicyActionsBlockDownload struct {
	Active    bool `json:"active"`
	Unscanned bool `json:"unscanned"`
}
