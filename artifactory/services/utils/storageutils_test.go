package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindRepositoryThatExists(t *testing.T) {
	si := buildFakeStorageInfo()

	result, err := si.FindRepositoryWithKey("repository-one")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "repository-one", result.RepoKey)
}

func TestFindRepositoryThatDoesNotExist(t *testing.T) {
	si := buildFakeStorageInfo()

	result, err := si.FindRepositoryWithKey("repository-three")

	assert.Error(t, err, "Failed to locate repository with key: repository-three")
	assert.Nil(t, result)
}

func buildFakeStorageInfo() StorageInfo {
	repositoryOne := RepositorySummary{RepoKey: "repository-one"}
	repositoryTwo := RepositorySummary{RepoKey: "repository-two"}

	return StorageInfo{
		BinariesSummary:         BinariesSummary{},
		RepositoriesSummaryList: []RepositorySummary{repositoryOne, repositoryTwo},
		FileStoreSummary:        FileStoreSummary{},
	}
}
