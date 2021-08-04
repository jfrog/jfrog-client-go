package tests

import (
	"github.com/jfrog/jfrog-client-go/access/services"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAccessProject(t *testing.T) {
	initAccessTest(t)
	t.Run("create", testAccessCreateProject)
	t.Run("update", testAccessUpdateProject)
	t.Run("delete", testAccessDeleteProject)
}

func testAccessCreateProject(t *testing.T) {
	projectParams := getTestProjectParams()
	err := testsAccessProjectService.CreateProject(projectParams)
	defer deleteProjectAndAssert(t, projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	createdProject, err := testsAccessProjectService.GetProject(projectParams.ProjectDetails.ProjectKey)
	assert.NotNil(t, createdProject)
	if !reflect.DeepEqual(projectParams.ProjectDetails, *createdProject) {
		t.Error("Unexpected project details built. Expected: `", projectParams.ProjectDetails, "` Got `", *createdProject, "`")
	}
}

func testAccessDeleteProject(t *testing.T) {
	projectParams := getTestProjectParams()
	err := testsAccessProjectService.CreateProject(projectParams)
	assert.NoError(t, err)
	err = testsAccessProjectService.DeleteProject(projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	updatedProject, err := testsAccessProjectService.GetProject(projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	assert.Nil(t, updatedProject)
}

func testAccessUpdateProject(t *testing.T) {
	projectParams := getTestProjectParams()
	err := testsAccessProjectService.CreateProject(projectParams)
	defer deleteProjectAndAssert(t, projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	projectParams.ProjectDetails.Description += "123"
	projectParams.ProjectDetails.StorageQuotaBytes += 123
	err = testsAccessProjectService.UpdateProject(projectParams)
	assert.NoError(t, err)
	updatedProject, err := testsAccessProjectService.GetProject(projectParams.ProjectDetails.ProjectKey)
	assert.NoError(t, err)
	if !reflect.DeepEqual(projectParams.ProjectDetails, *updatedProject) {
		t.Error("Unexpected project details built. Expected: `", projectParams.ProjectDetails, "` Got `", *updatedProject, "`")
	}
}

func deleteProjectAndAssert(t *testing.T, projectKey string) {
	err := testsAccessProjectService.DeleteProject(projectKey)
	assert.NoError(t, err)
}

func getTestProjectParams() services.ProjectParams {
	adminPriviligies := services.AdminPrivileges{
		ManageMembers:   true,
		ManageResources: true,
		IndexResources:  true,
	}
	projectDetails := services.Project{
		DisplayName:       "testProject",
		Description:       "My Test Project",
		AdminPrivileges:   &adminPriviligies,
		SoftLimit:         false,
		StorageQuotaBytes: 1073741825, // needs to be higher than 1073741824
		ProjectKey:        "tstprj",
	}
	return services.ProjectParams{
		ProjectDetails: projectDetails,
	}
}
