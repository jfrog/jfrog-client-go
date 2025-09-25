//go:build itest

package tests

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/stretchr/testify/assert"
)

func TestReplication(t *testing.T) {
	initArtifactoryTest(t)
	err := createReplication()
	if err != nil {
		t.Error(err)
	}
	err = getAndAssertReplication(t, GetReplicationConfig())
	if err != nil {
		t.Error(err)
	}
	err = deleteReplication()
	if err != nil {
		t.Error(err)
	}
	err = getAndAssertReplication(t, nil)
	assert.Error(t, err)
}

func createReplication() error {
	params := services.NewCreateReplicationParams()
	// Those fields are required
	params.Username = "anonymous"
	params.Password = "password"
	params.Url = "http://www.jfrog.com"
	params.CronExp = "0 0 14 * * ?"
	params.RepoKey = getRtTargetRepoKey()
	params.Enabled = true
	params.SocketTimeoutMillis = 100
	params.IncludePathPrefixPattern = "/include/path"
	return testsCreateReplicationService.CreateReplication(params)
}

func getAndAssertReplication(t *testing.T, expected []utils.ReplicationParams) error {
	replicationConf, err := testsReplicationGetService.GetReplication(getRtTargetRepoKey())
	if err != nil {
		return err
	}

	assert.Len(t, expected, 1, "Error in the test input. Probably a bug. Expecting only 1 replication. Got %d.", len(expected))
	assert.Len(t, replicationConf, 1, "Expected to fetch only 1 replication. Got %d.", len(replicationConf))

	// Artifactory may return the password encrypted. We therefore remove it,
	// before we can properly compare 'replicationConf' and 'expected'.
	replicationConf[0].Password = ""
	expected[0].Password = ""

	assert.ElementsMatch(t, replicationConf, expected)
	return nil
}

func deleteReplication() error {
	err := testsReplicationDeleteService.DeleteReplication(getRtTargetRepoKey())
	if err != nil {
		return err
	}
	return nil
}

func GetReplicationConfig() []utils.ReplicationParams {
	return []utils.ReplicationParams{
		{
			Url:                      "http://www.jfrog.com",
			Username:                 "anonymous",
			Password:                 "password",
			CronExp:                  "0 0 14 * * ?",
			RepoKey:                  getRtTargetRepoKey(),
			EnableEventReplication:   false,
			SocketTimeoutMillis:      100,
			Enabled:                  true,
			SyncDeletes:              false,
			SyncProperties:           false,
			SyncStatistics:           false,
			IncludePathPrefixPattern: "/include/path",
		},
	}
}
