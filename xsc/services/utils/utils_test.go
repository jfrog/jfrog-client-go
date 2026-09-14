package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXrayUrlToXscUrl(t *testing.T) {
	tests := []struct {
		testName      string
		xrayUrl       string
		xrayVersion   string
		expectedValue string
	}{
		{"after transition", "http://platform.jfrog.io/xray/", "3.107.13", "http://platform.jfrog.io/xray/api/v1/xsc/"},
		{"before transition", "http://platform.jfrog.io/xray/", "3.106.0", "http://platform.jfrog.io/xsc/api/v1/"},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actualValue := XrayUrlToXscUrl(test.xrayUrl, test.xrayVersion)
			if actualValue != test.expectedValue {
				t.Error(test.testName, "Expecting:", test.expectedValue, "Got:", actualValue)
			}
		})
	}
}

func TestGetGitRepoUrlKey(t *testing.T) {
	tests := []struct {
		testName   string
		gitRepoUrl string
		expected   string
	}{
		{"with_http", "http://git.com/jfrog/jfrog-client-go.git", "git.com/jfrog/jfrog-client-go.git"},
		{"with_https", "https://git.com/jfrog/jfrog-client-go.git", "git.com/jfrog/jfrog-client-go.git"},
		{"with_https_uppercase_scheme", "HTTPS://github.com/org/repo.git", "github.com/org/repo.git"},
		{"with_https_preserves_existing_case", "https://Git.COM/JFrog/jfrog-client-go.git", "Git.COM/JFrog/jfrog-client-go.git"},
		{"without_protocol", "git.com/jfrog/jfrog-client-go", "git.com/jfrog/jfrog-client-go.git"},
		{"without_protocol_and_port", "Git.COM:7999/JFrog/jfrog-client-go.git", "Git.COM:7999/JFrog/jfrog-client-go.git"},
		{"host_and_numeric_value", "github.com:443", "github.com:443.git"},
		{"scp", "git@Git.COM:JFrog/jfrog-client-go.git", "Git.COM/JFrog/jfrog-client-go.git"},
		{"ssh_with_port", "ssh://git@Git.COM:7999/JFrog/jfrog-client-go.git", "Git.COM/JFrog/jfrog-client-go.git"},
		{"azure_ssh", "git@ssh.dev.azure.com:v3/Org/Project/Repo", "dev.azure.com/Org/Project/_git/Repo.git"},
		{"azure_dev_azure_scp", "Org@dev.azure.com:v3/Org/Project/Repo", "dev.azure.com/Org/Project/_git/Repo.git"},
		{"azure_vs_ssh", "Org@vs-ssh.visualstudio.com:v3/Org/Project/Repo", "dev.azure.com/Org/Project/_git/Repo.git"},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			assert.Equal(t, test.expected, GetGitRepoUrlKey(test.gitRepoUrl))
		})
	}
}
