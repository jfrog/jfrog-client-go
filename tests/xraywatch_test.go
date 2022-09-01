package tests

import (
	"fmt"
	artifactoryServices "github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXrayWatch(t *testing.T) {
	initXrayTest(t)
	t.Run("testXrayWatchAll", testXrayWatchAll)
	t.Run("testXrayWatchSelectedRepos", testXrayWatchSelectedRepos)
	t.Run("testXrayWatchBuildsByPattern", testXrayWatchBuildsByPattern)
	t.Run("testXrayWatchUpdateMissingWatch", testXrayWatchUpdateMissingWatch)
	t.Run("testXrayWatchDeleteMissingWatch", testXrayWatchDeleteMissingWatch)
	t.Run("testXrayWatchGetMissingWatch", testXrayWatchGetMissingWatch)
}

func testXrayWatchAll(t *testing.T) {
	policy1Name := fmt.Sprintf("%s-%s", "policy1", getRunId())
	err := createDummyPolicy(policy1Name)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, testsXrayPolicyService.Delete(policy1Name))
	}()
	policy2Name := fmt.Sprintf("%s-%s", "policy2", getRunId())
	err = createDummyPolicy(policy2Name)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, testsXrayPolicyService.Delete(policy2Name))
	}()
	AllWatchName := fmt.Sprintf("%s-%s", "client-go-tests-watch-all-repos", getRunId())
	paramsAllRepos := utils.NewWatchParams()
	paramsAllRepos.Name = AllWatchName
	paramsAllRepos.Description = "All Repos"
	paramsAllRepos.Active = true

	paramsAllRepos.Repositories.Type = utils.WatchRepositoriesAll
	paramsAllRepos.Repositories.All.Filters.PackageTypes = []string{"NpM", "maven"}
	paramsAllRepos.Repositories.All.Filters.Names = []string{"example-name-1"}
	paramsAllRepos.Repositories.All.Filters.Paths = []string{"example-path-1"}
	paramsAllRepos.Repositories.All.Filters.MimeTypes = []string{"example-mime-type-1"}
	paramsAllRepos.Repositories.All.Filters.Properties = map[string]string{"some-key-1": "some-value-1", "some-key-2": "some-value-2"}

	paramsAllRepos.Repositories.ExcludePatterns = []string{"excludePath1", "excludePath2"}
	paramsAllRepos.Repositories.IncludePatterns = []string{"includePath1", "includePath2"}

	paramsAllRepos.Builds.Type = utils.WatchBuildAll
	paramsAllRepos.Builds.All.BinMgrID = "default"
	paramsAllRepos.Policies = []utils.AssignedPolicy{
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
	defer func() {
		assert.NoError(t, testsXrayWatchService.Delete(paramsAllRepos.Name))
	}()
	validateWatchGeneralSettings(t, paramsAllRepos)
	targetConfig, err := testsXrayWatchService.Get(paramsAllRepos.Name)
	assert.NoError(t, err)

	assert.Equal(t, []string{"excludePath1", "excludePath2"}, targetConfig.Repositories.ExcludePatterns)
	assert.Equal(t, []string{"includePath1", "includePath2"}, targetConfig.Repositories.IncludePatterns)
	assert.Equal(t, []string{"Maven", "Npm"}, targetConfig.Repositories.All.Filters.PackageTypes)
	assert.Equal(t, []string{"example-name-1"}, targetConfig.Repositories.All.Filters.Names)
	assert.Equal(t, []string{"example-path-1"}, targetConfig.Repositories.All.Filters.Paths)
	assert.Equal(t, []string{"example-mime-type-1"}, targetConfig.Repositories.All.Filters.MimeTypes)
	assert.Equal(t, map[string]string{"some-key-1": "some-value-1", "some-key-2": "some-value-2"}, targetConfig.Repositories.All.Filters.Properties)
	assert.Equal(t, utils.WatchRepositoriesAll, targetConfig.Repositories.Type)

	assert.Equal(t, utils.WatchBuildAll, targetConfig.Builds.Type)
	assert.Equal(t, "default", targetConfig.Builds.All.BinMgrID)

	targetConfig.Description = "Updated Description"
	targetConfig.Repositories.All.Filters.PackageTypes = []string{"generic", "pypi"}
	targetConfig.Repositories.All.Filters.Names = []string{"example-name-2"}
	targetConfig.Repositories.All.Filters.Paths = []string{"example-path-2"}
	targetConfig.Repositories.All.Filters.MimeTypes = []string{"example-mime-type-2"}
	targetConfig.Repositories.All.Filters.Properties = map[string]string{"some-key-2": "some-value-2", "some-key-4": "some-value-4"}

	targetConfig.Repositories.ExcludePatterns = []string{"excludePath3", "excludePath4"}
	targetConfig.Repositories.IncludePatterns = []string{"includePath3", "includePath4"}

	targetConfig.Builds.Type = utils.WatchBuildAll
	targetConfig.Builds.All.BinMgrID = "default"
	targetConfig.Policies = []utils.AssignedPolicy{
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
	assert.Equal(t, map[string]string{"some-key-2": "some-value-2", "some-key-4": "some-value-4"}, updatedTargetConfig.Repositories.All.Filters.Properties)
}

func testXrayWatchSelectedRepos(t *testing.T) {
	policy1Name := fmt.Sprintf("%s-%s", "policy1-pattern", getRunId())
	err := createDummyPolicy(policy1Name)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, testsXrayPolicyService.Delete(policy1Name))
	}()
	repo1Name := fmt.Sprintf("%s-%s", "repo1", getRunId())
	createRepoLocal(t, repo1Name)
	defer deleteRepo(t, repo1Name)
	repo2Name := fmt.Sprintf("%s-%s", "repo2", getRunId())
	createRepoRemote(t, repo2Name)
	defer deleteRepo(t, repo2Name)

	build1Name := fmt.Sprintf("%s-%s", "build1", getRunId())
	err = createAndIndexBuild(t, build1Name)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, deleteBuild(build1Name))
	}()
	build2Name := fmt.Sprintf("%s-%s", "build2", getRunId())
	err = createAndIndexBuild(t, build2Name)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, deleteBuild(build2Name))
	}()
	paramsSelectedRepos := utils.NewWatchParams()
	paramsSelectedRepos.Name = fmt.Sprintf("%s-%s", "client-go-tests-watch-selected-repos", getRunId())
	paramsSelectedRepos.Description = "Selected Repos"
	paramsSelectedRepos.Active = true
	paramsSelectedRepos.Policies = []utils.AssignedPolicy{
		{
			Name: policy1Name,
			Type: "security",
		},
	}

	var repos = map[string]utils.WatchRepository{}
	repo := utils.NewWatchRepository(repo1Name, "default", utils.WatchRepositoryLocal)
	repo.Filters.PackageTypes = []string{"npm", "maven"}
	repo.Filters.Names = []string{"example-name"}
	repo.Filters.Paths = []string{"example-path"}
	repo.Filters.MimeTypes = []string{"example-mime-type"}
	repo.Filters.Properties = map[string]string{"some-key": "some-value", "some-key1": "some-value1"}

	repos[repo1Name] = repo

	anotherRepo := utils.NewWatchRepository(repo2Name, "default", utils.WatchRepositoryRemote)
	anotherRepo.Filters.PackageTypes = []string{"nuget"}
	anotherRepo.Filters.Names = []string{"another-example-name"}
	anotherRepo.Filters.Paths = []string{"another-example-path"}
	anotherRepo.Filters.MimeTypes = []string{"another-example-mime-type"}
	anotherRepo.Filters.Properties = map[string]string{"another-key": "some-value", "another-key1": "another-value1"}

	repos[repo2Name] = anotherRepo

	paramsSelectedRepos.Repositories.Type = utils.WatchRepositoriesByName
	paramsSelectedRepos.Repositories.Repositories = repos
	paramsSelectedRepos.Repositories.ExcludePatterns = []string{"selectedExcludePath1", "selectedExcludePath2"}
	paramsSelectedRepos.Repositories.IncludePatterns = []string{"selectedIncludePath1", "selectedIncludePath2"}

	paramsSelectedRepos.Builds.Type = utils.WatchBuildByName
	paramsSelectedRepos.Builds.ByNames = map[string]utils.WatchBuildsByNameParams{}
	paramsSelectedRepos.Builds.ByNames[build1Name] = utils.WatchBuildsByNameParams{
		Name:     build1Name,
		BinMgrID: "default",
	}
	err = testsXrayWatchService.Create(paramsSelectedRepos)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, testsXrayWatchService.Delete(paramsSelectedRepos.Name))
	}()
	validateWatchGeneralSettings(t, paramsSelectedRepos)

	targetConfig, err := testsXrayWatchService.Get(paramsSelectedRepos.Name)
	assert.NoError(t, err)
	assert.Equal(t, []string{"selectedExcludePath1", "selectedExcludePath2"}, targetConfig.Repositories.ExcludePatterns)
	assert.Equal(t, []string{"selectedIncludePath1", "selectedIncludePath2"}, targetConfig.Repositories.IncludePatterns)
	assert.Equal(t, utils.WatchRepositoriesByName, targetConfig.Repositories.Type)

	assert.Equal(t, repo1Name, targetConfig.Repositories.Repositories[repo1Name].Name)
	assert.Equal(t, "default", targetConfig.Repositories.Repositories[repo1Name].BinMgrID)
	assert.Equal(t, utils.WatchRepositoryLocal, targetConfig.Repositories.Repositories[repo1Name].RepoType)
	assert.Equal(t, []string{"Maven", "Npm"}, targetConfig.Repositories.Repositories[repo1Name].Filters.PackageTypes)
	assert.Equal(t, []string{"example-name"}, targetConfig.Repositories.Repositories[repo1Name].Filters.Names)
	assert.Equal(t, []string{"example-path"}, targetConfig.Repositories.Repositories[repo1Name].Filters.Paths)
	assert.Equal(t, []string{"example-mime-type"}, targetConfig.Repositories.Repositories[repo1Name].Filters.MimeTypes)
	assert.Equal(t, map[string]string{"some-key": "some-value", "some-key1": "some-value1"}, targetConfig.Repositories.Repositories[repo1Name].Filters.Properties)

	assert.Equal(t, repo2Name, targetConfig.Repositories.Repositories[repo2Name].Name)
	assert.Equal(t, "default", targetConfig.Repositories.Repositories[repo2Name].BinMgrID)
	assert.Equal(t, utils.WatchRepositoryRemote, targetConfig.Repositories.Repositories[repo2Name].RepoType)
	assert.Equal(t, []string{"NuGet"}, targetConfig.Repositories.Repositories[repo2Name].Filters.PackageTypes)
	assert.Equal(t, []string{"another-example-name"}, targetConfig.Repositories.Repositories[repo2Name].Filters.Names)
	assert.Equal(t, []string{"another-example-path"}, targetConfig.Repositories.Repositories[repo2Name].Filters.Paths)
	assert.Equal(t, []string{"another-example-mime-type"}, targetConfig.Repositories.Repositories[repo2Name].Filters.MimeTypes)
	assert.Equal(t, map[string]string{"another-key": "some-value", "another-key1": "another-value1"}, targetConfig.Repositories.Repositories[repo2Name].Filters.Properties)

	assert.Equal(t, utils.WatchBuildByName, targetConfig.Builds.Type)
	assert.Empty(t, targetConfig.Builds.All.ExcludePatterns)
	assert.Empty(t, targetConfig.Builds.All.IncludePatterns)

	assert.Equal(t, build1Name, targetConfig.Builds.ByNames[build1Name].Name)
	assert.Equal(t, "default", targetConfig.Builds.ByNames[build1Name].BinMgrID)

	targetConfig.Repositories.ExcludePatterns = []string{"excludePath-2"}
	targetConfig.Repositories.IncludePatterns = []string{"includePath-2", "fake-2"}
	targetConfig.Builds.ByNames[build2Name] = utils.WatchBuildsByNameParams{
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

func testXrayWatchBuildsByPattern(t *testing.T) {
	policy1Name := fmt.Sprintf("%s-%s", "policy1-pattern", getRunId())
	err := createDummyPolicy(policy1Name)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, testsXrayPolicyService.Delete(policy1Name))
	}()
	paramsBuildsByPattern := utils.NewWatchParams()
	paramsBuildsByPattern.Name = fmt.Sprintf("%s-%s", "client-go-tests-watch-builds-by-pattern", getRunId())
	paramsBuildsByPattern.Description = "Builds By Pattern"
	paramsBuildsByPattern.Builds.Type = utils.WatchBuildAll
	paramsBuildsByPattern.Builds.All.ExcludePatterns = []string{"excludePath"}
	paramsBuildsByPattern.Builds.All.IncludePatterns = []string{"includePath", "fake"}
	paramsBuildsByPattern.Builds.All.BinMgrID = "default"
	paramsBuildsByPattern.Policies = []utils.AssignedPolicy{
		{
			Name: policy1Name,
			Type: "security",
		},
	}

	err = testsXrayWatchService.Create(paramsBuildsByPattern)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, testsXrayWatchService.Delete(paramsBuildsByPattern.Name))
	}()
	validateWatchGeneralSettings(t, paramsBuildsByPattern)

	targetConfig, err := testsXrayWatchService.Get(paramsBuildsByPattern.Name)
	if assert.NoError(t, err) {
		assert.Equal(t, utils.WatchBuildAll, targetConfig.Builds.Type)
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
}

