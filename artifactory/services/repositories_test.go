package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWithFilterUrl(t *testing.T) {
	tests := []struct {
		params   RepositoriesFilterParams
		expected string
	}{
		{RepositoriesFilterParams{RepoType: "git", PackageType: "npm", ProjectKey: "123"}, "api/repositories?packageType=npm&project=123&type=git"},
		{RepositoriesFilterParams{RepoType: "git", PackageType: "npm"}, "api/repositories?packageType=npm&type=git"},
		{RepositoriesFilterParams{RepoType: "git"}, "api/repositories?type=git"},
		{RepositoriesFilterParams{PackageType: "npm"}, "api/repositories?packageType=npm"},
		{RepositoriesFilterParams{ProjectKey: "123"}, "api/repositories?project=123"},
		{RepositoriesFilterParams{}, "api/repositories"},
	}

	for _, test := range tests {
		result := createWithFilterUrl(test.params)
		assert.Equal(t, test.expected, result, "For params %+v, expected %s, but got %s", test.params, test.expected, result)
	}
}
