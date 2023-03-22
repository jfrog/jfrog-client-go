package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func testWorkspaceValidationWhenPipelinesResourcesAreNotValid(t *testing.T) {

	if !assert.NotEmpty(t, *PipelinesAccessToken, "cannot run pipelines tests without access token configured") {
		return
	}
	filePath := "/Users/bhanur/go/src/pipeline-cli/jfrog-client-go/tests/testdata/pipelines.yml"

	// get pipelines.yaml content
	pipelineRes, err := getPipelinesResourceToValidate(filePath)
	fmt.Printf("Bytes received : %s \n", pipelineRes)
	if !assert.NoError(t, err) {
		return
	}
	err = testPipelinesWorkspaceService.ValidateWorkspace(pipelineRes)
	assert.NoError(t, err)
}

func testWorkspaceValidationWhenPipelinesResourcesAreValid(t *testing.T) {
	if !assert.NotEmpty(t, *PipelinesAccessToken, "cannot run pipelines tests without access token configured") {
		return
	}
	filePath := "/Users/bhanur/go/src/pipeline-cli/jfrog-client-go/tests/testdata/pipelines.yml"

	integrationName := "int_gh_pipe_cli"
	_, err := testsPipelinesIntegrationsService.CreateGithubIntegration(integrationName, *PipelinesVcsToken)
	assert.NoError(t, err)
	//defer deleteIntegrationAndAssert(t, integrationId)

	// get pipelines.yaml content
	pipelineRes, err := getPipelinesResourceToValidate(filePath)
	if !assert.NoError(t, err) {
		return
	}
	err = testPipelinesWorkspaceService.ValidateWorkspace(pipelineRes)
	assert.NoError(t, err)

	workspaces, wsErr := testPipelinesWorkspaceService.WorkspacePollSyncStatus()
	assert.NoError(t, wsErr)
	pipelineBranch, err := testPipelinesWorkspaceService.GetWorkspacePipelines(workspaces)
	assert.NoError(t, err)
	for pipName, branch := range pipelineBranch {
		err = testPipelinesRunService.TriggerPipelineRun(branch, pipName, false)
		assert.NoError(t, err)
	}
	pipelineNames := make([]string, len(pipelineBranch))

	i := 0
	for k := range pipelineBranch {
		pipelineNames[i] = k
		i++
	}
	pipeRunIDs, wsRunErr := testPipelinesWorkspaceService.WorkspaceRunIDs(pipelineNames)
	assert.NoError(t, wsRunErr)

	for _, runId := range pipeRunIDs {
		_, err2 := testPipelinesWorkspaceService.WorkspaceRunStatus(runId.LatestRunID)
		assert.NoError(t, err2)
		_, err3 := testPipelinesWorkspaceService.WorkspaceStepStatus(runId.LatestRunID)
		assert.NoError(t, err3)
	}

}

func getPipelinesResourceToValidate(filePath string) ([]byte, error) {

	ymlType := ""
	readFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	toJSON, err3 := convertYAMLToJSON(err, readFile)
	if err3 != nil {
		return nil, err3
	}
	vsc, err4 := convertJSONDataToMap(fInfo, toJSON)
	if err4 != nil {
		return nil, err4
	}
	var marErr error

	resMap, err5, done := splitDataToPipelinesAndResourcesMap(vsc, marErr, err, ymlType)
	if done {
		return nil, err5
	}

	data, resErr := getPayloadToValidatePipelineResource(resMap)
	if resErr != nil {
		return nil, resErr
	}

	return data.Bytes(), nil
}

func convertJSONDataToMap(file os.FileInfo, toJSON []byte) (map[string][]interface{}, error) {
	log.Info("validating pipeline resources ", file.Name())
	time.Sleep(1 * time.Second)
	vsc := make(map[string][]interface{})
	convErr := yaml.Unmarshal(toJSON, &vsc)
	if convErr != nil {
		return nil, convErr
	}
	return vsc, nil
}

func convertYAMLToJSON(err error, readFile []byte) ([]byte, error) {
	toJSON, err := yaml.YAMLToJSON(readFile)
	if err != nil {
		log.Error("Failed to convert to json")
		return nil, err
	}
	return toJSON, nil
}

func splitDataToPipelinesAndResourcesMap(vsc map[string][]interface{}, marErr error, err error, ymlType string) (map[string]string, error, bool) {
	resMap := make(map[string]string)
	var data []byte
	if v, ok := vsc["resources"]; ok {
		log.Info("resources found preparing to validate")
		data, marErr = json.Marshal(v)
		if marErr != nil {
			log.Error("failed to marshal to json")
			return nil, err, true
		}

		ymlType = "resources"
		resMap[ymlType] = string(data)

	}
	if vp, ok := vsc["pipelines"]; ok {
		log.Info("pipelines found preparing to validate")
		data, marErr = json.Marshal(vp)
		if marErr != nil {
			log.Error("failed to marshal to json")
			return nil, err, true
		}
		ymlType = "pipelines"
		fmt.Println(string(data))
		resMap[ymlType] = string(data)
	}
	//fmt.Printf("%+v \n", resMap)
	return resMap, nil, false
}

func getPayloadToValidatePipelineResource(resMap map[string]string) (*bytes.Buffer, error) {
	payload := getPayloadBasedOnYmlType(resMap)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(payload)
	if err != nil {
		log.Error("Failed to read stream to send payload to trigger pipelines")
		return nil, err
	}
	return buf, err
}

func getPayloadBasedOnYmlType(m map[string]string) *strings.Reader {
	var resReader, pipReader, valReader *strings.Reader
	for ymlType, _ := range m {
		if ymlType == "resources" {
			resReader = strings.NewReader(`{"fileName":"` + ymlType + `.yml","content": ` + m[ymlType] + `,"ymlType":"` + ymlType + `"}`)
		} else if ymlType == "pipelines" {
			//fmt.Printf("data : %+v \n", m[ymlType])
			pipReader = strings.NewReader(`{"fileName":"` + ymlType + `.yml","content":` + m[ymlType] + `,"ymlType":"` + ymlType + `"}`)
		}
	}
	if resReader != nil && pipReader != nil {
		resAll, err := io.ReadAll(resReader)
		if err != nil {
			return nil
		}
		pipAll, err := io.ReadAll(pipReader)
		if err != nil {
			return nil
		}
		valReader = strings.NewReader(`{"files":[` + string(resAll) + `,` + string(pipAll) + `]}`)
	} else if resReader != nil {
		all, err := io.ReadAll(resReader)
		if err != nil {
			return nil
		}
		valReader = strings.NewReader(`{"files":[` + string(all) + `]}`)
	} else if pipReader != nil {
		all, err := io.ReadAll(pipReader)
		if err != nil {
			return nil
		}
		valReader = strings.NewReader(`{"files":[` + string(all) + `]}`)
	}
	return valReader
}
