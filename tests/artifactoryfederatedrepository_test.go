package tests

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

func TestArtifactoryFederatedRepository(t *testing.T) {
	initArtifactoryTest(t)
	t.Run("federatedAlpineTest", federatedAlpineTest)
	t.Run("federatedBowerTest", federatedBowerTest)
	t.Run("federatedCargoTest", federatedCargoTest)
	t.Run("federatedChefTest", federatedChefTest)
	t.Run("federatedCocoapodsTest", federatedCocoapodsTest)
	t.Run("federatedComposerTest", federatedComposerTest)
	t.Run("federatedConanTest", federatedConanTest)
	t.Run("federatedCondaTest", federatedCondaTest)
	t.Run("federatedCranTest", federatedCranTest)
	t.Run("federatedDebianTest", federatedDebianTest)
	t.Run("federatedDockerTest", federatedDockerTest)
	t.Run("federatedGemsTest", federatedGemsTest)
	t.Run("federatedGenericTest", federatedGenericTest)
	t.Run("federatedGitlfsTest", federatedGitlfsTest)
	t.Run("federatedGoTest", federatedGoTest)
	t.Run("federatedGradleTest", federatedGradleTest)
	t.Run("federatedHelmTest", federatedHelmTest)
	t.Run("federatedIvyTest", federatedIvyTest)
	t.Run("federatedMavenTest", federatedMavenTest)
	t.Run("federatedNpmTest", federatedNpmTest)
	t.Run("federatedNugetTest", federatedNugetTest)
	t.Run("federatedOkgTest", federatedOpkgTest)
	t.Run("federatedPuppetTest", federatedPuppetTest)
	t.Run("federatedPypiTest", federatedPypiTest)
	t.Run("federatedRpmTest", federatedRpmTest)
	t.Run("federatedSbtTest", federatedSbtTest)
	t.Run("federatedVagrantTest", federatedVagrantTest)
	t.Run("federatedYumTest", federatedYumTest)
	t.Run("federatedCreateWithParamTest", federatedCreateWithParamTest)
	t.Run("getFederatedRepoDetailsTest", getFederatedRepoDetailsTest)
	t.Run("getAllFederatedRepoDetailsTest", getAllFederatedRepoDetailsTest)
}

func federatedAlpineTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	alp := services.NewAlpineFederatedRepositoryParams()
	alp.Key = repoKey
	alp.RepoLayoutRef = "simple-default"
	alp.Description = "Alpine Repo for jfrog-client-go federated-repository-test"
	alp.BlackedOut = &trueValue
	alp.XrayIndex = &falseValue
	alp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Alpine(alp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, alp)

	alp.Description += " - Updated"
	alp.Notes = "Repo has been updated"
	alp.ArchiveBrowsingEnabled = &falseValue
	alp.IncludesPattern = "dir1/*"
	alp.ExcludesPattern = "dir2/*"
	alp.BlackedOut = &falseValue
	alp.XrayIndex = &trueValue
	alp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Alpine(alp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, alp)
}

func federatedBowerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	blp := services.NewBowerFederatedRepositoryParams()
	blp.Key = repoKey
	blp.RepoLayoutRef = "bower-default"
	blp.Description = "Bower Repo for jfrog-client-go federated-repository-test"
	blp.BlackedOut = &trueValue
	blp.XrayIndex = &falseValue
	blp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Bower(blp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, blp)

	blp.Description += " - Updated"
	blp.Notes = "Repo has been updated"
	blp.ArchiveBrowsingEnabled = &falseValue
	blp.IncludesPattern = "dir1/*"
	blp.ExcludesPattern = "dir2/*"
	blp.BlackedOut = &falseValue
	blp.XrayIndex = &trueValue
	blp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Bower(blp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, blp)
}

func federatedCargoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewCargoFederatedRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "cargo-default"
	clp.Description = "Cran Repo for jfrog-client-go federated-repository-test"
	clp.IncludesPattern = "dir1/*"
	clp.ExcludesPattern = "dir2/*"
	clp.BlackedOut = &falseValue
	clp.XrayIndex = &trueValue
	clp.CargoAnonymousAccess = &trueValue
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Cargo(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo has been updated"
	clp.ArchiveBrowsingEnabled = &falseValue
	clp.ExcludesPattern = ""
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &falseValue
	clp.CargoAnonymousAccess = &falseValue
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Cargo(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func federatedChefTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewChefFederatedRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "simple-default"
	clp.Description = "Chef Repo for jfrog-client-go federated-repository-test"
	clp.DownloadRedirect = &falseValue
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &falseValue
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Chef(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo has been updated"
	clp.ArchiveBrowsingEnabled = &falseValue
	clp.IncludesPattern = "dir1/*"
	clp.ExcludesPattern = "dir2/*"
	clp.BlackedOut = &falseValue
	clp.XrayIndex = &trueValue
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Chef(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func federatedCocoapodsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewCocoapodsFederatedRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "simple-default"
	clp.Description = "Cocoapods Repo for jfrog-client-go federated-repository-test"
	clp.IncludesPattern = "*/**"
	clp.ExcludesPattern = "dir1/*"
	clp.BlackedOut = &falseValue
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Cocoapods(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo has been updated"
	clp.ArchiveBrowsingEnabled = &trueValue
	clp.ArchiveBrowsingEnabled = &trueValue
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &trueValue
	clp.DownloadRedirect = &falseValue
	clp.ExcludesPattern = "dir1/*,dir2/dir4/*,"
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Cocoapods(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func federatedComposerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewComposerFederatedRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "composer-default"
	clp.Description = "Composer Repo for jfrog-client-go federated-repository-test"
	clp.DownloadRedirect = &falseValue
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &trueValue
	clp.IncludesPattern = "dir1/*,"
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Composer(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo has been updated"
	clp.ArchiveBrowsingEnabled = &trueValue
	clp.ArchiveBrowsingEnabled = &trueValue
	clp.BlackedOut = &falseValue
	clp.XrayIndex = &falseValue
	clp.IncludesPattern = "*/**,"
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Composer(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func federatedConanTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewConanFederatedRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "conan-default"
	clp.Description = "Conan Repo for jfrog-client-go federated-repository-test"
	clp.IncludesPattern = "*/**"
	clp.ExcludesPattern = "ConanEx/*"
	clp.ArchiveBrowsingEnabled = &trueValue
	clp.XrayIndex = &trueValue
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Conan(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo has been updated"
	clp.ArchiveBrowsingEnabled = &falseValue
	clp.ExcludesPattern = ""
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &falseValue
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Conan(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func federatedCondaTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewCondaFederatedRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "simple-default"
	clp.Description = "Conda Repo for jfrog-client-go federated-repository-test"
	clp.IncludesPattern = "*/**"
	clp.ExcludesPattern = "CondaEx/*"
	clp.ArchiveBrowsingEnabled = &trueValue
	clp.XrayIndex = &trueValue
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Conda(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo has been updated"
	clp.ArchiveBrowsingEnabled = &falseValue
	clp.ExcludesPattern = ""
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &falseValue
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Conda(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func federatedCranTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewCranFederatedRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "simple-default"
	clp.Description = "Cran Repo for jfrog-client-go federated-repository-test"
	clp.IncludesPattern = "dir1/*"
	clp.ExcludesPattern = "dir2/*"
	clp.BlackedOut = &falseValue
	clp.XrayIndex = &trueValue
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Cran(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo has been updated"
	clp.ArchiveBrowsingEnabled = &falseValue
	clp.ExcludesPattern = ""
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &falseValue
	clp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Cran(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func federatedDebianTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	dlp := services.NewDebianFederatedRepositoryParams()
	dlp.Key = repoKey
	dlp.RepoLayoutRef = "simple-default"
	dlp.Description = "Debian Repo for jfrog-client-go federated-repository-test"
	dlp.IncludesPattern = "Debian1/*,dir3/*"
	dlp.ExcludesPattern = "dir3/*"
	dlp.DebianTrivialLayout = &trueValue
	dlp.BlackedOut = &falseValue
	dlp.XrayIndex = &trueValue
	dlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Debian(dlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, dlp)

	dlp.Description += " - Updated"
	dlp.Notes = "Repo has been updated"
	dlp.ArchiveBrowsingEnabled = &falseValue
	dlp.IncludesPattern = "*/**"
	dlp.ExcludesPattern = ""
	dlp.DebianTrivialLayout = &falseValue
	dlp.BlackedOut = &trueValue
	dlp.XrayIndex = &falseValue
	dlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Debian(dlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dlp)
}

func federatedDockerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	dlp := services.NewDockerFederatedRepositoryParams()
	dlp.Key = repoKey
	dlp.RepoLayoutRef = "simple-default"
	dlp.Description = "Docker Repo for jfrog-client-go federated-repository-test"
	dlp.IncludesPattern = "*/**"
	dlp.BlackedOut = &falseValue
	dlp.DockerApiVersion = "V1"
	dlp.MaxUniqueTags = 18
	dlp.BlockPushingSchema1 = &falseValue
	dlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Docker(dlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, dlp)

	dlp.Description += " - Updated"
	dlp.Notes = "Repo has been updated"
	dlp.ArchiveBrowsingEnabled = &trueValue
	dlp.ArchiveBrowsingEnabled = &trueValue
	dlp.BlackedOut = &trueValue
	dlp.XrayIndex = &trueValue
	dlp.IncludesPattern = "dir1/*,dir3/*"
	dlp.ExcludesPattern = "dir2/*"
	dlp.DockerApiVersion = "V2"
	dlp.MaxUniqueTags = 36
	dlp.BlockPushingSchema1 = &trueValue
	dlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Docker(dlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dlp)
}

func federatedGemsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGemsFederatedRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "simple-default"
	glp.Description = "Gems Repo for jfrog-client-go federated-repository-test"
	glp.IncludesPattern = "*/**"
	glp.ExcludesPattern = "dirEx/*"
	glp.BlackedOut = &trueValue
	glp.ArchiveBrowsingEnabled = &trueValue
	glp.XrayIndex = &trueValue
	glp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Gems(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, glp)

	glp.Description += " - Updated"
	glp.Notes = "Repo has been updated"
	glp.ArchiveBrowsingEnabled = &falseValue
	glp.ExcludesPattern = ""
	glp.BlackedOut = &falseValue
	glp.XrayIndex = &falseValue
	glp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Gems(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)
}

func federatedGenericTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGenericFederatedRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "simple-default"
	glp.Description = "Generic Repo for jfrog-client-go federated-repository-test"
	glp.XrayIndex = &trueValue
	glp.DownloadRedirect = &falseValue
	glp.ArchiveBrowsingEnabled = &falseValue
	glp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Generic(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, glp)

	glp.Description += " - Updated"
	glp.Notes = "Repo has been updated"
	glp.ArchiveBrowsingEnabled = &trueValue
	glp.ArchiveBrowsingEnabled = &falseValue
	glp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Generic(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)
}

func federatedGitlfsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGitlfsFederatedRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "simple-default"
	glp.Description = "Gitlfs Repo for jfrog-client-go federated-repository-test"
	glp.IncludesPattern = "dir1/*,dir3/*"
	glp.ExcludesPattern = "dir3/*"
	glp.BlackedOut = &falseValue
	glp.XrayIndex = &trueValue
	glp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Gitlfs(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, glp)

	glp.Description += " - Updated"
	glp.Notes = "Repo has been updated"
	glp.ArchiveBrowsingEnabled = &falseValue
	glp.ExcludesPattern = ""
	glp.BlackedOut = &trueValue
	glp.XrayIndex = &falseValue
	glp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Gitlfs(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)
}

func federatedGoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGoFederatedRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "go-default"
	glp.Description = "Go Repo for jfrog-client-go federated-repository-test"
	glp.XrayIndex = &trueValue
	glp.DownloadRedirect = &falseValue
	glp.PropertySets = []string{"artifactory"}
	glp.ArchiveBrowsingEnabled = &trueValue
	glp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Go(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, glp)

	glp.Description += " - Updated"
	glp.Notes = "Repo has been updated"
	glp.ArchiveBrowsingEnabled = &falseValue
	glp.PropertySets = []string{}
	glp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Go(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)
}

func federatedGradleTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGradleFederatedRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "maven-2-default"
	glp.Description = "Gradle Repo for jfrog-client-go federated-repository-test"
	glp.SuppressPomConsistencyChecks = &trueValue
	glp.HandleReleases = &trueValue
	glp.HandleSnapshots = &falseValue
	glp.XrayIndex = &trueValue
	glp.MaxUniqueSnapshots = 18
	glp.ChecksumPolicyType = "server-generated-checksums"
	glp.DownloadRedirect = &falseValue
	glp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Gradle(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, glp)

	glp.Description += " - Updated"
	glp.MaxUniqueSnapshots = 36
	glp.HandleReleases = nil
	glp.HandleSnapshots = &trueValue
	glp.ChecksumPolicyType = "client-checksums"
	glp.Notes = "Repo has been updated"
	glp.ArchiveBrowsingEnabled = &trueValue
	glp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Gradle(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)
}

func federatedHelmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	hlp := services.NewHelmFederatedRepositoryParams()
	hlp.Key = repoKey
	hlp.RepoLayoutRef = "simple-default"
	hlp.Description = "Helm Repo for jfrog-client-go federated-repository-test"
	hlp.IncludesPattern = "*/**"
	hlp.BlackedOut = &falseValue
	hlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Helm(hlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, hlp)

	hlp.Description += " - Updated"
	hlp.Notes = "Repo has been updated"
	hlp.ArchiveBrowsingEnabled = &trueValue
	hlp.ArchiveBrowsingEnabled = &trueValue
	hlp.BlackedOut = &trueValue
	hlp.XrayIndex = &trueValue
	hlp.IncludesPattern = "dir1/*,dir3/*"
	hlp.ExcludesPattern = "dir2/*"
	hlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Helm(hlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, hlp)
}

func federatedIvyTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	ilp := services.NewIvyFederatedRepositoryParams()
	ilp.Key = repoKey
	ilp.RepoLayoutRef = "ivy-default"
	ilp.Description = "Ivy Repo for jfrog-client-go federated-repository-test"
	ilp.IncludesPattern = "dir1/*,dir3/*"
	ilp.ExcludesPattern = "dir3/*"
	ilp.SuppressPomConsistencyChecks = &trueValue
	ilp.HandleReleases = &trueValue
	ilp.HandleSnapshots = &falseValue
	ilp.XrayIndex = &trueValue
	ilp.MaxUniqueSnapshots = 18
	ilp.ChecksumPolicyType = "server-generated-checksums"
	ilp.BlackedOut = &falseValue
	ilp.XrayIndex = &trueValue
	ilp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Ivy(ilp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, ilp)

	ilp.Description += " - Updated"
	ilp.MaxUniqueSnapshots = 36
	ilp.HandleReleases = nil
	ilp.HandleSnapshots = &trueValue
	ilp.ChecksumPolicyType = "client-checksums"
	ilp.Notes = "Repo has been updated"
	ilp.ArchiveBrowsingEnabled = &falseValue
	ilp.ExcludesPattern = ""
	ilp.BlackedOut = &trueValue
	ilp.XrayIndex = &falseValue
	ilp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Ivy(ilp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, ilp)
}

func federatedMavenTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	mlp := services.NewMavenFederatedRepositoryParams()
	mlp.Key = repoKey
	mlp.RepoLayoutRef = "maven-2-default"
	mlp.Description = "Maven Repo for jfrog-client-go federated-repository-test"
	mlp.SuppressPomConsistencyChecks = &trueValue
	mlp.HandleReleases = &trueValue
	mlp.HandleSnapshots = &falseValue
	mlp.XrayIndex = &trueValue
	mlp.MaxUniqueSnapshots = 18
	mlp.ChecksumPolicyType = "server-generated-checksums"
	mlp.DownloadRedirect = &falseValue
	mlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Maven(mlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, mlp)

	mlp.Description += " - Updated"
	mlp.MaxUniqueSnapshots = 36
	mlp.HandleReleases = nil
	mlp.HandleSnapshots = &trueValue
	mlp.ChecksumPolicyType = "client-checksums"
	mlp.Notes = "Repo has been updated"
	mlp.ArchiveBrowsingEnabled = &trueValue
	mlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Maven(mlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, mlp)
}

func federatedNpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nlp := services.NewNpmFederatedRepositoryParams()
	nlp.Key = repoKey
	nlp.RepoLayoutRef = "npm-default"
	nlp.Description = "Npm Repo for jfrog-client-go federated-repository-test"
	nlp.IncludesPattern = "dir1/*"
	nlp.ExcludesPattern = "dir2/*"
	nlp.XrayIndex = &trueValue
	nlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Npm(nlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, nlp)

	nlp.Description += " - Updated"
	nlp.Notes = "Repo has been updated"
	nlp.ArchiveBrowsingEnabled = &falseValue
	nlp.IncludesPattern = "dir3/*"
	nlp.ExcludesPattern = "dir4/*,dir5/*"
	nlp.BlackedOut = &trueValue
	nlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Npm(nlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nlp)
}

func federatedNugetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nlp := services.NewNugetFederatedRepositoryParams()
	nlp.Key = repoKey
	nlp.RepoLayoutRef = "nuget-default"
	nlp.Description = "Nuget Repo for jfrog-client-go federated-repository-test"
	nlp.IncludesPattern = "dir1/*"
	nlp.ExcludesPattern = "dir2/*"
	nlp.XrayIndex = &trueValue
	nlp.ForceNugetAuthentication = &falseValue
	nlp.MaxUniqueSnapshots = 24
	nlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Nuget(nlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, nlp)

	nlp.Description += " - Updated"
	nlp.Notes = "Repo has been updated"
	nlp.ArchiveBrowsingEnabled = &falseValue
	nlp.IncludesPattern = "dir3/*"
	nlp.ExcludesPattern = "dir4/*,dir5/*"
	nlp.BlackedOut = &trueValue
	nlp.ForceNugetAuthentication = &trueValue
	nlp.MaxUniqueSnapshots = 18
	nlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Nuget(nlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nlp)
}

func federatedOpkgTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	olp := services.NewOpkgFederatedRepositoryParams()
	olp.Key = repoKey
	olp.RepoLayoutRef = "simple-default"
	olp.Description = "Opkg Repo for jfrog-client-go federated-repository-test"
	olp.DownloadRedirect = &falseValue
	olp.BlackedOut = &trueValue
	olp.XrayIndex = &trueValue
	olp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Opkg(olp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, olp)

	olp.Description += " - Updated"
	olp.Notes = "Repo has been updated"
	olp.ArchiveBrowsingEnabled = &trueValue
	olp.ArchiveBrowsingEnabled = &trueValue
	olp.BlackedOut = &falseValue
	olp.XrayIndex = &falseValue
	olp.IncludesPattern = "dir1/*,"
	olp.ExcludesPattern = "dir3/*"
	olp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Opkg(olp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, olp)
}

func federatedPuppetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	plp := services.NewPuppetFederatedRepositoryParams()
	plp.Key = repoKey
	plp.RepoLayoutRef = "puppet-default"
	plp.Description = "puppet Repo for jfrog-client-go federated-repository-test"
	plp.BlackedOut = &falseValue
	plp.XrayIndex = &falseValue
	plp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Puppet(plp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, plp)
	plp.Description += " - Updated"
	plp.Notes = "Repo has been updated"
	plp.ArchiveBrowsingEnabled = &falseValue
	plp.IncludesPattern = "dir1/*"
	plp.ExcludesPattern = "dir2/*"
	plp.BlackedOut = &trueValue
	plp.XrayIndex = &trueValue
	plp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Puppet(plp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, plp)
}

func federatedPypiTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	plp := services.NewPypiFederatedRepositoryParams()
	plp.Key = repoKey
	plp.RepoLayoutRef = "simple-default"
	plp.Description = "Pypi Repo for jfrog-client-go federated-repository-test"
	plp.BlackedOut = &falseValue
	plp.XrayIndex = &falseValue
	plp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Pypi(plp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, plp)
	plp.Description += " - Updated"
	plp.Notes = "Repo has been updated"
	plp.ArchiveBrowsingEnabled = &falseValue
	plp.IncludesPattern = "dir1/*"
	plp.ExcludesPattern = "dir2/*"
	plp.BlackedOut = &trueValue
	plp.XrayIndex = &trueValue
	plp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Pypi(plp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, plp)
}

func federatedRpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	rlp := services.NewRpmFederatedRepositoryParams()
	rlp.Key = repoKey
	rlp.RepoLayoutRef = "simple-default"
	rlp.Description = "Rpm Repo for jfrog-client-go federated-repository-test"
	rlp.XrayIndex = &trueValue
	rlp.DownloadRedirect = &falseValue
	rlp.YumRootDepth = 6
	rlp.CalculateYumMetadata = &falseValue
	rlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Rpm(rlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, rlp)

	rlp.Description += " - Updated"
	rlp.Notes = "Repo has been updated"
	rlp.ArchiveBrowsingEnabled = &trueValue
	rlp.YumRootDepth = 18
	rlp.CalculateYumMetadata = &trueValue
	rlp.EnableFileListsIndexing = &falseValue
	rlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Rpm(rlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, rlp)
}

func federatedSbtTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	slp := services.NewSbtFederatedRepositoryParams()
	slp.Key = repoKey
	slp.RepoLayoutRef = "sbt-default"
	slp.Description = "Sbt Repo for jfrog-client-go federated-repository-test"
	slp.IncludesPattern = "dir1/*,dir2/*"
	slp.ExcludesPattern = "dir3/*"
	slp.SuppressPomConsistencyChecks = &trueValue
	slp.HandleReleases = &trueValue
	slp.HandleSnapshots = &falseValue
	slp.XrayIndex = &trueValue
	slp.MaxUniqueSnapshots = 18
	slp.ChecksumPolicyType = "server-generated-checksums"
	slp.BlackedOut = &falseValue
	slp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Sbt(slp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, slp)

	slp.Description += " - Updated"
	slp.MaxUniqueSnapshots = 36
	slp.HandleReleases = nil
	slp.HandleSnapshots = &trueValue
	slp.ChecksumPolicyType = "client-checksums"
	slp.Notes = "Repo has been updated"
	slp.ArchiveBrowsingEnabled = &trueValue
	slp.BlackedOut = &trueValue
	slp.XrayIndex = &trueValue
	slp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Sbt(slp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, slp)
}

func federatedVagrantTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	vlp := services.NewVagrantFederatedRepositoryParams()
	vlp.Key = repoKey
	vlp.RepoLayoutRef = "simple-default"
	vlp.Description = "Vagrant Repo for jfrog-client-go federated-repository-test"
	vlp.DownloadRedirect = &falseValue
	vlp.BlackedOut = &trueValue
	vlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Vagrant(vlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, vlp)

	vlp.Description += " - Updated"
	vlp.Notes = "Repo has been updated"
	vlp.ArchiveBrowsingEnabled = &trueValue
	vlp.ArchiveBrowsingEnabled = &trueValue
	vlp.BlackedOut = &falseValue
	vlp.XrayIndex = &trueValue
	vlp.IncludesPattern = "dir3/*,"
	vlp.ExcludesPattern = "dir1/*,dir2/*"
	vlp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Vagrant(vlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, vlp)
}

func federatedYumTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	ylp := services.NewYumFederatedRepositoryParams()
	ylp.Key = repoKey
	ylp.RepoLayoutRef = "simple-default"
	ylp.Description = "Yum Repo for jfrog-client-go federated-repository-test"
	ylp.IncludesPattern = "dir1/*"
	ylp.ExcludesPattern = "dir2/*"
	ylp.BlackedOut = &falseValue
	ylp.XrayIndex = &trueValue
	ylp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &trueValue},
	}

	err := testsCreateFederatedRepositoryService.Yum(ylp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// "yum" package type is converted to "rpm" by Artifactory, so we have to change it too to pass the validation.
	ylp.PackageType = "rpm"
	validateRepoConfig(t, repoKey, ylp)

	ylp.Description += " - Updated"
	ylp.Notes = "Repo has been updated"
	ylp.ArchiveBrowsingEnabled = &falseValue
	ylp.ExcludesPattern = ""
	ylp.BlackedOut = &trueValue
	ylp.XrayIndex = &falseValue
	ylp.Members = []services.FederatedRepositoryMemberParams{
		{Url: testsUpdateFederatedRepositoryService.ArtDetails.GetUrl() + "artifactory/" + repoKey, Enabled: &falseValue},
	}

	err = testsUpdateFederatedRepositoryService.Yum(ylp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, ylp)
}

func federatedCreateWithParamTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	params := services.NewFederatedRepositoryBaseParams()
	params.Key = repoKey
	err := testsRepositoriesService.CreateFederated(params)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, params)
}

func getFederatedRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGenericFederatedRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "simple-default"
	glp.Description = "Generic Repo for jfrog-client-go federated-repository-test"
	glp.XrayIndex = &trueValue
	glp.DownloadRedirect = &falseValue
	glp.ArchiveBrowsingEnabled = &falseValue

	err := testsCreateFederatedRepositoryService.Generic(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// Get repo details
	data := getRepo(t, repoKey)
	// Validate
	assert.Equal(t, data.Key, repoKey)
	assert.Equal(t, data.Description, glp.Description)
	assert.Equal(t, data.GetRepoType(), "federated")
	assert.Empty(t, data.Url)
	assert.Equal(t, data.PackageType, "generic")
}

func getAllFederatedRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGenericFederatedRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "simple-default"
	glp.Description = "Generic Repo for jfrog-client-go federated-repository-test"
	glp.XrayIndex = &trueValue
	glp.DownloadRedirect = &falseValue
	glp.ArchiveBrowsingEnabled = &falseValue

	err := testsCreateFederatedRepositoryService.Generic(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// Get repo details
	data := getAllRepos(t, "federated", "")
	assert.NotNil(t, data)
	repo := &services.RepositoryDetails{}
	for _, v := range *data {
		if v.Key == repoKey {
			repo = &v
			break
		}
	}
	// Validate
	assert.NotNil(t, repo, "Repo "+repoKey+" not found")
	assert.Equal(t, glp.Description, repo.Description)
	assert.Equal(t, "Generic", repo.PackageType)
}
