package usage

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/stretchr/testify/assert"
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
				t.Error(fmt.Errorf("expected %t, got %t", test.expectedResult, result))
			}
		})
	}
}

func TestReportUsageJson(t *testing.T) {
	type test struct {
		productId      string
		commandName    string
		serviceId      string
		serverSize     string
		expectedResult string
		jsonPatternNum int
	}

	jsonPatterns := []string{
		`{"productId":"%s","features":[{"featureId":"%s","attributes":{"serverSize":"%s","serviceId":"%s"}}]}`,
		`{"productId":"%s","features":[{"featureId":"%s"}]}`,
	}

	preTests := []test{
		{"jfrog-cli-go/1.26.0", "rt_transfer_files", "jfrt@01g8dj3wcw22y01atqp63n1haq", "6.08 GB", "{\"productId\":\"jfrog-cli-go/1.26.0\",\"features\":[{\"featureId\":\"rt_transfer_files\",\"attributes\":{\"serverSize\":\"6.08 GB\",\"serviceId\":\"jfrt@01g8dj3wcw22y01atqp63n1haq\"}}]}", 0},
		{"jfrog-client-go", "rt_download", "", "3.58 GB", "{\"productId\":\"jfrog-client-go\",\"features\":[{\"featureId\":\"rt_download\"}]}", 1},
		{"test", "rt_build", "jfrt@01g8dj3wcw22y01atqp63n1haq", "", "", 1},
		{"agent/1.25.0", "rt_go", "", "", "", 1},
	}

	var tests []test
	// Create the expected json
	for _, test := range preTests {
		if test.serverSize != "" && test.serviceId != "" {
			test.expectedResult = fmt.Sprintf(jsonPatterns[test.jsonPatternNum], test.productId, test.commandName, test.serverSize, test.serviceId)
		} else {
			test.expectedResult = fmt.Sprintf(jsonPatterns[test.jsonPatternNum], test.productId, test.commandName)
		}

		tests = append(tests, test)
	}

	for _, test := range tests {
		t.Run(test.commandName, func(t *testing.T) {
			body, err := reportUsageToJson(test.productId, test.commandName, test.serviceId, test.serverSize)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, string(body))
		})
	}
}
