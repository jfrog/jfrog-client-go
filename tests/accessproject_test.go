package tests

import (
	"reflect"
	"testing"

	"github.com/jfrog/jfrog-client-go/access/services"
	"github.com/stretchr/testify/assert"
)

func TestAccessProject(t *testing.T) {
	initAccessTest(t)
	t.Run("create-update-delete", testAccessProjectCreateUpdateDelete)
}

func TestAccessProjectGroups(t *testing.T) {
	initAccessTest(t)
	t.Run("groups-add-get-delete", testAccessProjectAddGetDeleteGroups)
}

func testAccessProjectAddGetDeleteGroups(t *testing.T) {
	testGroup := getTestProjectGroupParams("a-test-group")
	projectParams := getTestProjectParams()
	err := testsAccessProjectService.Create(projectParams)
	defer deleteProjectAndGroupAndAssert(t, projectParams.ProjectDetails.ProjectKey, testGroup.Name)
	assert.NoError(t, err)

	toBeAddedGroup := getTestGroupParams(true)
	toBeAddedGroup.GroupDetails.Name = testGroup.Name
	err = testGroupService.CreateGroup(toBeAddedGroup)
	assert.NoError(t, err)

	err = testsAccessProjectService.UpdateGroup(projectParams.ProjectDetails.ProjectKey, testGroup.Name, testGroup)
	assert.NoError(t, err)

	allGroups, err := testsAccessProjectService.GetGroups(projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	if assert.NotNil(t, allGroups) {
		assert.Equal(t, len(*allGroups), 1, "Expected 1 group in the project but got %d", len(*allGroups))
		assert.Contains(t, *allGroups, testGroup)
	}

	testGroup.Roles = append(testGroup.Roles, "Contributor")
	err = testsAccessProjectService.UpdateGroup(projectParams.ProjectDetails.ProjectKey, testGroup.Name, testGroup)
	assert.NoError(t, err)

	singleGroup, err := testsAccessProjectService.GetGroup(projectParams.ProjectDetails.ProjectKey, testGroup.Name)
	assert.NoError(t, err)
	assert.Equal(t, *singleGroup, testGroup, "Expected group %v but got %v", *singleGroup, testGroup)

	err = testsAccessProjectService.DeleteExistingGroup(projectParams.ProjectDetails.ProjectKey, testGroup.Name)
	assert.NoError(t, err)

	noGroups, err := testsAccessProjectService.GetGroups(projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	assert.Empty(t, noGroups)
}

func testAccessProjectCreateUpdateDelete(t *testing.T) {
	projectParams := getTestProjectParams()
	err := testsAccessProjectService.Create(projectParams)
	defer deleteProjectAndAssert(t, projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	projectParams.ProjectDetails.Description += "123"
	projectParams.ProjectDetails.StorageQuotaBytes += 123
	projectParams.ProjectDetails.SoftLimit = &trueValue
	projectParams.ProjectDetails.AdminPrivileges.ManageMembers = &falseValue
	projectParams.ProjectDetails.AdminPrivileges.ManageResources = &trueValue
	projectParams.ProjectDetails.AdminPrivileges.IndexResources = &falseValue
	assert.NoError(t, testsAccessProjectService.Update(projectParams))
	updatedProject, err := testsAccessProjectService.Get(projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	assert.NotNil(t, updatedProject)
	if assert.NotNil(t, updatedProject) && !reflect.DeepEqual(projectParams.ProjectDetails, *updatedProject) {
		t.Error("Unexpected project details built. Expected: `", projectParams.ProjectDetails, "` Got `", *updatedProject, "`")
	}
}

func deleteProjectAndGroupAndAssert(t *testing.T, projectKey string, groupName string) {
	deleteProjectAndAssert(t, projectKey)
	deleteGroupAndAssert(t, groupName)
}

func deleteProjectAndAssert(t *testing.T, projectKey string) {
	err := testsAccessProjectService.Delete(projectKey)
	assert.NoError(t, err)
}

func getTestProjectParams() services.ProjectParams {
	adminPrivileges := services.AdminPrivileges{
		ManageMembers:   &trueValue,
		ManageResources: &falseValue,
		IndexResources:  &trueValue,
	}
	runId := getRunId()
	runNumberSuffix := runId[len(runId)-3:]
	projectDetails := services.Project{
		DisplayName:       "testProject" + runNumberSuffix,
		Description:       "My Test Project",
		AdminPrivileges:   &adminPrivileges,
		SoftLimit:         &falseValue,
		StorageQuotaBytes: 1073741825,                 // needs to be higher than 1073741824
		ProjectKey:        "tstprj" + runNumberSuffix, // valid length: 2 <= ProjectKey <= 10
	}
	return services.ProjectParams{
		ProjectDetails: projectDetails,
	}
}

func getTestProjectGroupParams(groupName string) services.ProjectGroup {
	return services.ProjectGroup{
		Name:  groupName,
		Roles: []string{"Contributor", "Release Manager"},
	}
}
