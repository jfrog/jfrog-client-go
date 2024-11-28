package utils

import "testing"

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
