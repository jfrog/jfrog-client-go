package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"testing"
)

const VirtualRepo = "-virtual-"

func TestArtifactoryVirtualRepository(t *testing.T) {
	t.Run("virtualMavenTest", virtualMavenTest)
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
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, mvp) {
		t.Error("Validation after create failed for " + repoKey)
	}

	mvp.Description += " - Updated"
	mvp.Notes = "Repo been updated"
	mvp.RepoLayoutRef = "maven-1-default"
	mvp.ForceMavenAuthentication = nil
	mvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	mvp.ExcludesPattern = "**/****"

	err = testsUpdateVirtualRepositoryService.Maven(mvp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, mvp) {
		t.Error("Validation after update failed for " + repoKey)
	}

	err = testsDeleteRepositoryService.Delete(repoKey)
	if err != nil {
		t.Error("Failed to delete " + repoKey)
	}
	if isRepoExist(repoKey) {
		t.Error(repoKey + " still exists")
	}
}

func virtualGenericTest(t *testing.T) {
	repoKey := "generic" + VirtualRepo + timestamp
	gvp := services.NewGenericVirtualRepositoryParams()
	gvp.Key = repoKey
	gvp.RepoLayoutRef = "simple-default"
	gvp.Description = "Generic Repo for jfrog-client-go virtual-repository-test"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue

	err := testsCreateVirtualRepositoryService.Generic(gvp)
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, gvp) {
		t.Error("Validation after create failed for " + repoKey)
	}

	gvp.Description += " - Updated"
	gvp.Notes = "Repo been updated"
	gvp.RepoLayoutRef = "maven-1-default"
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	gvp.ExcludesPattern = "**/****,a/b/c/*"

	err = testsUpdateVirtualRepositoryService.Generic(gvp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, gvp) {
		t.Error("Validation after update failed for " + repoKey)
	}

	err = testsDeleteRepositoryService.Delete(repoKey)
	if err != nil {
		t.Error("Failed to delete " + repoKey)
	}
	if isRepoExist(repoKey) {
		t.Error(repoKey + " still exists")
	}
}
