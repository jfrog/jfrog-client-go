package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	PermissionTargetNamePrefix = "jfrog-client-go-tests-target"
)

func TestPermissionTarget(t *testing.T) {
	params := services.NewPermissionTargetParams()
	params.Name = fmt.Sprintf("%s-%d", PermissionTargetNamePrefix, time.Now().Unix())
	params.Repo.Repositories = []string{"ANY"}
	params.Repo.ExcludePatterns = []string{"dir/*"}
	params.Repo.Actions.Users = map[string][]string{
		"anonymous": {"read"},
	}
	params.Build.Repositories = []string{"artifactory-build-info"}
	params.Build.Actions.Users = map[string][]string{
		"anonymous": {"annotate"},
	}

	err := testsPermissionTargetService.Create(params)
	assert.NoError(t, err)
	// Fill in default values before validation
	params.Repo.IncludePatterns = []string{"**"}
	params.Build.Repositories = []string{"artifactory-build-info"}
	params.Build.IncludePatterns = []string{"**"}
	params.Build.ExcludePatterns = []string{}
	validatePermissionTarget(t, params)

	params.Repo.Actions.Users = nil
	params.Repo.Repositories = []string{"ANY REMOTE"}
	err = testsPermissionTargetService.Update(params)
	validatePermissionTarget(t, params)
	assert.NoError(t, err)
	err = testsPermissionTargetService.Delete(params.Name)
	assert.NoError(t, err)
}

func validatePermissionTarget(t *testing.T, params services.PermissionTargetParams) {
	targetConfig, err := getPermissionTarget(params.Name)
	assert.NoError(t, err)
	assert.Equal(t, params.Name, targetConfig.Name)
	assert.Equal(t, params.Repo, targetConfig.Repo)
	assert.Equal(t, params.Build, targetConfig.Build)
	return
}

func getPermissionTarget(targetName string) (targetParams *services.PermissionTargetParams, err error) {
	artDetails := GetRtDetails()
	artHttpDetails := artDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return
	}
	resp, body, _, err := client.SendGet(artDetails.GetUrl()+"api/v2/security/permissions/"+targetName, false, artHttpDetails)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	if err = json.Unmarshal(body, &targetParams); err != nil {
		return nil, errors.New("failed unmarshalling permission target " + targetName)
	}
	return
}
