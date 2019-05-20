package usage

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils"
	"testing"
)

func TestIsVersionCompatible(t *testing.T) {
	tests := []struct {
		artifactoryVersion string
		expectedResult     bool
	}{
		{"6.5.0", false},
		{"6.2.0", false},
		{"5.9.0", false},
		{"6.0.0", false},
		{"6.6.0", false},
		{"6.9.0", true},
		{utils.Development, true},
		{"6.10.2", true},
		{"6.15.2", true},
	}
	for _, test := range tests {
		t.Run(test.artifactoryVersion, func(t *testing.T) {
			result := isVersionCompatible(test.artifactoryVersion)
			if test.expectedResult != result {
				t.Error(fmt.Errorf("Expected %t, got %t", test.expectedResult, result))
			}
		})
	}
}

func TestReportUsageJson(t *testing.T) {
	type test struct {
		productId      string
		commandName    string
		expectedResult string
	}

	json := `{"productId":"%s","features":[{"featureId":"%s"}]}`
	preTests := []test{
		{"jfrog-cli-go/1.26.0", "rt_copy", ""},
		{"jfrog-client-go", "rt_download", ""},
		{"test", "rt_build", ""},
		{"agent/1.25.0", "rt_go", ""},
	}

	var tests []test
	// Create the expected json
	for _, test := range preTests {
		test.expectedResult = fmt.Sprintf(json, test.productId, test.commandName)
		tests = append(tests, test)
	}

	for _, test := range tests {
		t.Run(test.commandName, func(t *testing.T) {
			body, err := reportUsageToJson(test.productId, test.commandName)
			if err != nil {
				t.Error(err)
			}

			if string(body) != test.expectedResult {
				t.Error(fmt.Errorf("Expected %s, got %s", test.expectedResult, string(body)))
			}
		})
	}
}
