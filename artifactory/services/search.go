package services

import (
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
)

type SearchService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
}

func NewSearchService(client *rthttpclient.ArtifactoryHttpClient) *SearchService {
	return &SearchService{client: client}
}

func (s *SearchService) GetArtifactoryDetails() auth.ServiceDetails {
	return s.ArtDetails
}

func (s *SearchService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	s.ArtDetails = rt
}

func (s *SearchService) IsDryRun() bool {
	return false
}

func (s *SearchService) GetJfrogHttpClient() (*rthttpclient.ArtifactoryHttpClient, error) {
	return s.client, nil
}

func (s *SearchService) Search(searchParams SearchParams) ([]utils.ResultItem, error) {
	return SearchBySpecFiles(searchParams, s, utils.ALL)
}

type SearchParams struct {
	*utils.ArtifactoryCommonParams
}

func (s *SearchParams) GetFile() *utils.ArtifactoryCommonParams {
	return s.ArtifactoryCommonParams
}

func NewSearchParams() SearchParams {
	return SearchParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{}}
}

func SearchBySpecFiles(searchParams SearchParams, flags utils.CommonConf, requiredArtifactProps utils.RequiredArtifactProps) ([]utils.ResultItem, error) {
	var resultItems []utils.ResultItem
	var itemsFound []utils.ResultItem
	var err error

	switch searchParams.GetSpecType() {
	case utils.WILDCARD:
		itemsFound, err = utils.SearchBySpecWithPattern(searchParams.GetFile(), flags, requiredArtifactProps)
	case utils.BUILD:
		itemsFound, err = utils.SearchBySpecWithBuild(searchParams.GetFile(), flags)
	case utils.AQL:
		itemsFound, err = utils.SearchBySpecWithAql(searchParams.GetFile(), flags, requiredArtifactProps)
	}
	if err != nil {
		return resultItems, err
	}
	resultItems = append(resultItems, itemsFound...)

	return resultItems, err
}
