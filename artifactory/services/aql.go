package services

import (
	"io"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

type AqlService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
}

func NewAqlService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *AqlService {
	return &AqlService{artDetails: &artDetails, client: client}
}

func (s *AqlService) GetArtifactoryDetails() auth.ServiceDetails {
	return *s.artDetails
}

func (s *AqlService) IsDryRun() bool {
	return false
}

func (s *AqlService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return s.client
}

func (s *AqlService) ExecAql(aql string) (io.ReadCloser, error) {
	return s.exec(aql)
}

func (s *AqlService) exec(aql string) (io.ReadCloser, error) {
	return utils.ExecAql(aql, s)
}
