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

// Delete received release bundles from the edge nodes. On success, keep or delete the release bundle from the distribution service.
type DeleteReleaseBundleService struct {
	client      *rthttpclient.ArtifactoryHttpClient
	DistDetails auth.ServiceDetails
	DryRun      bool
}

func NewDeleteReleaseBundleService(client *rthttpclient.ArtifactoryHttpClient) *DeleteReleaseBundleService {
	return &DeleteReleaseBundleService{client: client}
}

func (dr *DeleteReleaseBundleService) GetDistDetails() auth.ServiceDetails {
	return dr.DistDetails
}

func (dr *DeleteReleaseBundleService) IsDryRun() bool {
	return dr.DryRun
}

func (dr *DeleteReleaseBundleService) DeleteDistribution(deleteDistributionParams DeleteDistributionParams) error {
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

	deleteDistribution := DeleteRemoteDistributionBody{
		DistributionBody: DistributionBody{
			DryRun:            dr.DryRun,
			DistributionRules: distributionRules,
		},
		OnSuccess: onSuccess,
	}

	return dr.execDeleteDistribute(deleteDistributionParams.Name, deleteDistributionParams.Version, deleteDistribution)
}

func (dr *DeleteReleaseBundleService) execDeleteDistribute(name, version string, deleteDistribution DeleteRemoteDistributionBody) error {
	dryRunStr := ""
	if dr.IsDryRun() {
		dryRunStr = "[Dry run] "
	}
	log.Info(dryRunStr + "Deleting: " + name + "/" + version)

	httpClientsDetails := dr.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(deleteDistribution)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := dr.DistDetails.GetUrl() + "api/v1/distribution/" + name + "/" + version + "/delete"
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := dr.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Distribution response: ", resp.Status)
	log.Debug(utils.IndentJson(body))
	return errorutils.CheckError(err)
}

type DeleteRemoteDistributionBody struct {
	DistributionBody
	OnSuccess OnSuccess `json:"on_success"`
}

type DeleteDistributionParams struct {
	DistributionParams
	DeleteFromDistribution bool
}

func NewDeleteReleaseBundleParams(name, version string) DeleteDistributionParams {
	return DeleteDistributionParams{
		DistributionParams: DistributionParams{
			Name:    name,
			Version: version,
		},
	}
}
