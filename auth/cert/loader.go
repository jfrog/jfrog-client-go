package cert

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
)

func loadCertificates(caCertPool *x509.CertPool, certificatesDirPath string) error {
	if !fileutils.IsPathExists(certificatesDirPath, false) {
		return nil
	}
	files, err := os.ReadDir(certificatesDirPath)
	err = errorutils.CheckError(err)
	if err != nil {
		return err
	}
	for _, file := range files {
		caCert, err := os.ReadFile(filepath.Join(certificatesDirPath, file.Name()))
		err = errorutils.CheckError(err)
		if err != nil {
			return err
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}
	return nil
}

func LoadCertificate(clientCertPath, clientCertKeyPath string) (certificate tls.Certificate, err error) {
	certificate, err = tls.LoadX509KeyPair(clientCertPath, clientCertKeyPath)
	if err != nil {
		if clientCertKeyPath == "" {
			err = errorutils.CheckErrorf("failed using the certificate located at %s. Reason: %s. Hint: A certificate key was not provided. Make sure that the certificate doesn't require a key", clientCertPath, err.Error())
			return
		}
		err = errorutils.CheckErrorf("failed loading client certificate: " + err.Error())
	}
	return
}

func GetTransportWithLoadedCert(certificatesDirPath string, insecureTls bool, transport *http.Transport) (*http.Transport, error) {
	// Remove once SystemCertPool supports windows
	caCertPool, err := loadSystemRoots()
	err = errorutils.CheckError(err)
	if err != nil {
		return nil, err
	}
	err = loadCertificates(caCertPool, certificatesDirPath)
	if err != nil {
		return nil, err
	}
	//#nosec G402 -- Skipping insecure tls verification was requested by the user.
	transport.TLSClientConfig = &tls.Config{
		RootCAs:            caCertPool,
		ClientSessionCache: tls.NewLRUClientSessionCache(1),
		InsecureSkipVerify: insecureTls,
	}
	transport.TLSClientConfig.BuildNameToCertificate()

	return transport, nil
}
