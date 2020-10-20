package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWatchRepositoryTypeAll(t *testing.T) {
	payloadBody := WatchBody{}
	allPayload := WatchParams{}
	allPayload.Repositories.Type = "all"
	err := configureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestWatchRepositoryTypeByName(t *testing.T) {
	payloadBody := WatchBody{}
	allPayload := WatchParams{}
	allPayload.Repositories.Type = "byname"
	err := configureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestWatchRepositoryTypeByEmpty(t *testing.T) {
	payloadBody := WatchBody{}
	allPayload := WatchParams{}
	allPayload.Repositories.Type = ""
	err := configureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestWatchRepositoryTypeBad(t *testing.T) {
	payloadBody := WatchBody{}
	allPayload := WatchParams{}
	allPayload.Repositories.Type = "bad"
	err := configureRepositories(&payloadBody, allPayload)
	assert.Error(t, err)
}

func TestWatchBuildTypeAll(t *testing.T) {
	payloadBody := WatchBody{}
	allPayload := WatchParams{}
	allPayload.Repositories.Type = "all"
	err := configureRepositories(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestWatchBuildTypeByName(t *testing.T) {
	payloadBody := WatchBody{}
	allPayload := WatchParams{}
	allPayload.Builds.Type = "byname"
	err := configureBuilds(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestWatchBuildTypeByEmpty(t *testing.T) {
	payloadBody := WatchBody{}
	allPayload := WatchParams{}
	allPayload.Builds.Type = ""
	err := configureBuilds(&payloadBody, allPayload)
	assert.NoError(t, err)
}

func TestWatchBuildTypeBad(t *testing.T) {
	payloadBody := WatchBody{}
	allPayload := WatchParams{}
	allPayload.Builds.Type = "bad"
	err := configureBuilds(&payloadBody, allPayload)
	assert.Error(t, err)
}
