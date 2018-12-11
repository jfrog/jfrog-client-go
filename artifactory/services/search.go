package services

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/httpclient"
)

type SearchService struct {
	client     *httpclient.HttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewSearchService(client *httpclient.HttpClient) *SearchService {
	return &SearchService{client: client}
}

func (s *SearchService) GetArtifactoryDetails() auth.ArtifactoryDetails {
	return s.ArtDetails
}

func (s *SearchService) SetArtifactoryDetails(rt auth.ArtifactoryDetails) {
	s.ArtDetails = rt
}

func (s *SearchService) IsDryRun() bool {
	return false
}

func (s *SearchService) GetJfrogHttpClient() *httpclient.HttpClient {
	return s.client
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
	return SearchParams{}
}

func SearchBySpecFiles(searchParams SearchParams, flags utils.CommonConf, requiredArtifactProps utils.RequiredArtifactProps) ([]utils.ResultItem, error) {
	var resultItems []utils.ResultItem
	var itemsFound []utils.ResultItem
	var err error

	switch searchParams.GetSpecType() {
	case utils.WILDCARD, utils.SIMPLE:
		itemsFound, e := utils.SearchBySpecWithPattern(searchParams.GetFile(), flags, requiredArtifactProps)
		if e != nil {
			err = e
			return resultItems, err
		}
		resultItems = append(resultItems, itemsFound...)
	case utils.BUILD:
		itemsFound, err = utils.SearchBySpecWithBuild(searchParams.GetFile(), flags)
		if err != nil {
			return resultItems, err
		}
		resultItems = append(resultItems, itemsFound...)
	case utils.AQL:
		itemsFound, err = utils.SearchBySpecWithAql(searchParams.GetFile(), flags, requiredArtifactProps)
		if err != nil {
			return resultItems, err
		}
		resultItems = append(resultItems, itemsFound...)
	}
	return resultItems, err
}
