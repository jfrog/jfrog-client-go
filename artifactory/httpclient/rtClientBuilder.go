package httpclient

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/httpclient"
)

func ArtifactoryClientBuilder() *artifactoryHttpClientBuilder {
	return &artifactoryHttpClientBuilder{}
}

type artifactoryHttpClientBuilder struct {
	certificatesDirPath string
	insecureTls         bool
	ArtDetails          *auth.ArtifactoryDetails
}

func (builder *artifactoryHttpClientBuilder) SetCertificatesPath(certificatesPath string) *artifactoryHttpClientBuilder {
	builder.certificatesDirPath = certificatesPath
	return builder
}

func (builder *artifactoryHttpClientBuilder) SetInsecureTls(insecureTls bool) *artifactoryHttpClientBuilder {
	builder.insecureTls = insecureTls
	return builder
}

func (builder *artifactoryHttpClientBuilder) SetArtDetails(rtDetails *auth.ArtifactoryDetails) *artifactoryHttpClientBuilder {
	builder.ArtDetails = rtDetails
	return builder
}

func (builder *artifactoryHttpClientBuilder) Build() (rtHttpClient *ArtifactoryHttpClient, err error) {
	rtHttpClient = &ArtifactoryHttpClient{ArtDetails: builder.ArtDetails}
	rtHttpClient.httpClient, err = httpclient.ClientBuilder().
		SetCertificatesPath(builder.certificatesDirPath).
		SetInsecureTls(builder.insecureTls).
		SetClientCertPath((*rtHttpClient.ArtDetails).GetClientCertPath()).
		SetClientCertKeyPath((*rtHttpClient.ArtDetails).GetClientCertKeyPath()).
		Build()
	return
}
