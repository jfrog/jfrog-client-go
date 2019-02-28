package httpclient

import (
	"crypto/tls"
	"github.com/jfrog/jfrog-client-go/artifactory/auth/cert"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"time"
)

func ClientBuilder() *httpClientBuilder {
	return &httpClientBuilder{}
}

type httpClientBuilder struct {
	certificatesDirPath string
	skipCertsVerify     bool
}

func (builder *httpClientBuilder) SetCertificatesPath(certificatesPath string) *httpClientBuilder {
	builder.certificatesDirPath = certificatesPath
	return builder
}

func (builder *httpClientBuilder) SetSkipCertsVerify(skipCertsVerify bool) *httpClientBuilder {
	builder.skipCertsVerify = skipCertsVerify
	return builder
}

func (builder *httpClientBuilder) Build() (*HttpClient, error) {
	if builder.certificatesDirPath == "" {
		transport := createDefaultHttpTransport()
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: builder.skipCertsVerify}
		return &HttpClient{Client: &http.Client{Transport: transport}}, nil
	}

	transport, err := cert.GetTransportWithLoadedCert(builder.certificatesDirPath, builder.skipCertsVerify, createDefaultHttpTransport())
	if err != nil {
		return nil, errorutils.CheckError(errors.New("Failed creating HttpClient: " + err.Error()))
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
