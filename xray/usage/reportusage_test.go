package usage

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsXrayVersionCompatible(t *testing.T) {
	tests := []struct {
		xrayVersion    string
		expectedResult bool
	}{
		{"1.2.0", false},
		{"2.9.0", false},
		{"2.0.0", false},
		{"3.79.3", false},
		{"3.80.0", true},
		{utils.Development, true},
		{"3.81.2", true},
		{"4.15.2", true},
	}
	for _, test := range tests {
		t.Run(test.xrayVersion, func(t *testing.T) {
			result := isVersionCompatible(test.xrayVersion)
			if test.expectedResult != result {
				t.Error(fmt.Errorf("expected %t, got %t", test.expectedResult, result))
			}
		})
	}
}

func TestXrayReportUsageJson(t *testing.T) {
	type reportUsageTestCase struct {
		productId  string
		EventId    string
		Attributes []ReportUsageAttribute
	}
	jsonPatterns := []string{
		`[{"product_name":"%s","event_name":"%s","origin":"API"}]`,
		`[{"product_name":"%s","event_name":"%s","origin":"API","data":{"%s":"%s"}}]`,
		`[{"product_name":"%s","event_name":"%s","origin":"API","data":{"%s":"%s","%s":"%s"}}]`,
	}

	cases := []reportUsageTestCase{
		{"jfrog-cli-go", "generic_audit", []ReportUsageAttribute{}},
		{"frogbot", "scan_pull_request", []ReportUsageAttribute{{AttributeName: "clientId", AttributeValue: "repo1"}}},
		{"jfrog-idea-plugin", "ci", []ReportUsageAttribute{{AttributeName: "buildNumber", AttributeValue: "1023456"}, {AttributeName: "clientId", AttributeValue: "user-hash"}}},
	}

	for _, test := range cases {
		// Create the expected json
		expectedResult := ""
		switch {
		case len(test.Attributes) == 1:
			expectedResult = fmt.Sprintf(jsonPatterns[1], test.productId, getExpectedEventName(test.productId, test.EventId), test.Attributes[0].AttributeName, test.Attributes[0].AttributeValue)
		case len(test.Attributes) == 2:
			expectedResult = fmt.Sprintf(jsonPatterns[2], test.productId, getExpectedEventName(test.productId, test.EventId), test.Attributes[0].AttributeName, test.Attributes[0].AttributeValue, test.Attributes[1].AttributeName, test.Attributes[1].AttributeValue)
		default:
			expectedResult = fmt.Sprintf(jsonPatterns[0], test.productId, getExpectedEventName(test.productId, test.EventId))
		}
		// Run test
		t.Run(test.EventId, func(t *testing.T) {
			body, err := reportUsageXrayToJson(CreateUsageEvents(test.productId, test.EventId, test.Attributes...))
			assert.NoError(t, err)
			assert.Equal(t, expectedResult, string(body))
		})
	}
}

func TestEcosystemReportUsageJson(t *testing.T) {
	type reportUsageTestCase struct {
		ProductId string
		AccountId string
		ClientId  string
		Features  []string
	}
	jsonPatterns := []string{
		`[{"productId":"%s","accountId":"%s","features":[]}]`,
		`[{"productId":"%s","accountId":"%s","features":["%s"]}]`,
		`[{"productId":"%s","accountId":"%s","clientId":"%s","features":["%s"]}]`,
		`[{"productId":"%s","accountId":"%s","clientId":"%s","features":["%s","%s"]}]`,
	}

	cases := []reportUsageTestCase{
		{"jfrog-cli-go", "platform.jfrog.io", "", []string{}},
		{"jfrog-cli-go", "platform.jfrog.io", "", []string{"generic_audit"}},
		{"frogbot", "platform.jfrog.io", "repo1", []string{"scan_pull_request"}},
		{"frogbot", "platform.jfrog.io", "repo1", []string{"scan_pull_request", "npm-dep"}},
	}

	// Create the expected json
	for _, test := range cases {
		// Create the expected json
		expectedResult := ""
		switch {
		case len(test.Features) == 1:
			if test.ClientId != "" {
				expectedResult = fmt.Sprintf(jsonPatterns[2], test.ProductId, test.AccountId, test.ClientId, test.Features[0])
			} else {
				expectedResult = fmt.Sprintf(jsonPatterns[1], test.ProductId, test.AccountId, test.Features[0])
			}
		case len(test.Features) == 2:
			if test.ClientId != "" {
				expectedResult = fmt.Sprintf(jsonPatterns[3], test.ProductId, test.AccountId, test.ClientId, test.Features[0], test.Features[1])
			} else {
				expectedResult = fmt.Sprintf(jsonPatterns[3], test.ProductId, test.AccountId, test.Features[0], test.Features[1])
			}
		default:
			if test.ClientId != "" {
				expectedResult = fmt.Sprintf(jsonPatterns[0], test.ProductId, test.AccountId, test.ClientId)
			} else {
				expectedResult = fmt.Sprintf(jsonPatterns[0], test.ProductId, test.AccountId)
			}
		}
		// Run test
		t.Run(strings.Join(test.Features, ","), func(t *testing.T) {
			if data, err := CreateUsageData(test.ProductId, test.AccountId, test.ClientId, test.Features...); len(test.Features) > 0 {
				assert.NoError(t, err)
				body, err := reportUsageEcosystemToJson(data)
				assert.NoError(t, err)
				assert.Equal(t, expectedResult, string(body))
			} else {
				assert.Error(t, err)
			}

		})
	}
}
