package services

import (
	"io"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

type AqlService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewAqlService(client *jfroghttpclient.JfrogHttpClient) *AqlService {
	return &AqlService{client: client}
}

func (s *AqlService) GetArtifactoryDetails() auth.ServiceDetails {
	return s.ArtDetails
}

func (s *AqlService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	s.ArtDetails = rt
}

func (s *AqlService) IsDryRun() bool {
	return false
}

func (s *AqlService) GetJfrogHttpClient() (*jfroghttpclient.JfrogHttpClient, error) {
	return s.client, nil
}

func (s *AqlService) ExecAql(aql string) (io.ReadCloser, error) {
	return s.exec(aql)
}

func (s *AqlService) exec(aql string) (io.ReadCloser, error) {
	return utils.ExecAql(aql, s)
}
