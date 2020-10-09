package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/httpclient"

	"github.com/stretchr/testify/assert"
)

func TestXrayWatch(t *testing.T) {
	policyName := fmt.Sprintf("%s-%d", "fake-policy", time.Now().Unix())
	err := addFakePolicy(policyName)
	assert.NoError(t, err)
	defer deletePolicy(policyName)

	buildName := fmt.Sprintf("%s-%d", "fake-build", time.Now().Unix())
	err = addFakeBuild(buildName)
	assert.NoError(t, err)
	defer deleteBuild(buildName)

	AllWatchName := fmt.Sprintf("%s-%d", "jfrog-client-go-tests-watch-all-repos", time.Now().Unix())
	paramsAllRepos := services.NewXrayWatchParams()
	paramsAllRepos.Name = AllWatchName
	paramsAllRepos.Description = "All Repos"
	paramsAllRepos.Active = true

	paramsAllRepos.Repositories.Type = services.WatchRepositoriesAll
	paramsAllRepos.Repositories.All.Filters.PackageTypes = []string{"npm", "maven"}
	paramsAllRepos.Repositories.All.Filters.Names = []string{"example-name"}
	paramsAllRepos.Repositories.All.Filters.Paths = []string{"example-path"}
	paramsAllRepos.Repositories.All.Filters.MimeTypes = []string{"example-mime-type"}
	paramsAllRepos.Repositories.All.Filters.Properties = map[string]string{"some-key": "some-value"}

	paramsAllRepos.Repositories.ExcludePatterns = []string{"excludePath1", "excludePath2"}
	paramsAllRepos.Repositories.IncludePatterns = []string{"includePath1", "includePath2"}

	paramsAllRepos.Builds.Type = services.WatchBuildAll
	paramsAllRepos.Builds.All.Bin_Mgr_ID = "default"
	paramsAllRepos.Policies = []services.XrayWatchPolicy{{
		Name: policyName,
		Type: "security",
	}}

	err = testsXrayWatchService.Create(paramsAllRepos)
	assert.NoError(t, err)
	defer testsXrayWatchService.Delete(paramsAllRepos.Name)

	validateWatchGeneralSettings(t, paramsAllRepos)
	targetConfig, err := testsXrayWatchService.Get(paramsAllRepos.Name)
	assert.Equal(t, []string{"excludePath1", "excludePath2"}, targetConfig.Repositories.ExcludePatterns)
	assert.Equal(t, []string{"includePath1", "includePath2"}, targetConfig.Repositories.IncludePatterns)
	assert.Equal(t, []string{"Maven", "Npm"}, targetConfig.Repositories.All.Filters.PackageTypes)
	assert.Equal(t, []string{"example-name"}, targetConfig.Repositories.All.Filters.Names)
	assert.Equal(t, []string{"example-path"}, targetConfig.Repositories.All.Filters.Paths)
	assert.Equal(t, []string{"example-mime-type"}, targetConfig.Repositories.All.Filters.MimeTypes)
	assert.Equal(t, map[string]string{"some-key": "some-value"}, targetConfig.Repositories.All.Filters.Properties)
	assert.Equal(t, services.WatchRepositoriesAll, targetConfig.Repositories.Type)

	assert.Equal(t, services.WatchBuildAll, targetConfig.Builds.Type)
	assert.Equal(t, "default", targetConfig.Builds.All.Bin_Mgr_ID)

	paramsSelectedRepos := services.NewXrayWatchParams()
	paramsSelectedRepos.Name = fmt.Sprintf("%s-%d", "jfrog-client-go-tests-watch-selected-repos", time.Now().Unix())
	paramsSelectedRepos.Description = "Selected Repos"
	paramsSelectedRepos.Active = true

	// Todo: update repository name
	// Repository must exist
	var repos = map[string]services.XrayWatchRepository{}
	repo := services.NewXrayWatchRepository("example-repo-local", "default")
	repo.Filters.PackageTypes = []string{"npm", "maven"}
	repo.Filters.Names = []string{"example-name"}
	repo.Filters.Paths = []string{"example-path"}
	repo.Filters.MimeTypes = []string{"example-mime-type"}
	repo.Filters.Properties = map[string]string{"some-key": "some-value"}

	repos["example-repo-local"] = repo

	anotherRepo := services.NewXrayWatchRepository("another-repo", "default")
	anotherRepo.Filters.PackageTypes = []string{"nuget"}
	anotherRepo.Filters.Names = []string{"another-example-name"}
	anotherRepo.Filters.Paths = []string{"another-example-path"}
	anotherRepo.Filters.MimeTypes = []string{"another-example-mime-type"}
	anotherRepo.Filters.Properties = map[string]string{"another-key": "some-value"}

	repos["another-repo"] = anotherRepo

	paramsSelectedRepos.Repositories.Type = services.WatchRepositoriesByName
	paramsSelectedRepos.Repositories.Repositories = repos
	paramsSelectedRepos.Repositories.ExcludePatterns = []string{"selectedExcludePath1", "selectedExcludePath2"}
	paramsSelectedRepos.Repositories.IncludePatterns = []string{"selectedIncludePath1", "selectedIncludePath2"}

	paramsSelectedRepos.Builds.Type = services.WatchBuildByName
	paramsSelectedRepos.Builds.ByNames = map[string]services.XrayWatchBuildsByNameParams{}
	paramsSelectedRepos.Builds.ByNames[buildName] = services.XrayWatchBuildsByNameParams{
		Name:       buildName,
		Bin_Mgr_ID: "default",
	}
	err = testsXrayWatchService.Create(paramsSelectedRepos)
	assert.NoError(t, err)
	defer testsXrayWatchService.Delete(paramsSelectedRepos.Name)
	validateWatchGeneralSettings(t, paramsSelectedRepos)
	targetConfig, err = testsXrayWatchService.Get(paramsSelectedRepos.Name)
	assert.Equal(t, []string{"selectedExcludePath1", "selectedExcludePath2"}, targetConfig.Repositories.ExcludePatterns)
	assert.Equal(t, []string{"selectedIncludePath1", "selectedIncludePath2"}, targetConfig.Repositories.IncludePatterns)
	assert.Equal(t, services.WatchRepositoriesByName, targetConfig.Repositories.Type)

	assert.Equal(t, "example-repo-local", targetConfig.Repositories.Repositories["example-repo-local"].Name)
	assert.Equal(t, "default", targetConfig.Repositories.Repositories["example-repo-local"].Bin_Mgr_ID)
	assert.Equal(t, []string{"Maven", "Npm"}, targetConfig.Repositories.Repositories["example-repo-local"].Filters.PackageTypes)
	assert.Equal(t, []string{"example-name"}, targetConfig.Repositories.Repositories["example-repo-local"].Filters.Names)
	assert.Equal(t, []string{"example-path"}, targetConfig.Repositories.Repositories["example-repo-local"].Filters.Paths)
	assert.Equal(t, []string{"example-mime-type"}, targetConfig.Repositories.Repositories["example-repo-local"].Filters.MimeTypes)
	assert.Equal(t, map[string]string{"some-key": "some-value"}, targetConfig.Repositories.Repositories["example-repo-local"].Filters.Properties)

	assert.Equal(t, "another-repo", targetConfig.Repositories.Repositories["another-repo"].Name)
	assert.Equal(t, "default", targetConfig.Repositories.Repositories["another-repo"].Bin_Mgr_ID)
	assert.Equal(t, []string{"NuGet"}, targetConfig.Repositories.Repositories["another-repo"].Filters.PackageTypes)
	assert.Equal(t, []string{"another-example-name"}, targetConfig.Repositories.Repositories["another-repo"].Filters.Names)
	assert.Equal(t, []string{"another-example-path"}, targetConfig.Repositories.Repositories["another-repo"].Filters.Paths)
	assert.Equal(t, []string{"another-example-mime-type"}, targetConfig.Repositories.Repositories["another-repo"].Filters.MimeTypes)
	assert.Equal(t, map[string]string{"another-key": "some-value"}, targetConfig.Repositories.Repositories["another-repo"].Filters.Properties)

	assert.Equal(t, services.WatchBuildByName, targetConfig.Builds.Type)
	assert.Empty(t, targetConfig.Builds.All.ExcludePatterns)
	assert.Empty(t, targetConfig.Builds.All.IncludePatterns)

	assert.Equal(t, buildName, targetConfig.Builds.ByNames[buildName].Name)
	assert.Equal(t, "default", targetConfig.Builds.ByNames[buildName].Bin_Mgr_ID)

	paramsBuildsByPattern := services.NewXrayWatchParams()
	paramsBuildsByPattern.Name = fmt.Sprintf("%s-%d", "jfrog-client-go-tests-watch-builds-by-pattern", time.Now().Unix())
	paramsBuildsByPattern.Description = "Builds By Pattern"
	paramsBuildsByPattern.Builds.Type = services.WatchBuildAll
	paramsBuildsByPattern.Builds.All.ExcludePatterns = []string{"excludePath"}
	paramsBuildsByPattern.Builds.All.IncludePatterns = []string{"includePath", "fake"}
	paramsBuildsByPattern.Builds.All.Bin_Mgr_ID = "default"
	err = testsXrayWatchService.Create(paramsBuildsByPattern)
	assert.NoError(t, err)
	defer testsXrayWatchService.Delete(paramsBuildsByPattern.Name)
	validateWatchGeneralSettings(t, paramsBuildsByPattern)
	targetConfig, err = testsXrayWatchService.Get(paramsBuildsByPattern.Name)
	assert.Equal(t, services.WatchBuildAll, targetConfig.Builds.Type)
	assert.Equal(t, []string{"excludePath"}, targetConfig.Builds.All.ExcludePatterns)
	assert.Equal(t, []string{"includePath", "fake"}, targetConfig.Builds.All.IncludePatterns)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
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

// func getPermissionTarget(targetName string) (targetParams *services.PermissionTargetParams, err error) {
// 	artDetails := GetRtDetails()
// 	artHttpDetails := artDetails.CreateHttpClientDetails()
// 	client, err := httpclient.ClientBuilder().Build()
// 	if err != nil {
// 		return
// 	}
// 	resp, body, _, err := client.SendGet(artDetails.GetUrl()+"api/v2/security/permissions/"+targetName, false, artHttpDetails)
// 	if err != nil || resp.StatusCode != http.StatusOK {
// 		return
// 	}
// 	if err = json.Unmarshal(body, &targetParams); err != nil {
// 		return nil, errors.New("failed unmarshalling permission target " + targetName)
// 	}
// 	return
// }

func addFakeBuild(buildName string) error {
	artDetails := GetRtDetails()
	artHTTPDetails := artDetails.CreateHttpClientDetails()

	utils.SetContentType("application/json", &artHTTPDetails.Headers)
	client, err := httpclient.ClientBuilder().Build()
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

	// TODO: Update URL
	resp, _, err := client.SendPut("https://artifactoryprovider.jfrog.io/artifactory/api/build", requestContentArtifactoryBuild, artHTTPDetails)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errors.New("Status is not OK or NoContent - " + strconv.Itoa(resp.StatusCode))
	}

	dataIndexBuild := struct {
		Names []string `json:"names"`
	}{
		Names: []string{buildName},
	}

	requestContentIndexBuild, err := json.Marshal(dataIndexBuild)

	// TODO: Update URL
	resp, _, err = client.SendPost("https://artifactoryprovider.jfrog.io/xray/api/v1/binMgr/builds", requestContentIndexBuild, artHTTPDetails)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func deleteBuildIndex(buildName string) error {
	artDetails := GetRtDetails()
	artHTTPDetails := artDetails.CreateHttpClientDetails()
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

	// TODO: Update URL
	resp, _, err := client.SendPut("https://artifactoryprovider.jfrog.io/xray/api/v1/binMgr/default/builds", requestContentIndexBuild, artHTTPDetails)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func deleteBuild(buildName string) error {
	deleteBuildIndex(buildName)

	artDetails := GetRtDetails()
	artHTTPDetails := artDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return nil
	}

	// TODO: Update URL
	resp, _, err := client.SendDelete("https://artifactoryprovider.jfrog.io/artifactory/api/build/"+buildName+"?deleteAll=1", nil, artHTTPDetails)
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

func addFakePolicy(policyName string) error {
	artDetails := GetRtDetails()
	artHTTPDetails := artDetails.CreateHttpClientDetails()

	utils.SetContentType("application/json", &artHTTPDetails.Headers)
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

	// TODO: Update URL
	resp, _, err := client.SendPost("https://artifactoryprovider.jfrog.io/xray/api/v2/policies", requestContent, artHTTPDetails)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("Status is not Created - " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}

func deletePolicy(policyName string) error {
	artDetails := GetRtDetails()
	artHTTPDetails := artDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &artHTTPDetails.Headers)
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return nil
	}
	// TODO: Update URL
	resp, _, err := client.SendDelete("https://artifactoryprovider.jfrog.io/xray/api/v2/policies/"+policyName, nil, artHTTPDetails)
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
