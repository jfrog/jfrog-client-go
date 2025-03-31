package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXrayUrlToXscUrl(t *testing.T) {
	tests := []struct {
		testName      string
		xrayUrl       string
		xrayVersion   string
		expectedValue string
	}{
		// jfrog-ignore for tests
		{"after transition", "http://platform.jfrog.io/xray/", "3.107.13", "http://platform.jfrog.io/xray/api/v1/xsc/"},
		// jfrog-ignore for tests
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
	expected := "git.com/jfrog/jfrog-client-go.git"
	tests := []struct {
		testName   string
		gitRepoUrl string
	}{
		{"with_http", "http://git.com/jfrog/jfrog-client-go.git"},
		{"with_https", "https://git.com/jfrog/jfrog-client-go.git"},
		{"without_protocol", "git.com/jfrog/jfrog-client-go"},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			assert.Equal(t, expected, GetGitRepoUrlKey(test.gitRepoUrl))
		})
	}
}
