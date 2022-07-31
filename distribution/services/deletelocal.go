package services

import (
	"net/http"

	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// Delete distributable release bundle from the distribution service.
type DeleteLocalReleaseBundleService struct {
	DeleteReleaseBundleService
}

func NewDeleteLocalDistributionService(client *jfroghttpclient.JfrogHttpClient) *DeleteLocalReleaseBundleService {
	return &DeleteLocalReleaseBundleService{DeleteReleaseBundleService: DeleteReleaseBundleService{client: client}}
}

func (dlr *DeleteLocalReleaseBundleService) GetDistDetails() auth.ServiceDetails {
	return dlr.DistDetails
}

func (dlr *DeleteLocalReleaseBundleService) DeleteDistribution(deleteDistributionParams DeleteDistributionParams) error {
	dlr.Sync = deleteDistributionParams.Sync
	dlr.MaxWaitMinutes = deleteDistributionParams.MaxWaitMinutes
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
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusNoContent); err != nil {
		return err
	}
	if dlr.Sync {
		err := dlr.waitForDeletion(name, version)
		if err != nil {
			return err
		}
	}

	log.Debug("Distribution response: ", resp.Status)
	return errorutils.CheckError(err)
}
