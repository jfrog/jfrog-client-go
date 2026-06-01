//go:build itest

package tests

import (
	"fmt"
	"sort"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroups(t *testing.T) {
	initArtifactoryTest(t)
	t.Run("create", testCreateGroup)
	t.Run("update", testUpdateGroup)
	t.Run("delete", testDeleteGroup)
	t.Run("add users", testAddUsersToGroup)
}

func testCreateGroup(t *testing.T) {
	groupParams := createGroup(t, "", false, true)

	createdGroup, err := testGroupService.GetGroup(groupParams)
	assert.NoError(t, err)
	assert.NotNil(t, createdGroup)
	assert.Equal(t, groupParams.GroupDetails, *createdGroup)

	allGroups, err := testGroupService.GetAllGroups()
	assert.NoError(t, err)
	assert.NotNil(t, allGroups)
	assert.Contains(t, *allGroups, groupParams.GroupDetails.Name)
}

func testUpdateGroup(t *testing.T) {
	groupParams := createGroup(t, "", false, true)
	groupParams.GroupDetails.Description = "Changed description"
	groupParams.GroupDetails.AutoJoin = utils.Pointer(true)
	groupParams.GroupDetails.AdminPrivileges = utils.Pointer(false)
	err := testGroupService.UpdateGroup(groupParams)
	require.NoError(t, err)
	group, err := testGroupService.GetGroup(groupParams)
	require.NoError(t, err)
	assert.Equal(t, groupParams.GroupDetails, *group)
}

func testAddUsersToGroup(t *testing.T) {
	// Create group
	groupParams := getTestGroupParams(true)
	assert.NoError(t, testGroupService.CreateGroup(groupParams))
	defer deleteGroupAndAssert(t, groupParams.GroupDetails.Name)

	// Create two new users
	userNames := []string{"Alice", "Bob"}
	for i, name := range userNames {
		UserParams := getTestUserParams(false, name)
		assert.NoError(t, testUserService.CreateUser(UserParams))
		defer deleteUserAndAssert(t, UserParams.UserDetails.Name)
		user, err := testUserService.GetUser(UserParams)
		if assert.NoError(t, err) {
			userNames[i] = user.Name
		}
	}

	// Add users to group
	groupParams.GroupDetails.UsersNames = userNames
	assert.NoError(t, testGroupService.UpdateGroup(groupParams))
	group, err := testGroupService.GetGroup(groupParams)
	assert.NoError(t, err)
	// Ignore usernames order
	sort.Strings(groupParams.GroupDetails.UsersNames)
	sort.Strings(group.UsersNames)
	assert.Equal(t, groupParams.GroupDetails, *group)
}

func testDeleteGroup(t *testing.T) {
	groupParams := getTestGroupParams(false)
	assert.NoError(t, testGroupService.CreateGroup(groupParams))
	deleteGroupAndAssert(t, groupParams.GroupDetails.Name)
	group, err := testGroupService.GetGroup(groupParams)
	assert.NoError(t, err)
	assert.Nil(t, group)

	allGroups, err := testGroupService.GetAllGroups()
	assert.NoError(t, err)
	assert.NotNil(t, allGroups)
	assert.NotContains(t, *allGroups, groupParams.GroupDetails.Name)
}

func getTestGroupParams(includeUsers bool) services.GroupParams {
	groupDetails := services.Group{
		Name:            fmt.Sprintf("test-%s", getRunId()),
		Description:     "hello",
		AutoJoin:        utils.Pointer(false),
		AdminPrivileges: utils.Pointer(true),
		Realm:           "internal",
		RealmAttributes: "",
	}
	return services.GroupParams{
		GroupDetails: groupDetails,
		IncludeUsers: includeUsers,
	}
}

func deleteGroupAndAssert(t *testing.T, groupName string) {
	assert.NoError(t, testGroupService.DeleteGroup(groupName))
}
