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

func TestConvertIntToStorageSizeString(t *testing.T) {
	tests := []struct {
		num    int
		output string
	}{
		{12546, "12.3KB"},
		{148576, "145.1KB"},
		{2587985, "2.5MB"},
		{12896547, "12.3MB"},
		{12896547785, "12.0GB"},
		{5248965785422365, "4773.9TB"},
	}

	for _, test := range tests {
		assert.Equal(t, test.output, ConvertIntToStorageSizeString(int64(test.num)))
	}
}
