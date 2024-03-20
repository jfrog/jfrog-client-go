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
	GetProjectKey() string
}

func CreateDistributeV1Body(distCommonParams []*DistributionCommonParams, dryRun, isAutoCreateRepo bool) ReleaseBundleDistributeV1Body {
	var distributionRules []DistributionRulesBody
	for i := range distCommonParams {
		distributionRule := DistributionRulesBody{
			SiteName:     distCommonParams[i].GetSiteName(),
			CityName:     distCommonParams[i].GetCityName(),
			CountryCodes: distCommonParams[i].GetCountryCodes(),
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
	content, err := json.Marshal(dr.GetDistributeBody())
	if err != nil {
		return "", errorutils.CheckError(err)
	}

	dryRunStr := ""
	if dr.IsDryRun() {
		dryRunStr = "[Dry run] "
	}
	log.Info(dryRunStr + "Distributing: " + name + "/" + version)

	requestFullUrl, err := clientUtils.BuildUrl(dr.ServiceDetails().GetUrl(), dr.GetRestApi(name, version), GetProjectQueryParam(dr.GetProjectKey()))
	if err != nil {
		return "", err
	}

	httpClientsDetails := dr.ServiceDetails().CreateHttpClientDetails()
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := dr.GetHttpClient().SendPost(requestFullUrl, content, &httpClientsDetails)
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

type DistributionStatus string

const (
	NotDistributed DistributionStatus = "Not distributed"
	InProgress     DistributionStatus = "In progress"
	InQueue        DistributionStatus = "In queue"
	Completed      DistributionStatus = "Completed"
	Failed         DistributionStatus = "Failed"
)

type DistributionStatusResponse struct {
	Id                json.Number              `json:"distribution_id"`
	FriendlyId        json.Number              `json:"distribution_friendly_id,omitempty"`
	Type              DistributionType         `json:"type,omitempty"`
	Name              string                   `json:"release_bundle_name,omitempty"`
	Version           string                   `json:"release_bundle_version,omitempty"`
	Status            DistributionStatus       `json:"status,omitempty"`
	DistributionRules []DistributionRulesBody  `json:"distribution_rules,omitempty"`
	Sites             []DistributionSiteStatus `json:"sites,omitempty"`
}

type DistributionType string

const (
	Distribute                 DistributionType = "distribute"
	DeleteReleaseBundleVersion DistributionType = "delete_release_bundle_version"
)

type DistributionSiteStatus struct {
	Status            DistributionStatus `json:"status,omitempty"`
	Error             string             `json:"general_error,omitempty"`
	TargetArtifactory TargetArtifactory  `json:"target_artifactory,omitempty"`
	TotalFiles        json.Number        `json:"total_files,omitempty"`
	TotalBytes        json.Number        `json:"total_bytes,omitempty"`
	DistributedBytes  json.Number        `json:"distributed_bytes,omitempty"`
	DistributedFiles  json.Number        `json:"distributed_files,omitempty"`
	FileErrors        []string           `json:"file_errors,omitempty"`
	FilesInProgress   []string           `json:"files_in_progress,omitempty"`
}

type TargetArtifactory struct {
	ServiceId string `json:"service_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}
