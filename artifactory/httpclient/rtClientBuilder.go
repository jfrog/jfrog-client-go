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
	skipCertsVerify     bool
	ArtDetails          *auth.ArtifactoryDetails
}

func (builder *artifactoryHttpClientBuilder) SetCertificatesPath(certificatesPath string) *artifactoryHttpClientBuilder {
	builder.certificatesDirPath = certificatesPath
	return builder
}

func (builder *artifactoryHttpClientBuilder) SetSkipCertsVerify(skipCertsVerify bool) *artifactoryHttpClientBuilder {
	builder.skipCertsVerify = skipCertsVerify
	return builder
}

func (builder *artifactoryHttpClientBuilder) SetArtDetails(rtDetails *auth.ArtifactoryDetails) *artifactoryHttpClientBuilder {
	builder.ArtDetails = rtDetails
	return builder
}

func (builder *artifactoryHttpClientBuilder) Build() (rtHttpClient *ArtifactoryHttpClient, err error) {
	rtHttpClient = &ArtifactoryHttpClient{ArtDetails: builder.ArtDetails}
	if builder.certificatesDirPath == "" {
		rtHttpClient.httpClient, err = httpclient.ClientBuilder().
			SetSkipCertsVerify(builder.skipCertsVerify).
			Build()
	} else {
		rtHttpClient.httpClient, err = httpclient.ClientBuilder().
			SetCertificatesPath(builder.certificatesDirPath).
			SetSkipCertsVerify(builder.skipCertsVerify).
			Build()
	}
	return
}
