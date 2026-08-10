//go:build itest

package tests

import (
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

const MavenCentralUrl = "https://repo.maven.apache.org"

func TestArtifactoryRemoteRepository(t *testing.T) {
	initRepositoryTest(t)
	t.Run("remoteAlpineTest", remoteAlpineTest)
	t.Run("remoteAnsibleTest", remoteAnsibleTest)
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
	t.Run("remoteTerraformTest", remoteTerraformTest)
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
		// Original repo params assigned on creation
		params.HardFail = utils.Pointer(true)
		params.Offline = utils.Pointer(true)
		params.StoreArtifactsLocally = utils.Pointer(true)
		params.SocketTimeoutMillis = utils.Pointer(500)
		params.RetrievalCachePeriodSecs = utils.Pointer(1000)
		params.MetadataRetrievalTimeoutSecs = utils.Pointer(500)
		params.MissedRetrievalCachePeriodSecs = utils.Pointer(3000)
		params.UnusedArtifactsCleanupPeriodHours = utils.Pointer(24)
		params.AssumedOfflinePeriodSecs = utils.Pointer(300)
		params.ShareConfiguration = utils.Pointer(true)
		params.SynchronizeProperties = utils.Pointer(true)
		params.BlockMismatchingMimeTypes = utils.Pointer(true)
		params.MismatchingMimeTypesOverrideList = "text/html,text/csv"
		params.AllowAnyHostAuth = utils.Pointer(true)
		params.EnableCookieManagement = utils.Pointer(true)
		params.BypassHeadRequests = utils.Pointer(true)
		params.ClientTlsCertificate = ""
	} else {
		// Repo params assigned on update
		params.HardFail = utils.Pointer(false)
		params.Offline = utils.Pointer(false)
		params.StoreArtifactsLocally = utils.Pointer(false)
		params.SocketTimeoutMillis = utils.Pointer(1000)
		params.RetrievalCachePeriodSecs = utils.Pointer(2000)
		params.MetadataRetrievalTimeoutSecs = utils.Pointer(1000)
		params.MissedRetrievalCachePeriodSecs = utils.Pointer(0)
		params.UnusedArtifactsCleanupPeriodHours = utils.Pointer(36)
		params.AssumedOfflinePeriodSecs = utils.Pointer(600)
		params.ShareConfiguration = utils.Pointer(false)
		params.SynchronizeProperties = utils.Pointer(false)
		params.BlockMismatchingMimeTypes = utils.Pointer(false)
		params.MismatchingMimeTypesOverrideList = ""
		params.AllowAnyHostAuth = utils.Pointer(false)
		params.EnableCookieManagement = utils.Pointer(false)
		params.BypassHeadRequests = utils.Pointer(false)
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
		params.MaxUniqueSnapshots = utils.Pointer(18)
		params.HandleReleases = utils.Pointer(true)
		params.HandleSnapshots = utils.Pointer(true)
		params.SuppressPomConsistencyChecks = utils.Pointer(true)
		params.RemoteRepoChecksumPolicyType = "ignore-and-generate"
		params.FetchJarsEagerly = utils.Pointer(true)
		params.FetchSourcesEagerly = utils.Pointer(true)
		params.RejectInvalidJars = utils.Pointer(true)
	} else {
		params.MaxUniqueSnapshots = utils.Pointer(36)
		params.HandleReleases = utils.Pointer(false)
		params.HandleSnapshots = utils.Pointer(false)
		params.SuppressPomConsistencyChecks = utils.Pointer(false)
		params.RemoteRepoChecksumPolicyType = "generate-if-absent"
		params.FetchJarsEagerly = utils.Pointer(false)
		params.FetchSourcesEagerly = utils.Pointer(false)
		params.RejectInvalidJars = utils.Pointer(false)
	}
}

func remoteAlpineTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	arp := services.NewAlpineRemoteRepositoryParams()
	arp.Key = repoKey
	arp.Url = "https://dl-cdn.alpinelinux.org/alpine"
	setRemoteRepositoryBaseParams(&arp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Alpine(arp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, arp)

	setRemoteRepositoryBaseParams(&arp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Alpine(arp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, arp)
	}
}

func remoteAnsibleTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	arp := services.NewAnsibleRemoteRepositoryParams()
	arp.Key = repoKey
	arp.Url = "https://galaxy.ansible.com"
	setRemoteRepositoryBaseParams(&arp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Ansible(arp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, arp)

	setRemoteRepositoryBaseParams(&arp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Ansible(arp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, arp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, brp)

	setRemoteRepositoryBaseParams(&brp.RemoteRepositoryBaseParams, true)
	setVcsRemoteRepositoryParams(&brp.VcsGitRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Bower(brp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, brp)
	}
}

func remoteCargoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	crp := services.NewCargoRemoteRepositoryParams()
	crp.Key = repoKey
	crp.Url = "https://github.com/rust-lang/crates.io-index"
	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, false)
	crp.GitRegistryUrl = "https://github.com/rust-lang/crates.io-index"
	crp.CargoAnonymousAccess = utils.Pointer(true)

	err := testsCreateRemoteRepositoryService.Cargo(crp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)
	crp.CargoAnonymousAccess = utils.Pointer(false)

	err = testsUpdateRemoteRepositoryService.Cargo(crp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, crp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Chef(crp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, crp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)
	setVcsRemoteRepositoryParams(&crp.VcsGitRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Cocoapods(crp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, crp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)
	setVcsRemoteRepositoryParams(&crp.VcsGitRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Composer(crp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, crp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Conan(crp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, crp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Conda(crp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, crp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, crp)

	setRemoteRepositoryBaseParams(&crp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Cran(crp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, crp)
	}
}

func remoteDebianTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	drp := services.NewDebianRemoteRepositoryParams()
	drp.Key = repoKey
	drp.Url = "https://archive.ubuntu.com/ubuntu/"
	setRemoteRepositoryBaseParams(&drp.RemoteRepositoryBaseParams, false)
	drp.ListRemoteFolderItems = utils.Pointer(true)

	err := testsCreateRemoteRepositoryService.Debian(drp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, drp)

	setRemoteRepositoryBaseParams(&drp.RemoteRepositoryBaseParams, true)
	drp.ListRemoteFolderItems = utils.Pointer(false)

	err = testsUpdateRemoteRepositoryService.Debian(drp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, drp)
	}
}

func remoteDockerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	drp := services.NewDockerRemoteRepositoryParams()
	drp.Key = repoKey
	drp.Url = "https://registry-1.docker.io/"
	setRemoteRepositoryBaseParams(&drp.RemoteRepositoryBaseParams, false)
	drp.ExternalDependenciesEnabled = utils.Pointer(true)
	drp.ExternalDependenciesPatterns = []string{"image/**"}
	drp.EnableTokenAuthentication = utils.Pointer(true)
	drp.BlockPushingSchema1 = utils.Pointer(true)

	err := testsCreateRemoteRepositoryService.Docker(drp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, drp)

	setRemoteRepositoryBaseParams(&drp.RemoteRepositoryBaseParams, true)
	drp.ExternalDependenciesEnabled = utils.Pointer(false)
	drp.ExternalDependenciesPatterns = nil
	drp.EnableTokenAuthentication = utils.Pointer(false)
	drp.BlockPushingSchema1 = utils.Pointer(false)
	// Docker prerequisite - artifacts must be stored locally in cache
	drp.StoreArtifactsLocally = utils.Pointer(true)
	err = testsUpdateRemoteRepositoryService.Docker(drp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, drp)
	}
}

func remoteGemsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGemsRemoteRepositoryParams()
	grp.Key = repoKey
	grp.Url = "https://rubygems.org/"
	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, false)
	grp.ListRemoteFolderItems = utils.Pointer(true)

	err := testsCreateRemoteRepositoryService.Gems(grp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)
	grp.ListRemoteFolderItems = utils.Pointer(false)

	err = testsUpdateRemoteRepositoryService.Gems(grp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, grp)
	}
}

func remoteGenericTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	grp := services.NewGenericRemoteRepositoryParams()
	grp.Key = repoKey
	grp.Url = MavenCentralUrl
	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, false)
	grp.ListRemoteFolderItems = utils.Pointer(true)

	err := testsCreateRemoteRepositoryService.Generic(grp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)
	grp.ListRemoteFolderItems = utils.Pointer(false)

	err = testsUpdateRemoteRepositoryService.Generic(grp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, grp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Gitlfs(grp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, grp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)
	grp.VcsGitProvider = "GITHUB"

	err = testsUpdateRemoteRepositoryService.Go(grp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, grp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)
	setJavaPackageManagersRemoteRepositoryParams(&grp.JavaPackageManagersRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Gradle(grp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, grp)
	}
}

func remoteHelmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	hrp := services.NewHelmRemoteRepositoryParams()
	hrp.Key = repoKey
	hrp.Url = "https://storage.googleapis.com/kubernetes-charts"
	setRemoteRepositoryBaseParams(&hrp.RemoteRepositoryBaseParams, false)
	hrp.ChartsBaseUrl = "charts"
	hrp.ExternalDependenciesEnabled = true
	hrp.ExternalDependenciesPatterns = []string{"https://github.com/**"}

	err := testsCreateRemoteRepositoryService.Helm(hrp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, hrp)

	setRemoteRepositoryBaseParams(&hrp.RemoteRepositoryBaseParams, true)
	hrp.ChartsBaseUrl = ""
	hrp.ExternalDependenciesEnabled = false
	hrp.ExternalDependenciesPatterns = []string{}

	err = testsUpdateRemoteRepositoryService.Helm(hrp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, hrp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, irp)

	setRemoteRepositoryBaseParams(&irp.RemoteRepositoryBaseParams, true)
	setJavaPackageManagersRemoteRepositoryParams(&irp.JavaPackageManagersRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Ivy(irp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, irp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, mrp)

	setRemoteRepositoryBaseParams(&mrp.RemoteRepositoryBaseParams, true)
	setJavaPackageManagersRemoteRepositoryParams(&mrp.JavaPackageManagersRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Maven(mrp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, mrp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, nrp)

	setRemoteRepositoryBaseParams(&nrp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Npm(nrp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, nrp)
	}
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
	nrp.ForceNugetAuthentication = utils.Pointer(true)
	nrp.SymbolServerUrl = "https://community.chocolatey.org"

	err := testsCreateRemoteRepositoryService.Nuget(nrp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, nrp)

	setRemoteRepositoryBaseParams(&nrp.RemoteRepositoryBaseParams, true)
	nrp.FeedContextPath = ""
	nrp.DownloadContextPath = ""
	nrp.V3FeedUrl = ""
	nrp.ForceNugetAuthentication = utils.Pointer(true)
	nrp.SymbolServerUrl = "https://community.chocolatey.org"

	err = testsUpdateRemoteRepositoryService.Nuget(nrp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, nrp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, orp)

	setRemoteRepositoryBaseParams(&orp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Opkg(orp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, orp)
	}
}

func remoteP2Test(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	prp := services.NewP2RemoteRepositoryParams()
	prp.Key = repoKey
	prp.Url = "https://repo.anaconda.com/pkgs/free"
	setRemoteRepositoryBaseParams(&prp.RemoteRepositoryBaseParams, false)
	prp.ListRemoteFolderItems = utils.Pointer(true)

	err := testsCreateRemoteRepositoryService.P2(prp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, prp)

	setRemoteRepositoryBaseParams(&prp.RemoteRepositoryBaseParams, true)
	prp.ListRemoteFolderItems = utils.Pointer(false)

	err = testsUpdateRemoteRepositoryService.P2(prp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, prp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, prp)

	setRemoteRepositoryBaseParams(&prp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Puppet(prp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, prp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, prp)

	setRemoteRepositoryBaseParams(&prp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Pypi(prp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, prp)
	}
}

func remoteRpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	rrp := services.NewRpmRemoteRepositoryParams()
	rrp.Key = repoKey
	rrp.Url = "https://mirror.centos.org/centos/"
	setRemoteRepositoryBaseParams(&rrp.RemoteRepositoryBaseParams, false)
	rrp.ListRemoteFolderItems = utils.Pointer(true)

	err := testsCreateRemoteRepositoryService.Rpm(rrp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, rrp)

	setRemoteRepositoryBaseParams(&rrp.RemoteRepositoryBaseParams, true)
	rrp.ListRemoteFolderItems = utils.Pointer(false)

	err = testsUpdateRemoteRepositoryService.Rpm(rrp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, rrp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, srp)

	setRemoteRepositoryBaseParams(&srp.RemoteRepositoryBaseParams, true)
	setJavaPackageManagersRemoteRepositoryParams(&srp.JavaPackageManagersRemoteRepositoryParams, true)

	err = testsUpdateRemoteRepositoryService.Sbt(srp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, srp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, srp)

	setRemoteRepositoryBaseParams(&srp.RemoteRepositoryBaseParams, true)

	err = testsUpdateRemoteRepositoryService.Swift(srp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, srp)
	}
}

func remoteTerraformTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	trp := services.NewTerraformRemoteRepositoryParams()
	trp.Key = repoKey
	trp.Url = "https://github.com"
	setRemoteRepositoryBaseParams(&trp.RemoteRepositoryBaseParams, false)

	err := testsCreateRemoteRepositoryService.Terraform(trp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, trp)

	setRemoteRepositoryBaseParams(&trp.RemoteRepositoryBaseParams, true)
	// Due to a bug on Artifactory side that prevents the update of "bypassHeadRequests" field to false on terraform we leave it unchanged.
	trp.BypassHeadRequests = utils.Pointer(true)
	// In terraform - Artifacts must be stored locally
	trp.StoreArtifactsLocally = utils.Pointer(true)
	err = testsUpdateRemoteRepositoryService.Terraform(trp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, trp)
	}
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
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, vrp)

	setRemoteRepositoryBaseParams(&vrp.RemoteRepositoryBaseParams, true)
	setVcsRemoteRepositoryParams(&vrp.VcsGitRemoteRepositoryParams, true)
	vrp.MaxUniqueSnapshots = 50

	err = testsUpdateRemoteRepositoryService.Vcs(vrp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, vrp)
	}
}

func remoteYumTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	yrp := services.NewYumRemoteRepositoryParams()
	yrp.Key = repoKey
	yrp.Url = "https://mirror.centos.org/centos/"
	setRemoteRepositoryBaseParams(&yrp.RemoteRepositoryBaseParams, false)
	yrp.ListRemoteFolderItems = utils.Pointer(true)

	err := testsCreateRemoteRepositoryService.Yum(yrp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	// "yum" package type is converted to "rpm" by Artifactory, so we have to change it too to pass the validation.
	yrp.PackageType = "rpm"
	validateRepoConfig(t, repoKey, yrp)

	setRemoteRepositoryBaseParams(&yrp.RemoteRepositoryBaseParams, true)
	yrp.ListRemoteFolderItems = utils.Pointer(false)

	err = testsUpdateRemoteRepositoryService.Yum(yrp)
	if assert.NoError(t, err, "Failed to update "+repoKey) {
		validateRepoConfig(t, repoKey, yrp)
	}
}

func remoteGenericSmartRemoteTest(t *testing.T) {
	repoKeyLocal := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKeyLocal
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Generic(glp)
	if !assert.NoError(t, err, "Failed to create "+repoKeyLocal) {
		return
	}
	deleteRepoOnTestDone(t, repoKeyLocal)
	validateRepoConfig(t, repoKeyLocal, glp)

	UserParams := getTestUserParams(false, "")
	UserParams.UserDetails.Admin = utils.Pointer(true)
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
		Enabled: utils.Pointer(true),
		Statistics: &services.ContentSynchronisationStatistics{
			Enabled: utils.Pointer(true),
		},
		Properties: &services.ContentSynchronisationProperties{
			Enabled: utils.Pointer(true),
		},
		Source: &services.ContentSynchronisationSource{
			OriginAbsenceDetection: utils.Pointer(true),
		},
	}
	grp.ListRemoteFolderItems = utils.Pointer(true)

	err = testsCreateRemoteRepositoryService.Generic(grp)
	if !assert.NoError(t, err, "Failed to create "+repoKeyRemote) {
		return
	}
	deleteRepoOnTestDone(t, repoKeyRemote)
	validateRepoConfig(t, repoKeyRemote, grp)

	setRemoteRepositoryBaseParams(&grp.RemoteRepositoryBaseParams, true)
	setAdditionalRepositoryBaseParams(&grp.AdditionalRepositoryBaseParams, true)
	grp.Username = UserParams.UserDetails.Name
	grp.Password = UserParams.UserDetails.Password
	grp.Proxy = ""
	grp.LocalAddress = ""
	grp.ContentSynchronisation = &services.ContentSynchronisation{
		Enabled: utils.Pointer(false),
		Statistics: &services.ContentSynchronisationStatistics{
			Enabled: utils.Pointer(false),
		},
		Properties: &services.ContentSynchronisationProperties{
			Enabled: utils.Pointer(false),
		},
		Source: &services.ContentSynchronisationSource{
			OriginAbsenceDetection: utils.Pointer(false),
		},
	}
	grp.ListRemoteFolderItems = utils.Pointer(false)

	err = testsUpdateRemoteRepositoryService.Generic(grp)
	if assert.NoError(t, err, "Failed to update "+repoKeyRemote) {
		validateRepoConfig(t, repoKeyRemote, grp)
	}
}

func remoteCreateWithParamTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	params := services.NewRemoteRepositoryBaseParams()
	params.Key = repoKey
	params.Url = "https://github.com/"
	err := testsRepositoriesService.Create(params, params.Key)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
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
	deleteRepoOnTestDone(t, repoKey)
	// Get repo details
	data := getRepo(t, repoKey)
	// Validate
	assert.Equal(t, repoKey, data.Key)
	assert.Equal(t, grp.Description, data.Description)
	assert.Equal(t, "remote", data.GetRepoType())
	assert.Equal(t, grp.Url, data.Url)
	assert.Equal(t, "generic", data.PackageType)
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
	deleteRepoOnTestDone(t, repoKey)
	// Get repo details
	data := getAllRepos(t, "remote")
	assert.NotNil(t, data)
	repo := &services.RepositoryDetails{}
	for _, v := range *data {
		if v.Key == repoKey {
			rRepo := v
			repo = &rRepo
			break
		}
	}
	// Validate
	assert.NotNil(t, repo, "Repo "+repoKey+" not found")
	repo.Description = strings.TrimSuffix(repo.Description, " (local file cache)")
	assert.Equal(t, grp.Description, repo.Description)
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
	deleteRepoOnTestDone(t, repoKey)

	// Validate repo exists
	exists = isRepoExists(t, repoKey)
	assert.True(t, exists)
}
