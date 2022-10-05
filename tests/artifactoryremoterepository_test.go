package tests

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

const ArtifactoryLocalFileCacheSuffix = " (local file cache)"
const MavenCentralUrl = "https://repo.maven.apache.org"

func TestArtifactoryRemoteRepository(t *testing.T) {
	initRepositoryTest(t)
	t.Run("remoteAlpineTest", remoteAlpineTest)
	t.Run("remoteBowerTest", remoteBowerTest)
	t.Run("remoteCargoTest", remoteCargoTest)
	t.Run("remoteChefTest", remoteChefTest)
	t.Run("remoteCocoapodsTest", remoteCocoapodsTest)
	t.Run("remoteComposerTest", remoteComposerTest)
	t.Run("remoteConanTest", remoteConanTest)
	t.Run("remoteCondaTest", remoteCondaTest)
	t.Run("remoteCranTest", remoteCranTest)
	t.Run("remoteDebianTest", remoteDebianTest)
	t.Run("remoteDockerTest", remoteDockerTest)
	t.Run("remoteGemsTest", remoteGemsTest)
	t.Run("remoteGenericTest", remoteGenericTest)
	t.Run("remoteGitlfsTest", remoteGitlfsTest)
	t.Run("remoteGoTest", remoteGoTest)
	t.Run("remoteGradleTest", remoteGradleTest)
	t.Run("remoteHelmTest", remoteHelmTest)
	t.Run("remoteIvyTest", remoteIvyTest)
	t.Run("remoteMavenTest", remoteMavenTest)
	t.Run("remoteNpmTest", remoteNpmTest)
	t.Run("remoteNugetTest", remoteNugetTest)
	t.Run("remoteOkgTest", remoteOpkgTest)
	t.Run("remoteP2Test", remoteP2Test)
	t.Run("remotePuppetTest", remotePuppetTest)
	t.Run("remotePypiTest", remotePypiTest)
	t.Run("remoteRpmTest", remoteRpmTest)
	t.Run("remoteSbtTest", remoteSbtTest)
	t.Run("remoteSwiftTest", remoteSwiftTest)
	t.Run("remoteVcsTest", remoteVcsTest)
	t.Run("remoteYumTest", remoteYumTest)
	t.Run("remoteGenericSmartRemoteTest", remoteGenericSmartRemoteTest)
	t.Run("remoteCreateWithParamTest", remoteCreateWithParamTest)
	t.Run("getRemoteRepoDetailsTest", getRemoteRepoDetailsTest)
	t.Run("getAllRemoteRepoDetailsTest", getAllRemoteRepoDetailsTest)
	t.Run("isRemoteRepoExistsTest", isRemoteRepoExistsTest)
}

func setRemoteRepositoryBaseParams(params *services.RemoteRepositoryBaseParams, isUpdate bool) {
	setRepositoryBaseParams(&params.RepositoryBaseParams, isUpdate)
	setAdditionalRepositoryBaseParams(&params.AdditionalRepositoryBaseParams, isUpdate)
	if !isUpdate {
		params.HardFail = &trueValue
		params.Offline = &trueValue
		params.StoreArtifactsLocally = &trueValue
		params.SocketTimeoutMillis = 500
		params.RetrievalCachePeriodSecs = 1000
		params.MetadataRetrievalTimeoutSecs = 500
		params.MissedRetrievalCachePeriodSecs = 3000
		params.UnusedArtifactsCleanupPeriodHours = 24
		params.AssumedOfflinePeriodSecs = 300
		params.ShareConfiguration = &trueValue
		params.SynchronizeProperties = &trueValue
		params.BlockMismatchingMimeTypes = &trueValue
		params.MismatchingMimeTypesOverrideList = "text/html,text/csv"
		params.AllowAnyHostAuth = &trueValue
		params.EnableCookieManagement = &trueValue
		params.BypassHeadRequests = &trueValue
		params.ClientTlsCertificate = ""
	} else {
		params.HardFail = &falseValue
		params.Offline = &falseValue
		params.StoreArtifactsLocally = &falseValue
		params.SocketTimeoutMillis = 1000
		params.RetrievalCachePeriodSecs = 2000
		params.MetadataRetrievalTimeoutSecs = 1000
		params.MissedRetrievalCachePeriodSecs = 5000
		params.UnusedArtifactsCleanupPeriodHours = 36
		params.AssumedOfflinePeriodSecs = 600
		params.ShareConfiguration = &falseValue
		params.SynchronizeProperties = &falseValue
		params.BlockMismatchingMimeTypes = &falseValue
		params.MismatchingMimeTypesOverrideList = ""
		params.AllowAnyHostAuth = &falseValue
		params.EnableCookieManagement = &falseValue
		params.BypassHeadRequests = &falseValue
		params.ClientTlsCertificate = ""
	}
}

