package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
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
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, gvp) {
		t.Error("Validation after create failed for " + repoKey)
	}

	gvp.Description += " - Updated"
	gvp.Notes = "Repo been updated"
	gvp.RepoLayoutRef = "maven-1-default"
	gvp.ForceMavenAuthentication = nil
	gvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	gvp.ExcludesPattern = "**/****"

	err = testsUpdateVirtualRepositoryService.Gradle(gvp)
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

func virtualP2Test(t *testing.T) {
	repoKey := "p2" + VirtualRepo + timestamp
	pvp := services.NewP2VirtualRepositoryParams()
	pvp.Key = repoKey
	pvp.RepoLayoutRef = "simple-default"
	pvp.Description = "P2 Repo for jfrog-client-go virtual-repository-test"
	pvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
	pvp.ExcludesPattern = "dir1/dir1.1/*"

	err := testsCreateVirtualRepositoryService.P2(pvp)
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, pvp) {
		t.Error("Validation after create failed for " + repoKey)
	}

	pvp.Description += " - Updated"
	pvp.Notes = "Repo been updated"
	pvp.RepoLayoutRef = "maven-1-default"
	pvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	pvp.ExcludesPattern = "dir2/*"

	err = testsUpdateVirtualRepositoryService.P2(pvp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, pvp) {
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
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, cvp) {
		t.Error("Validation after create failed for " + repoKey)
	}

	cvp.Description += " - Updated"
	cvp.Notes = "Repo been updated"
	cvp.RepoLayoutRef = "maven-1-default"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
	cvp.ExcludesPattern = "dir2/*"
	cvp.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue

	err = testsUpdateVirtualRepositoryService.Conda(cvp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, cvp) {
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
