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
type DeleteLocalDistributionService struct {
	client      *rthttpclient.ArtifactoryHttpClient
	DistDetails auth.CommonDetails
}

func NewDeleteLocalDistributionService(client *rthttpclient.ArtifactoryHttpClient) *DeleteLocalDistributionService {
	return &DeleteLocalDistributionService{client: client}
}

func (ds *DeleteLocalDistributionService) GetDistDetails() auth.CommonDetails {
	return ds.DistDetails
}

func (ds *DeleteLocalDistributionService) DeleteDistribution(deleteDistributionParams DeleteLocalDistributionParams) error {
	return ds.execDeleteLocalDistribution(deleteDistributionParams.Name, deleteDistributionParams.Version)
}

func (cbs *DeleteLocalDistributionService) execDeleteLocalDistribution(name, version string) error {
	httpClientsDetails := cbs.DistDetails.CreateHttpClientDetails()
	url := cbs.DistDetails.GetUrl() + "api/v1/distribution/" + name + "/" + version
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := cbs.client.SendDelete(url, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Distribution response: ", resp.Status)
	log.Output(utils.IndentJson(body))
	return errorutils.CheckError(err)
}

type DeleteLocalDistributionParams struct {
	Name    string
	Version string
}

func NewDeleteLocalDistributionParams(name, version string) DeleteLocalDistributionParams {
	return DeleteLocalDistributionParams{
		Name:    name,
		Version: version,
	}
}
