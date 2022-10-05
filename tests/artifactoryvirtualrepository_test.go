package tests

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

func TestArtifactoryVirtualRepository(t *testing.T) {
	initRepositoryTest(t)
	t.Run("virtualAlpineTest", virtualAlpineTest)
	t.Run("virtualBowerTest", virtualBowerTest)
	t.Run("virtualChefTest", virtualChefTest)
	t.Run("virtualConanTest", virtualConanTest)
	t.Run("virtualCondaTest", virtualCondaTest)
	t.Run("virtualCranTest", virtualCranTest)
	t.Run("virtualDebianTest", virtualDebianTest)
	t.Run("virtualDockerTest", virtualDockerTest)
	t.Run("virtualGemsTest", virtualGemsTest)
	t.Run("virtualGenericTest", virtualGenericTest)
	t.Run("virtualGitlfsTest", virtualGitlfsTest)
	t.Run("virtualGoTest", virtualGoTest)
	t.Run("virtualGradleTest", virtualGradleTest)
	t.Run("virtualHelmTest", virtualHelmTest)
	t.Run("virtualIvyTest", virtualIvyTest)
	t.Run("virtualMavenTest", virtualMavenTest)
	t.Run("virtualNpmTest", virtualNpmTest)
	t.Run("virtualNugetTest", virtualNugetTest)
	t.Run("virtualP2Test", virtualP2Test)
	t.Run("virtualPuppetTest", virtualPuppetTest)
	t.Run("virtualPypiTest", virtualPypiTest)
	t.Run("virtualRpmTest", virtualRpmTest)
	t.Run("virtualSbtTest", virtualSbtTest)
	t.Run("virtualSwiftTest", virtualSwiftTest)
	t.Run("virtualYumTest", virtualYumTest)
	t.Run("virtualCreateWithParamTest", virtualCreateWithParamTest)
	t.Run("getVirtualRepoDetailsTest", getVirtualRepoDetailsTest)
	t.Run("getAllVirtualRepoDetailsTest", getAllVirtualRepoDetailsTest)
	t.Run("isVirtualRepoExistsTest", isVirtualRepoExistsTest)
}

func setVirtualRepositoryBaseParams(params *services.VirtualRepositoryBaseParams, isUpdate bool) {
	setRepositoryBaseParams(&params.RepositoryBaseParams, isUpdate)
	if !isUpdate {
		params.Repositories = []string{getRtTargetRepoKey()}
		params.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &trueValue
		params.DefaultDeploymentRepo = getRtTargetRepoKey()
	} else {
		params.Repositories = nil
		params.ArtifactoryRequestsCanRetrieveRemoteArtifacts = &falseValue
		params.DefaultDeploymentRepo = ""
	}
}

func setCacheVirtualRepositoryParams(params *services.CommonCacheVirtualRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.VirtualRetrievalCachePeriodSecs = 300
	} else {
		params.VirtualRetrievalCachePeriodSecs = 0
	}
}

func setJavaPackageManagersVirtualRepositoryParams(params *services.CommonJavaVirtualRepositoryParams, isUpdate bool) {
	if !isUpdate {
		params.PomRepositoryReferencesCleanupPolicy = "nothing"
		params.KeyPair = ""

	} else {
		params.PomRepositoryReferencesCleanupPolicy = ""
		params.KeyPair = ""
	}
}

func virtualAlpineTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	avp := services.NewAlpineVirtualRepositoryParams()
	avp.Key = repoKey
	setVirtualRepositoryBaseParams(&avp.VirtualRepositoryBaseParams, false)
	setCacheVirtualRepositoryParams(&avp.CommonCacheVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Alpine(avp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, avp)

	setVirtualRepositoryBaseParams(&avp.VirtualRepositoryBaseParams, true)
	setCacheVirtualRepositoryParams(&avp.CommonCacheVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Alpine(avp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, avp)
}

func virtualBowerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	bvp := services.NewBowerVirtualRepositoryParams()
	bvp.Key = repoKey
	setVirtualRepositoryBaseParams(&bvp.VirtualRepositoryBaseParams, false)
	bvp.ExternalDependenciesEnabled = &trueValue
	bvp.ExternalDependenciesPatterns = []string{"**/*github*/**"}
	bvp.ExternalDependenciesRemoteRepo = ""

	err := testsCreateVirtualRepositoryService.Bower(bvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, bvp)

	setVirtualRepositoryBaseParams(&bvp.VirtualRepositoryBaseParams, true)
	bvp.ExternalDependenciesEnabled = &falseValue
	bvp.ExternalDependenciesPatterns = nil
	bvp.ExternalDependenciesRemoteRepo = ""

	err = testsUpdateVirtualRepositoryService.Bower(bvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, bvp)
}

func virtualChefTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cvp := services.NewChefVirtualRepositoryParams()
	cvp.Key = repoKey
	setVirtualRepositoryBaseParams(&cvp.VirtualRepositoryBaseParams, false)
	setCacheVirtualRepositoryParams(&cvp.CommonCacheVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Chef(cvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cvp)

	setVirtualRepositoryBaseParams(&cvp.VirtualRepositoryBaseParams, true)
	setCacheVirtualRepositoryParams(&cvp.CommonCacheVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Chef(cvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cvp)
}

func virtualConanTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cvp := services.NewConanVirtualRepositoryParams()
	cvp.Key = repoKey
	setVirtualRepositoryBaseParams(&cvp.VirtualRepositoryBaseParams, false)
	setCacheVirtualRepositoryParams(&cvp.CommonCacheVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Conan(cvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cvp)

	setVirtualRepositoryBaseParams(&cvp.VirtualRepositoryBaseParams, true)
	setCacheVirtualRepositoryParams(&cvp.CommonCacheVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Conan(cvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cvp)
}

func virtualCondaTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cvp := services.NewCondaVirtualRepositoryParams()
	cvp.Key = repoKey
	setVirtualRepositoryBaseParams(&cvp.VirtualRepositoryBaseParams, false)
	setCacheVirtualRepositoryParams(&cvp.CommonCacheVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Conda(cvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cvp)

	setVirtualRepositoryBaseParams(&cvp.VirtualRepositoryBaseParams, true)
	setCacheVirtualRepositoryParams(&cvp.CommonCacheVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Conda(cvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cvp)
}

func virtualCranTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cvp := services.NewCranVirtualRepositoryParams()
	cvp.Key = repoKey
	setVirtualRepositoryBaseParams(&cvp.VirtualRepositoryBaseParams, false)
	setCacheVirtualRepositoryParams(&cvp.CommonCacheVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Cran(cvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cvp)

	setVirtualRepositoryBaseParams(&cvp.VirtualRepositoryBaseParams, true)
	setCacheVirtualRepositoryParams(&cvp.CommonCacheVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Cran(cvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cvp)
}

func virtualDebianTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	dvp := services.NewDebianVirtualRepositoryParams()
	dvp.Key = repoKey
	setVirtualRepositoryBaseParams(&dvp.VirtualRepositoryBaseParams, false)
	setCacheVirtualRepositoryParams(&dvp.CommonCacheVirtualRepositoryParams, false)
	dvp.DebianDefaultArchitectures = "amd64, i386"
	dvp.OptionalIndexCompressionFormats = []string{"bz2", "lzma"}

	err := testsCreateVirtualRepositoryService.Debian(dvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, dvp)

	setVirtualRepositoryBaseParams(&dvp.VirtualRepositoryBaseParams, true)
	setCacheVirtualRepositoryParams(&dvp.CommonCacheVirtualRepositoryParams, true)
	dvp.DebianDefaultArchitectures = ""
	dvp.OptionalIndexCompressionFormats = nil

	err = testsUpdateVirtualRepositoryService.Debian(dvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dvp)
}

func virtualDockerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	dvp := services.NewDockerVirtualRepositoryParams()
	dvp.Key = repoKey
	setVirtualRepositoryBaseParams(&dvp.VirtualRepositoryBaseParams, false)
	dvp.ResolveDockerTagsByTimestamp = &trueValue

	err := testsCreateVirtualRepositoryService.Docker(dvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, dvp)

	setVirtualRepositoryBaseParams(&dvp.VirtualRepositoryBaseParams, true)
	dvp.ResolveDockerTagsByTimestamp = &falseValue

	err = testsUpdateVirtualRepositoryService.Docker(dvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dvp)
}

func virtualGemsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGemsVirtualRepositoryParams()
	gvp.Key = repoKey
	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, false)

	err := testsCreateVirtualRepositoryService.Gems(gvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gvp)

	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, true)

	err = testsUpdateVirtualRepositoryService.Gems(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)
}

func virtualGenericTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGenericVirtualRepositoryParams()
	gvp.Key = repoKey
	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, false)

	err := testsCreateVirtualRepositoryService.Generic(gvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gvp)

	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, true)

	err = testsUpdateVirtualRepositoryService.Generic(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)
}

func virtualGitlfsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGitlfsVirtualRepositoryParams()
	gvp.Key = repoKey
	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, false)

	err := testsCreateVirtualRepositoryService.Gitlfs(gvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gvp)

	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, true)

	err = testsUpdateVirtualRepositoryService.Gitlfs(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)
}

func virtualGoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGoVirtualRepositoryParams()
	gvp.Key = repoKey
	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, false)
	gvp.ExternalDependenciesEnabled = &trueValue
	gvp.ExternalDependenciesPatterns = []string{"**/*microsoft*/**", "**/*github*/**"}

	err := testsCreateVirtualRepositoryService.Go(gvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gvp)

	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, true)
	gvp.ExternalDependenciesEnabled = &falseValue
	gvp.ExternalDependenciesPatterns = nil

	err = testsUpdateVirtualRepositoryService.Go(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)
}

func virtualGradleTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGradleVirtualRepositoryParams()
	gvp.Key = repoKey
	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, false)
	setJavaPackageManagersVirtualRepositoryParams(&gvp.CommonJavaVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Gradle(gvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gvp)

	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, true)
	setJavaPackageManagersVirtualRepositoryParams(&gvp.CommonJavaVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Gradle(gvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gvp)
}

func virtualHelmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	hvp := services.NewHelmVirtualRepositoryParams()
	hvp.Key = repoKey
	setVirtualRepositoryBaseParams(&hvp.VirtualRepositoryBaseParams, false)
	setCacheVirtualRepositoryParams(&hvp.CommonCacheVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Helm(hvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, hvp)

	setVirtualRepositoryBaseParams(&hvp.VirtualRepositoryBaseParams, true)
	setCacheVirtualRepositoryParams(&hvp.CommonCacheVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Helm(hvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, hvp)
}

func virtualIvyTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	ivp := services.NewIvyVirtualRepositoryParams()
	ivp.Key = repoKey
	setVirtualRepositoryBaseParams(&ivp.VirtualRepositoryBaseParams, false)
	setJavaPackageManagersVirtualRepositoryParams(&ivp.CommonJavaVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Ivy(ivp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, ivp)

	setVirtualRepositoryBaseParams(&ivp.VirtualRepositoryBaseParams, true)
	setJavaPackageManagersVirtualRepositoryParams(&ivp.CommonJavaVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Ivy(ivp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, ivp)
}

func virtualMavenTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	mvp := services.NewMavenVirtualRepositoryParams()
	mvp.Key = repoKey
	setVirtualRepositoryBaseParams(&mvp.VirtualRepositoryBaseParams, false)
	setJavaPackageManagersVirtualRepositoryParams(&mvp.CommonJavaVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Maven(mvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, mvp)

	setVirtualRepositoryBaseParams(&mvp.VirtualRepositoryBaseParams, true)
	setJavaPackageManagersVirtualRepositoryParams(&mvp.CommonJavaVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Maven(mvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, mvp)
}

func virtualNpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nvp := services.NewNpmVirtualRepositoryParams()
	nvp.Key = repoKey
	setVirtualRepositoryBaseParams(&nvp.VirtualRepositoryBaseParams, false)
	setCacheVirtualRepositoryParams(&nvp.CommonCacheVirtualRepositoryParams, false)
	nvp.ExternalDependenciesEnabled = &trueValue
	nvp.ExternalDependenciesPatterns = []string{"**/*microsoft*/**"}
	nvp.ExternalDependenciesRemoteRepo = ""

	err := testsCreateVirtualRepositoryService.Npm(nvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, nvp)

	setVirtualRepositoryBaseParams(&nvp.VirtualRepositoryBaseParams, true)
	setCacheVirtualRepositoryParams(&nvp.CommonCacheVirtualRepositoryParams, true)
	nvp.ExternalDependenciesEnabled = &falseValue
	nvp.ExternalDependenciesPatterns = nil
	nvp.ExternalDependenciesRemoteRepo = ""

	err = testsUpdateVirtualRepositoryService.Npm(nvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nvp)
}

func virtualNugetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nvp := services.NewNugetVirtualRepositoryParams()
	nvp.Key = repoKey
	setVirtualRepositoryBaseParams(&nvp.VirtualRepositoryBaseParams, false)
	nvp.ForceNugetAuthentication = &trueValue

	err := testsCreateVirtualRepositoryService.Nuget(nvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, nvp)

	setVirtualRepositoryBaseParams(&nvp.VirtualRepositoryBaseParams, true)
	nvp.ForceNugetAuthentication = &falseValue

	err = testsUpdateVirtualRepositoryService.Nuget(nvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nvp)
}

func virtualP2Test(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	pvp := services.NewP2VirtualRepositoryParams()
	pvp.Key = repoKey
	setVirtualRepositoryBaseParams(&pvp.VirtualRepositoryBaseParams, false)

	err := testsCreateVirtualRepositoryService.P2(pvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, pvp)

	setVirtualRepositoryBaseParams(&pvp.VirtualRepositoryBaseParams, true)
	pvp.Repositories = nil

	err = testsUpdateVirtualRepositoryService.P2(pvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, pvp)
}

func virtualPuppetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	pvp := services.NewPuppetVirtualRepositoryParams()
	pvp.Key = repoKey
	setVirtualRepositoryBaseParams(&pvp.VirtualRepositoryBaseParams, false)

	err := testsCreateVirtualRepositoryService.Puppet(pvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, pvp)

	setVirtualRepositoryBaseParams(&pvp.VirtualRepositoryBaseParams, true)

	err = testsUpdateVirtualRepositoryService.Puppet(pvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, pvp)
}

func virtualPypiTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	pvp := services.NewPypiVirtualRepositoryParams()
	pvp.Key = repoKey
	setVirtualRepositoryBaseParams(&pvp.VirtualRepositoryBaseParams, false)

	err := testsCreateVirtualRepositoryService.Pypi(pvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, pvp)

	setVirtualRepositoryBaseParams(&pvp.VirtualRepositoryBaseParams, true)

	err = testsUpdateVirtualRepositoryService.Pypi(pvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, pvp)
}

func virtualRpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	rvp := services.NewRpmVirtualRepositoryParams()
	rvp.Key = repoKey
	setVirtualRepositoryBaseParams(&rvp.VirtualRepositoryBaseParams, false)
	setCacheVirtualRepositoryParams(&rvp.CommonCacheVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Rpm(rvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, rvp)

	setVirtualRepositoryBaseParams(&rvp.VirtualRepositoryBaseParams, true)
	setCacheVirtualRepositoryParams(&rvp.CommonCacheVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Rpm(rvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, rvp)
}

func virtualSbtTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	svp := services.NewSbtVirtualRepositoryParams()
	svp.Key = repoKey
	setVirtualRepositoryBaseParams(&svp.VirtualRepositoryBaseParams, false)
	setJavaPackageManagersVirtualRepositoryParams(&svp.CommonJavaVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Sbt(svp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, svp)

	setVirtualRepositoryBaseParams(&svp.VirtualRepositoryBaseParams, true)
	setJavaPackageManagersVirtualRepositoryParams(&svp.CommonJavaVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Sbt(svp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, svp)
}

func virtualSwiftTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	svp := services.NewSwiftVirtualRepositoryParams()
	svp.Key = repoKey
	setVirtualRepositoryBaseParams(&svp.VirtualRepositoryBaseParams, false)

	err := testsCreateVirtualRepositoryService.Swift(svp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, svp)

	setVirtualRepositoryBaseParams(&svp.VirtualRepositoryBaseParams, true)

	err = testsUpdateVirtualRepositoryService.Swift(svp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, svp)
}

func virtualYumTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	yvp := services.NewYumVirtualRepositoryParams()
	yvp.Key = repoKey
	setVirtualRepositoryBaseParams(&yvp.VirtualRepositoryBaseParams, false)
	setCacheVirtualRepositoryParams(&yvp.CommonCacheVirtualRepositoryParams, false)

	err := testsCreateVirtualRepositoryService.Yum(yvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// "yum" package type is converted to "rpm" by Artifactory, so we have to change it too to pass the validation.
	yvp.PackageType = "rpm"
	validateRepoConfig(t, repoKey, yvp)

	setVirtualRepositoryBaseParams(&yvp.VirtualRepositoryBaseParams, true)
	setCacheVirtualRepositoryParams(&yvp.CommonCacheVirtualRepositoryParams, true)

	err = testsUpdateVirtualRepositoryService.Yum(yvp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, yvp)
}

func virtualCreateWithParamTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	params := services.NewVirtualRepositoryBaseParams()
	params.Key = repoKey
	err := testsRepositoriesService.CreateVirtual(params)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, params)
}

func getVirtualRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGenericVirtualRepositoryParams()
	gvp.Key = repoKey
	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, false)

	err := testsCreateVirtualRepositoryService.Generic(gvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// Get repo details
	data := getRepo(t, repoKey)
	// Validate
	assert.Equal(t, data.Key, repoKey)
	assert.Equal(t, data.Description, gvp.Description)
	assert.Equal(t, data.GetRepoType(), "virtual")
	assert.Empty(t, data.Url)
	assert.Equal(t, data.PackageType, "generic")
}

func getAllVirtualRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gvp := services.NewGenericVirtualRepositoryParams()
	gvp.Key = repoKey
	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, false)

	err := testsCreateVirtualRepositoryService.Generic(gvp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// Get repo details
	data := getAllRepos(t, "virtual", "")
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
}

func isVirtualRepoExistsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()

	// Validate repo doesn't exist
	exists := isRepoExists(t, repoKey)
	assert.False(t, exists)

	// Create Repo
	gvp := services.NewGenericVirtualRepositoryParams()
	gvp.Key = repoKey
	setVirtualRepositoryBaseParams(&gvp.VirtualRepositoryBaseParams, false)
	err := testsCreateVirtualRepositoryService.Generic(gvp)
	assert.NoError(t, err)
	defer deleteRepo(t, repoKey)

	// Validate repo exists
	exists = isRepoExists(t, repoKey)
	assert.True(t, exists)
}
