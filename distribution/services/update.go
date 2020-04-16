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

type UpdateReleaseBundleService struct {
	client      *rthttpclient.ArtifactoryHttpClient
	DistDetails auth.ServiceDetails
	DryRun      bool
}

func NewUpdateReleaseBundleService(client *rthttpclient.ArtifactoryHttpClient) *UpdateReleaseBundleService {
	return &UpdateReleaseBundleService{client: client}
}

func (ur *UpdateReleaseBundleService) GetDistDetails() auth.ServiceDetails {
	return ur.DistDetails
}

func (ur *UpdateReleaseBundleService) UpdateReleaseBundle(createBundleParams UpdateReleaseBundleParams) error {
	releaseBundleBody, err := distrbutionServiceUtils.CreateBundleBody(createBundleParams.ReleaseBundleParams, ur.DryRun)
	if err != nil {
		return err
	}
	return ur.execUpdateReleaseBundle(createBundleParams.Name, createBundleParams.Version, createBundleParams.GpgPassphrase, releaseBundleBody)
}

func (ur *UpdateReleaseBundleService) execUpdateReleaseBundle(name, version, gpgPassphrase string, releaseBundle *distrbutionServiceUtils.ReleaseBundleBody) error {
	httpClientsDetails := ur.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(releaseBundle)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := ur.DistDetails.GetUrl() + "api/v1/release_bundle/" + name + "/" + version
	distrbutionServiceUtils.AddGpgPassphraseHeader(gpgPassphrase, &httpClientsDetails.Headers)
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := ur.client.SendPut(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Distribution response: ", resp.Status)
	return errorutils.CheckError(err)
}

type UpdateReleaseBundleParams struct {
	distrbutionServiceUtils.ReleaseBundleParams
}

func NewUpdateReleaseBundleParams(name, version string) UpdateReleaseBundleParams {
	return UpdateReleaseBundleParams{
		distrbutionServiceUtils.ReleaseBundleParams{
			Name:    name,
			Version: version,
		},
	}
}
