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

func (cb *CreateReleaseBundleService) GetDistDetails() auth.ServiceDetails {
	return cb.DistDetails
}

func (cb *CreateReleaseBundleService) CreateReleaseBundle(createBundleParams CreateReleaseBundleParams) error {
	releaseBundleBody, err := distrbutionServiceUtils.CreateBundleBody(createBundleParams.ReleaseBundleParams, cb.DryRun)
	if err != nil {
		return err
	}

	body := &createReleaseBundleBody{
		Name:              createBundleParams.Name,
		Version:           createBundleParams.Version,
		ReleaseBundleBody: *releaseBundleBody,
	}

	return cb.execCreateReleaseBundle(createBundleParams.GpgPassphrase, body)
}

func (cb *CreateReleaseBundleService) execCreateReleaseBundle(gpgPassprase string, releaseBundle *createReleaseBundleBody) error {
	httpClientsDetails := cb.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(releaseBundle)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := cb.DistDetails.GetUrl() + "api/v1/release_bundle"
	distrbutionServiceUtils.AddGpgPassphraseHeader(gpgPassprase, &httpClientsDetails.Headers)
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := cb.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Distribution response: ", resp.Status)
	log.Debug(utils.IndentJson(body))
	return errorutils.CheckError(err)
}

type createReleaseBundleBody struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	distrbutionServiceUtils.ReleaseBundleBody
}

type CreateReleaseBundleParams struct {
	distrbutionServiceUtils.ReleaseBundleParams
}

func NewCreateReleaseBundleParams(name, version string) CreateReleaseBundleParams {
	return CreateReleaseBundleParams{
		distrbutionServiceUtils.ReleaseBundleParams{
			Name:    name,
			Version: version,
		},
	}
}
