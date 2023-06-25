package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestPipelinesWorkspaceService(t *testing.T) {
	t.Run("test trigger pipelines workspace when resources are not valid", TestWorkspaceValidationWhenPipelinesResourcesAreNotValid)
	t.Run("test trigger pipelines workspace when resources are valid", TestWorkspaceValidationWhenPipelinesResourcesAreValid)
}

func TestWorkspaceValidationWhenPipelinesResourcesAreNotValid(t *testing.T) {

	if !assert.NotEmpty(t, *PipelinesAccessToken, "cannot run pipelines tests without access token configured") {
		return
	}
	filePath := "/Users/bhanur/go/src/pipeline-cli/jfrog-client-go/tests/testdata/pipelines.yml"

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
	filePath := "/Users/bhanur/go/src/pipeline-cli/jfrog-client-go/tests/testdata/pipelines.yml"

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
	for _, ws := range wsResp {
		fmt.Printf("%+v \n", ws)
		err = testPipelinesWorkspaceService.DeleteWorkspace(strconv.Itoa(ws.ID))
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
	// need to define handling values
	workSpaceValidation := WorkSpaceValidation{
		ProjectId:   "1",
		Files:       pipelineDefinitions,
		ProjectName: "default",
		Name:        "bhanu",
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
