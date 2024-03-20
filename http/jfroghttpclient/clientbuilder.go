package jfroghttpclient

import (
	"context"
	"net/http"
	"time"

	"github.com/jfrog/jfrog-client-go/http/httpclient"
)

func JfrogClientBuilder() *jfrogHttpClientBuilder {
	builder := &jfrogHttpClientBuilder{}
	builder.SetDialTimeout(httpclient.DefaultDialTimeout)
	return builder
}

type jfrogHttpClientBuilder struct {
	certificatesDirPath    string
	insecureTls            bool
	ctx                    context.Context
	retries                int
	retryWaitTimMilliSecs  int
	preRequestInterceptors []PreRequestInterceptorFunc
	clientCertPath         string
	clientCertKeyPath      string
	dialTimeout            time.Duration
	overallRequestTimeout  time.Duration
	httpClient             *http.Client
}

func (builder *jfrogHttpClientBuilder) SetCertificatesPath(certificatesPath string) *jfrogHttpClientBuilder {
	builder.certificatesDirPath = certificatesPath
	return builder
}

func (builder *jfrogHttpClientBuilder) SetInsecureTls(insecureTls bool) *jfrogHttpClientBuilder {
	builder.insecureTls = insecureTls
	return builder
}

func (builder *jfrogHttpClientBuilder) SetClientCertPath(clientCertPath string) *jfrogHttpClientBuilder {
	builder.clientCertPath = clientCertPath
	return builder
}

func (builder *jfrogHttpClientBuilder) SetClientCertKeyPath(clientCertKeyPath string) *jfrogHttpClientBuilder {
	builder.clientCertKeyPath = clientCertKeyPath
	return builder
}

func (builder *jfrogHttpClientBuilder) SetContext(ctx context.Context) *jfrogHttpClientBuilder {
	builder.ctx = ctx
	return builder
}

func (builder *jfrogHttpClientBuilder) SetRetries(retries int) *jfrogHttpClientBuilder {
	builder.retries = retries
	return builder
}

func (builder *jfrogHttpClientBuilder) SetRetryWaitMilliSecs(retryWaitMilliSecs int) *jfrogHttpClientBuilder {
	builder.retryWaitTimMilliSecs = retryWaitMilliSecs
	return builder
}

func (builder *jfrogHttpClientBuilder) AppendPreRequestInterceptor(interceptor PreRequestInterceptorFunc) *jfrogHttpClientBuilder {
	builder.preRequestInterceptors = append(builder.preRequestInterceptors, interceptor)
	return builder
}

func (builder *jfrogHttpClientBuilder) SetDialTimeout(dialTimeout time.Duration) *jfrogHttpClientBuilder {
	builder.dialTimeout = dialTimeout
	return builder
}

func (builder *jfrogHttpClientBuilder) SetOverallRequestTimeout(overallRequestTimeout time.Duration) *jfrogHttpClientBuilder {
	builder.overallRequestTimeout = overallRequestTimeout
	return builder
}

func (builder *jfrogHttpClientBuilder) SetHttpClient(httpClient *http.Client) *jfrogHttpClientBuilder {
	builder.httpClient = httpClient
	return builder
}

func (builder *jfrogHttpClientBuilder) Build() (rtHttpClient *JfrogHttpClient, err error) {
	rtHttpClient = &JfrogHttpClient{preRequestInterceptors: builder.preRequestInterceptors}
	rtHttpClient.httpClient, err = httpclient.ClientBuilder().
		SetCertificatesPath(builder.certificatesDirPath).
		SetInsecureTls(builder.insecureTls).
		SetClientCertPath(builder.clientCertPath).
		SetClientCertKeyPath(builder.clientCertKeyPath).
		SetContext(builder.ctx).
		SetDialTimeout(builder.dialTimeout).
		SetOverallRequestTimeout(builder.overallRequestTimeout).
		SetRetries(builder.retries).
		SetRetryWaitMilliSecs(builder.retryWaitTimMilliSecs).
		SetHttpClient(builder.httpClient).
		Build()
	return
}
