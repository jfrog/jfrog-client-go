package services

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

// VersionService returns the https client and Application details
type VersionService struct {
	client             *jfroghttpclient.JfrogHttpClient
	ApplicationDetails auth.ServiceDetails
}

// NewVersionService creates a new service to retrieve the version of Application
func NewVersionService(client *jfroghttpclient.JfrogHttpClient) *VersionService {
	return &VersionService{client: client}
}

// GetVersion endpoint currently not support by Application
// TODO Add support for this endpoint when application service will support this endpoint
func (vs *VersionService) GetVersion() (string, error) {
	return "1.0.0", nil
}
