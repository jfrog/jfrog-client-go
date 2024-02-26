package services

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"path"
)

const (
	distributionBaseApi = "api/v2/distribution/"
	distribute          = "distribute"
)

type DistributeReleaseBundleService struct {
	client           *jfroghttpclient.JfrogHttpClient
	LcDetails        auth.ServiceDetails
	DryRun           bool
	AutoCreateRepo   bool
	Sync             bool
	MaxWaitMinutes   int
	DistributeParams distribution.DistributionParams
	Modifications
	ProjectKey string
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
		ReleaseBundleDistributeV1Body: distribution.CreateDistributeV1Body(dr.DistributeParams, dr.DryRun, dr.AutoCreateRepo),
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
