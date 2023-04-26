package tests

import (
	"fmt"
	"testing"

	services "github.com/jfrog/jfrog-client-go/access/services/v2"
	"github.com/stretchr/testify/assert"
)

func TestAccessGroups(t *testing.T) {
	initAccessTest(t)
	t.Run("create", testCreateAccessGroup)
	t.Run("update", testUpdateAccessGroup)
	t.Run("delete", testDeleteAccessGroup)
}

func getTestAccessGroupParams() services.GroupParams {
	group := services.GroupDetails{
		Name:            fmt.Sprintf("test-%s", getRunId()),
		Description:     "hello",
		AutoJoin:        &falseValue,
		AdminPrivileges: &trueValue,
		Realm:           "internal",
		RealmAttributes: "",
		ExternalId:      "",
	}
	return services.GroupParams{GroupDetails: group}
}

func testCreateAccessGroup(t *testing.T) {
	groupParams := getTestAccessGroupParams()
	err := testAccessGroupService.Create(groupParams)
	defer deleteAccessGroupAndAssert(t, groupParams.GroupDetails.Name)
	assert.NoError(t, err)

	createdGroup, err := testAccessGroupService.Get(groupParams.Name)
	assert.NoError(t, err)
	assert.NotNil(t, createdGroup)
	assert.Equal(t, groupParams.GroupDetails, *createdGroup)

	allGroups, err := testAccessGroupService.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, allGroups)

	var groupNames []string
	for _, v := range allGroups {
		groupNames = append(groupNames, v.GroupName)
	}
	assert.Contains(t, groupNames, groupParams.GroupDetails.Name)

}

func testUpdateAccessGroup(t *testing.T) {
	groupParams := getTestAccessGroupParams()
	err := testAccessGroupService.Create(groupParams)
	defer deleteAccessGroupAndAssert(t, groupParams.Name)
	assert.NoError(t, err)
	groupParams.Description = "Changed description"
	groupParams.AutoJoin = &trueValue
	groupParams.AdminPrivileges = &falseValue
	err = testAccessGroupService.Update(groupParams)
	assert.NoError(t, err)
	group, err := testAccessGroupService.Get(groupParams.Name)
	assert.NoError(t, err)
	assert.Equal(t, groupParams.GroupDetails, *group)
}

func testDeleteAccessGroup(t *testing.T) {
	groupParams := getTestAccessGroupParams()
	assert.NoError(t, testAccessGroupService.Create(groupParams))

	deleteAccessGroupAndAssert(t, groupParams.Name)

	group, err := testAccessGroupService.Get(groupParams.Name)
	assert.NoError(t, err)
	assert.Nil(t, group)

	allGroups, err := testAccessGroupService.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, allGroups)

	var allGroupNames []string
	for _, v := range allGroups {
		allGroupNames = append(allGroupNames, v.GroupName)
	}
	assert.NotContains(t, allGroupNames, groupParams.Name)
}

func deleteAccessGroupAndAssert(t *testing.T, groupName string) {
	assert.NoError(t, testAccessGroupService.Delete(groupName))
}
