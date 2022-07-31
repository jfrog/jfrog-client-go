package services

import (
	"encoding/json"
	"fmt"
	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"time"
)

type OnSuccess string

const (
	Keep   OnSuccess = "keep"
	Delete OnSuccess = "delete"
)

// Delete received release bundles from the edge nodes. On success, keep or delete the release bundle from the distribution service.
type DeleteReleaseBundleService struct {
	client      *jfroghttpclient.JfrogHttpClient
	DistDetails auth.ServiceDetails
	DryRun      bool
	Sync        bool
	// Max time in minutes to wait for sync distribution to finish.
	MaxWaitMinutes int
}

func NewDeleteReleaseBundleService(client *jfroghttpclient.JfrogHttpClient) *DeleteReleaseBundleService {
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
	dr.Sync = deleteDistributionParams.Sync
	dr.MaxWaitMinutes = deleteDistributionParams.MaxWaitMinutes
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
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK, http.StatusAccepted); err != nil {
		return err
	}
	if dr.Sync {
		err := dr.waitForDeletion(name, version)
		if err != nil {
			return err
		}
	}
	log.Debug("Distribution response: ", resp.Status)
	log.Debug(utils.IndentJson(body))
	return errorutils.CheckError(err)
}

func (dr *DeleteReleaseBundleService) waitForDeletion(name, version string) error {
	distributeBundleService := NewDistributionStatusService(dr.client)
	distributeBundleService.DistDetails = dr.GetDistDetails()
	httpClientsDetails := distributeBundleService.DistDetails.CreateHttpClientDetails()
	maxWaitMinutes := defaultMaxWaitMinutes
	if dr.MaxWaitMinutes >= 1 {
		maxWaitMinutes = dr.MaxWaitMinutes
	}
	for timeElapsed := 0; timeElapsed < maxWaitMinutes*60; timeElapsed += defaultSyncSleepIntervalSeconds {
		if timeElapsed%60 == 0 {
			log.Info(fmt.Sprintf("Performing sync deletion of release bundle %s/%s...", name, version))
		}
		resp, _, _, err := dr.client.SendGet(dr.DistDetails.GetUrl()+"api/v1/release_bundle/"+name+"/"+version+"/distribution", true, &httpClientsDetails)
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusNotFound {
			log.Info("Deletion Completed!")
			return nil
		}
		if resp.StatusCode != http.StatusOK {
			return errorutils.CheckErrorf("Error while waiting to deletion: status code " + fmt.Sprint(resp.StatusCode) + ".")
		}
		time.Sleep(time.Second * defaultSyncSleepIntervalSeconds)
	}
	return errorutils.CheckErrorf("Timeout for sync deletion. ")
}

type DeleteRemoteDistributionBody struct {
	DistributionBody
	OnSuccess OnSuccess `json:"on_success"`
}

type DeleteDistributionParams struct {
	DistributionParams
	DeleteFromDistribution bool
	Sync                   bool
	// Max time in minutes to wait for sync distribution to finish.
	MaxWaitMinutes int
}

func NewDeleteReleaseBundleParams(name, version string) DeleteDistributionParams {
	return DeleteDistributionParams{
		DistributionParams: DistributionParams{
			Name:    name,
			Version: version,
		},
	}
}
