package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

const VirtualRepo = "-virtual-"

func TestArtifactoryVirtualRepository(t *testing.T) {
	t.Run("virtualMavenTest", virtualMavenTest)
	t.Run("virtualGradleTest", virtualGradleTest)
	t.Run("virtualP2Test", virtualP2Test)
	t.Run("virtualCondaTest", virtualCondaTest)
	t.Run("virtualGenericTest", virtualGenericTest)
}

func virtualMavenTest(t *testing.T) {
	repoKey := "maven" + VirtualRepo + timestamp
	mvp := services.NewMavenVirtualRepositoryParams()
	mvp.Key = repoKey
	mvp.RepoLayoutRef = "maven-2-default"
	mvp.Description = "Maven Repo for jfrog-client-go virtual-repository-test"
	mvp.PomRepositoryReferencesCleanupPolicy = "nothing"
	mvp.ForceMavenAuthentication = &trueValue
	mvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue

	err := testsCreateVirtualRepositoryService.Maven(mvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, mvp)

	mvp.Description += " - Updated"
	mvp.Notes = "Repo been updated"
	mvp.RepoLayoutRef = "maven-1-default"
	mvp.ForceMavenAuthentication = nil
	mvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	mvp.ExcludesPattern = "**/****"

	err = testsUpdateVirtualRepositoryService.Maven(mvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, mvp)

	deleteRepoAndValidate(t, repoKey)
}

func virtualGradleTest(t *testing.T) {
	repoKey := "gradle" + VirtualRepo + timestamp
	gvp := services.NewGradleVirtualRepositoryParams()
	gvp.Key = repoKey
	gvp.RepoLayoutRef = "gradle-default"
	gvp.Description = "Gradle Repo for jfrog-client-go virtual-repository-test"
	gvp.PomRepositoryReferencesCleanupPolicy = "nothing"
	gvp.ForceMavenAuthentication = &trueValue
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue

	err := testsCreateVirtualRepositoryService.Gradle(gvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, gvp)

	gvp.Description += " - Updated"
	gvp.Notes = "Repo been updated"
	gvp.RepoLayoutRef = "maven-1-default"
	gvp.ForceMavenAuthentication = nil
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	gvp.ExcludesPattern = "**/****"

	err = testsUpdateVirtualRepositoryService.Gradle(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)

	deleteRepoAndValidate(t, repoKey)
}

func virtualP2Test(t *testing.T) {
	repoKey := "p2" + VirtualRepo + timestamp
	pvp := services.NewP2VirtualRepositoryParams()
	pvp.Key = repoKey
	pvp.RepoLayoutRef = "simple-default"
	pvp.Description = "P2 Repo for jfrog-client-go virtual-repository-test"
	pvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	pvp.ExcludesPattern = "dir1/dir1.1/*"

	err := testsCreateVirtualRepositoryService.P2(pvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, pvp)

	pvp.Description += " - Updated"
	pvp.Notes = "Repo been updated"
	pvp.RepoLayoutRef = "maven-1-default"
	pvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	pvp.ExcludesPattern = "dir2/*"

	err = testsUpdateVirtualRepositoryService.P2(pvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, pvp)

	deleteRepoAndValidate(t, repoKey)
}

func virtualCondaTest(t *testing.T) {
	repoKey := "conda" + VirtualRepo + timestamp
	cvp := services.NewCondaVirtualRepositoryParams()
	cvp.Key = repoKey
	cvp.RepoLayoutRef = "simple-default"
	cvp.Description = "Conda Repo for jfrog-client-go virtual-repository-test"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	cvp.IncludesPattern = "**/*"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue

	err := testsCreateVirtualRepositoryService.Conda(cvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, cvp)

	cvp.Description += " - Updated"
	cvp.Notes = "Repo been updated"
	cvp.RepoLayoutRef = "maven-1-default"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	cvp.ExcludesPattern = "dir2/*"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue

	err = testsUpdateVirtualRepositoryService.Conda(cvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cvp)

	deleteRepoAndValidate(t, repoKey)
}

func virtualGenericTest(t *testing.T) {
	repoKey := "generic" + VirtualRepo + timestamp
	gvp := services.NewGenericVirtualRepositoryParams()
	gvp.Key = repoKey
	gvp.RepoLayoutRef = "simple-default"
	gvp.Description = "Generic Repo for jfrog-client-go virtual-repository-test"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue

	err := testsCreateVirtualRepositoryService.Generic(gvp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, gvp)

	gvp.Description += " - Updated"
	gvp.Notes = "Repo been updated"
	gvp.RepoLayoutRef = "maven-1-default"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	gvp.ExcludesPattern = "**/****,a/b/c/*"

	err = testsUpdateVirtualRepositoryService.Generic(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)

	deleteRepoAndValidate(t, repoKey)
}
