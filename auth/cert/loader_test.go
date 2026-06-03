package cert

import (
	"crypto/tls"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSystemRootsNeverReturnsNilPool(t *testing.T) {
	pool, err := loadSystemRoots()
	require.NoError(t, err)
	assert.NotNil(t, pool)
}

func TestGetTransportWithLoadedCert(t *testing.T) {
	transport := &http.Transport{}
	result, err := GetTransportWithLoadedCert(filepath.Join(t.TempDir(), "missing-certs"), false, transport)
	require.NoError(t, err)
	require.NotNil(t, result.TLSClientConfig)
	assert.NotNil(t, result.TLSClientConfig.RootCAs)
	assert.Equal(t, uint16(tls.VersionTLS12), result.TLSClientConfig.MinVersion)
	assert.False(t, result.TLSClientConfig.InsecureSkipVerify)
}

func TestGetTransportWithLoadedCertInsecureTls(t *testing.T) {
	transport := &http.Transport{}
	result, err := GetTransportWithLoadedCert("", true, transport)
	require.NoError(t, err)
	assert.True(t, result.TLSClientConfig.InsecureSkipVerify)
}
