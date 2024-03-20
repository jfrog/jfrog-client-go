package tests

import (
	"github.com/jfrog/jfrog-client-go/utils"
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
	projectParams := getTestProjectParams("tstprj", "testProject")
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
	projectParams := getTestProjectParams("tstprj", "testProject")
	assert.NoError(t, testsAccessProjectService.Create(projectParams))
	defer deleteProjectAndAssert(t, projectParams.ProjectDetails.ProjectKey)
	projectParams.ProjectDetails.Description += "123"
	projectParams.ProjectDetails.StorageQuotaBytes += 123
	projectParams.ProjectDetails.SoftLimit = utils.Pointer(true)
	projectParams.ProjectDetails.AdminPrivileges.ManageMembers = utils.Pointer(false)
	projectParams.ProjectDetails.AdminPrivileges.ManageResources = utils.Pointer(true)
	projectParams.ProjectDetails.AdminPrivileges.IndexResources = utils.Pointer(false)
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

func getTestProjectParams(projectKey string, projectName string) services.ProjectParams {
	adminPrivileges := services.AdminPrivileges{
		ManageMembers:   utils.Pointer(true),
		ManageResources: utils.Pointer(false),
		IndexResources:  utils.Pointer(true),
	}
	runId := getRunId()
	runNumberSuffix := runId[len(runId)-3:]
	projectDetails := services.Project{
		DisplayName:       projectName + runNumberSuffix,
		Description:       "My Test Project",
		AdminPrivileges:   &adminPrivileges,
		SoftLimit:         utils.Pointer(false),
		StorageQuotaBytes: 1073741825,                   // Needs to be higher than 1073741824
		ProjectKey:        projectKey + runNumberSuffix, // Valid length: 2 <= ProjectKey <= 10
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

func TestGetAllProjects(t *testing.T) {
	initAccessTest(t)
	preProjects, err := testsAccessProjectService.GetAll()
	assert.NoError(t, err)
	oldNumOfPrjs := len(preProjects)
	params1 := getTestProjectParams("tstprj", "projectForTest")
	params2 := getTestProjectParams("tstprj1", "projectTesting")
	params3 := getTestProjectParams("tstprj2", "testProject")
	params4 := getTestProjectParams("tstprj3", "It'sForTest")
	assert.NoError(t, testsAccessProjectService.Create(params1))
	assert.NoError(t, testsAccessProjectService.Create(params2))
	assert.NoError(t, testsAccessProjectService.Create(params3))
	assert.NoError(t, testsAccessProjectService.Create(params4))
	projects, err := testsAccessProjectService.GetAll()
	assert.NoError(t, err, "Failed to Unmarshal")
	assert.Equal(t, oldNumOfPrjs+4, len(projects))
	deleteProjectAndAssert(t, params1.ProjectDetails.ProjectKey)
	deleteProjectAndAssert(t, params2.ProjectDetails.ProjectKey)
	deleteProjectAndAssert(t, params3.ProjectDetails.ProjectKey)
	deleteProjectAndAssert(t, params4.ProjectDetails.ProjectKey)
}
