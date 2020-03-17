package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"testing"
)

const RemoteRepo = "-remote-"

func TestArtifactoryRemoteRepository(t *testing.T) {
	t.Run("remoteMavenTest", remoteMavenTest)
	t.Run("remoteGradleTest", remoteGradleTest)
	t.Run("remoteComposerTest", remoteComposerTest)
	t.Run("remoteVcsTest", remoteVcsTest)
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
	validateRepoConfig(t, repoKey, mrp)

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
	validateRepoConfig(t, repoKey, mrp)

	err = testsDeleteRepositoryService.Delete(repoKey)
	if err != nil {
		t.Error("Failed to delete " + repoKey)
	}
	if isRepoExist(repoKey) {
		t.Error(repoKey + " still exists")
	}
}

func remoteGradleTest(t *testing.T) {
	repoKey := "gradle" + RemoteRepo + timestamp
	grp := services.NewGradleRemoteRepositoryParams()
	grp.Key = repoKey
	grp.RepoLayoutRef = "gradle-default"
	grp.Url = "https://jcenter.bintray.com"
	grp.Description = "Gradle Repo for jfrog-client-go remote-repository-test"
	grp.SuppressPomConsistencyChecks = &trueValue
	grp.HandleReleases = &trueValue
	grp.HandleSnapshots = &trueValue
	grp.RemoteRepoChecksumPolicyType = "ignore-and-generate"
	grp.AssumedOfflinePeriodSecs = 2345
	grp.StoreArtifactsLocally = &falseValue
	grp.ShareConfiguration = &falseValue

	err := testsCreateRemoteRepositoryService.Gradle(grp)
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	validateRepoConfig(t, repoKey, grp)

	grp.Description += " - Updated"
	grp.HandleReleases = nil
	grp.HandleSnapshots = &falseValue
	grp.Notes = "Repo been updated"
	grp.AssumedOfflinePeriodSecs = 2000
	grp.EnableCookieManagement = &trueValue
	grp.RemoteRepoChecksumPolicyType = "generate-if-absent"
	grp.FetchJarsEagerly = &falseValue
	grp.SocketTimeoutMillis = 666

	err = testsUpdateRemoteRepositoryService.Gradle(grp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	validateRepoConfig(t, repoKey, grp)

	err = testsDeleteRepositoryService.Delete(repoKey)
	if err != nil {
		t.Error("Failed to delete " + repoKey)
	}
	if isRepoExist(repoKey) {
		t.Error(repoKey + " still exists")
	}
}

func remoteComposerTest(t *testing.T) {
	repoKey := "composer" + RemoteRepo + timestamp
	crp := services.NewComposerRemoteRepositoryParams()
	crp.Key = repoKey
	crp.RepoLayoutRef = "composer-default"
	crp.Url = "https://github.com/"
	crp.Description = "Composer Repo for jfrog-client-go remote-repository-test"
	crp.AssumedOfflinePeriodSecs = 2345
	crp.StoreArtifactsLocally = &falseValue
	crp.ShareConfiguration = &falseValue
	crp.ComposerRegistryUrl = "https://composer.registry.com/"
	crp.IncludesPattern = "dir1/*, dir2/dir2.1/*"
	crp.BypassHeadRequests = &trueValue

	err := testsCreateRemoteRepositoryService.Composer(crp)
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	validateRepoConfig(t, repoKey, crp)

	crp.Description += " - Updated"
	crp.Notes = "Repo been updated"
	crp.AssumedOfflinePeriodSecs = 2000
	crp.EnableCookieManagement = &trueValue
	crp.SocketTimeoutMillis = 666
	crp.IncludesPattern = "**/*"
	crp.BypassHeadRequests = &falseValue

	err = testsUpdateRemoteRepositoryService.Composer(crp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	validateRepoConfig(t, repoKey, crp)

	err = testsDeleteRepositoryService.Delete(repoKey)
	if err != nil {
		t.Error("Failed to delete " + repoKey)
	}
	if isRepoExist(repoKey) {
		t.Error(repoKey + " still exists")
	}
}

func remoteVcsTest(t *testing.T) {
	repoKey := "vcs" + RemoteRepo + timestamp
	vrp := services.NewVcsRemoteRepositoryParams()
	vrp.Key = repoKey
	vrp.RepoLayoutRef = "composer-default"
	vrp.Url = "https://github.com/"
	vrp.Description = "Vcs Repo for jfrog-client-go remote-repository-test"
	vrp.AssumedOfflinePeriodSecs = 2345
	vrp.StoreArtifactsLocally = &falseValue
	vrp.ShareConfiguration = &falseValue
	vrp.VcsGitDownloadUrl = "https://github.com/download.git"
	vrp.VcsGitProvider = "BITBUCKET"
	vrp.VcsType = "GIT"
	vrp.IncludesPattern = "dir1/*, dir2/dir2.1/*"
	vrp.BypassHeadRequests = &trueValue
	vrp.SocketTimeoutMillis = 1111

	err := testsCreateRemoteRepositoryService.Vcs(vrp)
	if err != nil {
		t.Error("Failed to create " + repoKey)
	}
	validateRepoConfig(t, repoKey, vrp)

	vrp.Description += " - Updated"
	vrp.Notes = "Repo been updated"
	vrp.AssumedOfflinePeriodSecs = 2000
	vrp.EnableCookieManagement = &trueValue
	vrp.SocketTimeoutMillis = 666
	vrp.IncludesPattern = "**/*"
	vrp.BypassHeadRequests = &falseValue
	vrp.VcsGitProvider = "OLDSTASH"
	vrp.SocketTimeoutMillis = 1110

	err = testsUpdateRemoteRepositoryService.Vcs(vrp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	validateRepoConfig(t, repoKey, vrp)

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
	validateRepoConfig(t, repoKey, grp)

	grp.Description += " - Updated"
	grp.Notes = "Repo been updated"
	grp.AssumedOfflinePeriodSecs = 2000
	grp.EnableCookieManagement = &trueValue
	grp.SocketTimeoutMillis = 666

	err = testsUpdateRemoteRepositoryService.Generic(grp)
	if err != nil {
		t.Error("Failed to update " + repoKey)
	}
	validateRepoConfig(t, repoKey, grp)

	err = testsDeleteRepositoryService.Delete(repoKey)
	if err != nil {
		t.Error("Failed to delete " + repoKey)
	}
	if isRepoExist(repoKey) {
		t.Error(repoKey + " still exists")
	}
}
