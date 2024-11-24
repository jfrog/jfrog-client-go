package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGitRepoUrlKey(t *testing.T) {
	expected := "github.com/jfrog/jfrog-client-go.git"
	tests := []struct {
		testName   string
		gitRepoUrl string
	}{
		{"with_http", "http://github.com/jfrog/jfrog-client-go.git"},
		{"with_https", "https://github.com/jfrog/jfrog-client-go.git"},
		{"with_ssh", "git@github.com:jfrog/jfrog-client-go.git"},
		{"with_ssh_bb", "ssh://git@git.com/jfrog/jfrog-client-go.git"},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			assert.Equal(t, expected, GetGitRepoUrlKey(test.gitRepoUrl))
		})
	}
}
