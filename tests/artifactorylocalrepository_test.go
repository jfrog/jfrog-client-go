package tests

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

func TestArtifactoryLocalRepository(t *testing.T) {
	initRepositoryTest(t)
	t.Run("localAlpineTest", localAlpineTest)
	t.Run("localBowerTest", localBowerTest)
	t.Run("localCargoTest", localCargoTest)
	t.Run("localChefTest", localChefTest)
	t.Run("localCocoapodsTest", localCocoapodsTest)
	t.Run("localComposerTest", localComposerTest)
	t.Run("localConanTest", localConanTest)
	t.Run("localCondaTest", localCondaTest)
	t.Run("localCranTest", localCranTest)
	t.Run("localDebianTest", localDebianTest)
	t.Run("localDockerTest", localDockerTest)
	t.Run("localGemsTest", localGemsTest)
	t.Run("localGenericTest", localGenericTest)
	t.Run("localGitlfsTest", localGitlfsTest)
	t.Run("localGoTest", localGoTest)
	t.Run("localGradleTest", localGradleTest)
	t.Run("localHelmTest", localHelmTest)
	t.Run("localIvyTest", localIvyTest)
	t.Run("localMavenTest", localMavenTest)
	t.Run("localNpmTest", localNpmTest)
	t.Run("localNugetTest", localNugetTest)
	t.Run("localOkgTest", localOpkgTest)
	t.Run("localPuppetTest", localPuppetTest)
	t.Run("localPypiTest", localPypiTest)
	t.Run("localRpmTest", localRpmTest)
	t.Run("localSbtTest", localSbtTest)
	t.Run("localSwiftTest", localSwiftTest)
	t.Run("localVagrantTest", localVagrantTest)
	t.Run("localYumTest", localYumTest)
	t.Run("localCreateWithParamTest", localCreateWithParamTest)
	t.Run("getLocalRepoDetailsTest", getLocalRepoDetailsTest)
	t.Run("getAllLocalRepoDetailsTest", getAllLocalRepoDetailsTest)
	t.Run("isLocalRepoExistsTest", isLocalRepoExistsTest)
}

func setLocalRepositoryBaseParams(params *services.LocalRepositoryBaseParams, isUpdate bool) {
	setRepositoryBaseParams(&params.RepositoryBaseParams, isUpdate)
	setAdditionalRepositoryBaseParams(&params.AdditionalRepositoryBaseParams, isUpdate)
	if !isUpdate {
		params.ArchiveBrowsingEnabled = &trueValue
	} else {
		params.ArchiveBrowsingEnabled = &falseValue
	}
}

func localAlpineTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	alp := services.NewAlpineLocalRepositoryParams()
	alp.Key = repoKey
	setLocalRepositoryBaseParams(&alp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Alpine(alp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, alp)

	setLocalRepositoryBaseParams(&alp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Alpine(alp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, alp)
}

func localBowerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	blp := services.NewBowerLocalRepositoryParams()
	blp.Key = repoKey
	setLocalRepositoryBaseParams(&blp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Bower(blp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, blp)

	setLocalRepositoryBaseParams(&blp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Bower(blp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, blp)
}

func localCargoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewCargoLocalRepositoryParams()
	clp.Key = repoKey
	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, false)
	setCargoRepositoryParams(&clp.CargoRepositoryParams, false)

	err := testsCreateLocalRepositoryService.Cargo(clp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, true)
	setCargoRepositoryParams(&clp.CargoRepositoryParams, true)

	err = testsUpdateLocalRepositoryService.Cargo(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func localChefTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewChefLocalRepositoryParams()
	clp.Key = repoKey
	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Chef(clp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Chef(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func localCocoapodsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewCocoapodsLocalRepositoryParams()
	clp.Key = repoKey
	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Cocoapods(clp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Cocoapods(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func localComposerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewComposerLocalRepositoryParams()
	clp.Key = repoKey
	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Composer(clp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Composer(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func localConanTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewConanLocalRepositoryParams()
	clp.Key = repoKey
	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Conan(clp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Conan(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func localCondaTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewCondaLocalRepositoryParams()
	clp.Key = repoKey
	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Conda(clp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Conda(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func localCranTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	clp := services.NewCranLocalRepositoryParams()
	clp.Key = repoKey
	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Cran(clp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, clp)

	setLocalRepositoryBaseParams(&clp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Cran(clp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, clp)
}

func localDebianTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	dlp := services.NewDebianLocalRepositoryParams()
	dlp.Key = repoKey
	setLocalRepositoryBaseParams(&dlp.LocalRepositoryBaseParams, false)
	setDebianRepositoryParams(&dlp.DebianRepositoryParams, false)

	err := testsCreateLocalRepositoryService.Debian(dlp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, dlp)

	setLocalRepositoryBaseParams(&dlp.LocalRepositoryBaseParams, true)
	setDebianRepositoryParams(&dlp.DebianRepositoryParams, true)

	err = testsUpdateLocalRepositoryService.Debian(dlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dlp)
}

func localDockerTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	dlp := services.NewDockerLocalRepositoryParams()
	dlp.Key = repoKey
	setLocalRepositoryBaseParams(&dlp.LocalRepositoryBaseParams, false)
	setDockerRepositoryParams(&dlp.DockerRepositoryParams, false)

	err := testsCreateLocalRepositoryService.Docker(dlp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, dlp)

	setLocalRepositoryBaseParams(&dlp.LocalRepositoryBaseParams, true)
	setDockerRepositoryParams(&dlp.DockerRepositoryParams, true)

	err = testsUpdateLocalRepositoryService.Docker(dlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, dlp)
}

func localGemsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGemsLocalRepositoryParams()
	glp.Key = repoKey
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Gems(glp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, glp)

	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Gems(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)
}

func localGenericTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKey
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Generic(glp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, glp)

	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Generic(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)
}

func localGitlfsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGitlfsLocalRepositoryParams()
	glp.Key = repoKey
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Gitlfs(glp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, glp)

	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Gitlfs(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)
}

func localGoTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGoLocalRepositoryParams()
	glp.Key = repoKey
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Go(glp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, glp)

	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Go(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)
}

func localGradleTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGradleLocalRepositoryParams()
	glp.Key = repoKey
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, false)
	setJavaPackageManagersRepositoryParams(&glp.JavaPackageManagersRepositoryParams, false)

	err := testsCreateLocalRepositoryService.Gradle(glp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, glp)

	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, true)
	setJavaPackageManagersRepositoryParams(&glp.JavaPackageManagersRepositoryParams, true)

	err = testsUpdateLocalRepositoryService.Gradle(glp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, glp)
}

func localHelmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	hlp := services.NewHelmLocalRepositoryParams()
	hlp.Key = repoKey
	setLocalRepositoryBaseParams(&hlp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Helm(hlp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, hlp)

	setLocalRepositoryBaseParams(&hlp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Helm(hlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, hlp)
}

func localIvyTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	ilp := services.NewIvyLocalRepositoryParams()
	ilp.Key = repoKey
	setLocalRepositoryBaseParams(&ilp.LocalRepositoryBaseParams, false)
	setJavaPackageManagersRepositoryParams(&ilp.JavaPackageManagersRepositoryParams, false)

	err := testsCreateLocalRepositoryService.Ivy(ilp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, ilp)

	setLocalRepositoryBaseParams(&ilp.LocalRepositoryBaseParams, true)
	setJavaPackageManagersRepositoryParams(&ilp.JavaPackageManagersRepositoryParams, true)

	err = testsUpdateLocalRepositoryService.Ivy(ilp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, ilp)
}

func localMavenTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	mlp := services.NewMavenLocalRepositoryParams()
	mlp.Key = repoKey
	setLocalRepositoryBaseParams(&mlp.LocalRepositoryBaseParams, false)
	setJavaPackageManagersRepositoryParams(&mlp.JavaPackageManagersRepositoryParams, false)

	err := testsCreateLocalRepositoryService.Maven(mlp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, mlp)

	setLocalRepositoryBaseParams(&mlp.LocalRepositoryBaseParams, true)
	setJavaPackageManagersRepositoryParams(&mlp.JavaPackageManagersRepositoryParams, true)

	err = testsUpdateLocalRepositoryService.Maven(mlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, mlp)
}

func localNpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nlp := services.NewNpmLocalRepositoryParams()
	nlp.Key = repoKey
	setLocalRepositoryBaseParams(&nlp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Npm(nlp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, nlp)

	setLocalRepositoryBaseParams(&nlp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Npm(nlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nlp)
}

func localNugetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	nlp := services.NewNugetLocalRepositoryParams()
	nlp.Key = repoKey
	setLocalRepositoryBaseParams(&nlp.LocalRepositoryBaseParams, false)
	setNugetRepositoryParams(&nlp.NugetRepositoryParams, false)

	err := testsCreateLocalRepositoryService.Nuget(nlp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, nlp)

	setLocalRepositoryBaseParams(&nlp.LocalRepositoryBaseParams, true)
	setNugetRepositoryParams(&nlp.NugetRepositoryParams, true)

	err = testsUpdateLocalRepositoryService.Nuget(nlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, nlp)
}

func localOpkgTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	olp := services.NewOpkgLocalRepositoryParams()
	olp.Key = repoKey
	setLocalRepositoryBaseParams(&olp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Opkg(olp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, olp)

	setLocalRepositoryBaseParams(&olp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Opkg(olp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, olp)
}

func localPuppetTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	plp := services.NewPuppetLocalRepositoryParams()
	plp.Key = repoKey
	setLocalRepositoryBaseParams(&plp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Puppet(plp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, plp)

	setLocalRepositoryBaseParams(&plp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Puppet(plp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, plp)
}

func localPypiTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	plp := services.NewPypiLocalRepositoryParams()
	plp.Key = repoKey
	setLocalRepositoryBaseParams(&plp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Pypi(plp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, plp)

	setLocalRepositoryBaseParams(&plp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Pypi(plp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, plp)
}

func localRpmTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	rlp := services.NewRpmLocalRepositoryParams()
	rlp.Key = repoKey
	setLocalRepositoryBaseParams(&rlp.LocalRepositoryBaseParams, false)
	setRpmRepositoryParams(&rlp.RpmRepositoryParams, false)

	err := testsCreateLocalRepositoryService.Rpm(rlp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, rlp)

	setLocalRepositoryBaseParams(&rlp.LocalRepositoryBaseParams, true)
	setRpmRepositoryParams(&rlp.RpmRepositoryParams, true)

	err = testsUpdateLocalRepositoryService.Rpm(rlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, rlp)
}

func localSbtTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	slp := services.NewSbtLocalRepositoryParams()
	slp.Key = repoKey
	setLocalRepositoryBaseParams(&slp.LocalRepositoryBaseParams, false)
	setJavaPackageManagersRepositoryParams(&slp.JavaPackageManagersRepositoryParams, false)

	err := testsCreateLocalRepositoryService.Sbt(slp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, slp)

	setLocalRepositoryBaseParams(&slp.LocalRepositoryBaseParams, true)
	setJavaPackageManagersRepositoryParams(&slp.JavaPackageManagersRepositoryParams, true)

	err = testsUpdateLocalRepositoryService.Sbt(slp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, slp)
}

func localSwiftTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	slp := services.NewSwiftLocalRepositoryParams()
	slp.Key = repoKey
	setLocalRepositoryBaseParams(&slp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Swift(slp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, slp)

	setLocalRepositoryBaseParams(&slp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Swift(slp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, slp)
}

func localVagrantTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	vlp := services.NewVagrantLocalRepositoryParams()
	vlp.Key = repoKey
	setLocalRepositoryBaseParams(&vlp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Vagrant(vlp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, vlp)

	setLocalRepositoryBaseParams(&vlp.LocalRepositoryBaseParams, true)

	err = testsUpdateLocalRepositoryService.Vagrant(vlp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, vlp)
}

func localYumTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	ylp := services.NewYumLocalRepositoryParams()
	ylp.Key = repoKey
	setLocalRepositoryBaseParams(&ylp.LocalRepositoryBaseParams, false)
	yumRootDepth := 6
	ylp.YumRootDepth = &yumRootDepth
	ylp.CalculateYumMetadata = &trueValue
	ylp.EnableFileListsIndexing = &trueValue
	ylp.YumGroupFileNames = "filename"

	err := testsCreateLocalRepositoryService.Yum(ylp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// "yum" package type is converted to "rpm" by Artifactory, so we have to change it too to pass the validation.
	ylp.PackageType = "rpm"
	validateRepoConfig(t, repoKey, ylp)

	setLocalRepositoryBaseParams(&ylp.LocalRepositoryBaseParams, true)
	*ylp.YumRootDepth = 18
	ylp.CalculateYumMetadata = &falseValue
	ylp.EnableFileListsIndexing = &falseValue
	ylp.YumGroupFileNames = ""

	err = testsUpdateLocalRepositoryService.Yum(ylp)
	assert.NoError(t, err, "Failed to update "+repoKey)
	validateRepoConfig(t, repoKey, ylp)
}

func localCreateWithParamTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	params := services.NewLocalRepositoryBaseParams()
	params.Key = repoKey
	err := testsRepositoriesService.CreateLocal(params)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	validateRepoConfig(t, repoKey, params)
}

func getLocalRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKey
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Generic(glp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// Get repo details
	data := getRepo(t, repoKey)
	// Validate
	assert.Equal(t, data.Key, repoKey)
	assert.Equal(t, data.Description, glp.Description)
	assert.Equal(t, data.GetRepoType(), "local")
	assert.Empty(t, data.Url)
	assert.Equal(t, data.PackageType, "generic")
}

func isLocalRepoExistsTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	// Validate repo doesn't exist
	exists := isRepoExists(t, repoKey)
	assert.False(t, exists)
	// Create Repo
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKey
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, false)
	err := testsCreateLocalRepositoryService.Generic(glp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// Validate repo exists
	exists = isRepoExists(t, repoKey)
	assert.True(t, exists)
}

func getAllLocalRepoDetailsTest(t *testing.T) {
	// Create Repo
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKey
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Generic(glp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	defer deleteRepo(t, repoKey)
	// Get repo details
	data := getAllRepos(t, "local", "")
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
