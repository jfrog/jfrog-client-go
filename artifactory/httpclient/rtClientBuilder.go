package httpclient

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/httpclient"
)

func ArtifactoryClientBuilder() *artifactoryHttpClientBuilder {
	return &artifactoryHttpClientBuilder{}
}

type artifactoryHttpClientBuilder struct {
	certificatesDirPath string
	insecureTls         bool
	CommonDetails       *auth.CommonDetails
}

func (builder *artifactoryHttpClientBuilder) SetCertificatesPath(certificatesPath string) *artifactoryHttpClientBuilder {
	builder.certificatesDirPath = certificatesPath
	return builder
}

func (builder *artifactoryHttpClientBuilder) SetInsecureTls(insecureTls bool) *artifactoryHttpClientBuilder {
	builder.insecureTls = insecureTls
	return builder
}

func (builder *artifactoryHttpClientBuilder) SetCommonDetails(rtDetails *auth.CommonDetails) *artifactoryHttpClientBuilder {
	builder.CommonDetails = rtDetails
	return builder
}

func (builder *artifactoryHttpClientBuilder) Build() (rtHttpClient *ArtifactoryHttpClient, err error) {
	rtHttpClient = &ArtifactoryHttpClient{ArtDetails: builder.CommonDetails}
	rtHttpClient.httpClient, err = httpclient.ClientBuilder().
		SetCertificatesPath(builder.certificatesDirPath).
		SetInsecureTls(builder.insecureTls).
		SetClientCertPath((*rtHttpClient.ArtDetails).GetClientCertPath()).
		SetClientCertKeyPath((*rtHttpClient.ArtDetails).GetClientCertKeyPath()).
		Build()
	return
}
