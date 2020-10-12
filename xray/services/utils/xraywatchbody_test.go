package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXrayRepositoryTypeAll(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Repositories.Type = "all"
	err := ConfigureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayRepositoryTypeByName(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Repositories.Type = "byname"
	err := ConfigureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayRepositoryTypeByEmpty(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Repositories.Type = ""
	err := ConfigureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayRepositoryTypeBad(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Repositories.Type = "bad"
	err := ConfigureRepositories(&payloadBody, allPayload)
	assert.Error(t, err)
}

func TestXrayBuildTypeAll(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Repositories.Type = "all"
	err := ConfigureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayBuildTypeByName(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Builds.Type = "byname"
	err := ConfigureBuilds(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayBuildTypeByEmpty(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Builds.Type = ""
	err := ConfigureBuilds(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayBuildTypeBad(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Builds.Type = "bad"
	err := ConfigureBuilds(&payloadBody, allPayload)
	assert.Error(t, err)
}
