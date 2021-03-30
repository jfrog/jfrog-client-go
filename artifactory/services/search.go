package services

import (
	"errors"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/version"
)

type SearchService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
}

func NewSearchService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *SearchService {
	return &SearchService{artDetails: &artDetails, client: client}
}

func (s *SearchService) GetArtifactoryDetails() auth.ServiceDetails {
	return *s.artDetails
}

func (s *SearchService) IsDryRun() bool {
	return false
}

func (s *SearchService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return s.client
}

func (s *SearchService) Search(searchParams SearchParams) (*content.ContentReader, error) {
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

func SearchBySpecFiles(searchParams SearchParams, flags utils.CommonConf, requiredArtifactProps utils.RequiredArtifactProps) (*content.ContentReader, error) {
	artifactoryVersionStr, err := flags.GetArtifactoryDetails().GetVersion()
	if err != nil {
		return nil, err
	}
	artifactoryVersion := version.NewVersion(artifactoryVersionStr)
	err = utils.ValidateTransitiveSearchAllowed(searchParams.ArtifactoryCommonParams, artifactoryVersion)
	if err != nil {
		return nil, err
	}
	switch searchParams.GetSpecType() {
	case utils.WILDCARD:
		return utils.SearchBySpecWithPattern(searchParams.GetFile(), flags, requiredArtifactProps)
	case utils.BUILD:
		return utils.SearchBySpecWithBuild(searchParams.GetFile(), flags)
	case utils.AQL:
		return utils.SearchBySpecWithAql(searchParams.GetFile(), flags, requiredArtifactProps)
	default:
		return nil, errorutils.CheckError(errors.New("Error at SearchBySpecFiles: Unknown spec type"))
	}
}
