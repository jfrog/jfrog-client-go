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
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type UpdateReleaseBundleService struct {
	client      *jfroghttpclient.JfrogHttpClient
	DistDetails auth.ServiceDetails
	DryRun      bool
}

func NewUpdateReleaseBundleService(client *jfroghttpclient.JfrogHttpClient) *UpdateReleaseBundleService {
	return &UpdateReleaseBundleService{client: client}
}

func (ur *UpdateReleaseBundleService) GetDistDetails() auth.ServiceDetails {
	return ur.DistDetails
}

func (ur *UpdateReleaseBundleService) UpdateReleaseBundle(createBundleParams UpdateReleaseBundleParams) (*clientutils.Sha256Summary, error) {
	releaseBundleBody, err := distributionServiceUtils.CreateBundleBody(createBundleParams.ReleaseBundleParams, ur.DryRun)
	if err != nil {
		return nil, err
	}
	summary, err := ur.execUpdateReleaseBundle(createBundleParams.Name, createBundleParams.Version, createBundleParams.GpgPassphrase, releaseBundleBody)
	if createBundleParams.SignImmediately {
		return summary, err
	}
	return nil, err
}

func (ur *UpdateReleaseBundleService) execUpdateReleaseBundle(name, version, gpgPassphrase string, releaseBundle *distributionServiceUtils.ReleaseBundleBody) (*clientutils.Sha256Summary, error) {
	summary := clientutils.NewSha256Summary()
	httpClientsDetails := ur.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(releaseBundle)
	if err != nil {
		return summary, errorutils.CheckError(err)
	}

	dryRunStr := ""
	if releaseBundle.DryRun {
		dryRunStr = "[Dry run] "
	}
	log.Info(dryRunStr + "Updating: " + name + "/" + version)

	url := ur.DistDetails.GetUrl() + "api/v1/release_bundle/" + name + "/" + version
	distributionServiceUtils.AddGpgPassphraseHeader(gpgPassphrase, &httpClientsDetails.Headers)
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := ur.client.SendPut(url, content, &httpClientsDetails)
	if err != nil {
		return summary, err
	}
	if resp.StatusCode != http.StatusOK {
		return summary, errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	summary.SetSucceeded(true)
	summary.SetSha256(resp.Header.Get("X-Checksum-Sha256"))

	log.Debug("Distribution response: ", resp.Status)
	log.Debug(utils.IndentJson(body))
	return summary, nil
}

type UpdateReleaseBundleParams struct {
	distributionServiceUtils.ReleaseBundleParams
}

func NewUpdateReleaseBundleParams(name, version string) UpdateReleaseBundleParams {
	return UpdateReleaseBundleParams{
		distributionServiceUtils.ReleaseBundleParams{
			Name:    name,
			Version: version,
		},
	}
}
