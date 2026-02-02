//go:build itest

package tests

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"sort"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"

	"github.com/jfrog/jfrog-client-go/access/services"
	rtservices "github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	projectKey := createRandomProject(t).ProjectDetails.ProjectKey

	testGroup := getTestProjectGroupParams("a-test-group")

	createGroup(t, testGroup.Name, true, false)

	require.NoError(t, testsAccessProjectService.UpdateGroup(projectKey, testGroup.Name, testGroup))

	allGroups, err := testsAccessProjectService.GetGroups(projectKey)
	if assert.NoError(t, err) &&
		assert.NotNil(t, allGroups, "Expected 1 group in the project but got 0") {
		assert.Equal(t, len(*allGroups), 1, "Expected 1 group in the project but got %d", len(*allGroups))
		assert.Contains(t, *allGroups, testGroup)
	}

	testGroup.Roles = append(testGroup.Roles, "Viewer")
	assert.NoError(t, testsAccessProjectService.UpdateGroup(projectKey, testGroup.Name, testGroup))
	// Sort roles for comparison
	sort.Slice(testGroup.Roles, func(i, j int) bool {
		return testGroup.Roles[i] < testGroup.Roles[j]
	})
	singleGroup, err := testsAccessProjectService.GetGroup(projectKey, testGroup.Name)
	// Sort roles for comparison
	if singleGroup != nil {
		sort.Slice(singleGroup.Roles, func(i, j int) bool {
			return singleGroup.Roles[i] < singleGroup.Roles[j]
		})
	}
	if assert.NoError(t, err) &&
		assert.NotNil(t, singleGroup, "Expected group %s but got nil", testGroup.Name) {
		assert.Equal(t, testGroup, *singleGroup, "Expected group %v but got %v", testGroup, *singleGroup)
	}

	assert.NoError(t, testsAccessProjectService.DeleteExistingGroup(projectKey, testGroup.Name))

	noGroups, err := testsAccessProjectService.GetGroups(projectKey)
	require.NoError(t, err)
	assert.Empty(t, noGroups)
}

func testAccessProjectCreateUpdateDelete(t *testing.T) {
	projectParams := createRandomProject(t)
	originalDetails := projectParams.ProjectDetails
	projectParams.ProjectDetails.Description += "123"
	projectParams.ProjectDetails.StorageQuotaBytes += 123
	projectParams.ProjectDetails.SoftLimit = utils.Pointer(true)
	projectParams.ProjectDetails.AdminPrivileges.ManageMembers = utils.Pointer(false)
	projectParams.ProjectDetails.AdminPrivileges.ManageResources = utils.Pointer(true)
	projectParams.ProjectDetails.AdminPrivileges.IndexResources = utils.Pointer(false)
	require.NoError(t, testsAccessProjectService.Update(projectParams))
	updatedProject, err := testsAccessProjectService.Get(projectParams.ProjectDetails.ProjectKey)
	require.NoError(t, err)
	require.NotNilf(t, updatedProject, "Expected project %s but got nil", projectParams.ProjectDetails.ProjectKey)
	assert.Falsef(t, reflect.DeepEqual(originalDetails, *updatedProject), "Expected project %v not to be %v", projectParams.ProjectDetails, *updatedProject)
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

func createGroup(t *testing.T, groupName string, includeUsers bool, adminPrivileges bool) rtservices.GroupParams {
	groupParams := getTestGroupParams(includeUsers)
	if groupName != "" {
		groupParams.GroupDetails.Name = groupName
	}
	groupParams.GroupDetails.AdminPrivileges = utils.Pointer(adminPrivileges)
	require.NoError(t, testGroupService.CreateGroup(groupParams))
	t.Cleanup(func() {
		err := testGroupService.DeleteGroup(groupParams.GroupDetails.Name)
		if err != nil {
			log.Warn(fmt.Sprintf("Failed to delete group %s: %+v", groupName, err))
		}
	})
	return groupParams
}

func createRandomProject(t *testing.T) services.ProjectParams {
	projectKey := fmt.Sprintf("tstprj%d%s", time.Now().Unix(), randomString(t, 6))
	return createProject(t, projectKey, projectKey)
}

func createProject(t *testing.T, projectKey string, projectName string) services.ProjectParams {
	projectParams := getTestProjectParams(projectKey, projectName)
	require.NoError(t, testsAccessProjectService.Create(projectParams))
	t.Cleanup(func() {
		err := testsAccessProjectService.Delete(projectKey)
		if err != nil {
			log.Warn(fmt.Sprintf("Failed to delete project %s: %+v", projectKey, err))
		}
	})
	return projectParams
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
	require.NoError(t, err)
	oldNumOfPrjs := len(preProjects)

	createRandomProject(t)
	createRandomProject(t)
	createRandomProject(t)
	createRandomProject(t)

	projects, err := testsAccessProjectService.GetAll()
	require.NoError(t, err, "Failed to Unmarshal")
	assert.Equal(t, oldNumOfPrjs+4, len(projects))
}
