package services

import (
	"encoding/json"
	"errors"
	"net/http"

	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	distributionServiceUtils "github.com/jfrog/jfrog-client-go/distribution/services/utils"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type CreateReleaseBundleService struct {
	UpdateReleaseBundleService
}

func NewCreateReleaseBundleService(client *jfroghttpclient.JfrogHttpClient) *CreateReleaseBundleService {
	return &CreateReleaseBundleService{UpdateReleaseBundleService{client: client}}
}

func (cb *CreateReleaseBundleService) GetDistDetails() auth.ServiceDetails {
	return cb.DistDetails
}

func (cb *CreateReleaseBundleService) CreateReleaseBundle(createBundleParams CreateReleaseBundleParams) error {
	releaseBundleBody, err := distributionServiceUtils.CreateBundleBody(createBundleParams.ReleaseBundleParams, cb.DryRun)
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

func (cb *CreateReleaseBundleService) execCreateReleaseBundle(gpgPassphrase string, releaseBundle *createReleaseBundleBody) error {
	httpClientsDetails := cb.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(releaseBundle)
	if err != nil {
		return errorutils.CheckError(err)
	}
	dryRunStr := ""
	if releaseBundle.DryRun {
		dryRunStr = "[Dry run] "
	}
	log.Info(dryRunStr + "Creating: " + releaseBundle.Name + "/" + releaseBundle.Version)

	url := cb.DistDetails.GetUrl() + "api/v1/release_bundle"
	distributionServiceUtils.AddGpgPassphraseHeader(gpgPassphrase, &httpClientsDetails.Headers)
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := cb.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if !(resp.StatusCode == http.StatusCreated || (resp.StatusCode == http.StatusOK && releaseBundle.DryRun)) {
		return errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Distribution response: ", resp.Status)
	log.Debug(utils.IndentJson(body))
	return errorutils.CheckError(err)
}

type createReleaseBundleBody struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	distributionServiceUtils.ReleaseBundleBody
}

type CreateReleaseBundleParams struct {
	distributionServiceUtils.ReleaseBundleParams
}

func NewCreateReleaseBundleParams(name, version string) CreateReleaseBundleParams {
	return CreateReleaseBundleParams{
		distributionServiceUtils.ReleaseBundleParams{
			Name:    name,
			Version: version,
		},
	}
}