func testXrayWatchUpdateMissingWatch(t *testing.T) {
	paramsMissingWatch := utils.NewWatchParams()
	paramsMissingWatch.Name = fmt.Sprintf("%s-%s", "client-go-tests-watch-missing", getRunId())
	paramsMissingWatch.Description = "Missing Watch"
	paramsMissingWatch.Builds.Type = utils.WatchBuildAll
	paramsMissingWatch.Policies = []utils.AssignedPolicy{}

	err := testsXrayWatchService.Update(paramsMissingWatch)
	assert.EqualError(t, err, "server response: 404 Not Found\n{\n  \"error\": \"Failed to update Watch: Watch was not found\"\n}")
}

func testXrayWatchDeleteMissingWatch(t *testing.T) {
	err := testsXrayWatchService.Delete("client-go-tests-watch-builds-missing")
	assert.EqualError(t, err, "server response: 404 Not Found\n{\n  \"error\": \"Failed to delete Watch: Watch was not found\"\n}")
}

func testXrayWatchGetMissingWatch(t *testing.T) {
	_, err := testsXrayWatchService.Get("client-go-tests-watch-builds-missing")
	assert.EqualError(t, err, "server response: 404 Not Found\n{\n  \"error\": \"Watch was not found\"\n}")
}

func validateWatchGeneralSettings(t *testing.T, params utils.WatchParams) {
	targetConfig, err := testsXrayWatchService.Get(params.Name)
	if assert.NoError(t, err) {
		assert.Equal(t, params.Name, targetConfig.Name)
		assert.Equal(t, params.Description, targetConfig.Description)
		assert.Equal(t, params.Active, targetConfig.Active)
		assert.ElementsMatch(t, params.Policies, targetConfig.Policies)
	}
}

