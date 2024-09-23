package tests

import (
	"net/http"
	"testing"

	"github.com/jfrog/jfrog-client-go/access"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/distribution"
	"github.com/jfrog/jfrog-client-go/evidence"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/lifecycle"
	"github.com/jfrog/jfrog-client-go/metadata"
	"github.com/jfrog/jfrog-client-go/pipelines"
	"github.com/jfrog/jfrog-client-go/xray"
	"github.com/jfrog/jfrog-client-go/xsc"
	"github.com/stretchr/testify/assert"
)

func TestUsesCustomHttpClient(t *testing.T) {
	t.Run("artifactory", func(t *testing.T) { usesCustomHttpClient(t, access.New) })
	t.Run("access", func(t *testing.T) { usesCustomHttpClient(t, access.New) })
	t.Run("distribution", func(t *testing.T) { usesCustomHttpClient(t, distribution.New) })
	t.Run("evidence", func(t *testing.T) { usesCustomHttpClient(t, evidence.New) })
	t.Run("lifecycle", func(t *testing.T) { usesCustomHttpClient(t, lifecycle.New) })
	t.Run("metadata", func(t *testing.T) { usesCustomHttpClient(t, metadata.NewManager) })
	t.Run("pipelines", func(t *testing.T) { usesCustomHttpClient(t, pipelines.New) })
	t.Run("xray", func(t *testing.T) { usesCustomHttpClient(t, xray.New) })
	t.Run("xsc", func(t *testing.T) { usesCustomHttpClient(t, xsc.New) })
}

type Manager interface {
	Client() *jfroghttpclient.JfrogHttpClient
}

func usesCustomHttpClient[K Manager](t *testing.T, createManager func(config config.Config) (K, error)) {
	client := http.DefaultClient
	config, err := config.NewConfigBuilder().
		SetServiceDetails(GetRtDetails()).
		SetHttpClient(client).
		Build()
	if err != nil {
		t.Error(err)
	}
	m, err := createManager(config)
	if err != nil {
		t.Error(err)
	}
	actualClient := m.Client().GetHttpClient().GetClient()
	assert.Equal(t, client, actualClient, "Expected the client to be the same")
}
