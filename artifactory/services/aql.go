package services

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/utils/httpclient"
)

type AqlService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewAqlService(client *rthttpclient.ArtifactoryHttpClient) *AqlService {
	return &AqlService{client: client}
}

func (s *AqlService) GetArtifactoryDetails() auth.ArtifactoryDetails {
	return s.ArtDetails
}

func (s *AqlService) SetArtifactoryDetails(rt auth.ArtifactoryDetails) {
	s.ArtDetails = rt
}

func (s *AqlService) IsDryRun() bool {
	return false
}

func (s *AqlService) GetJfrogHttpClient() (*rthttpclient.ArtifactoryHttpClient, error) {
	return s.client, nil
}

func (s *AqlService) ExecAql(aql string) ([]byte, error) {
	return s.exec(aql)
}

func (s *AqlService) exec(aql string) ([]byte, error) {
	return utils.ExecAql(aql, s)
}
