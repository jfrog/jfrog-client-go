package services

import (
	"errors"
	"net/http"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// Delete distributable release bundle from the distribution service.
type DeleteLocalReleaseBundleService struct {
	DeleteReleaseBundleService
}

func NewDeleteLocalDistributionService(client *rthttpclient.ArtifactoryHttpClient) *DeleteLocalReleaseBundleService {
	return &DeleteLocalReleaseBundleService{DeleteReleaseBundleService: DeleteReleaseBundleService{client: client}}
}

func (dlr *DeleteLocalReleaseBundleService) GetDistDetails() auth.ServiceDetails {
	return dlr.DistDetails
}

func (dlr *DeleteLocalReleaseBundleService) DeleteDistribution(deleteDistributionParams DeleteDistributionParams) error {
	return dlr.execDeleteLocalDistribution(deleteDistributionParams.Name, deleteDistributionParams.Version)
}

func (dlr *DeleteLocalReleaseBundleService) execDeleteLocalDistribution(name, version string) error {
	if dlr.IsDryRun() {
		log.Info("[Dry run] Deleting release bundle:", name, version)
		return nil
	}
	log.Info("Deleting release bundle:", name, version)
	httpClientsDetails := dlr.DistDetails.CreateHttpClientDetails()
	url := dlr.DistDetails.GetUrl() + "api/v1/release_bundle/" + name + "/" + version
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := dlr.client.SendDelete(url, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Distribution response: ", resp.Status)
	return errorutils.CheckError(err)
}
