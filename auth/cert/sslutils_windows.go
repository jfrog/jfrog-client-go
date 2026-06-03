//go:build windows
// +build windows

package cert

import (
	"crypto/x509"
)

func loadSystemRoots() (*x509.CertPool, error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	if pool == nil {
		pool = x509.NewCertPool()
	}
	return pool, nil
}
