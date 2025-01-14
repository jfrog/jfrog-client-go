package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientUtils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const defaultMaxWaitMinutes = 60                     // 1 hour
const DefaultDistributeSyncSleepIntervalSeconds = 10 // 10 seconds

type DistributeReleaseBundleV1Service struct {
	client         *jfroghttpclient.JfrogHttpClient
	DistDetails    auth.ServiceDetails
	DryRun         bool
	Sync           bool
	AutoCreateRepo bool
	// Max time in minutes to wait for sync distribution to finish.
	MaxWaitMinutes   int
	DistributeParams distribution.DistributionParams
}

func (dr *DistributeReleaseBundleV1Service) GetHttpClient() *jfroghttpclient.JfrogHttpClient {
	return dr.client
}

func (dr *DistributeReleaseBundleV1Service) ServiceDetails() auth.ServiceDetails {
	return dr.DistDetails
}

func (dr *DistributeReleaseBundleV1Service) IsDryRun() bool {
	return dr.DryRun
}

func (dr *DistributeReleaseBundleV1Service) IsSync() bool {
	return dr.Sync
}

func (dr *DistributeReleaseBundleV1Service) GetMaxWaitMinutes() int {
	return dr.MaxWaitMinutes
}

func (dr *DistributeReleaseBundleV1Service) GetRestApi(name, version string) string {
	return "api/v1/distribution/" + name + "/" + version
}

func (dr *DistributeReleaseBundleV1Service) GetDistributeBody() any {
	return distribution.CreateDistributeV1Body(dr.DistributeParams.DistributionRules, dr.DryRun, dr.AutoCreateRepo)
}

func (dr *DistributeReleaseBundleV1Service) GetDistributionParams() distribution.DistributionParams {
	return dr.DistributeParams
}

func (dr *DistributeReleaseBundleV1Service) GetProjectKey() string {
	return ""
}

func NewDistributeReleaseBundleV1Service(client *jfroghttpclient.JfrogHttpClient) *DistributeReleaseBundleV1Service {
	return &DistributeReleaseBundleV1Service{client: client}
}

func (dr *DistributeReleaseBundleV1Service) Distribute() error {
	trackerId, err := distribution.DoDistribute(dr)
	if err != nil || !dr.IsSync() || dr.IsDryRun() {
		return err
	}

	// Sync distribution
	return dr.waitForDistribution(&dr.DistributeParams, trackerId)
}

func (dr *DistributeReleaseBundleV1Service) waitForDistribution(distributeParams *distribution.DistributionParams, trackerId json.Number) error {
	distributeBundleService := NewDistributionStatusService(dr.GetHttpClient())
	distributeBundleService.DistDetails = dr.ServiceDetails()
	distributionStatusParams := DistributionStatusParams{
		Name:      distributeParams.Name,
		Version:   distributeParams.Version,
		TrackerId: trackerId.String(),
	}
	maxWaitMinutes := defaultMaxWaitMinutes
	if dr.GetMaxWaitMinutes() >= 1 {
		maxWaitMinutes = dr.GetMaxWaitMinutes()
	}
	distributingMessage := fmt.Sprintf("Sync: Distributing %s/%s...", distributeParams.Name, distributeParams.Version)
	retryExecutor := &clientUtils.RetryExecutor{
		MaxRetries:               maxWaitMinutes * 60 / DefaultDistributeSyncSleepIntervalSeconds,
		RetriesIntervalMilliSecs: DefaultDistributeSyncSleepIntervalSeconds * 1000,
		ErrorMessage:             "",
		LogMsgPrefix:             distributingMessage,
		ExecutionHandler: func() (bool, error) {
			response, err := distributeBundleService.GetStatus(distributionStatusParams)
			if err != nil {
				return false, errorutils.CheckError(err)
			}
			if (*response)[0].Status == distribution.Failed {
				bytes, err := json.Marshal(response)
				if err != nil {
					return false, errorutils.CheckError(err)
				}
				return false, errorutils.CheckErrorf("Distribution failed: " + clientUtils.IndentJson(bytes))
			}
			if (*response)[0].Status == distribution.Completed {
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
