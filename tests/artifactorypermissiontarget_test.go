package tests

import (
	"fmt"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"

	"github.com/stretchr/testify/assert"
)

const (
	PermissionTargetNamePrefix = JfrogRepoPrefix + "-client-go-tests-target"
)

func TestPermissionTarget(t *testing.T) {
	initArtifactoryTest(t)
	params := services.NewPermissionTargetParams()
	params.Name = fmt.Sprintf("%s-%s", PermissionTargetNamePrefix, randomRunNumber)
	params.Repo = &services.PermissionTargetSection{}
	params.Repo.Repositories = []string{"ANY"}
	params.Repo.ExcludePatterns = []string{"dir/*"}
	params.Repo.Actions = &services.Actions{}
	params.Repo.Actions.Users = map[string][]string{
		"anonymous": {"read"},
	}
	params.Build = &services.PermissionTargetSection{}
	params.Build.Repositories = []string{"artifactory-build-info"}
	params.Build.Actions = &services.Actions{}
	params.Build.Actions.Users = map[string][]string{
		"anonymous": {"annotate"},
	}

	err := testsPermissionTargetService.Create(params)
	assert.NoError(t, err)
	// Fill in default values before validation
	params.Repo.IncludePatterns = []string{"**"}
	params.Build.Repositories = []string{"artifactory-build-info"}
	params.Build.IncludePatterns = []string{"**"}
	params.Build.ExcludePatterns = []string{}
	validatePermissionTarget(t, params)

	params.Repo.Actions.Users = nil
	params.Repo.Repositories = []string{"ANY REMOTE"}
	err = testsPermissionTargetService.Update(params)
	validatePermissionTarget(t, params)
	assert.NoError(t, err)
	err = testsPermissionTargetService.Delete(params.Name)
	assert.NoError(t, err)
	targetParams, err := getPermissionTarget(params.Name)
	assert.NoError(t, err)
	assert.Nil(t, targetParams)
}

func validatePermissionTarget(t *testing.T, params services.PermissionTargetParams) {
	targetConfig, err := getPermissionTarget(params.Name)
	assert.NoError(t, err)
	assert.Equal(t, params.Name, targetConfig.Name)
	assert.Equal(t, params.Repo, targetConfig.Repo)
	assert.Equal(t, params.Build, targetConfig.Build)
	assert.Equal(t, params.ReleaseBundle, targetConfig.ReleaseBundle)
}

func getPermissionTarget(targetName string) (targetParams *services.PermissionTargetParams, err error) {
	return testsPermissionTargetService.Get(targetName)
}

// Assert empty inner structs remain nil unless explicitly set.
func TestPermissionTargetEmptyFields(t *testing.T) {
	initArtifactoryTest(t)
	params := services.NewPermissionTargetParams()
	params.Name = fmt.Sprintf("%s-%s", PermissionTargetNamePrefix, randomRunNumber)

	assert.Nil(t, params.Repo)
	params.Repo = &services.PermissionTargetSection{}
	params.Repo.Repositories = []string{"ANY"}
	params.Repo.IncludePatterns = []string{"**"}
	params.Repo.ExcludePatterns = []string{"dir/*"}
	params.Repo.Actions = &services.Actions{}
	params.Repo.Actions.Users = map[string][]string{
		"anonymous": {"read"},
	}

	assert.Nil(t, params.Build)
	assert.Nil(t, params.ReleaseBundle)
	assert.NoError(t, testsPermissionTargetService.Create(params))
	validatePermissionTarget(t, params)
	err := testsPermissionTargetService.Delete(params.Name)
	assert.NoError(t, err)
}
