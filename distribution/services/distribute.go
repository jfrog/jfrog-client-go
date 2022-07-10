package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	distributionUtils "github.com/jfrog/jfrog-client-go/distribution/services/utils"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const defaultMaxWaitMinutes = 60           // 1 hour
const defaultSyncSleepIntervalSeconds = 10 // 10 seconds

type DistributeReleaseBundleService struct {
	client         *jfroghttpclient.JfrogHttpClient
	DistDetails    auth.ServiceDetails
	DryRun         bool
	Sync           bool
	AutoCreateRepo bool
	// Max time in minutes to wait for sync distribution to finish.
	MaxWaitMinutes int
}

func NewDistributeReleaseBundleService(client *jfroghttpclient.JfrogHttpClient) *DistributeReleaseBundleService {
	return &DistributeReleaseBundleService{client: client}
}

func (dr *DistributeReleaseBundleService) GetDistDetails() auth.ServiceDetails {
	return dr.DistDetails
}

func (dr *DistributeReleaseBundleService) Distribute(distributeParams DistributionParams) error {
	var distributionRules []DistributionRulesBody
	for _, spec := range distributeParams.DistributionRules {
		distributionRule := DistributionRulesBody{
			SiteName:     spec.GetSiteName(),
			CityName:     spec.GetCityName(),
			CountryCodes: spec.GetCountryCodes(),
		}
		distributionRules = append(distributionRules, distributionRule)
	}
	distribution := &DistributionBody{
		DryRun:            dr.DryRun,
		DistributionRules: distributionRules,
		AutoCreateRepo:    dr.AutoCreateRepo,
	}

	trackerId, err := dr.execDistribute(distributeParams.Name, distributeParams.Version, distribution)
	if err != nil || !dr.Sync || dr.DryRun {
		return err
	}

	// Sync distribution
	return dr.waitForDistribution(&distributeParams, trackerId)
}

func (dr *DistributeReleaseBundleService) execDistribute(name, version string, distribution *DistributionBody) (json.Number, error) {
	httpClientsDetails := dr.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(distribution)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	dryRunStr := ""
	if distribution.DryRun {
		dryRunStr = "[Dry run] "
	}
	log.Info(dryRunStr + "Distributing: " + name + "/" + version)

	url := dr.DistDetails.GetUrl() + "api/v1/distribution/" + name + "/" + version
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := dr.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return "", err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusAccepted); err != nil {
		return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, utils.IndentJson(body)))
	}
	response := distributionResponseBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", errorutils.CheckError(err)
	}

	log.Debug("Distribution response: ", resp.Status)
	log.Debug(utils.IndentJson(body))
	return response.TrackerId, nil
}

func (dr *DistributeReleaseBundleService) waitForDistribution(distributeParams *DistributionParams, trackerId json.Number) error {
	distributeBundleService := NewDistributionStatusService(dr.client)
	distributeBundleService.DistDetails = dr.GetDistDetails()
	distributionStatusParams := DistributionStatusParams{
		Name:      distributeParams.Name,
		Version:   distributeParams.Version,
		TrackerId: trackerId.String(),
	}
	maxWaitMinutes := defaultMaxWaitMinutes
	if dr.MaxWaitMinutes >= 1 {
		maxWaitMinutes = dr.MaxWaitMinutes
	}
	distributingMessage := fmt.Sprintf("Sync: Distributing %s/%s...", distributeParams.Name, distributeParams.Version)
	retryExecutor := &utils.RetryExecutor{
		MaxRetries:               maxWaitMinutes * 60 / defaultMaxWaitMinutes,
		RetriesIntervalMilliSecs: defaultSyncSleepIntervalSeconds * 1000,
		ErrorMessage:             "",
		LogMsgPrefix:             distributingMessage,
		ExecutionHandler: func() (bool, error) {
			response, err := distributeBundleService.GetStatus(distributionStatusParams)
			if err != nil {
				return false, errorutils.CheckError(err)
			}
			if (*response)[0].Status == Failed {
				bytes, err := json.Marshal(response)
				if err != nil {
					return false, errorutils.CheckError(err)
				}
				return false, errorutils.CheckErrorf("Distribution failed: " + utils.IndentJson(bytes))
			}
			if (*response)[0].Status == Completed {
				log.Info("Distribution Completed!")
				return false, nil
			}
			// Keep trying to get an answer
			log.Info(distributingMessage)
			return true, nil
		},
	}
	return retryExecutor.Execute()
}

type DistributionBody struct {
	DryRun            bool                    `json:"dry_run"`
	DistributionRules []DistributionRulesBody `json:"distribution_rules"`
	AutoCreateRepo    bool                    `json:"auto_create_missing_repositories,omitempty"`
}

type DistributionRulesBody struct {
	SiteName     string   `json:"site_name,omitempty"`
	CityName     string   `json:"city_name,omitempty"`
	CountryCodes []string `json:"country_codes,omitempty"`
}

type distributionResponseBody struct {
	TrackerId json.Number `json:"id"`
}

type DistributionParams struct {
	DistributionRules []*distributionUtils.DistributionCommonParams
	Name              string
	Version           string
}

func NewDistributeReleaseBundleParams(name, version string) DistributionParams {
	return DistributionParams{
		Name:    name,
		Version: version,
	}
}
