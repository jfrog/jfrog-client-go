package tests

import (
	"github.com/jfrog/jfrog-client-go/access/services"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
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
	projectParams := getTestProjectParams()
	assert.NoError(t, testsAccessProjectService.Create(projectParams))

	testGroup := getTestProjectGroupParams("a-test-group")
	defer deleteProjectAndGroupAndAssert(t, projectParams.ProjectDetails.ProjectKey, testGroup.Name)

	toBeAddedGroup := getTestGroupParams(true)
	toBeAddedGroup.GroupDetails.Name = testGroup.Name
	assert.NoError(t, testGroupService.CreateGroup(toBeAddedGroup))
	assert.NoError(t, testsAccessProjectService.UpdateGroup(projectParams.ProjectDetails.ProjectKey, testGroup.Name, testGroup))

	allGroups, err := testsAccessProjectService.GetGroups(projectParams.ProjectDetails.ProjectKey)
	if assert.NoError(t, err) &&
		assert.NotNil(t, allGroups, "Expected 1 group in the project but got 0") {
		assert.Equal(t, len(*allGroups), 1, "Expected 1 group in the project but got %d", len(*allGroups))
		assert.Contains(t, *allGroups, testGroup)
	}

	testGroup.Roles = append(testGroup.Roles, "Viewer")
	assert.NoError(t, testsAccessProjectService.UpdateGroup(projectParams.ProjectDetails.ProjectKey, testGroup.Name, testGroup))

	singleGroup, err := testsAccessProjectService.GetGroup(projectParams.ProjectDetails.ProjectKey, testGroup.Name)
	if assert.NoError(t, err) &&
		assert.NotNil(t, singleGroup, "Expected group %s but got nil", testGroup.Name) {
		assert.Equal(t, testGroup, *singleGroup, "Expected group %v but got %v", testGroup, *singleGroup)
	}

	assert.NoError(t, testsAccessProjectService.DeleteExistingGroup(projectParams.ProjectDetails.ProjectKey, testGroup.Name))

	noGroups, err := testsAccessProjectService.GetGroups(projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	assert.Empty(t, noGroups)
}

func testAccessProjectCreateUpdateDelete(t *testing.T) {
	projectParams := getTestProjectParams()
	assert.NoError(t, testsAccessProjectService.Create(projectParams))
	defer deleteProjectAndAssert(t, projectParams.ProjectDetails.ProjectKey)
	projectParams.ProjectDetails.Description += "123"
	projectParams.ProjectDetails.StorageQuotaBytes += 123
	projectParams.ProjectDetails.SoftLimit = &trueValue
	projectParams.ProjectDetails.AdminPrivileges.ManageMembers = &falseValue
	projectParams.ProjectDetails.AdminPrivileges.ManageResources = &trueValue
	projectParams.ProjectDetails.AdminPrivileges.IndexResources = &falseValue
	assert.NoError(t, testsAccessProjectService.Update(projectParams))
	updatedProject, err := testsAccessProjectService.Get(projectParams.ProjectDetails.ProjectKey)
	if assert.NoError(t, err) &&
		assert.NotNil(t, updatedProject, "Expected project %s but got nil", projectParams.ProjectDetails.ProjectKey) &&
		!reflect.DeepEqual(projectParams.ProjectDetails, *updatedProject) {
		t.Error("Unexpected project details built. Expected: `", projectParams.ProjectDetails, "` Got `", *updatedProject, "`")
	}
}

func deleteProjectAndGroupAndAssert(t *testing.T, projectKey string, groupName string) {
	deleteProjectAndAssert(t, projectKey)
	deleteGroupAndAssert(t, groupName)
}

func deleteProjectAndAssert(t *testing.T, projectKey string) {
	assert.NoError(t, testsAccessProjectService.Delete(projectKey))
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
