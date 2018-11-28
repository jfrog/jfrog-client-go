package artifactory

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth/cert"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"net/http"
)

func CreateArtifactoryHttpClient(config Config) (*httpclient.HttpClient, error) {
	if config.GetCertifactesPath() == "" {
		return httpclient.NewDefaultHttpClient(), nil
	}

	transport, err := cert.GetTransportWithLoadedCert(config.GetCertifactesPath())
	if err != nil {
		return nil, err
	}
	return httpclient.NewHttpClient(&http.Client{Transport: transport}), nil
}
