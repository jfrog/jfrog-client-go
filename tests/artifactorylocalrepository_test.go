package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"testing"
)

const LocalRepo = "-local-"

func TestArtifactoryLocalRepository(t *testing.T) {
	t.Run("localMavenTest", localMavenTest)
	t.Run("localGenericTest", localGenericTest)
}

func localMavenTest(t *testing.T) {
	repoKey := "maven" + LocalRepo + timestamp
	mlp := services.NewMavenLocalRepositoryParams()
	mlp.Key = repoKey
	mlp.RepoLayoutRef = "maven-2-default"
	mlp.Description = "Maven Repo for jfrog-client-go local-repository-test"
	mlp.SuppressPomConsistencyChecks = &trueValue
	mlp.HandleReleases = &trueValue
	mlp.HandleSnapshots = &falseValue
	mlp.XrayIndex = &trueValue
	mlp.MaxUniqueSnapshots = 18
	mlp.ChecksumPolicyType = "server-generated-checksums"
	mlp.DownloadRedirect = &falseValue

	err := testsCreateLocalRepositoryService.Maven(mlp)
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, mlp) {
		t.Error("Validation after create failed for " + repoKey)
	}

	mlp.Description += " - Updated"
	mlp.MaxUniqueSnapshots = 36
	mlp.HandleReleases = nil
	mlp.HandleSnapshots = &trueValue
	mlp.ChecksumPolicyType = "client-checksums"
	mlp.Notes = "Repo been updated"
	mlp.ArchiveBrowsingEnabled = &trueValue

	err = testsUpdateLocalRepositoryService.Maven(mlp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, mlp) {
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

func localGenericTest(t *testing.T) {
	repoKey := "generic" + LocalRepo + timestamp
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "simple-default"
	glp.Description = "Generic Repo for jfrog-client-go local-repository-test"
	glp.XrayIndex = &trueValue
	glp.DownloadRedirect = &falseValue
	glp.ArchiveBrowsingEnabled = &falseValue

	err := testsCreateLocalRepositoryService.Generic(glp)
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, glp) {
		t.Error("Validation after create failed for " + repoKey)
	}

	glp.Description += " - Updated"
	glp.Notes = "Repo been updated"
	glp.ArchiveBrowsingEnabled = &trueValue
	glp.ArchiveBrowsingEnabled = &falseValue
	glp.BlockPushingSchema1 = nil

	err = testsUpdateLocalRepositoryService.Generic(glp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, glp) {
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
