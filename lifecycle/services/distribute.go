package services

import (
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
	DistributeParams distribution.DistributionParams
	PathMapping
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

func (dr *DistributeReleaseBundleService) IsAutoCreateRepo() bool {
	return dr.AutoCreateRepo
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

func NewDistributeReleaseBundleService(client *jfroghttpclient.JfrogHttpClient) *DistributeReleaseBundleService {
	return &DistributeReleaseBundleService{client: client}
}

func (dr *DistributeReleaseBundleService) Distribute() error {
	_, err := distribution.DoDistribute(dr)
	return err
}

func (dr *DistributeReleaseBundleService) createDistributeBody() ReleaseBundleDistributeBody {
	return ReleaseBundleDistributeBody{
		ReleaseBundleDistributeV1Body: distribution.CreateDistributeV1Body(dr.DistributeParams, dr.DryRun, dr.AutoCreateRepo),
		Modifications: Modifications{
			PathMappings: distribution.CreatePathMappings(dr.Pattern, dr.Target),
		},
	}
}

type ReleaseBundleDistributeBody struct {
	distribution.ReleaseBundleDistributeV1Body
	Modifications `json:"modifications"`
}

type Modifications struct {
	PathMappings []distribution.PathMapping `json:"mappings"`
}

type PathMapping struct {
	Pattern string
	Target  string
}
