package tests

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/pipelines/services"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPipelinesWorkspaceService(t *testing.T) {
	t.Run("test trigger pipelines workspace when resources are not valid", TestWorkspaceValidationWhenPipelinesResourcesAreNotValid)
	t.Run("test trigger pipelines workspace when resources are valid", TestWorkspaceValidationWhenPipelinesResourcesAreValid)
	t.Run("test workspace sync when it fails for the first time and succeeds later", TestWorkspaceValidationFailureAndSucceedsWhenValidIntegrationCreated)
}

func TestWorkspaceValidationWhenPipelinesResourcesAreNotValid(t *testing.T) {

	if !assert.NotEmpty(t, *PipelinesAccessToken, "cannot run pipelines tests without access token configured") {
		return
	}
	filePath := "testdata/pipelines.yml"

	// get pipelines.yaml content
	pipelineRes, err := getWorkspaceRunPayload([]string{filePath})
	if !assert.NoError(t, err) {
		return
	}
	fmt.Printf("Bytes received : %s \n", pipelineRes)
	if !assert.NoError(t, err) {
		return
	}
	err = testPipelinesWorkspaceService.ValidateWorkspace(pipelineRes)
	assert.NoError(t, err)
}

func TestWorkspaceValidationWhenPipelinesResourcesAreValid(t *testing.T) {
	if !assert.NotEmpty(t, *PipelinesAccessToken, "cannot run pipelines tests without access token configured") {
		return
	}
	filePath := "testdata/pipelines.yml"

	// get pipelines.yaml content
	pipelineRes, err := getWorkspaceRunPayload([]string{filePath})
	if !assert.NoError(t, err) {
		return
	}
	err = testPipelinesWorkspaceService.ValidateWorkspace(pipelineRes)
	assert.NoError(t, err)

	wsResp, err := testPipelinesWorkspaceService.GetWorkspace()
	assert.NoError(t, err)
	if len(wsResp) < 1 {
		assert.Fail(t, "No workspace created")
	}
	syncStatusResp, err := testPipelinesWorkspaceService.WorkspacePollSyncStatus()
	assert.NoError(t, err)
	for _, ws := range syncStatusResp {
		fmt.Printf("%+v \n", ws)
		err = testPipelinesWorkspaceService.DeleteWorkspace("default")
	}
	assert.NoError(t, err)
}

func TestWorkspaceValidationFailureAndSucceedsWhenValidIntegrationCreated(t *testing.T) {
	if !assert.NotEmpty(t, *PipelinesAccessToken, "cannot run pipelines tests without access token configured") {
		return
	}
	filePath := "testdata/pipelines-integration.yml"

	// Get pipelines-integration.yaml content
	pipelineRes, err := getWorkspaceRunPayload([]string{filePath})
	if !assert.NoError(t, err) {
		return
	}
	// Call workspace validation
	err = testPipelinesWorkspaceService.ValidateWorkspace(pipelineRes)
	assert.Error(t, err)

	// Create valid integration required for pipelines
	id, err := testsPipelinesIntegrationsService.CreateArtifactoryIntegration("int_workspace_artifactory", testsDummyRtUrl, testsDummyUser, testsDummyApiKey)
	if !assert.NoError(t, err) {
		return
	}
	defer deleteIntegrationAndAssert(t, id)
	getIntegrationAndAssert(t, id, "int_workspace_artifactory", services.ArtifactoryName)

	// Validate workspace pipelines should be successful here
	err = testPipelinesWorkspaceService.ValidateWorkspace(pipelineRes)
	assert.NoError(t, err)

	wsResp, err := testPipelinesWorkspaceService.GetWorkspace()
	assert.NoError(t, err)
	if len(wsResp) < 1 {
		assert.Fail(t, "No workspace created")
	}
	_, err = testPipelinesWorkspaceService.WorkspacePollSyncStatus()
	if !assert.NoError(t, err) {
		return
	}

	err = testPipelinesWorkspaceService.WorkspaceSync("default")
	assert.NoError(t, err)
	syncStatusRespSuccess, err := testPipelinesWorkspaceService.WorkspacePollSyncStatus()
	assert.NoError(t, err)
	for _, ws := range syncStatusRespSuccess {
		fmt.Printf("%+v \n", ws)
		err = testPipelinesWorkspaceService.DeleteWorkspace("default")
	}
	assert.NoError(t, err)
}

type PipelineDefinition struct {
	FileName string `json:"fileName,omitempty"`
	Content  string `json:"content,omitempty"`
	YmlType  string `json:"ymlType,omitempty"`
}

type WorkSpaceValidation struct {
	ProjectId   string               `json:"-"`
	Files       []PipelineDefinition `json:"files,omitempty"`
	ProjectName string               `json:"projectName,omitempty"`
	Name        string               `json:"name,omitempty"`
}

func getWorkspaceRunPayload(resources []string) ([]byte, error) {
	var pipelineDefinitions []PipelineDefinition
	for _, pathToFile := range resources {
		fileContent, fileInfo, err := getFileContentAndBaseName(pathToFile)
		if err != nil {
			return nil, err
		}
		pipeDefinition := PipelineDefinition{
			FileName: fileInfo.Name(),
			Content:  string(fileContent),
			YmlType:  "pipelines",
		}
		pipelineDefinitions = append(pipelineDefinitions, pipeDefinition)
	}
	workSpaceValidation := WorkSpaceValidation{
		Files:       pipelineDefinitions,
		ProjectName: "default",
		Name:        "",
	}
	return json.Marshal(workSpaceValidation)
}

func getFileContentAndBaseName(pathToFile string) ([]byte, os.FileInfo, error) {
	fileContent, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, nil, err
	}
	fileInfo, err := os.Stat(pathToFile)
	if err != nil {
		return nil, nil, err
	}
	return fileContent, fileInfo, nil
}
