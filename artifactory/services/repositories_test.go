package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUrlWithFilter(t *testing.T) {
	rs := RepositoriesService{}
	testCases := []struct {
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

	for _, testCase := range testCases {
		result := rs.createUrlWithFilter(testCase.params)
		assert.Equal(t, testCase.expected, result, "For params %+v, expected %s, but got %s", testCase.params, testCase.expected, result)
	}
}
