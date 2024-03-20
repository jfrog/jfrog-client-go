package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientUtils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
)

const (
	distributionBaseApi = "api/v2/distribution/"
	distribute          = "distribute"
	trackers            = "trackers"
)

type DistributeReleaseBundleService struct {
	client           *jfroghttpclient.JfrogHttpClient
	LcDetails        auth.ServiceDetails
	DryRun           bool
	AutoCreateRepo   bool
	Sync             bool
	MaxWaitMinutes   int
	DistributeParams distribution.DistributionParams
	ProjectKey       string
	Modifications
}

type DistributeReleaseBundleParams struct {
	Sync              bool
	AutoCreateRepo    bool
	MaxWaitMinutes    int
	DistributionRules []*distribution.DistributionCommonParams
	PathMappings      []PathMapping
	ProjectKey        string
}

func (dr *DistributeReleaseBundleService) GetHttpClient() *jfroghttpclient.JfrogHttpClient {
	return dr.client
}

func (dr *DistributeReleaseBundleService) ServiceDetails() auth.ServiceDetails {
	return dr.LcDetails
}

func (dr *DistributeReleaseBundleService) IsDryRun() bool {
	return dr.DryRun
}

func (dr *DistributeReleaseBundleService) IsSync() bool {
	return dr.Sync
}

func (dr *DistributeReleaseBundleService) GetMaxWaitMinutes() int {
	return dr.MaxWaitMinutes
}

func (dr *DistributeReleaseBundleService) GetRestApi(name, version string) string {
	return path.Join(distributionBaseApi, distribute, name, version)
}

func (dr *DistributeReleaseBundleService) GetDistributeBody() any {
	return dr.createDistributeBody()
}

func (dr *DistributeReleaseBundleService) GetDistributionParams() distribution.DistributionParams {
	return dr.DistributeParams
}

func (dr *DistributeReleaseBundleService) GetProjectKey() string {
	return dr.ProjectKey
}

func NewDistributeReleaseBundleService(client *jfroghttpclient.JfrogHttpClient) *DistributeReleaseBundleService {
	return &DistributeReleaseBundleService{client: client}
}

func (dr *DistributeReleaseBundleService) Distribute() error {
	trackerId, err := distribution.DoDistribute(dr)
	if err != nil || !dr.IsSync() || dr.IsDryRun() {
		return err
	}

	// Sync distribution
	return dr.waitForDistributionOperationCompletion(&dr.DistributeParams, trackerId)
}

func (dr *DistributeReleaseBundleService) createDistributeBody() ReleaseBundleDistributeBody {
	return ReleaseBundleDistributeBody{
		ReleaseBundleDistributeV1Body: distribution.CreateDistributeV1Body(dr.DistributeParams.DistributionRules, dr.DryRun, dr.AutoCreateRepo),
		Modifications:                 dr.Modifications,
	}
}

type ReleaseBundleDistributeBody struct {
	distribution.ReleaseBundleDistributeV1Body
	Modifications `json:"modifications"`
}

type Modifications struct {
	PathMappings []utils.PathMapping `json:"mappings"`
}

type PathMapping struct {
	Pattern string
	Target  string
}

func (rbs *ReleaseBundlesService) getReleaseBundleDistributions(rbDetails ReleaseBundleDetails, projectKey string) (distributionsResp GetDistributionsResponse, body []byte, err error) {
	restApi := GetReleaseBundleDistributionsApi(rbDetails)
	requestFullUrl, err := clientUtils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, distribution.GetProjectQueryParam(projectKey))
	if err != nil {
		return
	}
	httpClientsDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	resp, body, _, err := rbs.client.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		return
	}
	log.Debug("Artifactory response:", resp.Status)
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusAccepted, http.StatusOK); err != nil {
		return
	}
	err = errorutils.CheckError(json.Unmarshal(body, &distributionsResp))
	return
}

func GetReleaseBundleDistributionsApi(rbDetails ReleaseBundleDetails) string {
	return path.Join(distributionBaseApi, trackers, rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion)
}

type GetDistributionsResponse []struct {
	FriendlyId           json.Number `json:"distribution_tracker_friendly_id"`
	Type                 string      `json:"type"`
	ReleaseBundleName    string      `json:"release_bundle_name"`
	ReleaseBundleVersion string      `json:"release_bundle_version"`
	Repository           string      `json:"storing_repository"`
	Status               RbStatus    `json:"status"`
	DistributedBy        string      `json:"distributed_by"`
	Created              string      `json:"created"`
	StartTime            string      `json:"start_time"`
	FinishTime           string      `json:"finish_time"`
	Targets              []string    `json:"targets"`
}
