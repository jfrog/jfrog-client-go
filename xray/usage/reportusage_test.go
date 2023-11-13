package usage

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsXrayVersionCompatible(t *testing.T) {
	tests := []struct {
		xrayVersion string
		compatible  bool
	}{
		{"1.2.0", false},
		{"2.9.0", false},
		{"2.0.0", false},
		{"3.80.3", false},
		{"3.81.4", false},
		{utils.Development, true},
		{"3.83.0", true},
		{"3.83.3", true},
		{"4.15.2", true},
	}
	for _, test := range tests {
		t.Run(test.xrayVersion, func(t *testing.T) {
			err := utils.ValidateMinimumVersion(utils.Xray, test.xrayVersion, minXrayReportUsageVersion)
			if test.compatible {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

type reportUsageTestCase struct {
	productId   string
	EventId     string
	Attributes  []ReportUsageAttribute
	jsonPattern string
}

var reportCases = []reportUsageTestCase{
	{"jfrog-cli-go", "generic_audit", []ReportUsageAttribute{}, `[{"product_name":"%s","event_name":"%s","origin":"API_CLI"}]`},
	{"frogbot", "scan_pull_request", []ReportUsageAttribute{{AttributeName: "clientId", AttributeValue: "repo1"}}, `[{"data":{"%s":"%s"},"product_name":"%s","event_name":"%s","origin":"API_CLI"}]`},
	{"jfrog-idea-plugin", "ci", []ReportUsageAttribute{{AttributeName: "buildNumber", AttributeValue: "1023456"}, {AttributeName: "clientId", AttributeValue: "user-hash"}}, `[{"data":{"%s":"%s","%s":"%s"},"product_name":"%s","event_name":"%s","origin":"API_CLI"}]`},
}

func TestXrayUsageEventToJson(t *testing.T) {
	for _, test := range reportCases {
		// Create the expected json
		var expectedResult string
		switch len(test.Attributes) {
		case 1:
			expectedResult = fmt.Sprintf(test.jsonPattern, test.Attributes[0].AttributeName, test.Attributes[0].AttributeValue, test.productId, GetExpectedXrayEventName(test.productId, test.EventId))
		case 2:
			expectedResult = fmt.Sprintf(test.jsonPattern, test.Attributes[0].AttributeName, test.Attributes[0].AttributeValue, test.Attributes[1].AttributeName, test.Attributes[1].AttributeValue, test.productId, GetExpectedXrayEventName(test.productId, test.EventId))
		default:
			expectedResult = fmt.Sprintf(test.jsonPattern, test.productId, GetExpectedXrayEventName(test.productId, test.EventId))
		}
		// Run test
		t.Run(test.EventId, func(t *testing.T) {
			body, err := json.Marshal([]ReportXrayEventData{CreateUsageEvent(test.productId, test.EventId, test.Attributes...)})
			assert.NoError(t, err)
			assert.Equal(t, expectedResult, string(body))
		})
	}
}
