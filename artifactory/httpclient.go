package artifactory

import (
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/auth/cert"
	"github.com/jfrog/jfrog-client-go/httpclient"
)

func CreateArtifactoryHttpClient(config Config) (*httpclient.HttpClient, error) {
	if config.GetCertifactesPath() == "" {
		return &httpclient.HttpClient{
			Client: &http.Client{
				Timeout: config.GetTimeout(),
			},
		}, nil
	}

	transport, err := cert.GetTransportWithLoadedCert(config.GetCertifactesPath())
	if err != nil {
		return nil, err
	}
	return httpclient.NewHttpClient(&http.Client{
		Transport: transport,
		Timeout:   config.GetTimeout(),
	}), nil
}
