package tests

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

func TestGroups(t *testing.T) {
	t.Run("create", testCreateGroup)
	t.Run("update", testUpdateGroup)
	t.Run("delete", testDeleteGroup)
}

func testCreateGroup(t *testing.T) {
	groupParams := getTestGroupParams(false)
	err := testGroupService.CreateGroup(groupParams)
	defer testGroupService.DeleteGroup(groupParams.GroupDetails.Name)
	assert.NoError(t, err)
	createdGroup, _, err := testGroupService.GetGroup(groupParams)
	assert.NotNil(t, createdGroup)
	assert.Equal(t, groupParams.GroupDetails, *createdGroup)
}

func testUpdateGroup(t *testing.T) {
	groupParams := getTestGroupParams(false)
	err := testGroupService.CreateGroup(groupParams)
	defer testGroupService.DeleteGroup(groupParams.GroupDetails.Name)
	assert.NoError(t, err)
	groupParams.GroupDetails.Description = "Changed description"
	err = testGroupService.UpdateGroup(groupParams)
	assert.NoError(t, err)
	grp, _, err := testGroupService.GetGroup(groupParams)
	assert.NoError(t, err)
	assert.Equal(t, groupParams.GroupDetails, *grp)
}

func testAddUsersToGroup(t *testing.T) {
	groupParams := getTestGroupParams(true)
	err := testGroupService.CreateGroup(groupParams)
	defer testGroupService.DeleteGroup(groupParams.GroupDetails.Name)
	assert.NoError(t, err)
	groupParams.GroupDetails.UsersNames = []string{"Alice", "Bob"}
	err = testGroupService.UpdateGroup(groupParams)
	assert.NoError(t, err)
	grp, _, err := testGroupService.GetGroup(groupParams)
	assert.NoError(t, err)
	assert.Equal(t, groupParams.GroupDetails, *grp)
}

func testDeleteGroup(t *testing.T) {
	groupParams := getTestGroupParams(false)
	err := testGroupService.CreateGroup(groupParams)
	assert.NoError(t, err)
	err = testGroupService.DeleteGroup(groupParams.GroupDetails.Name)
	assert.NoError(t, err)
	group, notExists, err := testGroupService.GetGroup(groupParams)
	assert.True(t, notExists)
	assert.Nil(t, group)
}

func getTestGroupParams(includeUsers bool) services.GroupParams {
	groupDetails := services.Group{
		Name:            fmt.Sprintf("test%d", rand.Int()),
		Description:     "hello",
		AutoJoin:        false,
		AdminPrivileges: true,
		Realm:           "internal",
		RealmAttributes: "",
	}
	return services.GroupParams{
		GroupDetails: groupDetails,
		IncludeUsers: includeUsers,
	}
}
