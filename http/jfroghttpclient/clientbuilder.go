package jfroghttpclient

import (
	"context"
	"github.com/jfrog/jfrog-client-go/http/httpclient"

	"github.com/jfrog/jfrog-client-go/auth"
)

func JfrogClientBuilder() *jfrogHttpClientBuilder {
	return &jfrogHttpClientBuilder{}
}

type jfrogHttpClientBuilder struct {
	certificatesDirPath string
	insecureTls         bool
	ctx                 context.Context
	ServiceDetails      *auth.ServiceDetails
}

func (builder *jfrogHttpClientBuilder) SetCertificatesPath(certificatesPath string) *jfrogHttpClientBuilder {
	builder.certificatesDirPath = certificatesPath
	return builder
}

func (builder *jfrogHttpClientBuilder) SetInsecureTls(insecureTls bool) *jfrogHttpClientBuilder {
	builder.insecureTls = insecureTls
	return builder
}

func (builder *jfrogHttpClientBuilder) SetServiceDetails(rtDetails *auth.ServiceDetails) *jfrogHttpClientBuilder {
	builder.ServiceDetails = rtDetails
	return builder
}

func (builder *jfrogHttpClientBuilder) SetContext(ctx context.Context) *jfrogHttpClientBuilder {
	builder.ctx = ctx
	return builder
}

func (builder *jfrogHttpClientBuilder) Build() (rtHttpClient *JfrogHttpClient, err error) {
	rtHttpClient = &JfrogHttpClient{JfrogServiceDetails: builder.ServiceDetails}
	rtHttpClient.httpClient, err = httpclient.ClientBuilder().
		SetCertificatesPath(builder.certificatesDirPath).
		SetInsecureTls(builder.insecureTls).
		SetClientCertPath((*rtHttpClient.JfrogServiceDetails).GetClientCertPath()).
		SetClientCertKeyPath((*rtHttpClient.JfrogServiceDetails).GetClientCertKeyPath()).
		SetContext(builder.ctx).
		Build()
	return
}
