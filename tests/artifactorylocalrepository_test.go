package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArtifactoryLocalRepository(t *testing.T) {
	t.Run("localMavenTest", localMavenTest)
	t.Run("localGradleTest", localGradleTest)
	t.Run("localIvyTest", localIvyTest)
	t.Run("localSbtTest", localSbtTest)
	t.Run("localHelmTest", localHelmTest)
	t.Run("localRpmTest", localRpmTest)
	t.Run("localNugetTest", localNugetTest)
	t.Run("localCranTest", localCranTest)
	t.Run("localGemsTest", localGemsTest)
	t.Run("localNpmTest", localNpmTest)
	t.Run("localBowerTest", localBowerTest)
	t.Run("localDebianTest", localDebianTest)
	t.Run("localPypiTest", localPypiTest)
	t.Run("localDockerTest", localDockerTest)
	t.Run("localGitlfsTest", localGitlfsTest)
	t.Run("localGoTest", localGoTest)
	t.Run("localYumTest", localYumTest)
	t.Run("localConanTest", localConanTest)
	t.Run("localChefTest", localChefTest)
	t.Run("localPuppetTest", localPuppetTest)
	t.Run("localCocoapodsTest", localCocoapodsTest)
	t.Run("localOkgTest", localOpkgTest)
	t.Run("localComposerTest", localComposerTest)
	t.Run("localvagrantTest", localVagrantTest)
	t.Run("localGenericTest", localGenericTest)
}

func generateLocalRepoKey(pkgType string) (repoKey string) {
	return GenerateRepoKeyForRepoServiceTest(pkgType, "local")
}

