package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXrayRepositoryTypeAll(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Repositories.Type = "all"
	err := configureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayRepositoryTypeByName(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Repositories.Type = "byname"
	err := configureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayRepositoryTypeByEmpty(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Repositories.Type = ""
	err := configureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayRepositoryTypeBad(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Repositories.Type = "bad"
	err := configureRepositories(&payloadBody, allPayload)
	assert.Error(t, err)
}

func TestXrayBuildTypeAll(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Repositories.Type = "all"
	err := configureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayBuildTypeByName(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Builds.Type = "byname"
	err := configureBuilds(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayBuildTypeByEmpty(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Builds.Type = ""
	err := configureBuilds(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestXrayBuildTypeBad(t *testing.T) {
	payloadBody := XrayWatchBody{}
	allPayload := XrayWatchParams{}
	allPayload.Builds.Type = "bad"
	err := configureBuilds(&payloadBody, allPayload)
	assert.Error(t, err)
}
