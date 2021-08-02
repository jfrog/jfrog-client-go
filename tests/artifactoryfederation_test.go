package tests

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/utils/version"
	"github.com/stretchr/testify/assert"
)

func TestArtifactoryFederation(t *testing.T) {
	initArtifactoryTest(t)
	rtVersion, err := GetRtDetails().GetVersion()
	if err != nil {
		t.Error(err)
	}
	if !version.NewVersion(rtVersion).AtLeast("7.18.3") {
		t.Skip("Skipping artifactory test. Federated repositories are only supported by Artifactory 7.18.3 or higher.")
	}
	t.Run("localConvertLocalToFederatedTest", localConvertLocalToFederatedTest)
	t.Run("localConvertNonExistentLocalToFederatedTest", localConvertNonExistentLocalToFederatedTest)
}

func localConvertLocalToFederatedTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKey
	glp.RepoLayoutRef = "simple-default"
	glp.Description = "Generic Repo for jfrog-client-go federation-test"
	glp.XrayIndex = &trueValue
	glp.DownloadRedirect = &falseValue
	glp.ArchiveBrowsingEnabled = &falseValue

	err := testsCreateLocalRepositoryService.Generic(glp)
	assert.NoError(t, err, "Failed to create "+repoKey)
	defer deleteRepo(t, repoKey)

	err = testsFederationService.ConvertLocalToFederated(repoKey)
	assert.NoError(t, err, "Failed to convert "+repoKey)
}

func localConvertNonExistentLocalToFederatedTest(t *testing.T) {
	repoKey := GenerateRepoKeyForRepoServiceTest()
	err := testsFederationService.ConvertLocalToFederated(repoKey)
	assert.Error(t, err, "Failed to not convert "+repoKey)
}
