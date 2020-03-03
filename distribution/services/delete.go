package services

import (
	"encoding/json"
	"errors"
	"net/http"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type OnSuccess string

const (
	Keep   OnSuccess = "keep"
	Delete           = "delete"
)

type DeleteDistributionService struct {
	client      *rthttpclient.ArtifactoryHttpClient
	DistDetails auth.CommonDetails
	DryRun      bool
}

func NewDeleteDistributionService(client *rthttpclient.ArtifactoryHttpClient) *DeleteDistributionService {
	return &DeleteDistributionService{client: client}
}

func (ds *DeleteDistributionService) GetDistDetails() auth.CommonDetails {
	return ds.DistDetails
}

func (ds *DeleteDistributionService) DeleteDistribution(deleteDistributionParams DeleteDistributionParams) error {
	var distributionRules []DistributionRulesBody
	for _, rule := range deleteDistributionParams.DistributionRules {
		distributionRule := DistributionRulesBody{
			SiteName:     rule.GetSiteName(),
			CityName:     rule.GetCityName(),
			CountryCodes: rule.GetCountryCodes(),
		}
		distributionRules = append(distributionRules, distributionRule)
	}

	var onSuccess OnSuccess
	if deleteDistributionParams.DeleteFromDistribution {
		onSuccess = Delete
	} else {
		onSuccess = Keep
	}

	deleteDistribution := DeleteDistributionBody{
		DistributionBody: DistributionBody{
			DryRun:            ds.DryRun,
			DistributionRules: distributionRules,
		},
		OnSuccess: onSuccess,
	}

	return ds.execDeleteDistribute(deleteDistributionParams.Name, deleteDistributionParams.Version, deleteDistribution)
}

func (cbs *DeleteDistributionService) execDeleteDistribute(name, version string, deleteDistribution DeleteDistributionBody) error {
	httpClientsDetails := cbs.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(deleteDistribution)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := cbs.DistDetails.GetUrl() + "api/v1/distribution/" + name + "/" + version + "/delete"
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := cbs.client.SendPost(url, content, &httpClientsDetails)
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

type DeleteDistributionBody struct {
	DistributionBody
	OnSuccess OnSuccess `json:"on_success"`
}

type DeleteDistributionParams struct {
	DistributionParams
	DeleteFromDistribution bool
}

func NewDeleteDistributionParams(name, version string) DeleteDistributionParams {
	return DeleteDistributionParams{
		DistributionParams: DistributionParams{
			Name:    name,
			Version: version,
		},
	}
}
