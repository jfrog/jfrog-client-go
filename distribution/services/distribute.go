package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	distributionUtils "github.com/jfrog/jfrog-client-go/distribution/services/utils"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const defaultMaxWaitMinutes = 60    // 1 hour
const defaultSyncSleepInterval = 10 // 10 seconds

type DistributeReleaseBundleService struct {
	client      *jfroghttpclient.JfrogHttpClient
	DistDetails auth.ServiceDetails
	DryRun      bool
	Sync        bool
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
	}

	trackerId, err := dr.execDistribute(distributeParams.Name, distributeParams.Version, distribution)
	if err != nil || !dr.Sync {
		return err
	}

	// Sync distribution
	return dr.waitForDistribution(&distributeParams, trackerId)
}

func (dr *DistributeReleaseBundleService) execDistribute(name, version string, distribution *DistributionBody) (int, error) {
	httpClientsDetails := dr.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(distribution)
	if err != nil {
		return 0, errorutils.CheckError(err)
	}
	url := dr.DistDetails.GetUrl() + "api/v1/distribution/" + name + "/" + version
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := dr.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return 0, errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	response := distributionResponseBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	log.Debug("Distribution response: ", resp.Status)
	log.Debug(utils.IndentJson(body))
	return response.TrackerId, errorutils.CheckError(err)
}

func (dr *DistributeReleaseBundleService) waitForDistribution(distributeParams *DistributionParams, trackerId int) error {
	distributeBundleService := NewDistributionStatusService(dr.client)
	distributeBundleService.DistDetails = dr.GetDistDetails()
	distributionStatusParams := DistributionStatusParams{
		Name:      distributeParams.Name,
		Version:   distributeParams.Version,
		TrackerId: strconv.Itoa(trackerId),
	}
	maxWaitMinutes := defaultMaxWaitMinutes
	if dr.MaxWaitMinutes >= 1 {
		maxWaitMinutes = dr.MaxWaitMinutes
	}
	distributingMessage := fmt.Sprintf("Distributing %s/%s...", distributeParams.Name, distributeParams.Version)
	for timeElapsed := 0; timeElapsed < maxWaitMinutes*60; timeElapsed += defaultSyncSleepInterval {
		if timeElapsed%60 == 0 {
			log.Info(distributingMessage)
		}
		response, err := distributeBundleService.GetStatus(distributionStatusParams)
		if err != nil {
			return err
		}

		if (*response)[0].Status == Failed {
			bytes, err := json.Marshal(response)
			if err != nil {
				return errorutils.CheckError(err)
			}
			return errorutils.CheckError(errors.New("Distribution failed: " + utils.IndentJson(bytes)))
		}
		if (*response)[0].Status == Completed {
			return nil
		}
		time.Sleep(time.Second * defaultSyncSleepInterval)
	}
	return errorutils.CheckError(errors.New("Timeout for sync distribution"))
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

type distributionResponseBody struct {
	TrackerId int `json:"id"`
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
