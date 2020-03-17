package tests

import (
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

var (
	// TrimSuffix cannot be constants
	// we can declare them as top-level variables
	repoKey string = strings.TrimSuffix(RtTargetRepo, "/")
)

func TestReplication(t *testing.T) {
	err := createReplication()
	if err != nil {
		t.Error(err)
	}
	err = showReplication(t, GetReplicationConfge())
	if err != nil {
		t.Error(err)
	}
	err = deleteReplication(t)
	if err != nil {
		t.Error(err)
	}
	err = showReplication(t, nil)
	assert.Error(t, err)
}

func createReplication() error {
	params := services.NewPushReplicationParams()
	//those fields are required
	params.Username = "anonymous"
	params.Password = "password"
	params.URL = "http://www.jfrog.com"
	params.CronExp = "0 0 14 * * ?"
	params.RepoKey = repoKey
	params.Enabled = true
	params.SocketTimeoutMillis = 100
	return testsReplicationService.Push(params)
}

func showReplication(t *testing.T, expected []services.PushReplicationParams) error {
	replicationConf, err := testsReplicationShowService.Show(repoKey)
	if err != nil {
		return err
	}
	assert.ElementsMatch(t, replicationConf, expected)
	return nil
}

func deleteReplication(t *testing.T) error {
	err := testsReplicationDeleteService.Delete(repoKey)
	if err != nil {
		return err
	}
	return nil
}

func GetReplicationConfge() []services.PushReplicationParams {
	return []services.PushReplicationParams{
		{
			URL:      "http://www.jfrog.com",
			Username: "anonymous",
			Password: "password",
			CommonReplicationParams: services.CommonReplicationParams{
				CronExp:                "0 0 14 * * ?",
				RepoKey:                repoKey,
				EnableEventReplication: false,
				SocketTimeoutMillis:    100,
				Enabled:                true,
				SyncDeletes:            false,
				SyncProperties:         false,
				SyncStatistics:         false,
				PathPrefix:             "",
			},
		},
	}
}
