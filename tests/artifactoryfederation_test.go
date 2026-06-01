//go:build itest

package tests

import (
	"testing"

	"github.com/jfrog/gofrog/version"
	"github.com/jfrog/jfrog-client-go/utils"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

func TestArtifactoryFederation(t *testing.T) {
	initRepositoryTest(t)
	rtVersion, err := GetRtDetails().GetVersion()
	if err != nil {
		t.Error(err)
	}
	if !version.NewVersion(rtVersion).AtLeast("7.18.3") {
		t.Skip("Skipping artifactory test. Federated repositories are only supported by Artifactory 7.18.3 or higher.")
	}
	t.Run("localConvertLocalToFederatedTest", localConvertLocalToFederatedTest)
	t.Run("localConvertNonExistentLocalToFederatedTest", localConvertNonExistentLocalToFederatedTest)
	t.Run("localTriggerFederatedFullSyncAllTest", localTriggerFederatedFullSyncAllTest)
	t.Run("localTriggerFederatedFullSyncMirrorTest", localTriggerFederatedFullSyncMirrorTest)
}

func localConvertLocalToFederatedTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "simple-default"
	glp.Description = "Generic Repo for jfrog-client-go federation-test"
	glp.XrayIndex = utils.Pointer(true)
	glp.DownloadRedirect = utils.Pointer(false)
	glp.ArchiveBrowsingEnabled = utils.Pointer(false)

	err := testsCreateLocalRepositoryService.Generic(glp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)

	err = testsFederationService.ConvertLocalToFederated(repoKey)
	assert.NoError(t, err, "Failed to convert "+repoKey)
}

func localConvertNonExistentLocalToFederatedTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	err := testsFederationService.ConvertLocalToFederated(repoKey)
	assert.Error(t, err, "Failed to not convert "+repoKey)
}

func localTriggerFederatedFullSyncAllTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gfp := services.NewGenericFederatedRepositoryParams()
	gfp.Key = repoKey
	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Generic(gfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, gfp)

	err = testsFederationService.TriggerFederatedFullSyncAll(repoKey)
	assert.NoError(t, err, "Failed to trigger full synchonisation "+repoKey)
}

func localTriggerFederatedFullSyncMirrorTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	gfp := services.NewGenericFederatedRepositoryParams()
	gfp.Key = repoKey
	setFederatedRepositoryBaseParams(&gfp.FederatedRepositoryBaseParams, false)

	err := testsCreateFederatedRepositoryService.Generic(gfp)
	if !assert.NoError(t, err, "Failed to create "+repoKey) {
		return
	}
	deleteRepoOnTestDone(t, repoKey)
	validateRepoConfig(t, repoKey, gfp)

	mirror := gfp.Members[0].Url
	err = testsFederationService.TriggerFederatedFullSyncMirror(repoKey, mirror)
	assert.NoError(t, err, "Failed to trigger synchonisation "+repoKey+" for "+mirror)
}
