package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"testing"
)

const RemoteRepo = "-remote-"

func TestArtifactoryRemoteRepository(t *testing.T) {
	t.Run("remoteMavenTest", remoteMavenTest)
	t.Run("remoteGenericTest", remoteGenericTest)
}

func remoteMavenTest(t *testing.T) {
	repoKey := "maven" + RemoteRepo + timestamp
	mrp := services.NewMavenRemoteRepositoryParams()
	mrp.Key = repoKey
	mrp.RepoLayoutRef = "maven-2-default"
	mrp.Url = "https://jcenter.bintray.com"
	mrp.Description = "Maven Repo for jfrog-client-go remote-repository-test"
	mrp.SuppressPomConsistencyChecks = &trueValue
	mrp.HandleReleases = &trueValue
	mrp.HandleSnapshots = &trueValue
	mrp.RemoteRepoChecksumPolicyType = "ignore-and-generate"
	mrp.AssumedOfflinePeriodSecs = 2345
	mrp.StoreArtifactsLocally = &falseValue
	mrp.ShareConfiguration = &falseValue

	err := testsCreateRemoteRepositoryService.Maven(mrp)
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, mrp) {
		t.Error("Validation after create failed for " + repoKey)
	}

	mrp.Description += " - Updated"
	mrp.HandleReleases = nil
	mrp.HandleSnapshots = &falseValue
	mrp.Notes = "Repo been updated"
	mrp.AssumedOfflinePeriodSecs = 2000
	mrp.EnableCookieManagement = &trueValue
	mrp.FetchJarsEagerly = &falseValue
	mrp.SocketTimeoutMillis = 666

	err = testsUpdateRemoteRepositoryService.Maven(mrp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, mrp) {
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

func remoteGenericTest(t *testing.T) {
	repoKey := "generic" + RemoteRepo + timestamp
	grp := services.NewGenericRemoteRepositoryParams()
	grp.Key = repoKey
	grp.RepoLayoutRef = "simple-default"
	grp.Url = "https://jcenter.bintray.com"
	grp.Description = "Generic Repo for jfrog-client-go remote-repository-test"
	grp.AssumedOfflinePeriodSecs = 2345
	grp.StoreArtifactsLocally = &falseValue
	grp.ShareConfiguration = &falseValue

	err := testsCreateRemoteRepositoryService.Generic(grp)
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, grp) {
		t.Error("Validation after create failed for " + repoKey)
	}

	grp.Description += " - Updated"
	grp.Notes = "Repo been updated"
	grp.AssumedOfflinePeriodSecs = 2000
	grp.EnableCookieManagement = &trueValue
	grp.SocketTimeoutMillis = 666

	err = testsUpdateRemoteRepositoryService.Generic(grp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	if !validateRepoConfig(t, repoKey, grp) {
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
