package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

const LocalRepo = "-local-"

func TestArtifactoryLocalRepository(t *testing.T) {
	t.Run("localMavenTest", localMavenTest)
	t.Run("localGradleTest", localGradleTest)
	t.Run("localRpmTest", localRpmTest)
	t.Run("localGoTest", localGoTest)
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
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, mlp)

	mlp.Description += " - Updated"
	mlp.MaxUniqueSnapshots = 36
	mlp.HandleReleases = nil
	mlp.HandleSnapshots = &trueValue
	mlp.ChecksumPolicyType = "client-checksums"
	mlp.Notes = "Repo been updated"
	mlp.ArchiveBrowsingEnabled = &trueValue

	err = testsUpdateLocalRepositoryService.Maven(mlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, mlp)

	deleteRepoAndValidate(t, repoKey)
}

func localGradleTest(t *testing.T) {
	repoKey := "gradle" + LocalRepo + timestamp
	glp := services.NewGradleLocalRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "gradle-default"
	glp.Description = "Gradle Repo for jfrog-client-go local-repository-test"
	glp.SuppressPomConsistencyChecks = &trueValue
	glp.HandleReleases = &trueValue
	glp.HandleSnapshots = &falseValue
	glp.XrayIndex = &trueValue
	glp.MaxUniqueSnapshots = 18
	glp.ChecksumPolicyType = "server-generated-checksums"
	glp.DownloadRedirect = &falseValue

	err := testsCreateLocalRepositoryService.Gradle(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, glp)

	glp.Description += " - Updated"
	glp.MaxUniqueSnapshots = 36
	glp.HandleReleases = nil
	glp.HandleSnapshots = &trueValue
	glp.ChecksumPolicyType = "client-checksums"
	glp.Notes = "Repo been updated"
	glp.ArchiveBrowsingEnabled = &trueValue

	err = testsUpdateLocalRepositoryService.Gradle(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)

	deleteRepoAndValidate(t, repoKey)
}

func localRpmTest(t *testing.T) {
	repoKey := "rpm" + LocalRepo + timestamp
	rlp := services.NewRpmLocalRepositoryParams()
	rlp.Key = repoKey
	rlp.RepoLayoutRef = "simple-default"
	rlp.Description = "Rpm Repo for jfrog-client-go local-repository-test"
	rlp.XrayIndex = &trueValue
	rlp.DownloadRedirect = &falseValue
	rlp.YumRootDepth = 6
	rlp.CalculateYumMetadata = &falseValue

	err := testsCreateLocalRepositoryService.Rpm(rlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, rlp)

	rlp.Description += " - Updated"
	rlp.Notes = "Repo been updated"
	rlp.ArchiveBrowsingEnabled = &trueValue
	rlp.YumRootDepth = 18
	rlp.CalculateYumMetadata = &trueValue
	rlp.EnableFileListsIndexing = &falseValue

	err = testsUpdateLocalRepositoryService.Rpm(rlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, rlp)

	deleteRepoAndValidate(t, repoKey)
}

func localGoTest(t *testing.T) {
	repoKey := "go" + LocalRepo + timestamp
	glp := services.NewGoLocalRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "go-default"
	glp.Description = "Go Repo for jfrog-client-go local-repository-test"
	glp.XrayIndex = &trueValue
	glp.DownloadRedirect = &falseValue
	glp.PropertySets = []string{"artifactory"}
	glp.ArchiveBrowsingEnabled = &trueValue

	err := testsCreateLocalRepositoryService.Go(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, glp)

	glp.Description += " - Updated"
	glp.Notes = "Repo been updated"
	glp.ArchiveBrowsingEnabled = &falseValue
	glp.PropertySets = []string{}

	err = testsUpdateLocalRepositoryService.Go(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)

	deleteRepoAndValidate(t, repoKey)
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
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, glp)

	glp.Description += " - Updated"
	glp.Notes = "Repo been updated"
	glp.ArchiveBrowsingEnabled = &trueValue
	glp.ArchiveBrowsingEnabled = &falseValue
	glp.BlockPushingSchema1 = nil

	err = testsUpdateLocalRepositoryService.Generic(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)

	deleteRepoAndValidate(t, repoKey)
}
