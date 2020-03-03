package services

import (
	"encoding/json"
	"errors"
	"net/http"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	distributionUtils "github.com/jfrog/jfrog-client-go/distribution/services/utils"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DistributionService struct {
	client      *rthttpclient.ArtifactoryHttpClient
	DistDetails auth.CommonDetails
	DryRun      bool
}

func NewDistributeService(client *rthttpclient.ArtifactoryHttpClient) *DistributionService {
	return &DistributionService{client: client}
}

func (ds *DistributionService) GetDistDetails() auth.CommonDetails {
	return ds.DistDetails
}

func (ds *DistributionService) Distribute(distributeParams DistributionParams) error {
	var distributionRules []DistributionRulesBody
	for _, spec := range distributeParams.DistributionRules {
		distributionRule := DistributionRulesBody{
			SiteName:     spec.GetSiteName(),
			CityName:     spec.GetCityName(),
			CountryCodes: spec.GetCountryCodes(),
		}
		distributionRules = append(distributionRules, distributionRule)
	}
	distribution := DistributionBody{
		DryRun:            ds.DryRun,
		DistributionRules: distributionRules,
	}

	return ds.execDistribute(distributeParams.Name, distributeParams.Version, distribution)
}

func (cbs *DistributionService) execDistribute(name, version string, distribution DistributionBody) error {
	httpClientsDetails := cbs.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(distribution)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := cbs.DistDetails.GetUrl() + "api/v1/distribution/" + name + "/" + version
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

type DistributionBody struct {
	DryRun            bool                    `json:"dry_run"`
	DistributionRules []DistributionRulesBody `json:"distribution_rules"`
}

type DistributionRulesBody struct {
	SiteName     string   `json:"site_name,omitempty"`
	CityName     string   `json:"city_name,omitempty"`
	CountryCodes []string `json:"country_codes,omitempty"`
}

type DistributionParams struct {
	DistributionRules []*distributionUtils.DistributionCommonParams
	Name              string
	Version           string
}

func NewDistributeParams(name, version string) DistributionParams {
	return DistributionParams{
		Name:    name,
		Version: version,
	}
}
