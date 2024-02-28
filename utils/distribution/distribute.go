package distribution

import (
	"encoding/json"
	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientUtils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

type DistributeReleaseBundleExecutor interface {
	GetHttpClient() *jfroghttpclient.JfrogHttpClient
	ServiceDetails() auth.ServiceDetails
	IsDryRun() bool
	GetRestApi(name, version string) string
	GetDistributeBody() any
	GetDistributionParams() DistributionParams
}

func CreateDistributeV1Body(distributeParams DistributionParams, dryRun, isAutoCreateRepo bool) ReleaseBundleDistributeV1Body {
	var distributionRules []DistributionRulesBody
	for _, spec := range distributeParams.DistributionRules {
		distributionRule := DistributionRulesBody{
			SiteName:     spec.GetSiteName(),
			CityName:     spec.GetCityName(),
			CountryCodes: spec.GetCountryCodes(),
		}
		distributionRules = append(distributionRules, distributionRule)
	}
	body := ReleaseBundleDistributeV1Body{
		DryRun:            dryRun,
		DistributionRules: distributionRules,
		AutoCreateRepo:    isAutoCreateRepo,
	}
	return body
}

func DoDistribute(dr DistributeReleaseBundleExecutor) (trackerId json.Number, err error) {
	distributeParams := dr.GetDistributionParams()
	return execDistribute(dr, distributeParams.Name, distributeParams.Version)
}

func execDistribute(dr DistributeReleaseBundleExecutor, name, version string) (json.Number, error) {
	httpClientsDetails := dr.ServiceDetails().CreateHttpClientDetails()
	content, err := json.Marshal(dr.GetDistributeBody())
	if err != nil {
		return "", errorutils.CheckError(err)
	}

	dryRunStr := ""
	if dr.IsDryRun() {
		dryRunStr = "[Dry run] "
	}
	log.Info(dryRunStr + "Distributing: " + name + "/" + version)

	url := dr.ServiceDetails().GetUrl() + dr.GetRestApi(name, version)
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := dr.GetHttpClient().SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return "", err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusAccepted); err != nil {
		return "", err
	}
	response := DistributionResponseBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", errorutils.CheckError(err)
	}

	log.Debug("Distribution response:", resp.Status)
	log.Debug(clientUtils.IndentJson(body))
	return response.TrackerId, nil
}

func NewDistributeReleaseBundleParams(name, version string) DistributionParams {
	return DistributionParams{
		Name:    name,
		Version: version,
	}
}

type DistributionParams struct {
	DistributionRules []*DistributionCommonParams
	Name              string
	Version           string
}

type ReleaseBundleDistributeV1Body struct {
	DryRun            bool                    `json:"dry_run"`
	DistributionRules []DistributionRulesBody `json:"distribution_rules"`
	AutoCreateRepo    bool                    `json:"auto_create_missing_repositories,omitempty"`
}

type DistributionRulesBody struct {
	SiteName     string   `json:"site_name,omitempty"`
	CityName     string   `json:"city_name,omitempty"`
	CountryCodes []string `json:"country_codes,omitempty"`
}

type DistributionResponseBody struct {
	TrackerId json.Number `json:"id"`
}
