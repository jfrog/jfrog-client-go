package tests

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils"

	"github.com/jfrog/jfrog-client-go/artifactory/services"

	"github.com/stretchr/testify/assert"
)

const (
	PermissionTargetNamePrefix = "client-go-tests-target"
)

func TestPermissionTarget(t *testing.T) {
	initArtifactoryTest(t)
	params := services.NewPermissionTargetParams()
	params.Name = fmt.Sprintf("%s-%s", PermissionTargetNamePrefix, getRunId())
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
	assert.NotNil(t, targetConfig)
	if targetConfig == nil {
		return
	}
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
	params.Name = fmt.Sprintf("%s-%s", PermissionTargetNamePrefix, getRunId())

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

// This test directly copies the example in the documentation:
// - Creating and Updating Permission Targets
func TestDocumentationExampleCreateUpdateAndDeletePermissionTarget(t *testing.T) {
	initArtifactoryTest(t)

	// Preamble to get dependent entities setup
	user1 := createRandomUser(t)
	defer deleteUserAndAssert(t, user1)

	user2 := createRandomUser(t)
	defer deleteUserAndAssert(t, user2)

	localRepo1 := createRandomRepo(t)
	defer deleteRepo(t, localRepo1)

	localRepo2 := createRandomRepo(t)
	defer deleteRepo(t, localRepo2)

	group1 := createRandomGroup(t)
	defer deleteGroupAndAssert(t, group1)

	group2 := createRandomGroup(t)
	defer deleteGroupAndAssert(t, group2)

	// Example code from Documentation
	params := services.NewPermissionTargetParams()
	params.Name = "java-developers"
	params.Repo = &services.PermissionTargetSection{}
	params.Repo.Repositories = []string{"ANY REMOTE", localRepo1, localRepo2}
	params.Repo.ExcludePatterns = []string{"dir/*"}
	params.Repo.Actions = &services.Actions{}
	params.Repo.Actions.Users = map[string][]string{
		user1: {"read", "write"},
		user2: {"write", "annotate", "read"},
	}
	params.Repo.Actions.Groups = map[string][]string{
		group1: {"manage", "read", "annotate"},
	}
	// This is the default value that cannot be changed
	params.Build = &services.PermissionTargetSection{}
	params.Build.Repositories = []string{"artifactory-build-info"}
	params.Build.Actions = &services.Actions{}
	params.Build.Actions.Groups = map[string][]string{
		group1: {"manage", "read", "write", "annotate", "delete"},
		group2: {"read"},
	}

	// Creating the Permission Target
	err := testsPermissionTargetService.Create(params)
	assert.NoError(t, err)

	// Update the permission target
	err = testsPermissionTargetService.Update(params)
	assert.NoError(t, err)

	// Fetch a permission target
	_, err = testsPermissionTargetService.Get("java-developers")
	assert.NoError(t, err)

	// Fetch all permission targets
	_, err = testsPermissionTargetService.GetAll()
	assert.NoError(t, err)

	// Delete the permission target
	err = testsPermissionTargetService.Delete("java-developers")
	assert.NoError(t, err)
}

func createRandomUser(t *testing.T) string {
	name := fmt.Sprintf("test-%s-%s", timestampStr, randomString(t, 16))
	userDetails := services.User{
		Name:                     name,
		Email:                    name + "@jfrog.com",
		Password:                 "Password1*",
		Admin:                    utils.Pointer(false),
		Realm:                    "internal",
		ShouldInvite:             utils.Pointer(false),
		ProfileUpdatable:         utils.Pointer(true),
		DisableUIAccess:          utils.Pointer(false),
		InternalPasswordDisabled: utils.Pointer(false),
		WatchManager:             utils.Pointer(false),
		ReportsManager:           utils.Pointer(false),
		PolicyManager:            utils.Pointer(false),
	}

	err := testUserService.CreateUser(services.UserParams{
		UserDetails:     userDetails,
		ReplaceIfExists: true,
	})

	assert.NoError(t, err)

	return name
}

func createRandomRepo(t *testing.T) string {
	repoKey := fmt.Sprintf("test-%s-%s", timestampStr, randomString(t, 16))
	glp := services.NewGenericLocalRepositoryParams()
	glp.Key = repoKey
	setLocalRepositoryBaseParams(&glp.LocalRepositoryBaseParams, false)

	err := testsCreateLocalRepositoryService.Generic(glp)
	assert.NoError(t, err)

	return repoKey
}

func createRandomGroup(t *testing.T) string {
	name := fmt.Sprintf("test-%s-%s", timestampStr, randomString(t, 16))

	groupDetails := services.Group{
		Name:            name,
		Description:     "hello",
		AutoJoin:        utils.Pointer(false),
		AdminPrivileges: utils.Pointer(false),
		Realm:           "internal",
		RealmAttributes: "",
	}

	groupParams := services.GroupParams{
		GroupDetails: groupDetails,
		IncludeUsers: false,
	}

	err := testGroupService.CreateGroup(groupParams)
	assert.NoError(t, err)

	return name
}

func randomString(t *testing.T, length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	assert.NoError(t, err)
	return fmt.Sprintf("%x", b)[:length]
}