func localMavenTest(t *testing.T) {
	repoKey := generateLocalRepoKey("maven")
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
	repoKey := generateLocalRepoKey("gradle")
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

func localIvyTest(t *testing.T) {
	repoKey := generateLocalRepoKey("ivy")
	ilp := services.NewIvyLocalRepositoryParams()
	ilp.Key = repoKey
	ilp.RepoLayoutRef = "ivy-default"
	ilp.Description = "Ivy Repo for jfrog-client-go local-repository-test"
	ilp.IncludesPattern = "dir1/*,dir3/*"
	ilp.ExcludesPattern = "dir3/*"
	ilp.DownloadRedirect = &trueValue
	ilp.BlackedOut = &falseValue
	ilp.XrayIndex = &trueValue

	err := testsCreateLocalRepositoryService.Ivy(ilp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, ilp)

	ilp.Description += " - Updated"
	ilp.Notes = "Repo been updated"
	ilp.ArchiveBrowsingEnabled = &falseValue
	ilp.ExcludesPattern = ""
	ilp.BlackedOut = &trueValue
	ilp.XrayIndex = &falseValue

	err = testsUpdateLocalRepositoryService.Ivy(ilp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, ilp)

	deleteRepoAndValidate(t, repoKey)
}

func localSbtTest(t *testing.T) {
	repoKey := generateLocalRepoKey("sbt")
	slp := services.NewSbtLocalRepositoryParams()
	slp.Key = repoKey
	slp.RepoLayoutRef = "sbt-default"
	slp.Description = "Sbt Repo for jfrog-client-go local-repository-test"
	slp.IncludesPattern = "dir1/*,dir2/*"
	slp.ExcludesPattern = "dir3/*"
	slp.DownloadRedirect = &trueValue
	slp.BlackedOut = &falseValue

	err := testsCreateLocalRepositoryService.Sbt(slp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, slp)

	slp.Description += " - Updated"
	slp.Notes = "Repo been updated"
	slp.ArchiveBrowsingEnabled = &trueValue
	slp.BlackedOut = &trueValue
	slp.XrayIndex = &trueValue

	err = testsUpdateLocalRepositoryService.Sbt(slp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, slp)

	deleteRepoAndValidate(t, repoKey)
}

func localHelmTest(t *testing.T) {
	repoKey := generateLocalRepoKey("helm")
	hlp := services.NewHelmLocalRepositoryParams()
	hlp.Key = repoKey
	hlp.RepoLayoutRef = "simple-default"
	hlp.Description = "Helm Repo for jfrog-client-go local-repository-test"
	hlp.IncludesPattern = "*/**"
	hlp.DownloadRedirect = &trueValue
	hlp.BlackedOut = &falseValue

	err := testsCreateLocalRepositoryService.Helm(hlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, hlp)

	hlp.Description += " - Updated"
	hlp.Notes = "Repo been updated"
	hlp.ArchiveBrowsingEnabled = &trueValue
	hlp.ArchiveBrowsingEnabled = &trueValue
	hlp.BlackedOut = &trueValue
	hlp.XrayIndex = &trueValue
	hlp.IncludesPattern = "dir1/*,dir3/*"
	hlp.ExcludesPattern = "dir2/*"

	err = testsUpdateLocalRepositoryService.Helm(hlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, hlp)

	deleteRepoAndValidate(t, repoKey)
}

func localRpmTest(t *testing.T) {
	repoKey := generateLocalRepoKey("rpm")
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

func localNugetTest(t *testing.T) {
	repoKey := generateLocalRepoKey("nuget")
	nlp := services.NewNugetLocalRepositoryParams()
	nlp.Key = repoKey
	nlp.RepoLayoutRef = "nuget-default"
	nlp.Description = "Nuget Repo for jfrog-client-go local-repository-test"
	nlp.IncludesPattern = "dir1/*"
	nlp.ExcludesPattern = "dir2/*"
	nlp.DownloadRedirect = &trueValue
	nlp.XrayIndex = &trueValue
	nlp.ForceNugetAuthentication = &falseValue
	nlp.MaxUniqueSnapshots = 24

	err := testsCreateLocalRepositoryService.Nuget(nlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, nlp)

	nlp.Description += " - Updated"
	nlp.Notes = "Repo been updated"
	nlp.ArchiveBrowsingEnabled = &falseValue
	nlp.IncludesPattern = "dir3/*"
	nlp.ExcludesPattern = "dir4/*,dir5/*"
	nlp.BlackedOut = &trueValue
	nlp.ForceNugetAuthentication = &trueValue
	nlp.MaxUniqueSnapshots = 18

	err = testsUpdateLocalRepositoryService.Nuget(nlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nlp)

	deleteRepoAndValidate(t, repoKey)
}

func localCranTest(t *testing.T) {
	repoKey := generateLocalRepoKey("cran")
	clp := services.NewCranLocalRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "simple-default"
	clp.Description = "Cran Repo for jfrog-client-go local-repository-test"
	clp.IncludesPattern = "dir1/*"
	clp.ExcludesPattern = "dir2/*"
	clp.DownloadRedirect = &trueValue
	clp.BlackedOut = &falseValue
	clp.XrayIndex = &trueValue

	err := testsCreateLocalRepositoryService.Cran(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo been updated"
	clp.ArchiveBrowsingEnabled = &falseValue
	clp.ExcludesPattern = ""
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &falseValue

	err = testsUpdateLocalRepositoryService.Cran(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)

	deleteRepoAndValidate(t, repoKey)
}

func localGemsTest(t *testing.T) {
	repoKey := generateLocalRepoKey("gems")
	glp := services.NewGemsLocalRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "simple-default"
	glp.Description = "Gems Repo for jfrog-client-go local-repository-test"
	glp.IncludesPattern = "*/**"
	glp.ExcludesPattern = "dirEx/*"
	glp.DownloadRedirect = &trueValue
	glp.BlackedOut = &trueValue
	glp.ArchiveBrowsingEnabled = &trueValue
	glp.XrayIndex = &trueValue

	err := testsCreateLocalRepositoryService.Gems(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, glp)

	glp.Description += " - Updated"
	glp.Notes = "Repo been updated"
	glp.ArchiveBrowsingEnabled = &falseValue
	glp.ExcludesPattern = ""
	glp.BlackedOut = &falseValue
	glp.XrayIndex = &falseValue

	err = testsUpdateLocalRepositoryService.Gems(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)

	deleteRepoAndValidate(t, repoKey)
}

func localNpmTest(t *testing.T) {
	repoKey := generateLocalRepoKey("npm")
	nlp := services.NewNpmLocalRepositoryParams()
	nlp.Key = repoKey
	nlp.RepoLayoutRef = "npm-default"
	nlp.Description = "Npm Repo for jfrog-client-go local-repository-test"
	nlp.IncludesPattern = "dir1/*"
	nlp.ExcludesPattern = "dir2/*"
	nlp.DownloadRedirect = &trueValue
	nlp.XrayIndex = &trueValue

	err := testsCreateLocalRepositoryService.Npm(nlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, nlp)

	nlp.Description += " - Updated"
	nlp.Notes = "Repo been updated"
	nlp.ArchiveBrowsingEnabled = &falseValue
	nlp.IncludesPattern = "dir3/*"
	nlp.ExcludesPattern = "dir4/*,dir5/*"
	nlp.BlackedOut = &trueValue

	err = testsUpdateLocalRepositoryService.Npm(nlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nlp)

	deleteRepoAndValidate(t, repoKey)
}

func localBowerTest(t *testing.T) {
	repoKey := generateLocalRepoKey("bower")
	blp := services.NewBowerLocalRepositoryParams()
	blp.Key = repoKey
	blp.RepoLayoutRef = "bower-default"
	blp.Description = "Boer Repo for jfrog-client-go local-repository-test"
	blp.DownloadRedirect = &trueValue
	blp.BlackedOut = &trueValue
	blp.XrayIndex = &falseValue

	err := testsCreateLocalRepositoryService.Bower(blp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, blp)

	blp.Description += " - Updated"
	blp.Notes = "Repo been updated"
	blp.ArchiveBrowsingEnabled = &falseValue
	blp.IncludesPattern = "dir1/*"
	blp.ExcludesPattern = "dir2/*"
	blp.BlackedOut = &falseValue
	blp.XrayIndex = &trueValue

	err = testsUpdateLocalRepositoryService.Bower(blp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, blp)

	deleteRepoAndValidate(t, repoKey)
}

func localDebianTest(t *testing.T) {
	repoKey := generateLocalRepoKey("debian")
	dlp := services.NewDebianLocalRepositoryParams()
	dlp.Key = repoKey
	dlp.RepoLayoutRef = "simple-default"
	dlp.Description = "Debian Repo for jfrog-client-go local-repository-test"
	dlp.IncludesPattern = "Debian1/*,dir3/*"
	dlp.ExcludesPattern = "dir3/*"
	dlp.DebianTrivialLayout = &trueValue
	dlp.DownloadRedirect = &trueValue
	dlp.BlackedOut = &falseValue
	dlp.XrayIndex = &trueValue

	err := testsCreateLocalRepositoryService.Debian(dlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, dlp)

	dlp.Description += " - Updated"
	dlp.Notes = "Repo been updated"
	dlp.ArchiveBrowsingEnabled = &falseValue
	dlp.IncludesPattern = "*/**"
	dlp.ExcludesPattern = ""
	dlp.DebianTrivialLayout = &falseValue
	dlp.BlackedOut = &trueValue
	dlp.XrayIndex = &falseValue

	err = testsUpdateLocalRepositoryService.Debian(dlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dlp)

	deleteRepoAndValidate(t, repoKey)
}

func localPypiTest(t *testing.T) {
	repoKey := generateLocalRepoKey("pypi")
	plp := services.NewPypiLocalRepositoryParams()
	plp.Key = repoKey
	plp.RepoLayoutRef = "simple-default"
	plp.Description = "Pypi Repo for jfrog-client-go local-repository-test"

	plp.BlackedOut = &falseValue
	plp.XrayIndex = &falseValue

	err := testsCreateLocalRepositoryService.Pypi(plp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, plp)
	plp.Description += " - Updated"
	plp.Notes = "Repo been updated"
	plp.ArchiveBrowsingEnabled = &falseValue
	plp.IncludesPattern = "dir1/*"
	plp.ExcludesPattern = "dir2/*"
	plp.BlackedOut = &trueValue
	plp.XrayIndex = &trueValue
	plp.DownloadRedirect = &trueValue

	err = testsUpdateLocalRepositoryService.Pypi(plp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, plp)

	deleteRepoAndValidate(t, repoKey)
}

func localDockerTest(t *testing.T) {
	repoKey := generateLocalRepoKey("docker")
	dlp := services.NewDockerLocalRepositoryParams()
	dlp.Key = repoKey
	dlp.RepoLayoutRef = "simple-default"
	dlp.Description = "Docker Repo for jfrog-client-go local-repository-test"
	dlp.IncludesPattern = "*/**"
	dlp.DownloadRedirect = &trueValue
	dlp.BlackedOut = &falseValue
	dlp.DockerApiVersion = "V1"
	dlp.MaxUniqueTags = 18

	err := testsCreateLocalRepositoryService.Docker(dlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, dlp)

	dlp.Description += " - Updated"
	dlp.Notes = "Repo been updated"
	dlp.ArchiveBrowsingEnabled = &trueValue
	dlp.ArchiveBrowsingEnabled = &trueValue
	dlp.BlackedOut = &trueValue
	dlp.XrayIndex = &trueValue
	dlp.IncludesPattern = "dir1/*,dir3/*"
	dlp.ExcludesPattern = "dir2/*"
	dlp.DockerApiVersion = "V2"
	dlp.MaxUniqueTags = 36

	err = testsUpdateLocalRepositoryService.Docker(dlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dlp)

	deleteRepoAndValidate(t, repoKey)
}

func localGitlfsTest(t *testing.T) {
	repoKey := generateLocalRepoKey("gitlfs")
	glp := services.NewGitlfsLocalRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "simple-default"
	glp.Description = "Gitlfs Repo for jfrog-client-go local-repository-test"
	glp.IncludesPattern = "dir1/*,dir3/*"
	glp.ExcludesPattern = "dir3/*"
	glp.DownloadRedirect = &trueValue
	glp.BlackedOut = &falseValue
	glp.XrayIndex = &trueValue

	err := testsCreateLocalRepositoryService.Gitlfs(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, glp)

	glp.Description += " - Updated"
	glp.Notes = "Repo been updated"
	glp.ArchiveBrowsingEnabled = &falseValue
	glp.ExcludesPattern = ""
	glp.BlackedOut = &trueValue
	glp.XrayIndex = &falseValue

	err = testsUpdateLocalRepositoryService.Gitlfs(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)

	deleteRepoAndValidate(t, repoKey)
}

func localGoTest(t *testing.T) {
	repoKey := generateLocalRepoKey("go")
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

func localYumTest(t *testing.T) {
	repoKey := generateLocalRepoKey("yum")
	ylp := services.NewYumLocalRepositoryParams()
	ylp.Key = repoKey
	ylp.RepoLayoutRef = "simple-default"
	ylp.Description = "Yum Repo for jfrog-client-go local-repository-test"
	ylp.IncludesPattern = "dir1/*"
	ylp.ExcludesPattern = "dir2/*"
	ylp.DownloadRedirect = &trueValue
	ylp.BlackedOut = &falseValue
	ylp.XrayIndex = &trueValue

	err := testsCreateLocalRepositoryService.Yum(ylp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	// "yum" package type is converted to "rpm" by Artifactory, so we have to change it too to pass the validation.
	ylp.PackageType = "rpm"
	validateRepoConfig(t, repoKey, ylp)

	ylp.Description += " - Updated"
	ylp.Notes = "Repo been updated"
	ylp.ArchiveBrowsingEnabled = &falseValue
	ylp.ExcludesPattern = ""
	ylp.BlackedOut = &trueValue
	ylp.XrayIndex = &falseValue

	err = testsUpdateLocalRepositoryService.Yum(ylp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, ylp)

	deleteRepoAndValidate(t, repoKey)
}

func localConanTest(t *testing.T) {
	repoKey := generateLocalRepoKey("conan")
	clp := services.NewConanLocalRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "conan-default"
	clp.Description = "Conan Repo for jfrog-client-go local-repository-test"
	clp.IncludesPattern = "*/**"
	clp.ExcludesPattern = "ConanEx/*"
	clp.DownloadRedirect = &trueValue
	clp.ArchiveBrowsingEnabled = &trueValue
	clp.XrayIndex = &trueValue

	err := testsCreateLocalRepositoryService.Conan(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo been updated"
	clp.ArchiveBrowsingEnabled = &falseValue
	clp.ExcludesPattern = ""
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &falseValue

	err = testsUpdateLocalRepositoryService.Conan(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)

	deleteRepoAndValidate(t, repoKey)
}

func localChefTest(t *testing.T) {
	repoKey := generateLocalRepoKey("chef")
	clp := services.NewChefLocalRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "simple-default"
	clp.Description = "Chef Repo for jfrog-client-go local-repository-test"
	clp.DownloadRedirect = &falseValue
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &falseValue

	err := testsCreateLocalRepositoryService.Chef(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo been updated"
	clp.ArchiveBrowsingEnabled = &falseValue
	clp.IncludesPattern = "dir1/*"
	clp.ExcludesPattern = "dir2/*"
	clp.BlackedOut = &falseValue
	clp.XrayIndex = &trueValue
	clp.DownloadRedirect = &trueValue

	err = testsUpdateLocalRepositoryService.Chef(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)

	deleteRepoAndValidate(t, repoKey)
}

func localPuppetTest(t *testing.T) {
	repoKey := generateLocalRepoKey("puppet")
	plp := services.NewPuppetLocalRepositoryParams()
	plp.Key = repoKey
	plp.RepoLayoutRef = "puppet-default"
	plp.Description = "puppet Repo for jfrog-client-go local-repository-test"

	plp.BlackedOut = &falseValue
	plp.XrayIndex = &falseValue

	err := testsCreateLocalRepositoryService.Puppet(plp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, plp)
	plp.Description += " - Updated"
	plp.Notes = "Repo been updated"
	plp.ArchiveBrowsingEnabled = &falseValue
	plp.IncludesPattern = "dir1/*"
	plp.ExcludesPattern = "dir2/*"
	plp.BlackedOut = &trueValue
	plp.XrayIndex = &trueValue
	plp.DownloadRedirect = &trueValue

	err = testsUpdateLocalRepositoryService.Puppet(plp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, plp)

	deleteRepoAndValidate(t, repoKey)
}

func localCocoapodsTest(t *testing.T) {
	repoKey := generateLocalRepoKey("cocoapods")
	clp := services.NewCocoapodsLocalRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "simple-default"
	clp.Description = "Cocoapods Repo for jfrog-client-go local-repository-test"
	clp.IncludesPattern = "*/**"
	clp.ExcludesPattern = "dir1/*"
	clp.BlackedOut = &falseValue

	err := testsCreateLocalRepositoryService.Cocoapods(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo been updated"
	clp.ArchiveBrowsingEnabled = &trueValue
	clp.ArchiveBrowsingEnabled = &trueValue
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &trueValue
	clp.DownloadRedirect = &falseValue
	clp.ExcludesPattern = "dir1/*,dir2/dir4/*,"

	err = testsUpdateLocalRepositoryService.Cocoapods(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)

	deleteRepoAndValidate(t, repoKey)
}

func localOpkgTest(t *testing.T) {
	repoKey := generateLocalRepoKey("opkg")
	olp := services.NewOpkgLocalRepositoryParams()
	olp.Key = repoKey
	olp.RepoLayoutRef = "simple-default"
	olp.Description = "Opkg Repo for jfrog-client-go local-repository-test"
	olp.DownloadRedirect = &falseValue
	olp.BlackedOut = &trueValue
	olp.XrayIndex = &trueValue

	err := testsCreateLocalRepositoryService.Opkg(olp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, olp)

	olp.Description += " - Updated"
	olp.Notes = "Repo been updated"
	olp.ArchiveBrowsingEnabled = &trueValue
	olp.ArchiveBrowsingEnabled = &trueValue
	olp.BlackedOut = &falseValue
	olp.XrayIndex = &falseValue
	olp.IncludesPattern = "dir1/*,"
	olp.ExcludesPattern = "dir3/*"

	err = testsUpdateLocalRepositoryService.Opkg(olp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, olp)

	deleteRepoAndValidate(t, repoKey)
}

func localComposerTest(t *testing.T) {
	repoKey := generateLocalRepoKey("composer")
	clp := services.NewComposerLocalRepositoryParams()
	clp.Key = repoKey
	clp.RepoLayoutRef = "composer-default"
	clp.Description = "Composer Repo for jfrog-client-go local-repository-test"
	clp.DownloadRedirect = &falseValue
	clp.BlackedOut = &trueValue
	clp.XrayIndex = &trueValue
	clp.IncludesPattern = "dir1/*,"

	err := testsCreateLocalRepositoryService.Composer(clp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, clp)

	clp.Description += " - Updated"
	clp.Notes = "Repo been updated"
	clp.ArchiveBrowsingEnabled = &trueValue
	clp.ArchiveBrowsingEnabled = &trueValue
	clp.BlackedOut = &falseValue
	clp.XrayIndex = &falseValue
	clp.IncludesPattern = "*/**,"

	err = testsUpdateLocalRepositoryService.Composer(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)

	deleteRepoAndValidate(t, repoKey)
}

func localVagrantTest(t *testing.T) {
	repoKey := generateLocalRepoKey("vagrant")
	vlp := services.NewVagrantLocalRepositoryParams()
	vlp.Key = repoKey
	vlp.RepoLayoutRef = "simple-default"
	vlp.Description = "Vagrant Repo for jfrog-client-go local-repository-test"
	vlp.DownloadRedirect = &falseValue
	vlp.BlackedOut = &trueValue

	err := testsCreateLocalRepositoryService.Vagrant(vlp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	validateRepoConfig(t, repoKey, vlp)

	vlp.Description += " - Updated"
	vlp.Notes = "Repo been updated"
	vlp.ArchiveBrowsingEnabled = &trueValue
	vlp.ArchiveBrowsingEnabled = &trueValue
	vlp.BlackedOut = &falseValue
	vlp.XrayIndex = &trueValue
	vlp.IncludesPattern = "dir3/*,"
	vlp.ExcludesPattern = "dir1/*,dir2/*"

	err = testsUpdateLocalRepositoryService.Vagrant(vlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, vlp)

	deleteRepoAndValidate(t, repoKey)
}

func localGenericTest(t *testing.T) {
	repoKey := generateLocalRepoKey("generic")
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
