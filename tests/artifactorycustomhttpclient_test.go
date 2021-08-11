package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestGetArtifactoryVersionWithProxyShouldFail(t *testing.T) {
	initArtifactoryTest(t)
	rtDetails := GetRtDetails()

	proxyUrl, err := url.Parse("http://invalidproxy:12345")
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	}

	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(rtDetails).
		SetDryRun(false).
		SetHttpClient(client).
		Build()
	if err != nil {
		t.Error(err)
	}

	rtManager, err := artifactory.New(serviceConfig)
	if err != nil {
		t.Error(err)
	}

	_, err = rtManager.GetVersion()
	assert.Error(t, err, "Should fail with invalid proxy")
}