func setVcsRemoteRepositoryParams(params *services.VcsGitRemoteRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.VcsType = "GIT"
		params.VcsGitProvider = "CUSTOM"
		params.VcsGitDownloadUrl = "https://github.com/download.git"
	} else {
		params.VcsType = "GIT"
		params.VcsGitProvider = ""
		params.VcsGitDownloadUrl = ""
	}
}

func setJavaPackageManagersRemoteRepositoryParams(params *services.JavaPackageManagersRemoteRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.MaxUniqueSnapshots = 18
		params.HandleReleases = &trueValue
		params.HandleSnapshots = &trueValue
		params.SuppressPomConsistencyChecks = &trueValue
		params.RemoteRepoChecksumPolicyType = "ignore-and-generate"
		params.FetchJarsEagerly = &trueValue
		params.FetchSourcesEagerly = &trueValue
		params.RejectInvalidJars = &trueValue
	} else {
		params.MaxUniqueSnapshots = 36
		params.HandleReleases = &falseValue
		params.HandleSnapshots = &falseValue
		params.SuppressPomConsistencyChecks = &falseValue
		params.RemoteRepoChecksumPolicyType = "generate-if-absent"
		params.FetchJarsEagerly = &falseValue
		params.FetchSourcesEagerly = &falseValue
		params.RejectInvalidJars = &falseValue
	}
}

func remoteAlpineTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	arp := services.NewAlpineRemoteRepositoryParams()
	arp.Key = repoKey
	arp.Url = "http://dl-cdn.alpinelinux.org/alpine"
	setRemoteRepositoryBaseParams(&arp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Alpine(arp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	arp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, arp)

	setRemoteRepositoryBaseParams(&arp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Alpine(arp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, arp)
}

func remoteBowerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	brp := services.NewBowerRemoteRepositoryParams()
	brp.Key = repoKey
	brp.Url = "https://github.com/"
	setRemoteRepositoryBaseParams(&brp.RemoteRepositoryBaseParams, false)
	setVcsRemoteRepositoryParams(&brp.VcsGitRemoteRepositoryParams, false)
	brp.BowerRegistryUrl = "https://registry.bower.io"

	err := testsCreateRemoteRepositoryService.Bower(brp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	brp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, brp)

	setRemoteRepositoryBaseParams(&brp.RemoteRepositoryBaseParams, true)
	setVcsRemoteRepositoryParams(&brp.VcsGitRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Bower(brp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, brp)
}

func remoteCargoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewCargoRemoteRepositoryParams()
	crp.Key = repoKey
	crp.Url = "https://github.com/rust-lang/crates.io-index"
	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, false)
	crp.GitRegistryUrl = "https://github.com/rust-lang/crates.io-index"
	crp.CargoAnonymousAccess = &trueValue

	err := testsCreateRemoteRepositoryService.Cargo(crp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	crp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)
	crp.CargoAnonymousAccess = &falseValue

	err = testsUpdateRemoteRepositoryService.Cargo(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteChefTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewChefRemoteRepositoryParams()
	crp.Key = repoKey
	crp.Url = "https://supermarket.chef.io"
	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Chef(crp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	crp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Chef(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteCocoapodsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewCocoapodsRemoteRepositoryParams()
	crp.Key = repoKey
	crp.Url = "https://github.com/"
	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, false)
	setVcsRemoteRepositoryParams(&crp.VcsGitRemoteRepositoryParams, false)
	crp.PodsSpecsRepoUrl = "https://github.com/CocoaPods/Specs"

	err := testsCreateRemoteRepositoryService.Cocoapods(crp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	crp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)
	setVcsRemoteRepositoryParams(&crp.VcsGitRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Cocoapods(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteComposerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewComposerRemoteRepositoryParams()
	crp.Key = repoKey
	crp.Url = "https://github.com/"
	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, false)
	setVcsRemoteRepositoryParams(&crp.VcsGitRemoteRepositoryParams, false)
	crp.ComposerRegistryUrl = "https://composer.registry.com/"

	err := testsCreateRemoteRepositoryService.Composer(crp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	crp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)
	setVcsRemoteRepositoryParams(&crp.VcsGitRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Composer(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteConanTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewConanRemoteRepositoryParams()
	crp.Key = repoKey
	crp.Url = "https://conan.bintray.com"
	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Conan(crp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	crp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Conan(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteCondaTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewCondaRemoteRepositoryParams()
	crp.Key = repoKey
	crp.Url = "https://repo.anaconda.com/pkgs/free"
	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Conda(crp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	crp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Conda(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteCranTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewCranRemoteRepositoryParams()
	crp.Key = repoKey
	crp.Url = "https://cran.r-project.org/"
	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Cran(crp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	crp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Cran(crp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, crp)
}

func remoteDebianTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	drp := services.NewDebianRemoteRepositoryParams()
	drp.Key = repoKey
	drp.Url = "http://archive.ubuntu.com/ubuntu/"
	setRemoteRepositoryBaseParams(&drp.RemoteRepositoryBaseParams, false)
	drp.ListRemoteFolderItems = &trueValue

	err := testsCreateRemoteRepositoryService.Debian(drp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	drp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, drp)

	setRemoteRepositoryBaseParams(&drp.RemoteRepositoryBaseParams, true)
	drp.ListRemoteFolderItems = &falseValue

	err = testsUpdateRemoteRepositoryService.Debian(drp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, drp)
}

func remoteDockerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	drp := services.NewDockerRemoteRepositoryParams()
	drp.Key = repoKey
	drp.Url = "https://registry-1.docker.io/"
	setRemoteRepositoryBaseParams(&drp.RemoteRepositoryBaseParams, false)
	drp.ExternalDependenciesEnabled = &trueValue
	drp.ExternalDependenciesPatterns = []string{"image/**"}
	drp.EnableTokenAuthentication = &trueValue
	drp.BlockPullingSchema1 = &trueValue

	err := testsCreateRemoteRepositoryService.Docker(drp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	drp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, drp)

	setRemoteRepositoryBaseParams(&drp.RemoteRepositoryBaseParams, true)
	drp.ExternalDependenciesEnabled = &falseValue
	drp.ExternalDependenciesPatterns = nil
	drp.EnableTokenAuthentication = &falseValue
	drp.BlockPullingSchema1 = &falseValue
	// Docker prerequisite - artifacts must be stored locally in cache
	drp.StoreArtifactsLocally = &trueValue
	err = testsUpdateRemoteRepositoryService.Docker(drp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, drp)
}

func remoteGemsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGemsRemoteRepositoryParams()
	grp.Key = repoKey
	grp.Url = "https://rubygems.org/"
	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, false)
	grp.ListRemoteFolderItems = &trueValue

	err := testsCreateRemoteRepositoryService.Gems(grp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	grp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)
	grp.ListRemoteFolderItems = &falseValue

	err = testsUpdateRemoteRepositoryService.Gems(grp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, grp)
}

func remoteGenericTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGenericRemoteRepositoryParams()
	grp.Key = repoKey
	grp.Url = MavenCentralUrl
	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, false)
	grp.ListRemoteFolderItems = &trueValue

	err := testsCreateRemoteRepositoryService.Generic(grp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	grp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)
	grp.ListRemoteFolderItems = &falseValue

	err = testsUpdateRemoteRepositoryService.Generic(grp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, grp)
}

func remoteGitlfsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGitlfsRemoteRepositoryParams()
	grp.Key = repoKey
	grp.Url = "https://github.com/"
	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Gitlfs(grp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	grp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Gitlfs(grp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, grp)
}

func remoteGoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGoRemoteRepositoryParams()
	grp.Key = repoKey
	grp.Url = "https://gocenter.io/"
	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, false)
	grp.VcsGitProvider = "ARTIFACTORY"

	err := testsCreateRemoteRepositoryService.Go(grp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	grp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)
	grp.VcsGitProvider = "GITHUB"

	err = testsUpdateRemoteRepositoryService.Go(grp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, grp)
}

func remoteGradleTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGradleRemoteRepositoryParams()
	grp.Key = repoKey
	grp.Url = MavenCentralUrl
	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, false)
	setJavaPackageManagersRemoteRepositoryParams(&grp.JavaPackageManagersRemoteRepositoryParams, false)

	err := testsCreateRemoteRepositoryService.Gradle(grp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	grp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)
	setJavaPackageManagersRemoteRepositoryParams(&grp.JavaPackageManagersRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Gradle(grp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, grp)
}

func remoteHelmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	hrp := services.NewHelmRemoteRepositoryParams()
	hrp.Key = repoKey
	hrp.Url = "https://storage.googleapis.com/kubernetes-charts"
	setRemoteRepositoryBaseParams(&hrp.RemoteRepositoryBaseParams, false)
	hrp.ChartsBaseUrl = "charts"

	err := testsCreateRemoteRepositoryService.Helm(hrp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	hrp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, hrp)

	setRemoteRepositoryBaseParams(&hrp.RemoteRepositoryBaseParams, true)
	hrp.ChartsBaseUrl = ""

	err = testsUpdateRemoteRepositoryService.Helm(hrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, hrp)
}

func remoteIvyTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	irp := services.NewIvyRemoteRepositoryParams()
	irp.Key = repoKey
	irp.Url = MavenCentralUrl
	setRemoteRepositoryBaseParams(&irp.RemoteRepositoryBaseParams, false)
	setJavaPackageManagersRemoteRepositoryParams(&irp.JavaPackageManagersRemoteRepositoryParams, false)

	err := testsCreateRemoteRepositoryService.Ivy(irp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	irp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, irp)

	setRemoteRepositoryBaseParams(&irp.RemoteRepositoryBaseParams, true)
	setJavaPackageManagersRemoteRepositoryParams(&irp.JavaPackageManagersRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Ivy(irp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, irp)
}

func remoteMavenTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	mrp := services.NewMavenRemoteRepositoryParams()
	mrp.Key = repoKey
	mrp.Url = MavenCentralUrl
	setRemoteRepositoryBaseParams(&mrp.RemoteRepositoryBaseParams, false)
	setJavaPackageManagersRemoteRepositoryParams(&mrp.JavaPackageManagersRemoteRepositoryParams, false)

	err := testsCreateRemoteRepositoryService.Maven(mrp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	mrp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, mrp)

	setRemoteRepositoryBaseParams(&mrp.RemoteRepositoryBaseParams, true)
	setJavaPackageManagersRemoteRepositoryParams(&mrp.JavaPackageManagersRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Maven(mrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, mrp)
}

func remoteNpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nrp := services.NewNpmRemoteRepositoryParams()
	nrp.Key = repoKey
	nrp.Url = "https://registry.npmjs.org"
	setRemoteRepositoryBaseParams(&nrp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Npm(nrp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	nrp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, nrp)

	setRemoteRepositoryBaseParams(&nrp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Npm(nrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nrp)
}

func remoteNugetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nrp := services.NewNugetRemoteRepositoryParams()
	nrp.Key = repoKey
	nrp.Url = "https://www.nuget.org/"
	setRemoteRepositoryBaseParams(&nrp.RemoteRepositoryBaseParams, false)
	nrp.FeedContextPath = "api/v1"
	nrp.DownloadContextPath = "api/v1/package"
	nrp.V3FeedUrl = "https://api.nuget.org/v3/index.json"
	nrp.ForceNugetAuthentication = &trueValue

	err := testsCreateRemoteRepositoryService.Nuget(nrp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	nrp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, nrp)

	setRemoteRepositoryBaseParams(&nrp.RemoteRepositoryBaseParams, true)
	nrp.FeedContextPath = ""
	nrp.DownloadContextPath = ""
	nrp.V3FeedUrl = ""
	nrp.ForceNugetAuthentication = &trueValue

	err = testsUpdateRemoteRepositoryService.Nuget(nrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nrp)
}

func remoteOpkgTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	orp := services.NewOpkgRemoteRepositoryParams()
	orp.Key = repoKey
	orp.Url = "https://opkg.com/download.git"
	setRemoteRepositoryBaseParams(&orp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Opkg(orp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	orp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, orp)

	setRemoteRepositoryBaseParams(&orp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Opkg(orp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, orp)
}

func remoteP2Test(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	prp := services.NewP2RemoteRepositoryParams()
	prp.Key = repoKey
	prp.Url = "https://repo.anaconda.com/pkgs/free"
	setRemoteRepositoryBaseParams(&prp.RemoteRepositoryBaseParams, false)
	prp.ListRemoteFolderItems = &trueValue

	err := testsCreateRemoteRepositoryService.P2(prp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	prp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, prp)

	setRemoteRepositoryBaseParams(&prp.RemoteRepositoryBaseParams, true)
	prp.ListRemoteFolderItems = &falseValue

	err = testsUpdateRemoteRepositoryService.P2(prp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, prp)
}

func remotePuppetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	prp := services.NewPuppetRemoteRepositoryParams()
	prp.Key = repoKey
	prp.Url = "https://forgeapi.puppetlabs.com/"
	setRemoteRepositoryBaseParams(&prp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Puppet(prp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	prp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, prp)

	setRemoteRepositoryBaseParams(&prp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Puppet(prp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, prp)
}

func remotePypiTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	prp := services.NewPypiRemoteRepositoryParams()
	prp.Key = repoKey
	prp.Url = "https://files.pythonhosted.org"
	setRemoteRepositoryBaseParams(&prp.RemoteRepositoryBaseParams, false)
	prp.PypiRegistryUrl = "https://pypi.org"
	prp.PypiRepositorySuffix = "simple"

	err := testsCreateRemoteRepositoryService.Pypi(prp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	prp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, prp)

	setRemoteRepositoryBaseParams(&prp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Pypi(prp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, prp)
}

func remoteRpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	rrp := services.NewRpmRemoteRepositoryParams()
	rrp.Key = repoKey
	rrp.Url = "http://mirror.centos.org/centos/"
	setRemoteRepositoryBaseParams(&rrp.RemoteRepositoryBaseParams, false)
	rrp.ListRemoteFolderItems = &trueValue

	err := testsCreateRemoteRepositoryService.Rpm(rrp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	rrp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, rrp)

	setRemoteRepositoryBaseParams(&rrp.RemoteRepositoryBaseParams, true)
	rrp.ListRemoteFolderItems = &falseValue

	err = testsUpdateRemoteRepositoryService.Rpm(rrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, rrp)
}

func remoteSbtTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	srp := services.NewSbtRemoteRepositoryParams()
	srp.Key = repoKey
	srp.Url = MavenCentralUrl
	setRemoteRepositoryBaseParams(&srp.RemoteRepositoryBaseParams, false)
	setJavaPackageManagersRemoteRepositoryParams(&srp.JavaPackageManagersRemoteRepositoryParams, false)

	err := testsCreateRemoteRepositoryService.Sbt(srp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	srp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, srp)

	setRemoteRepositoryBaseParams(&srp.RemoteRepositoryBaseParams, true)
	setJavaPackageManagersRemoteRepositoryParams(&srp.JavaPackageManagersRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Sbt(srp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, srp)
}

func remoteSwiftTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	srp := services.NewSwiftRemoteRepositoryParams()
	srp.Key = repoKey
	srp.Url = "https://github.com"
	setRemoteRepositoryBaseParams(&srp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Swift(srp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	srp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, srp)

	setRemoteRepositoryBaseParams(&srp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Swift(srp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, srp)
}

func remoteVcsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	vrp := services.NewVcsRemoteRepositoryParams()
	vrp.Key = repoKey
	vrp.Url = "https://github.com/"
	setRemoteRepositoryBaseParams(&vrp.RemoteRepositoryBaseParams, false)
	setVcsRemoteRepositoryParams(&vrp.VcsGitRemoteRepositoryParams, false)
	vrp.MaxUniqueSnapshots = 25

	err := testsCreateRemoteRepositoryService.Vcs(vrp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	vrp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, vrp)

	setRemoteRepositoryBaseParams(&vrp.RemoteRepositoryBaseParams, true)
	setVcsRemoteRepositoryParams(&vrp.VcsGitRemoteRepositoryParams, true)
	vrp.MaxUniqueSnapshots = 50

	err = testsUpdateRemoteRepositoryService.Vcs(vrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, vrp)
}

func remoteYumTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	yrp := services.NewYumRemoteRepositoryParams()
	yrp.Key = repoKey
	yrp.Url = "http://mirror.centos.org/centos/"
	setRemoteRepositoryBaseParams(&yrp.RemoteRepositoryBaseParams, false)
	yrp.ListRemoteFolderItems = &trueValue

	err := testsCreateRemoteRepositoryService.Yum(yrp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// "yum" package type is converted to "rpm" by Artifactory, so we have to change it too to pass the validation.
	yrp.PackageType = "rpm"
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	yrp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKey, yrp)

	setRemoteRepositoryBaseParams(&yrp.RemoteRepositoryBaseParams, true)
	yrp.ListRemoteFolderItems = &falseValue

	err = testsUpdateRemoteRepositoryService.Yum(yrp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, yrp)
}

