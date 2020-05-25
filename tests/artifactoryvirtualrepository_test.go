package tests

import (
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

var trimmedRtTargetRepo = strings.TrimSuffix(RtTargetRepo, "/")
var repos = []string{trimmedRtTargetRepo}

func TestArtifactoryVirtualRepository(t *testing.T) {
	t.Run("virtualMavenTest", virtualMavenTest)
	t.Run("virtualGradleTest", virtualGradleTest)
	t.Run("virtualIvyTest", virtualIvyTest)
	t.Run("virtualSbtTest", virtualSbtTest)
	t.Run("virtualHelmTest", virtualHelmTest)
	t.Run("virtualRpmTest", virtualRpmTest)
	t.Run("virtualNugetTest", virtualNugetTest)
	t.Run("virtualCranTest", virtualCranTest)
	t.Run("virtualGemsTest", virtualGemsTest)
	t.Run("virtualNpmTest", virtualNpmTest)
	t.Run("virtualBowerTest", virtualBowerTest)
	t.Run("virtualDebianTest", virtualDebianTest)
	t.Run("virtualPypiTest", virtualPypiTest)
	t.Run("virtualDockerTest", virtualDockerTest)
	t.Run("virtualGitlfsTest", virtualGitlfsTest)
	t.Run("virtualGoTest", virtualGoTest)
	t.Run("virtualYumTest", virtualYumTest)
	t.Run("virtualConanTest", virtualConanTest)
	t.Run("virtualChefTest", virtualChefTest)
	t.Run("virtualP2Test", virtualP2Test)
	t.Run("virtualPuppetTest", virtualPuppetTest)
	t.Run("virtualCondaTest", virtualCondaTest)
	t.Run("virtualGenericTest", virtualGenericTest)
	t.Run("getVirtualRepoDetailsTest", getVirtualRepoDetailsTest)
}

func virtualMavenTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	mvp := services.NewMavenVirtualRepositoryParams()
	mvp.Key = repoKey
	mvp.RepoLayoutRef = "maven-1-default"
	mvp.Repositories = repos
	mvp.Description = "Maven Repo for jfrog-client-go virtual-repository-test"
	mvp.PomRepositoryReferencesCleanupPolicy = "nothing"
	mvp.ForceMavenAuthentication = &trueValue
	mvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue

	err := testsCreateVirtualRepositoryService.Maven(mvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, mvp)

	mvp.Description += " - Updated"
	mvp.Notes = "Repo been updated"
	mvp.DefaultDeploymentRepo = trimmedRtTargetRepo
	mvp.RepoLayoutRef = "maven-2-default"
	mvp.ForceMavenAuthentication = nil
	mvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	mvp.ExcludesPattern = "**/****"

	err = testsUpdateVirtualRepositoryService.Maven(mvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, mvp)
}

func virtualGradleTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGradleVirtualRepositoryParams()
	gvp.Key = repoKey
	gvp.RepoLayoutRef = "simple-default"
	gvp.Description = "Gradle Repo for jfrog-client-go virtual-repository-test"
	gvp.PomRepositoryReferencesCleanupPolicy = "nothing"
	gvp.ForceMavenAuthentication = &trueValue
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue

	err := testsCreateVirtualRepositoryService.Gradle(gvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gvp)

	gvp.Description += " - Updated"
	gvp.Notes = "Repo been updated"
	gvp.RepoLayoutRef = "gradle-default"
	gvp.ForceMavenAuthentication = nil
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	gvp.ExcludesPattern = "**/****"
	gvp.Repositories = repos
	gvp.DefaultDeploymentRepo = trimmedRtTargetRepo

	err = testsUpdateVirtualRepositoryService.Gradle(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)
}

func virtualIvyTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	ivp := services.NewIvyVirtualRepositoryParams()
	ivp.Key = repoKey
	ivp.RepoLayoutRef = "ivy-default"
	ivp.Description = "Ivy Repo for jfrog-client-go virtual-repository-test"
	ivp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	ivp.DefaultDeploymentRepo = trimmedRtTargetRepo
	ivp.Repositories = repos
	ivp.IncludesPattern = "onlyDir/*"

	err := testsCreateVirtualRepositoryService.Ivy(ivp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, ivp)

	ivp.Description += " - Updated"
	ivp.Notes = "Repo been updated"
	ivp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	ivp.ExcludesPattern = "a/b/c/*"
	ivp.IncludesPattern = "**/*"

	err = testsUpdateVirtualRepositoryService.Ivy(ivp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, ivp)
}

func virtualSbtTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	svp := services.NewSbtVirtualRepositoryParams()
	svp.Key = repoKey
	svp.RepoLayoutRef = "sbt-default"
	svp.Description = "Sbt Repo for jfrog-client-go virtual-repository-test"
	svp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	svp.Repositories = repos
	svp.DefaultDeploymentRepo = trimmedRtTargetRepo

	err := testsCreateVirtualRepositoryService.Sbt(svp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, svp)

	svp.Description += " - Updated"
	svp.Notes = "Repo been updated"
	svp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	svp.IncludesPattern = "incDir/*"
	svp.ExcludesPattern = "exDir/*"

	err = testsUpdateVirtualRepositoryService.Sbt(svp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, svp)
}

func virtualHelmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	hvp := services.NewHelmVirtualRepositoryParams()
	hvp.Key = repoKey
	hvp.RepoLayoutRef = "simple-default"
	hvp.Description = "Helm Repo for jfrog-client-go virtual-repository-test"
	hvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	hvp.ExcludesPattern = "dir1/dir11/*"
	hvp.Repositories = repos
	hvp.DefaultDeploymentRepo = trimmedRtTargetRepo

	err := testsCreateVirtualRepositoryService.Helm(hvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, hvp)

	hvp.Description += " - Updated"
	hvp.Notes = "Repo been updated"
	hvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	hvp.ExcludesPattern = "dir2/*"
	hvp.IncludesPattern = "includeDir/*"
	hvp.VirtualRetrievalCachePeriodSecs = 666

	err = testsUpdateVirtualRepositoryService.Helm(hvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, hvp)
}

func virtualRpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	rvp := services.NewRpmVirtualRepositoryParams()
	rvp.Key = repoKey
	rvp.RepoLayoutRef = "simple-default"
	rvp.Description = "Rpm Repo for jfrog-client-go virtual-repository-test"
	rvp.ExcludesPattern = "dir1/dir11/*"
	rvp.VirtualRetrievalCachePeriodSecs = 5555
	rvp.Repositories = repos

	err := testsCreateVirtualRepositoryService.Rpm(rvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, rvp)

	rvp.Description += " - Updated"
	rvp.Notes = "Repo been updated"
	rvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	rvp.ExcludesPattern = "dir2/*"
	rvp.IncludesPattern = "includeDir/*"
	rvp.VirtualRetrievalCachePeriodSecs = 1818
	rvp.DefaultDeploymentRepo = trimmedRtTargetRepo

	err = testsUpdateVirtualRepositoryService.Rpm(rvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, rvp)
}

func virtualNugetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nvp := services.NewNugetVirtualRepositoryParams()
	nvp.Key = repoKey
	nvp.RepoLayoutRef = "nuget-default"
	nvp.Description = "Nuget Repo for jfrog-client-go virtual-repository-test"
	nvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	nvp.IncludesPattern = "**/*"
	nvp.ExcludesPattern = "*/ex/*"

	err := testsCreateVirtualRepositoryService.Nuget(nvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, nvp)

	nvp.Description += " - Updated"
	nvp.Notes = "Repo been updated"
	nvp.Repositories = repos
	nvp.DefaultDeploymentRepo = trimmedRtTargetRepo
	nvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	nvp.ExcludesPattern = "nugetExclude/*"
	nvp.ForceNugetAuthentication = &trueValue

	err = testsUpdateVirtualRepositoryService.Nuget(nvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nvp)
}

func virtualCranTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cvp := services.NewCranVirtualRepositoryParams()
	cvp.Key = repoKey
	cvp.RepoLayoutRef = "simple-default"
	cvp.Description = "Cran Repo for jfrog-client-go virtual-repository-test"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	cvp.ExcludesPattern = "dir1/dir11/*"
	cvp.VirtualRetrievalCachePeriodSecs = 5555

	err := testsCreateVirtualRepositoryService.Cran(cvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cvp)

	cvp.Description += " - Updated"
	cvp.Notes = "Repo been updated"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	cvp.ExcludesPattern = "dir2/*"
	cvp.IncludesPattern = "includeDir/*"
	cvp.VirtualRetrievalCachePeriodSecs = 1818

	err = testsUpdateVirtualRepositoryService.Cran(cvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cvp)
}

func virtualGemsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGemsVirtualRepositoryParams()
	gvp.Key = repoKey
	gvp.RepoLayoutRef = "simple-default"
	gvp.Description = "Gems Repo for jfrog-client-go virtual-repository-test"
	gvp.Repositories = repos
	gvp.DefaultDeploymentRepo = trimmedRtTargetRepo
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	gvp.ExcludesPattern = "dir1/"

	err := testsCreateVirtualRepositoryService.Gems(gvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gvp)

	gvp.Description += " - Updated"
	gvp.Notes = "Repo been updated"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	gvp.ExcludesPattern = "**/****,a/b/c/*"
	gvp.DefaultDeploymentRepo = ""

	err = testsUpdateVirtualRepositoryService.Gems(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)
}

func virtualNpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nvp := services.NewNpmVirtualRepositoryParams()
	nvp.Key = repoKey
	nvp.RepoLayoutRef = "npm-default"
	nvp.Description = "Npm Repo for jfrog-client-go virtual-repository-test"
	nvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	nvp.DefaultDeploymentRepo = trimmedRtTargetRepo
	nvp.Repositories = repos
	nvp.IncludesPattern = "includeNpm/*"
	nvp.ExternalDependenciesEnabled = &trueValue
	nvp.ExternalDependenciesPatterns = []string{"**/*microsoft*/**"}
	nvp.VirtualRetrievalCachePeriodSecs = 1818

	err := testsCreateVirtualRepositoryService.Npm(nvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, nvp)

	nvp.Description += " - Updated"
	nvp.Notes = "Repo been updated"
	nvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	nvp.ExcludesPattern = "npmEx/*"
	nvp.IncludesPattern = "**/*"
	nvp.DefaultDeploymentRepo = ""
	nvp.ExternalDependenciesPatterns = append(nvp.ExternalDependenciesPatterns, "**/*github*/**")
	nvp.VirtualRetrievalCachePeriodSecs = 1500

	err = testsUpdateVirtualRepositoryService.Npm(nvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nvp)
}

func virtualBowerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	bvp := services.NewBowerVirtualRepositoryParams()
	bvp.Key = repoKey
	bvp.RepoLayoutRef = "bower-default"
	bvp.Description = "Bower Repo for jfrog-client-go virtual-repository-test"
	bvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	bvp.Repositories = repos
	bvp.IncludesPattern = "bowerInc/*"
	bvp.ExternalDependenciesEnabled = &trueValue
	bvp.ExternalDependenciesPatterns = []string{"**/*github*/**"}

	err := testsCreateVirtualRepositoryService.Bower(bvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, bvp)

	bvp.Description += " - Updated"
	bvp.Notes = "Repo been updated"
	bvp.DefaultDeploymentRepo = trimmedRtTargetRepo
	bvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	bvp.ExcludesPattern = "dir1/*"
	bvp.IncludesPattern = "**/*"
	bvp.DefaultDeploymentRepo = ""
	bvp.ExternalDependenciesPatterns = append(bvp.ExternalDependenciesPatterns, "**/*microsoft*/**")

	err = testsUpdateVirtualRepositoryService.Bower(bvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, bvp)
}

func virtualDebianTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	dvp := services.NewDebianVirtualRepositoryParams()
	dvp.Key = repoKey
	dvp.RepoLayoutRef = "simple-default"
	dvp.Description = "Debian Repo for jfrog-client-go virtual-repository-test"
	dvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	dvp.ExcludesPattern = "dir1/dir2/*"

	err := testsCreateVirtualRepositoryService.Debian(dvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, dvp)

	dvp.Description += " - Updated"
	dvp.Notes = "Repo been updated"
	dvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	dvp.ExcludesPattern = "dirEx/*"
	dvp.Repositories = repos
	dvp.DefaultDeploymentRepo = trimmedRtTargetRepo

	err = testsUpdateVirtualRepositoryService.Debian(dvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dvp)
}

func virtualPypiTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	pvp := services.NewPypiVirtualRepositoryParams()
	pvp.Key = repoKey
	pvp.RepoLayoutRef = "simple-default"
	pvp.Description = "Pypi Repo for jfrog-client-go virtual-repository-test"
	pvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	pvp.ExcludesPattern = "dir1/dir2/*"

	err := testsCreateVirtualRepositoryService.Pypi(pvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, pvp)

	pvp.Description += " - Updated"
	pvp.Notes = "Repo been updated"
	pvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	pvp.ExcludesPattern = "dirEx/*"
	pvp.Repositories = repos
	pvp.DefaultDeploymentRepo = trimmedRtTargetRepo

	err = testsUpdateVirtualRepositoryService.Pypi(pvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, pvp)
}

func virtualDockerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	dvp := services.NewDockerVirtualRepositoryParams()
	dvp.Key = repoKey
	dvp.RepoLayoutRef = "simple-default"
	dvp.Description = "Docker Repo for jfrog-client-go virtual-repository-test"
	dvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	dvp.IncludesPattern = "**/*"
	dvp.ExcludesPattern = "*/ex/*"

	err := testsCreateVirtualRepositoryService.Docker(dvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, dvp)

	dvp.Description += " - Updated"
	dvp.Notes = "Repo been updated"
	dvp.Repositories = repos
	dvp.DefaultDeploymentRepo = trimmedRtTargetRepo
	dvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	dvp.ExcludesPattern = "docker1/*"

	err = testsUpdateVirtualRepositoryService.Docker(dvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dvp)
}

func virtualGitlfsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGitlfsVirtualRepositoryParams()
	gvp.Key = repoKey
	gvp.RepoLayoutRef = "simple-default"
	gvp.Description = "Gitlfs Repo for jfrog-client-go virtual-repository-test"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	gvp.ExcludesPattern = "dir1/dir1.1/*"

	err := testsCreateVirtualRepositoryService.Gitlfs(gvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gvp)

	gvp.Description += " - Updated"
	gvp.Notes = "Repo been updated"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	gvp.ExcludesPattern = "dir2/*"
	gvp.Repositories = repos
	gvp.DefaultDeploymentRepo = trimmedRtTargetRepo

	err = testsUpdateVirtualRepositoryService.Gitlfs(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)
}

func virtualGoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGoVirtualRepositoryParams()
	gvp.Key = repoKey
	gvp.RepoLayoutRef = "go-default"
	gvp.Description = "Go Repo for jfrog-client-go virtual-repository-test"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	gvp.DefaultDeploymentRepo = trimmedRtTargetRepo
	gvp.Repositories = repos
	gvp.IncludesPattern = "includeGo/*"
	gvp.ExternalDependenciesEnabled = &trueValue
	gvp.ExternalDependenciesPatterns = []string{"**/*microsoft*/**", "**/*github*/**"}

	err := testsCreateVirtualRepositoryService.Go(gvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gvp)

	gvp.Description += " - Updated"
	gvp.Notes = "Repo been updated"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	gvp.ExcludesPattern = "go1/go2/go3/*"
	gvp.IncludesPattern = "**/*"
	gvp.DefaultDeploymentRepo = ""
	gvp.ExternalDependenciesPatterns = append(gvp.ExternalDependenciesPatterns, "**/gopkg.in/**", "**/go.googlesource.com/**")

	err = testsUpdateVirtualRepositoryService.Go(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)
}

func virtualYumTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	yvp := services.NewYumVirtualRepositoryParams()
	yvp.Key = repoKey
	yvp.RepoLayoutRef = "simple-default"
	yvp.Description = "Yum Repo for jfrog-client-go virtual-repository-test"
	yvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	yvp.DefaultDeploymentRepo = trimmedRtTargetRepo
	yvp.Repositories = repos
	yvp.IncludesPattern = "onlyDir/*"

	err := testsCreateVirtualRepositoryService.Yum(yvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// "yum" package type is converted to "rpm" by Artifactory, so we have to change it too to pass the validation.
	yvp.PackageType = "rpm"
	validateRepoConfig(t, repoKey, yvp)

	yvp.Description += " - Updated"
	yvp.Notes = "Repo been updated"
	yvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	yvp.ExcludesPattern = "a/b/c/*"
	yvp.IncludesPattern = "**/*"
	yvp.DefaultDeploymentRepo = ""

	err = testsUpdateVirtualRepositoryService.Yum(yvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, yvp)
}

func virtualConanTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cvp := services.NewConanVirtualRepositoryParams()
	cvp.Key = repoKey
	cvp.RepoLayoutRef = "conan-default"
	cvp.Description = "Conan Repo for jfrog-client-go virtual-repository-test"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	cvp.ExcludesPattern = "dir1/dir11/*"
	cvp.VirtualRetrievalCachePeriodSecs = 1818

	err := testsCreateVirtualRepositoryService.Conan(cvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cvp)

	cvp.Description += " - Updated"
	cvp.Notes = "Repo been updated"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	cvp.ExcludesPattern = "dir2/*"
	cvp.IncludesPattern = "includeDir/*"
	cvp.VirtualRetrievalCachePeriodSecs = 5555

	err = testsUpdateVirtualRepositoryService.Conan(cvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cvp)
}

func virtualChefTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cvp := services.NewChefVirtualRepositoryParams()
	cvp.Key = repoKey
	cvp.RepoLayoutRef = "simple-default"
	cvp.Description = "Chef Repo for jfrog-client-go virtual-repository-test"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	cvp.IncludesPattern = "**/*"
	cvp.ExcludesPattern = "chef/ex/*"
	cvp.Repositories = repos
	cvp.DefaultDeploymentRepo = trimmedRtTargetRepo

	err := testsCreateVirtualRepositoryService.Chef(cvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cvp)

	cvp.Description += " - Updated"
	cvp.Notes = "Repo been updated"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	cvp.ExcludesPattern = "dir2/*"

	err = testsUpdateVirtualRepositoryService.Chef(cvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cvp)
}

func virtualPuppetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	pvp := services.NewPuppetVirtualRepositoryParams()
	pvp.Key = repoKey
	pvp.RepoLayoutRef = "puppet-default"
	pvp.Description = "Puppet Repo for jfrog-client-go virtual-repository-test"
	pvp.Repositories = repos
	pvp.DefaultDeploymentRepo = trimmedRtTargetRepo
	pvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	pvp.ExcludesPattern = "dir1/*"

	err := testsCreateVirtualRepositoryService.Puppet(pvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, pvp)

	pvp.Description += " - Updated"
	pvp.Notes = "Repo been updated"
	pvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	pvp.ExcludesPattern = "dir2/*"
	pvp.IncludesPattern = "dir1/*"
	pvp.DefaultDeploymentRepo = ""

	err = testsUpdateVirtualRepositoryService.Puppet(pvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, pvp)
}

func virtualP2Test(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	pvp := services.NewP2VirtualRepositoryParams()
	pvp.Key = repoKey
	pvp.RepoLayoutRef = "simple-default"
	pvp.Description = "P2 Repo for jfrog-client-go virtual-repository-test"
	pvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	pvp.ExcludesPattern = "dir1/dir1.1/*"

	err := testsCreateVirtualRepositoryService.P2(pvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, pvp)

	pvp.Description += " - Updated"
	pvp.Notes = "Repo been updated"
	pvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	pvp.ExcludesPattern = "dir2/*"

	err = testsUpdateVirtualRepositoryService.P2(pvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, pvp)
}

func virtualCondaTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cvp := services.NewCondaVirtualRepositoryParams()
	cvp.Key = repoKey
	cvp.RepoLayoutRef = "simple-default"
	cvp.Description = "Conda Repo for jfrog-client-go virtual-repository-test"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	cvp.IncludesPattern = "**/*"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue

	err := testsCreateVirtualRepositoryService.Conda(cvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cvp)

	cvp.Description += " - Updated"
	cvp.Notes = "Repo been updated"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	cvp.ExcludesPattern = "dir2/*"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue

	err = testsUpdateVirtualRepositoryService.Conda(cvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cvp)
}

func virtualGenericTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGenericVirtualRepositoryParams()
	gvp.Key = repoKey
	gvp.RepoLayoutRef = "simple-default"
	gvp.Description = "Generic Repo for jfrog-client-go virtual-repository-test"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue

	err := testsCreateVirtualRepositoryService.Generic(gvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gvp)

	gvp.Description += " - Updated"
	gvp.Notes = "Repo been updated"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	gvp.ExcludesPattern = "**/****,a/b/c/*"

	err = testsUpdateVirtualRepositoryService.Generic(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)
}

func getVirtualRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGoVirtualRepositoryParams()
	gvp.Key = repoKey
	gvp.Description = "Repo for jfrog-client-go virtual-repository-test"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue

	err := testsCreateVirtualRepositoryService.Go(gvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// Get repo details
	data := getRepo(t, repoKey)
	// Validate
	assert.Equal(t, data.Key, repoKey)
	assert.Equal(t, data.Description, gvp.Description)
	assert.Equal(t, data.Rclass, "virtual")
	assert.Empty(t, data.Url)
	assert.Equal(t, data.PackageType, "go")
}
