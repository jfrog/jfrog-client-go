package httpclient

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/jfrog/jfrog-client-go/auth/cert"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

var DefaultHttpTimeout = 30 * time.Second

func ClientBuilder() *httpClientBuilder {
	builder := &httpClientBuilder{}
	builder.SetTimeout(DefaultHttpTimeout)
	return builder
}

type httpClientBuilder struct {
	certificatesDirPath string
	clientCertPath      string
	clientCertKeyPath   string
	insecureTls         bool
	ctx                 context.Context
	timeout             time.Duration
	retries             int
	retryWaitMilliSecs  int
	httpClient          *http.Client
}

func (builder *httpClientBuilder) SetCertificatesPath(certificatesPath string) *httpClientBuilder {
	builder.certificatesDirPath = certificatesPath
	return builder
}

func (builder *httpClientBuilder) SetClientCertPath(certificatePath string) *httpClientBuilder {
	builder.clientCertPath = certificatePath
	return builder
}

func (builder *httpClientBuilder) SetClientCertKeyPath(certificatePath string) *httpClientBuilder {
	builder.clientCertKeyPath = certificatePath
	return builder
}

func (builder *httpClientBuilder) SetInsecureTls(insecureTls bool) *httpClientBuilder {
	builder.insecureTls = insecureTls
	return builder
}

func (builder *httpClientBuilder) SetHttpClient(httpClient *http.Client) *httpClientBuilder {
	builder.httpClient = httpClient
	return builder
}

func (builder *httpClientBuilder) SetContext(ctx context.Context) *httpClientBuilder {
	builder.ctx = ctx
	return builder
}

func (builder *httpClientBuilder) SetTimeout(timeout time.Duration) *httpClientBuilder {
	builder.timeout = timeout
	return builder
}

func (builder *httpClientBuilder) SetRetries(retries int) *httpClientBuilder {
	builder.retries = retries
	return builder
}

func (builder *httpClientBuilder) SetRetryWaitMilliSecs(retryWaitMilliSecs int) *httpClientBuilder {
	builder.retryWaitMilliSecs = retryWaitMilliSecs
	return builder
}

func (builder *httpClientBuilder) AddClientCertToTransport(transport *http.Transport) error {
	if builder.clientCertPath != "" {
		certificate, err := cert.LoadCertificate(builder.clientCertPath, builder.clientCertKeyPath)
		if err != nil {
			return err
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{certificate}
	}
	return nil
}

func (builder *httpClientBuilder) Build() (*HttpClient, error) {
	if builder.httpClient != nil {
		// Using a custom http.Client, pass-though.
		return &HttpClient{client: builder.httpClient, ctx: builder.ctx, retries: builder.retries, retryWaitMilliSecs: builder.retryWaitMilliSecs}, nil
	}

	var err error
	var transport *http.Transport

	if builder.certificatesDirPath == "" {
		transport = builder.createDefaultHttpTransport()
		//#nosec G402 -- Insecure TLS allowed here.
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: builder.insecureTls}
	} else {
		transport, err = cert.GetTransportWithLoadedCert(builder.certificatesDirPath, builder.insecureTls, builder.createDefaultHttpTransport())
		if err != nil {
			return nil, errorutils.CheckErrorf("Failed creating HttpClient: " + err.Error())
		}
	}
	err = builder.AddClientCertToTransport(transport)
	return &HttpClient{client: &http.Client{Transport: transport}, ctx: builder.ctx, retries: builder.retries, retryWaitMilliSecs: builder.retryWaitMilliSecs}, err
}

func (builder *httpClientBuilder) createDefaultHttpTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   builder.timeout,
			KeepAlive: 20 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