func remoteGenericSmartRemoteTest(t *testing.T) {
	repoKeyLocal := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKeyLocal
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Generic(glp)
	assert.NoError(t, err, "Failed to create "+repoKeyLocal)
	defer deleteRepo(t, repoKeyLocal)
	validateRepoConfig(t, repoKeyLocal, glp)

	UserParams := getTestUserParams(false, "")
	UserParams.UserDetails.Admin = &trueValue
	err = testUserService.CreateUser(UserParams)
	defer deleteUserAndAssert(t, UserParams.UserDetails.Name)
	assert.NoError(t, err)

	repoKeyRemote := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGenericRemoteRepositoryParams()
	grp.Key = repoKeyRemote
	grp.Url = testsCreateRemoteRepositoryService.ArtDetails.GetUrl() + glp.Key
	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, false)
	setAdditionalRepositoryBaseParams(&grp.AdditionalRepositoryBaseParams, false)
	grp.Username = UserParams.UserDetails.Name
	grp.Password = UserParams.UserDetails.Password
	grp.Proxy = ""
	grp.LocalAddress = ""
	grp.ContentSynchronisation = &services.ContentSynchronisation{
		Enabled: &trueValue,
		Statistics: &services.ContentSynchronisationStatistics{
			Enabled: &trueValue,
		},
		Properties: &services.ContentSynchronisationProperties{
			Enabled: &trueValue,
		},
		Source: &services.ContentSynchronisationSource{
			OriginAbsenceDetection: &trueValue,
		},
	}
	grp.ListRemoteFolderItems = &trueValue

	err = testsCreateRemoteRepositoryService.Generic(grp)
	assert.NoError(t, err, "Failed to create "+repoKeyRemote)
	defer deleteRepo(t, repoKeyRemote)
	// The local file cache suffix is added by Artifactory, so we add it here to pass the validation
	grp.Description += ArtifactoryLocalFileCacheSuffix
	validateRepoConfig(t, repoKeyRemote, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)
	setAdditionalRepositoryBaseParams(&grp.AdditionalRepositoryBaseParams, true)
	grp.Username = UserParams.UserDetails.Name
	grp.Password = UserParams.UserDetails.Password
	grp.Proxy = ""
	grp.LocalAddress = ""
	grp.ContentSynchronisation = &services.ContentSynchronisation{
		Enabled: &falseValue,
		Statistics: &services.ContentSynchronisationStatistics{
			Enabled: &falseValue,
		},
		Properties: &services.ContentSynchronisationProperties{
			Enabled: &falseValue,
		},
		Source: &services.ContentSynchronisationSource{
			OriginAbsenceDetection: &falseValue,
		},
	}
	grp.ListRemoteFolderItems = &falseValue

	err = testsUpdateRemoteRepositoryService.Generic(grp)
	assert.NoError(t, err, "Failed to update "+repoKeyRemote)
	validateRepoConfig(t, repoKeyRemote, grp)
}

func remoteCreateWithParamTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	params := services.NewRemoteRepositoryBaseParams()
	params.Key = repoKey
	params.Url = "https://github.com/"
	err := testsRepositoriesService.CreateRemote(params)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, params)
}

func getRemoteRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGenericRemoteRepositoryParams()
	grp.Key = repoKey
	grp.Url = MavenCentralUrl
	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Generic(grp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// Get repo details
	data := getRepo(t, repoKey)
	// Validate
	assert.Equal(t, data.Key, repoKey)
	assert.Equal(t, data.Description, grp.Description+" (local file cache)")
	assert.Equal(t, data.GetRepoType(), "remote")
	assert.Equal(t, data.Url, grp.Url)
	assert.Equal(t, data.PackageType, "generic")
}

func getAllRemoteRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGenericRemoteRepositoryParams()
	grp.Key = repoKey
	grp.Url = MavenCentralUrl
	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Generic(grp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// Get repo details
	data := getAllRepos(t, "remote", "")
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
	assert.Equal(t, grp.Description+" (local file cache)", repo.Description)
	assert.Equal(t, "Generic", repo.PackageType)
	assert.Equal(t, grp.Url, repo.Url)
}

func isRemoteRepoExistsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()

	// Validate repo doesn't exist
	exists := isRepoExists(t, repoKey)
	assert.False(t, exists)

	// Create Repo
	grp := services.NewGenericRemoteRepositoryParams()
	grp.Key = repoKey
	grp.Url = MavenCentralUrl
	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, false)
	err := testsCreateRemoteRepositoryService.Generic(grp)
	assert.NoError(t, err)
	defer deleteRepo(t, repoKey)

	// Validate repo exists
	exists = isRepoExists(t, repoKey)
	assert.True(t, exists)
}
