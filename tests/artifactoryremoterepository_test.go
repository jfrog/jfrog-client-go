package tests

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

const ArtifactoryLocalFileCacheSuffix = " (local file cache)"

func TestArtifactoryRemoteRepository(t *testing.T) {
	t.Run("remoteMavenTest", remoteMavenTest)
	t.Run("remoteGradleTest", remoteGradleTest)
	t.Run("remoteIvyTest", remoteIvyTest)
	t.Run("remoteSbtTest", remoteSbtTest)
	t.Run("remoteHelmTest", remoteHelmTest)
	t.Run("remoteRpmTest", remoteRpmTest)
	t.Run("remoteNugetTest", remoteNugetTest)
	t.Run("remoteCranTest", remoteCranTest)
	t.Run("remoteGemsTest", remoteGemsTest)
	t.Run("remoteNpmTest", remoteNpmTest)
	t.Run("remoteBowerTest", remoteBowerTest)
	t.Run("remoteDebianTest", remoteDebianTest)
	t.Run("remotePypiTest", remotePypiTest)
	t.Run("remoteDockerTest", remoteDockerTest)
	t.Run("remoteGitlfsTest", remoteGitlfsTest)
	t.Run("remoteGoTest", remoteGoTest)
	t.Run("remoteYumTest", remoteYumTest)
	t.Run("remoteConanTest", remoteConanTest)
	t.Run("remoteChefTest", remoteChefTest)
	t.Run("remotePuppetTest", remotePuppetTest)
	t.Run("remoteComposerTest", remoteComposerTest)
	t.Run("remoteCocoapodsTest", remoteCocoapodsTest)
	t.Run("remoteOpkgTest", remoteOpkgTest)
	t.Run("remoteCondaTest", remoteCondaTest)
	t.Run("remoteP2Test", remoteP2Test)
	t.Run("remoteVcsTest", remoteVcsTest)
	t.Run("remoteGenericTest", remoteGenericTest)
	t.Run("getRemoteRepoDetailsTest", getRemoteRepoDetailsTest)
}

func remoteMavenTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
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
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
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
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, mrp)
}

func remoteGradleTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
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
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
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
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, grp)
}

func remoteIvyTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	irp := services.NewIvyRemoteRepositoryParams()
	irp.Key = repoKey
	irp.RepoLayoutRef = "ivy-default"
	irp.Url = "https://jcenter.bintray.com"
	irp.Description = "Ivy Repo for jfrog-client-go remote-repository-test"
	irp.AssumedOfflinePeriodSecs = 8080
	irp.StoreArtifactsLocally = &trueValue
	irp.ShareConfiguration = &trueValue

	err := testsCreateRemoteRepositoryService.Ivy(irp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	irp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, irp)

	irp.Description += " - Updated"
	irp.Notes = "Repo been updated"
	irp.AssumedOfflinePeriodSecs = 9090
	irp.EnableCookieManagement = &trueValue
	irp.SocketTimeoutMillis = 1818
	irp.ShareConfiguration = &falseValue

	err = testsUpdateRemoteRepositoryService.Ivy(irp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, irp)
}

func remoteSbtTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	srp := services.NewSbtRemoteRepositoryParams()
	srp.Key = repoKey
	srp.RepoLayoutRef = "sbt-default"
	srp.Url = "https://jcenter.bintray.com"
	srp.Description = "Sbt Repo for jfrog-client-go remote-repository-test"
	srp.AssumedOfflinePeriodSecs = 9999
	srp.StoreArtifactsLocally = &trueValue
	srp.ShareConfiguration = &falseValue

	err := testsCreateRemoteRepositoryService.Sbt(srp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	srp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, srp)

	srp.Notes = "Repo been updated"
	srp.AssumedOfflinePeriodSecs = 1818
	srp.EnableCookieManagement = &falseValue
	srp.SocketTimeoutMillis = 1111
	srp.ShareConfiguration = &trueValue
	srp.StoreArtifactsLocally = &falseValue
	srp.ShareConfiguration = &falseValue
	srp.AllowAnyHostAuth = &trueValue

	err = testsUpdateRemoteRepositoryService.Sbt(srp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, srp)
}

func remoteHelmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	hrp := services.NewHelmRemoteRepositoryParams()
	hrp.Key = repoKey
	hrp.RepoLayoutRef = "simple-default"
	hrp.Url = "https://storage.googleapis.com/kubernetes-charts"
	hrp.Description = "Helm Repo for jfrog-client-go remote-repository-test"
	hrp.AssumedOfflinePeriodSecs = 5432
	hrp.StoreArtifactsLocally = &falseValue
	hrp.ShareConfiguration = &falseValue
	hrp.BlackedOut = &trueValue
	hrp.IncludesPattern = "*/**"

	err := testsCreateRemoteRepositoryService.Helm(hrp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, hrp)

	hrp.Description += " - Updated"
	hrp.Notes = "Repo been updated"
	hrp.AssumedOfflinePeriodSecs = 2000
	hrp.EnableCookieManagement = &trueValue
	hrp.BypassHeadRequests = &trueValue
	hrp.IncludesPattern = "dir1/*,dir5/*"
	hrp.SocketTimeoutMillis = 666
	hrp.BlackedOut = &falseValue

	err = testsUpdateRemoteRepositoryService.Helm(hrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, hrp)
}

func remoteRpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	rrp := services.NewRpmRemoteRepositoryParams()
	rrp.Key = repoKey
	rrp.RepoLayoutRef = "simple-default"
	rrp.Url = "http://mirror.centos.org/centos/"
	rrp.Description = "Rpm Repo for jfrog-client-go remote-repository-test"
	rrp.AssumedOfflinePeriodSecs = 5555
	rrp.ListRemoteFolderItems = &falseValue
	rrp.StoreArtifactsLocally = &trueValue
	rrp.ShareConfiguration = &trueValue

	err := testsCreateRemoteRepositoryService.Rpm(rrp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	rrp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, rrp)

	rrp.Notes = "Repo been updated"
	rrp.AssumedOfflinePeriodSecs = 1818
	rrp.ListRemoteFolderItems = &trueValue
	rrp.AssumedOfflinePeriodSecs = 2525
	rrp.EnableCookieManagement = &trueValue
	rrp.SocketTimeoutMillis = 1010
	rrp.ShareConfiguration = &falseValue

	err = testsUpdateRemoteRepositoryService.Rpm(rrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, rrp)
}

func remoteNugetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nrp := services.NewNugetRemoteRepositoryParams()
	nrp.Key = repoKey
	nrp.RepoLayoutRef = "nuget-default"
	nrp.Url = "https://www.nuget.org/"
	nrp.Description = "NuGet Repo for jfrog-client-go remote-repository-test"
	nrp.AssumedOfflinePeriodSecs = 3600
	nrp.StoreArtifactsLocally = &falseValue
	nrp.ShareConfiguration = &falseValue
	nrp.BypassHeadRequests = &trueValue
	nrp.DownloadContextPath = "api/v1"
	nrp.V3FeedUrl = "https://api.nuget.org/v3/index.json"

	err := testsCreateRemoteRepositoryService.Nuget(nrp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, nrp)

	nrp.Notes = "Repo been updated"
	nrp.AssumedOfflinePeriodSecs = 1818
	nrp.AssumedOfflinePeriodSecs = 2525
	nrp.EnableCookieManagement = &trueValue
	nrp.SocketTimeoutMillis = 1010
	nrp.ShareConfiguration = &trueValue
	nrp.BlackedOut = &trueValue
	nrp.DownloadContextPath = "api/v2"
	nrp.ForceNugetAuthentication = &trueValue
	nrp.DownloadContextPath = "https://api.nuget.org/v3/index.json"

	err = testsUpdateRemoteRepositoryService.Nuget(nrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nrp)
}

func remoteCranTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewCranRemoteRepositoryParams()
	crp.Key = repoKey
	crp.RepoLayoutRef = "simple-default"
	crp.Url = "https://cran.r-project.org/"
	crp.Description = "Cran Repo for jfrog-client-go remote-repository-test"
	crp.AssumedOfflinePeriodSecs = 8080
	crp.StoreArtifactsLocally = &trueValue
	crp.ShareConfiguration = &trueValue

	err := testsCreateRemoteRepositoryService.Cran(crp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	crp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, crp)

	crp.Description += " - Updated"
	crp.Notes = "Repo been updated"
	crp.AssumedOfflinePeriodSecs = 9090
	crp.EnableCookieManagement = &trueValue
	crp.SocketTimeoutMillis = 1818
	crp.ShareConfiguration = &falseValue

	err = testsUpdateRemoteRepositoryService.Cran(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteGemsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGemsRemoteRepositoryParams()
	grp.Key = repoKey
	grp.RepoLayoutRef = "simple-default"
	grp.Url = "https://rubygems.org/"
	grp.Description = "Gems Repo for jfrog-client-go remote-repository-test"
	grp.AssumedOfflinePeriodSecs = 8080
	grp.StoreArtifactsLocally = &trueValue
	grp.ShareConfiguration = &trueValue
	grp.BlockMismatchingMimeTypes = &trueValue
	grp.IncludesPattern = "**/*"
	grp.ExcludesPattern = "dirEx/*"

	err := testsCreateRemoteRepositoryService.Gems(grp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	grp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, grp)

	grp.Description += " - Updated"
	grp.Notes = "Repo been updated"
	grp.AssumedOfflinePeriodSecs = 5555
	grp.EnableCookieManagement = &trueValue
	grp.SocketTimeoutMillis = 2131
	grp.Offline = &trueValue
	grp.AllowAnyHostAuth = &trueValue
	grp.ShareConfiguration = &falseValue

	err = testsUpdateRemoteRepositoryService.Gems(grp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, grp)
}

func remoteNpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nrp := services.NewNpmRemoteRepositoryParams()
	nrp.Key = repoKey
	nrp.RepoLayoutRef = "npm-default"
	nrp.Url = "https://registry.npmjs.org"
	nrp.Description = "Npm Repo for jfrog-client-go remote-repository-test"
	nrp.AssumedOfflinePeriodSecs = 6060
	nrp.StoreArtifactsLocally = &trueValue
	nrp.ShareConfiguration = &trueValue
	nrp.IncludesPattern = "goDir1/*"
	nrp.ListRemoteFolderItems = &trueValue
	nrp.Offline = &trueValue
	nrp.RetrievalCachePeriodSecs = 999

	err := testsCreateRemoteRepositoryService.Npm(nrp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	nrp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, nrp)

	nrp.Notes = "Repo been updated"
	nrp.AssumedOfflinePeriodSecs = 9090
	nrp.EnableCookieManagement = &trueValue
	nrp.SocketTimeoutMillis = 1111
	nrp.ShareConfiguration = &trueValue
	nrp.StoreArtifactsLocally = &falseValue
	nrp.ShareConfiguration = &falseValue
	nrp.ExcludesPattern = "goDir2/*,dir3/dir4/*"
	nrp.AllowAnyHostAuth = &trueValue
	nrp.Offline = &falseValue

	err = testsUpdateRemoteRepositoryService.Npm(nrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nrp)
}

func remoteBowerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	brp := services.NewBowerRemoteRepositoryParams()
	brp.Key = repoKey
	brp.RepoLayoutRef = "bower-default"
	brp.Url = "https://github.com/"
	brp.Description = "Bower Repo for jfrog-client-go remote-repository-test"
	brp.SocketTimeoutMillis = 5555
	brp.BypassHeadRequests = &falseValue
	brp.BlockMismatchingMimeTypes = &trueValue
	brp.Offline = &trueValue
	brp.BypassHeadRequests = &trueValue
	brp.IncludesPattern = "**/*"

	err := testsCreateRemoteRepositoryService.Bower(brp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	brp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, brp)

	brp.Description += " - Updated"
	brp.Notes = "Repo been updated"
	brp.AssumedOfflinePeriodSecs = 2000
	brp.ShareConfiguration = &trueValue
	brp.BypassHeadRequests = &trueValue
	brp.IncludesPattern = "BowerEx/Dir1/*"
	brp.SocketTimeoutMillis = 666
	brp.BlackedOut = &trueValue
	brp.BlockMismatchingMimeTypes = &falseValue
	brp.BypassHeadRequests = &falseValue

	err = testsUpdateRemoteRepositoryService.Bower(brp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, brp)
}

func remoteDebianTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	drp := services.NewDebianRemoteRepositoryParams()
	drp.Key = repoKey
	drp.RepoLayoutRef = "simple-default"
	drp.Url = "http://archive.ubuntu.com/ubuntu/"
	drp.Description = "Debian Repo for jfrog-client-go remote-repository-test"
	drp.AssumedOfflinePeriodSecs = 6060
	drp.StoreArtifactsLocally = &trueValue
	drp.ShareConfiguration = &trueValue
	drp.IncludesPattern = "goDir1/*"
	drp.ListRemoteFolderItems = &trueValue
	drp.Offline = &trueValue
	drp.RetrievalCachePeriodSecs = 999

	err := testsCreateRemoteRepositoryService.Debian(drp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	drp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, drp)

	drp.Notes = "Repo been updated"
	drp.AssumedOfflinePeriodSecs = 9090
	drp.EnableCookieManagement = &trueValue
	drp.SocketTimeoutMillis = 1111
	drp.ShareConfiguration = &trueValue
	drp.StoreArtifactsLocally = &falseValue
	drp.ShareConfiguration = &falseValue
	drp.ExcludesPattern = "goDir2/*,dir3/dir4/*"
	drp.AllowAnyHostAuth = &trueValue
	drp.ListRemoteFolderItems = &falseValue

	err = testsUpdateRemoteRepositoryService.Debian(drp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, drp)
}

func remotePypiTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewCondaRemoteRepositoryParams()
	crp.Key = repoKey
	crp.RepoLayoutRef = "simple-default"
	crp.Url = "https://repo.anaconda.com/pkgs/free"
	crp.Description = "Conda Repo for jfrog-client-go remote-repository-test"
	crp.AssumedOfflinePeriodSecs = 1800
	crp.StoreArtifactsLocally = &falseValue
	crp.ShareConfiguration = &trueValue

	err := testsCreateRemoteRepositoryService.Conda(crp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	crp.Description += " - Updated"
	crp.Notes = "Repo been updated"
	crp.AssumedOfflinePeriodSecs = 2222
	crp.EnableCookieManagement = &trueValue
	crp.SocketTimeoutMillis = 1818
	crp.ShareConfiguration = &falseValue

	err = testsUpdateRemoteRepositoryService.Conda(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteDockerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	drp := services.NewDockerRemoteRepositoryParams()
	drp.Key = repoKey
	drp.RepoLayoutRef = "simple-default"
	drp.Url = "https://registry-1.docker.io/"
	drp.Description = "Docker Repo for jfrog-client-go remote-repository-test"
	drp.AssumedOfflinePeriodSecs = 8080
	drp.StoreArtifactsLocally = &trueValue
	drp.ShareConfiguration = &trueValue
	drp.SocketTimeoutMillis = 1200
	drp.UnusedArtifactsCleanupPeriodHours = 72

	err := testsCreateRemoteRepositoryService.Docker(drp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	drp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, drp)

	drp.Notes = "Repo been updated"
	drp.AssumedOfflinePeriodSecs = 2202
	drp.EnableCookieManagement = &trueValue
	drp.SocketTimeoutMillis = 1800
	drp.ShareConfiguration = &falseValue
	drp.StoreArtifactsLocally = &falseValue
	drp.ShareConfiguration = &falseValue
	drp.UnusedArtifactsCleanupPeriodHours = 48

	err = testsUpdateRemoteRepositoryService.Docker(drp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, drp)
}

func remoteGitlfsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGitlfsRemoteRepositoryParams()
	grp.Key = repoKey
	grp.RepoLayoutRef = "simple-default"
	grp.Url = "https://github.com/"
	grp.Description = "Gitlfs Repo for jfrog-client-go remote-repository-test"
	grp.AssumedOfflinePeriodSecs = 5555
	grp.StoreArtifactsLocally = &trueValue
	grp.ShareConfiguration = &trueValue
	grp.BypassHeadRequests = &trueValue
	grp.BlockMismatchingMimeTypes = &falseValue
	grp.ListRemoteFolderItems = &trueValue

	err := testsCreateRemoteRepositoryService.Gitlfs(grp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	grp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, grp)

	grp.Description += " - Updated"
	grp.Notes = "Repo been updated"
	grp.AssumedOfflinePeriodSecs = 9090
	grp.EnableCookieManagement = &trueValue
	grp.SocketTimeoutMillis = 1818
	grp.ShareConfiguration = &falseValue
	grp.ListRemoteFolderItems = &falseValue
	grp.BlockMismatchingMimeTypes = &trueValue

	err = testsUpdateRemoteRepositoryService.Gitlfs(grp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, grp)
}

func remoteGoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGoRemoteRepositoryParams()
	grp.Key = repoKey
	grp.RepoLayoutRef = "go-default"
	grp.Url = "https://gocenter.io/"
	grp.Description = "Go Repo for jfrog-client-go remote-repository-test"
	grp.AssumedOfflinePeriodSecs = 6060
	grp.StoreArtifactsLocally = &trueValue
	grp.ShareConfiguration = &trueValue
	grp.IncludesPattern = "goDir1/*"

	err := testsCreateRemoteRepositoryService.Go(grp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	grp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, grp)

	grp.Notes = "Repo been updated"
	grp.AssumedOfflinePeriodSecs = 9090
	grp.EnableCookieManagement = &trueValue
	grp.SocketTimeoutMillis = 1111
	grp.ShareConfiguration = &trueValue
	grp.StoreArtifactsLocally = &falseValue
	grp.ShareConfiguration = &falseValue
	grp.ExcludesPattern = "goDir2/*,dir3/dir4/*"
	grp.AllowAnyHostAuth = &trueValue

	err = testsUpdateRemoteRepositoryService.Go(grp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, grp)
}

func remoteYumTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	yrp := services.NewYumRemoteRepositoryParams()
	yrp.Key = repoKey
	yrp.RepoLayoutRef = "simple-default"
	yrp.Url = "http://mirror.centos.org/centos/"
	yrp.Description = "Yum Repo for jfrog-client-go remote-repository-test"
	yrp.SocketTimeoutMillis = 5555
	yrp.BypassHeadRequests = &falseValue
	yrp.BlockMismatchingMimeTypes = &trueValue
	yrp.BlackedOut = &trueValue
	yrp.IncludesPattern = "*/**"

	err := testsCreateRemoteRepositoryService.Yum(yrp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// "yum" package type is converted to "rpm" by Artifactory, so we have to change it too to pass the validation.
	yrp.PackageType = "rpm"
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	yrp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, yrp)

	yrp.Description += " - Updated"
	yrp.Notes = "Repo been updated"
	yrp.AssumedOfflinePeriodSecs = 2000
	yrp.ShareConfiguration = &trueValue
	yrp.BypassHeadRequests = &trueValue
	yrp.IncludesPattern = "dir1/*,dir5/*"
	yrp.SocketTimeoutMillis = 666
	yrp.Offline = &trueValue
	yrp.BlockMismatchingMimeTypes = &falseValue

	err = testsUpdateRemoteRepositoryService.Yum(yrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, yrp)
}

func remoteConanTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewConanRemoteRepositoryParams()
	crp.Key = repoKey
	crp.RepoLayoutRef = "conan-default"
	crp.Url = "https://conan.bintray.com"
	crp.Description = "Conan Repo for jfrog-client-go remote-repository-test"
	crp.AssumedOfflinePeriodSecs = 1800
	crp.SynchronizeProperties = &trueValue
	crp.StoreArtifactsLocally = &falseValue
	crp.ShareConfiguration = &trueValue
	crp.BlockMismatchingMimeTypes = &falseValue
	crp.BypassHeadRequests = &trueValue

	err := testsCreateRemoteRepositoryService.Conan(crp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	crp.Description += " - Updated"
	crp.Notes = "Repo been updated"
	crp.AssumedOfflinePeriodSecs = 2222
	crp.EnableCookieManagement = &trueValue
	crp.SocketTimeoutMillis = 1818
	crp.ShareConfiguration = &falseValue
	crp.AssumedOfflinePeriodSecs = 4321
	crp.BypassHeadRequests = &falseValue
	crp.BlockMismatchingMimeTypes = &trueValue

	err = testsUpdateRemoteRepositoryService.Conan(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteChefTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewChefRemoteRepositoryParams()
	crp.Key = repoKey
	crp.RepoLayoutRef = "simple-default"
	crp.Url = "https://supermarket.chef.io"
	crp.Description = "Chef Repo for jfrog-client-go remote-repository-test"
	crp.AssumedOfflinePeriodSecs = 2345
	crp.StoreArtifactsLocally = &falseValue
	crp.ShareConfiguration = &falseValue
	crp.BypassHeadRequests = &trueValue
	crp.SynchronizeProperties = &trueValue
	crp.IncludesPattern = "**/*"
	crp.ExcludesPattern = "dir1/dir2/dir3/*"
	crp.AllowAnyHostAuth = &trueValue

	err := testsCreateRemoteRepositoryService.Chef(crp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	crp.Description += " - Updated"
	crp.Notes = "Repo been updated"
	crp.AssumedOfflinePeriodSecs = 2406
	crp.EnableCookieManagement = &trueValue
	crp.SocketTimeoutMillis = 1989
	crp.SynchronizeProperties = &falseValue
	crp.BypassHeadRequests = &falseValue
	crp.IncludesPattern = "**/*"
	crp.ExcludesPattern = "dir1/dir2/dir3/dir4/*,dir1/dir2/dir3/dir5/*"
	crp.AllowAnyHostAuth = &falseValue

	err = testsUpdateRemoteRepositoryService.Chef(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remotePuppetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	prp := services.NewPuppetRemoteRepositoryParams()
	prp.Key = repoKey
	prp.RepoLayoutRef = "puppet-default"
	prp.Url = "https://forgeapi.puppetlabs.com/"
	prp.Description = "Puppet Repo for jfrog-client-go remote-repository-test"
	prp.AssumedOfflinePeriodSecs = 999
	prp.StoreArtifactsLocally = &trueValue
	prp.ShareConfiguration = &trueValue
	prp.AssumedOfflinePeriodSecs = 1803
	prp.Offline = &trueValue

	err := testsCreateRemoteRepositoryService.Puppet(prp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	prp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, prp)

	prp.Notes = "Repo been updated"
	prp.AssumedOfflinePeriodSecs = 2202
	prp.EnableCookieManagement = &trueValue
	prp.SocketTimeoutMillis = 1800
	prp.ShareConfiguration = &falseValue
	prp.StoreArtifactsLocally = &falseValue
	prp.ShareConfiguration = &falseValue
	prp.BlockMismatchingMimeTypes = &falseValue
	prp.SynchronizeProperties = &trueValue

	err = testsUpdateRemoteRepositoryService.Puppet(prp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, prp)
}

func remoteComposerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
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
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	crp.Description += " - Updated"
	crp.Notes = "Repo been updated"
	crp.AssumedOfflinePeriodSecs = 2000
	crp.EnableCookieManagement = &trueValue
	crp.SocketTimeoutMillis = 666
	crp.IncludesPattern = "**/*"
	crp.BypassHeadRequests = &falseValue

	err = testsUpdateRemoteRepositoryService.Composer(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteVcsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
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
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
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
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, vrp)
}

func remoteCocoapodsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewCocoapodsRemoteRepositoryParams()
	crp.Key = repoKey
	crp.RepoLayoutRef = "simple-default"
	crp.Url = "https://github.com/"
	crp.Description = "Cocoapods Repo for jfrog-client-go remote-repository-test"
	crp.AssumedOfflinePeriodSecs = 3801
	crp.StoreArtifactsLocally = &trueValue
	crp.ShareConfiguration = &trueValue

	err := testsCreateRemoteRepositoryService.Cocoapods(crp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	crp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, crp)

	crp.Notes = "Repo been updated"
	crp.AssumedOfflinePeriodSecs = 1111
	crp.AssumedOfflinePeriodSecs = 3799
	crp.EnableCookieManagement = &trueValue
	crp.SocketTimeoutMillis = 1818
	crp.ShareConfiguration = &falseValue

	err = testsUpdateRemoteRepositoryService.Cocoapods(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteOpkgTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	orp := services.NewOpkgRemoteRepositoryParams()
	orp.Key = repoKey
	orp.RepoLayoutRef = "simple-default"
	orp.Url = "https://opkg.com/download.git"
	orp.Description = "Opkg Repo for jfrog-client-go remote-repository-test"
	orp.AssumedOfflinePeriodSecs = 1500
	orp.StoreArtifactsLocally = &falseValue
	orp.ShareConfiguration = &trueValue
	orp.ListRemoteFolderItems = &falseValue

	err := testsCreateRemoteRepositoryService.Opkg(orp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, orp)

	orp.Description += " - Updated"
	orp.Notes = "Repo been updated"
	orp.AssumedOfflinePeriodSecs = 2222
	orp.EnableCookieManagement = &trueValue
	orp.SocketTimeoutMillis = 1818
	orp.ShareConfiguration = &falseValue
	orp.ListRemoteFolderItems = &trueValue

	err = testsUpdateRemoteRepositoryService.Opkg(orp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, orp)
}

func remoteCondaTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewCondaRemoteRepositoryParams()
	crp.Key = repoKey
	crp.RepoLayoutRef = "simple-default"
	crp.Url = "https://repo.anaconda.com/pkgs/free"
	crp.Description = "Conda Repo for jfrog-client-go remote-repository-test"
	crp.AssumedOfflinePeriodSecs = 1800
	crp.StoreArtifactsLocally = &falseValue
	crp.ShareConfiguration = &trueValue

	err := testsCreateRemoteRepositoryService.Conda(crp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	crp.Description += " - Updated"
	crp.Notes = "Repo been updated"
	crp.AssumedOfflinePeriodSecs = 2222
	crp.EnableCookieManagement = &trueValue
	crp.SocketTimeoutMillis = 1818
	crp.ShareConfiguration = &falseValue

	err = testsUpdateRemoteRepositoryService.Conda(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteP2Test(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	prp := services.NewP2RemoteRepositoryParams()
	prp.Key = repoKey
	prp.RepoLayoutRef = "simple-default"
	prp.Url = "https://repo.anaconda.com/pkgs/free"
	prp.Description = "P2 Repo for jfrog-client-go remote-repository-test"
	prp.AssumedOfflinePeriodSecs = 999
	prp.StoreArtifactsLocally = &trueValue
	prp.ShareConfiguration = &trueValue

	err := testsCreateRemoteRepositoryService.P2(prp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	prp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, prp)

	prp.Notes = "Repo been updated"
	prp.AssumedOfflinePeriodSecs = 2202
	prp.EnableCookieManagement = &trueValue
	prp.SocketTimeoutMillis = 1800
	prp.ShareConfiguration = &falseValue
	prp.StoreArtifactsLocally = &falseValue
	prp.ShareConfiguration = &falseValue

	err = testsUpdateRemoteRepositoryService.P2(prp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, prp)
}

func remoteGenericTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGenericRemoteRepositoryParams()
	grp.Key = repoKey
	grp.RepoLayoutRef = "simple-default"
	grp.Url = "https://jcenter.bintray.com"
	grp.Description = "Generic Repo for jfrog-client-go remote-repository-test"
	grp.AssumedOfflinePeriodSecs = 2345
	grp.StoreArtifactsLocally = &falseValue
	grp.ShareConfiguration = &falseValue

	err := testsCreateRemoteRepositoryService.Generic(grp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, grp)

	grp.Description += " - Updated"
	grp.Notes = "Repo been updated"
	grp.AssumedOfflinePeriodSecs = 2000
	grp.EnableCookieManagement = &trueValue
	grp.SocketTimeoutMillis = 666

	err = testsUpdateRemoteRepositoryService.Generic(grp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, grp)
}

func getRemoteRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGenericRemoteRepositoryParams()
	grp.Key = repoKey
	grp.RepoLayoutRef = "simple-default"
	grp.Url = "https://jcenter.bintray.com"
	grp.Description = "Generic Repo for jfrog-client-go remote-repository-test"

	err := testsCreateRemoteRepositoryService.Generic(grp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)
	// Get repo details
	data := getRepo(t, repoKey)
	// Validate
	assert.Equal(t, data.Key, repoKey)
	assert.Equal(t, data.Description, grp.Description+" (local file cache)")
	assert.Equal(t, data.Rclass, "remote")
	assert.Equal(t, data.Url, grp.Url)
	assert.Equal(t, data.PackageType, "generic")
}