func createRepoLocal(t *testing.T, repoKey string) {
	glp := artifactoryServices.NewGenericLocalRepositoryParams()
	glp.Key = repoKey
	glp.XrayIndex = &trueValue

	err := testsCreateLocalRepositoryService.Generic(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
}

func createRepoRemote(t *testing.T, repoKey string) {
	nrp := artifactoryServices.NewNpmRemoteRepositoryParams()
	nrp.Key = repoKey
	nrp.RepoLayoutRef = "npm-default"
	nrp.Url = "https://registry.npmjs.org"
	nrp.XrayIndex = &trueValue

	err := testsCreateRemoteRepositoryService.Npm(nrp)
	assert.NoError(t, err, "Failed to create "+repoKey)
}

func createDummyPolicy(policyName string) error {
	params := utils.PolicyParams{
		Name:        policyName,
		Description: "example policy",
		Type:        utils.Security,
		Rules: []utils.PolicyRule{{
			Name:     "sec_rule",
			Criteria: *utils.CreateSeverityPolicyCriteria(utils.Medium),
			Actions: &utils.PolicyAction{
				Webhooks: []string{},
				BlockDownload: utils.PolicyBlockDownload{
					Active:    &trueValue,
					Unscanned: &falseValue,
				},
				BlockReleaseBundleDistribution: &trueValue,
				FailBuild:                      &trueValue,
				NotifyDeployer:                 &trueValue,
				NotifyWatchRecipients:          &trueValue,
			},
			Priority: 1,
		}},
	}
	err := testsXrayPolicyService.Create(params)
	return err
}

func createAndIndexBuild(t *testing.T, buildName string) error {
	err := createDummyBuild(buildName)
	assert.NoError(t, err)
	err = testXrayBinMgrService.AddBuildsToIndexing([]string{buildName})
	return err
}
