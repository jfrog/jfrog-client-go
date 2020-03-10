package services

import (
	"encoding/json"
	"errors"
	"net/http"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	distrbutionServiceUtils "github.com/jfrog/jfrog-client-go/distribution/services/utils"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type CreateReleaseBundleService struct {
	UpdateReleaseBundleService
}

func NewCreateReleseBundleService(client *rthttpclient.ArtifactoryHttpClient) *CreateReleaseBundleService {
	return &CreateReleaseBundleService{UpdateReleaseBundleService{client: client}}
}

func (cbs *CreateReleaseBundleService) GetDistDetails() auth.CommonDetails {
	return cbs.DistDetails
}

func (cbs *CreateReleaseBundleService) CreateReleaseBundle(createBundleParams CreateUpdateReleaseBundleParams) error {
	CreateUpdateReleaseBundleBody, err := CreateBundleBody(createBundleParams, cbs.DryRun)
	if err != nil {
		return err
	}
	createReleaseBundleBody := &CreateReleaseBundleBody{
		Name:              createBundleParams.Name,
		Version:           createBundleParams.Version,
		ReleaseBundleBody: *CreateUpdateReleaseBundleBody,
	}

	return cbs.execCreateReleaseBundle(createReleaseBundleBody)
}

func (cbs *CreateReleaseBundleService) execCreateReleaseBundle(releaseBundle *CreateReleaseBundleBody) error {
	httpClientsDetails := cbs.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(releaseBundle)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := cbs.DistDetails.GetUrl() + "api/v1/release_bundle"
	distrbutionServiceUtils.SetGpgPassphrase(cbs.GpgPassphrase, &httpClientsDetails.Headers)
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := cbs.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Artifactory response: ", resp.Status)
	log.Output(utils.IndentJson(body))
	return errorutils.CheckError(err)
}

type CreateReleaseBundleBody struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	ReleaseBundleBody
}
