package tests

import (
	"github.com/jfrog/gofrog/version"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

func TestArtifactoryFederatedRepository(t *testing.T) {
	initRepositoryTest(t)
	rtVersion, err := GetRtDetails().GetVersion()
	if err != nil {
		t.Error(err)
	}
	if !version.NewVersion(rtVersion).AtLeast("7.18.3") {
		t.Skip("Skipping artifactory test. Federated repositories are only supported by Artifactory 7.18.3 or higher.")
	}
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
	t.Run("federatedSwiftTest", federatedSwiftTest)
	t.Run("federatedVagrantTest", federatedVagrantTest)
	t.Run("federatedYumTest", federatedYumTest)
	t.Run("federatedCreateWithParamTest", federatedCreateWithParamTest)
	t.Run("getFederatedRepoDetailsTest", getFederatedRepoDetailsTest)
	t.Run("getAllFederatedRepoDetailsTest", getAllFederatedRepoDetailsTest)
}

func setFederatedRepositoryBaseParams(params *services.FederatedRepositoryBaseParams, isUpdate bool) {
	setRepositoryBaseParams(&params.RepositoryBaseParams, isUpdate)
	setAdditionalRepositoryBaseParams(&params.AdditionalRepositoryBaseParams, isUpdate)
	memberUrl := testsCreateFederatedRepositoryService.ArtDetails.GetUrl() + params.Key
	if !isUpdate {
		params.ArchiveBrowsingEnabled = &trueValue
		params.Members = []services.FederatedRepositoryMember{
			{Url: memberUrl, Enabled: &trueValue},
		}
	} else {
		params.ArchiveBrowsingEnabled = &falseValue
		params.Members = []services.FederatedRepositoryMember{
			{Url: memberUrl, Enabled: &falseValue},
		}
	}
}

func federatedAlpineTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	afp := services.NewAlpineFederatedRepositoryParams()
	afp.Key = repoKey
	setFederatedRepositoryBaseParams(&afp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Alpine(afp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, afp)

	setFederatedRepositoryBaseParams(&afp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Alpine(afp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, afp)
}

func federatedBowerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	bfp := services.NewBowerFederatedRepositoryParams()
	bfp.Key = repoKey
	setFederatedRepositoryBaseParams(&bfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Bower(bfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, bfp)

	setFederatedRepositoryBaseParams(&bfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Bower(bfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, bfp)
}

func federatedCargoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cfp := services.NewCargoFederatedRepositoryParams()
	cfp.Key = repoKey
	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, false)
	setCargoRepositoryParams(&cfp.CargoRepositoryParams, false)

	err := testsCreateFederatedRepositoryService.Cargo(cfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cfp)

	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, true)
	setCargoRepositoryParams(&cfp.CargoRepositoryParams, true)

	err = testsUpdateFederatedRepositoryService.Cargo(cfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cfp)
}

func federatedChefTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cfp := services.NewChefFederatedRepositoryParams()
	cfp.Key = repoKey
	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Chef(cfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cfp)

	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Chef(cfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cfp)
}

func federatedCocoapodsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cfp := services.NewCocoapodsFederatedRepositoryParams()
	cfp.Key = repoKey
	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Cocoapods(cfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cfp)

	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Cocoapods(cfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cfp)
}

func federatedComposerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cfp := services.NewComposerFederatedRepositoryParams()
	cfp.Key = repoKey
	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Composer(cfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cfp)

	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Composer(cfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cfp)
}

func federatedConanTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cfp := services.NewConanFederatedRepositoryParams()
	cfp.Key = repoKey
	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Conan(cfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cfp)

	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Conan(cfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cfp)
}

func federatedCondaTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cfp := services.NewCondaFederatedRepositoryParams()
	cfp.Key = repoKey
	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Conda(cfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cfp)

	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Conda(cfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cfp)
}

func federatedCranTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	cfp := services.NewCranFederatedRepositoryParams()
	cfp.Key = repoKey
	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Cran(cfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, cfp)

	setFederatedRepositoryBaseParams(&cfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Cran(cfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, cfp)
}

func federatedDebianTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	dfp := services.NewDebianFederatedRepositoryParams()
	dfp.Key = repoKey
	setFederatedRepositoryBaseParams(&dfp.FederatedRepositoryBaseParams, false)
	setDebianRepositoryParams(&dfp.DebianRepositoryParams, false)

	err := testsCreateFederatedRepositoryService.Debian(dfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, dfp)

	setFederatedRepositoryBaseParams(&dfp.FederatedRepositoryBaseParams, true)
	setDebianRepositoryParams(&dfp.DebianRepositoryParams, true)

	err = testsUpdateFederatedRepositoryService.Debian(dfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dfp)
}

func federatedDockerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	dfp := services.NewDockerFederatedRepositoryParams()
	dfp.Key = repoKey
	setFederatedRepositoryBaseParams(&dfp.FederatedRepositoryBaseParams, false)
	setDockerRepositoryParams(&dfp.DockerRepositoryParams, false)

	err := testsCreateFederatedRepositoryService.Docker(dfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, dfp)

	setFederatedRepositoryBaseParams(&dfp.FederatedRepositoryBaseParams, true)
	setDockerRepositoryParams(&dfp.DockerRepositoryParams, true)

	err = testsUpdateFederatedRepositoryService.Docker(dfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dfp)
}

func federatedGemsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gfp := services.NewGemsFederatedRepositoryParams()
	gfp.Key = repoKey
	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Gems(gfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gfp)

	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Gems(gfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gfp)
}

func federatedGenericTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gfp := services.NewGenericFederatedRepositoryParams()
	gfp.Key = repoKey
	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Generic(gfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gfp)

	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Generic(gfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gfp)
}

func federatedGitlfsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gfp := services.NewGitlfsFederatedRepositoryParams()
	gfp.Key = repoKey
	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Gitlfs(gfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gfp)

	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Gitlfs(gfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gfp)
}

func federatedGoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gfp := services.NewGoFederatedRepositoryParams()
	gfp.Key = repoKey
	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Go(gfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gfp)

	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Go(gfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gfp)
}

func federatedGradleTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gfp := services.NewGradleFederatedRepositoryParams()
	gfp.Key = repoKey
	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, false)
	setJavaPackageManagersRepositoryParams(&gfp.JavaPackageManagersRepositoryParams, false)

	err := testsCreateFederatedRepositoryService.Gradle(gfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, gfp)

	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, true)
	setJavaPackageManagersRepositoryParams(&gfp.JavaPackageManagersRepositoryParams, true)

	err = testsUpdateFederatedRepositoryService.Gradle(gfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, gfp)
}

func federatedHelmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	hfp := services.NewHelmFederatedRepositoryParams()
	hfp.Key = repoKey
	setFederatedRepositoryBaseParams(&hfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Helm(hfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, hfp)

	setFederatedRepositoryBaseParams(&hfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Helm(hfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, hfp)
}

func federatedIvyTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	ifp := services.NewIvyFederatedRepositoryParams()
	ifp.Key = repoKey
	setFederatedRepositoryBaseParams(&ifp.FederatedRepositoryBaseParams, false)
	setJavaPackageManagersRepositoryParams(&ifp.JavaPackageManagersRepositoryParams, false)

	err := testsCreateFederatedRepositoryService.Ivy(ifp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, ifp)

	setFederatedRepositoryBaseParams(&ifp.FederatedRepositoryBaseParams, true)
	setJavaPackageManagersRepositoryParams(&ifp.JavaPackageManagersRepositoryParams, true)

	err = testsUpdateFederatedRepositoryService.Ivy(ifp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, ifp)
}

func federatedMavenTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	mfp := services.NewMavenFederatedRepositoryParams()
	mfp.Key = repoKey
	setFederatedRepositoryBaseParams(&mfp.FederatedRepositoryBaseParams, false)
	setJavaPackageManagersRepositoryParams(&mfp.JavaPackageManagersRepositoryParams, false)

	err := testsCreateFederatedRepositoryService.Maven(mfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, mfp)

	setFederatedRepositoryBaseParams(&mfp.FederatedRepositoryBaseParams, true)
	setJavaPackageManagersRepositoryParams(&mfp.JavaPackageManagersRepositoryParams, true)

	err = testsUpdateFederatedRepositoryService.Maven(mfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, mfp)
}

func federatedNpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nfp := services.NewNpmFederatedRepositoryParams()
	nfp.Key = repoKey
	setFederatedRepositoryBaseParams(&nfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Npm(nfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, nfp)

	setFederatedRepositoryBaseParams(&nfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Npm(nfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nfp)
}

func federatedNugetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nfp := services.NewNugetFederatedRepositoryParams()
	nfp.Key = repoKey
	setFederatedRepositoryBaseParams(&nfp.FederatedRepositoryBaseParams, false)
	setNugetRepositoryParams(&nfp.NugetRepositoryParams, false)

	err := testsCreateFederatedRepositoryService.Nuget(nfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, nfp)

	setFederatedRepositoryBaseParams(&nfp.FederatedRepositoryBaseParams, true)
	setNugetRepositoryParams(&nfp.NugetRepositoryParams, true)

	err = testsUpdateFederatedRepositoryService.Nuget(nfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nfp)
}

func federatedOpkgTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	ofp := services.NewOpkgFederatedRepositoryParams()
	ofp.Key = repoKey
	setFederatedRepositoryBaseParams(&ofp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Opkg(ofp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, ofp)

	setFederatedRepositoryBaseParams(&ofp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Opkg(ofp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, ofp)
}

func federatedPuppetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	pfp := services.NewPuppetFederatedRepositoryParams()
	pfp.Key = repoKey
	setFederatedRepositoryBaseParams(&pfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Puppet(pfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, pfp)

	setFederatedRepositoryBaseParams(&pfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Puppet(pfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, pfp)
}

func federatedPypiTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	pfp := services.NewPypiFederatedRepositoryParams()
	pfp.Key = repoKey
	setFederatedRepositoryBaseParams(&pfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Pypi(pfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, pfp)

	setFederatedRepositoryBaseParams(&pfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Pypi(pfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, pfp)
}

func federatedRpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	rfp := services.NewRpmFederatedRepositoryParams()
	rfp.Key = repoKey
	setFederatedRepositoryBaseParams(&rfp.FederatedRepositoryBaseParams, false)
	setRpmRepositoryParams(&rfp.RpmRepositoryParams, false)

	err := testsCreateFederatedRepositoryService.Rpm(rfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, rfp)

	setFederatedRepositoryBaseParams(&rfp.FederatedRepositoryBaseParams, true)
	setRpmRepositoryParams(&rfp.RpmRepositoryParams, true)

	err = testsUpdateFederatedRepositoryService.Rpm(rfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, rfp)
}

func federatedSbtTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	sfp := services.NewSbtFederatedRepositoryParams()
	sfp.Key = repoKey
	setFederatedRepositoryBaseParams(&sfp.FederatedRepositoryBaseParams, false)
	setJavaPackageManagersRepositoryParams(&sfp.JavaPackageManagersRepositoryParams, false)

	err := testsCreateFederatedRepositoryService.Sbt(sfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, sfp)

	setFederatedRepositoryBaseParams(&sfp.FederatedRepositoryBaseParams, true)
	setJavaPackageManagersRepositoryParams(&sfp.JavaPackageManagersRepositoryParams, true)

	err = testsUpdateFederatedRepositoryService.Sbt(sfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, sfp)
}

func federatedSwiftTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	sfp := services.NewSwiftFederatedRepositoryParams()
	sfp.Key = repoKey
	setFederatedRepositoryBaseParams(&sfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Swift(sfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, sfp)

	setFederatedRepositoryBaseParams(&sfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Swift(sfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, sfp)
}

func federatedVagrantTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	vfp := services.NewVagrantFederatedRepositoryParams()
	vfp.Key = repoKey
	setFederatedRepositoryBaseParams(&vfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Vagrant(vfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, vfp)

	setFederatedRepositoryBaseParams(&vfp.FederatedRepositoryBaseParams, true)

	err = testsUpdateFederatedRepositoryService.Vagrant(vfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, vfp)
}

func federatedYumTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	yfp := services.NewYumFederatedRepositoryParams()
	yfp.Key = repoKey
	setFederatedRepositoryBaseParams(&yfp.FederatedRepositoryBaseParams, false)
	setRpmRepositoryParams(&yfp.RpmRepositoryParams, false)

	err := testsCreateFederatedRepositoryService.Yum(yfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// "yum" package type is converted to "rpm" by Artifactory, so we have to change it too to pass the validation.
	yfp.PackageType = "rpm"
	validateRepoConfig(t, repoKey, yfp)

	setFederatedRepositoryBaseParams(&yfp.FederatedRepositoryBaseParams, true)
	setRpmRepositoryParams(&yfp.RpmRepositoryParams, true)

	err = testsUpdateFederatedRepositoryService.Yum(yfp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, yfp)
}

func federatedCreateWithParamTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	params := services.NewFederatedRepositoryBaseParams()
	params.Key = repoKey
	err := testsRepositoriesService.CreateFederated(params)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, params)
}

func getFederatedRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gfp := services.NewGenericFederatedRepositoryParams()
	gfp.Key = repoKey
	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Generic(gfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// Get repo details
	data := getRepo(t, repoKey)
	// Validate
	assert.Equal(t, data.Key, repoKey)
	assert.Equal(t, data.Description, gfp.Description)
	assert.Equal(t, data.GetRepoType(), "federated")
	assert.Empty(t, data.Url)
	assert.Equal(t, data.PackageType, "generic")
}

func getAllFederatedRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gfp := services.NewGenericFederatedRepositoryParams()
	gfp.Key = repoKey
	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Generic(gfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
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
	assert.Equal(t, gfp.Description, repo.Description)
	assert.Equal(t, "Generic", repo.PackageType)
}
