package httpclient

import (
	"crypto/tls"
	"errors"
	"github.com/jfrog/jfrog-client-go/artifactory/auth/cert"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net"
	"net/http"
	"time"
)

func ClientBuilder() *httpClientBuilder {
	return &httpClientBuilder{}
}

type httpClientBuilder struct {
	certificatesDirPath      string
	clientCertificatePath    string
	clientCertificateKeyPath string
	insecureTls              bool
}

func (builder *httpClientBuilder) SetCertificatesPath(certificatesPath string) *httpClientBuilder {
	builder.certificatesDirPath = certificatesPath
	return builder
}

func (builder *httpClientBuilder) SetClientCertificatePath(certificatePath string) *httpClientBuilder {
	builder.clientCertificatePath = certificatePath
	return builder
}

func (builder *httpClientBuilder) SetClientCertificateKeyPath(certificatePath string) *httpClientBuilder {
	builder.clientCertificateKeyPath = certificatePath
	return builder
}

func (builder *httpClientBuilder) SetInsecureTls(insecureTls bool) *httpClientBuilder {
	builder.insecureTls = insecureTls
	return builder
}

func (builder *httpClientBuilder) AddClientCertificateToTransport(transport *http.Transport) error {
	if builder.clientCertificatePath != "" {
		cert, err := tls.LoadX509KeyPair(builder.clientCertificatePath, builder.clientCertificateKeyPath)
		if err != nil {
			return errorutils.CheckError(errors.New("Failed loading client certificate: " + err.Error()))
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
	}

	return nil
}

func (builder *httpClientBuilder) Build() (*HttpClient, error) {
	if builder.certificatesDirPath == "" {
		transport := createDefaultHttpTransport()
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: builder.insecureTls}
		err := builder.AddClientCertificateToTransport(transport)
		if err != nil {
			return nil, err
		}
		return &HttpClient{Client: &http.Client{Transport: transport}}, nil
	}

	transport, err := cert.GetTransportWithLoadedCert(builder.certificatesDirPath, builder.insecureTls, createDefaultHttpTransport())
	if err != nil {
		return nil, errorutils.CheckError(errors.New("Failed creating HttpClient: " + err.Error()))
	}
	err = builder.AddClientCertificateToTransport(transport)
	if err != nil {
		return nil, err
	}
	return &HttpClient{Client: &http.Client{Transport: transport}}, nil
}

func createDefaultHttpTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 20 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
