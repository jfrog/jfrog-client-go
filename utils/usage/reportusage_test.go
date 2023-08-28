package usage

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type reportUsageTestCase struct {
	ProductId   string
	AccountId   string
	ClientId    string
	Features    []string
	JsonPattern string
}

var cases = []reportUsageTestCase{
	{"jfrog-cli-go", "platform.jfrog.io", "", []string{}, `[{"productId":"%s","accountId":"%s","features":[]}]`},
	{"jfrog-cli-go", "platform.jfrog.io", "", []string{"generic_audit"}, `[{"productId":"%s","accountId":"%s","features":["%s"]}]`},
	{"frogbot", "platform.jfrog.io", "repo1", []string{"scan_pull_request"}, `[{"productId":"%s","accountId":"%s","clientId":"%s","features":["%s"]}]`},
	{"frogbot", "platform.jfrog.io", "repo1", []string{"scan_pull_request", "npm-dep"}, `[{"productId":"%s","accountId":"%s","clientId":"%s","features":["%s","%s"]}]`},
}

func TestEcosystemReportUsageToJson(t *testing.T) {
	// Create the expected json
	for _, test := range cases {
		// Create the expected json
		var expectedResult string
		switch len(test.Features) {
		case 1:
			if test.ClientId != "" {
				expectedResult = fmt.Sprintf(test.JsonPattern, test.ProductId, test.AccountId, test.ClientId, test.Features[0])
			} else {
				expectedResult = fmt.Sprintf(test.JsonPattern, test.ProductId, test.AccountId, test.Features[0])
			}
		case 2:
			if test.ClientId != "" {
				expectedResult = fmt.Sprintf(test.JsonPattern, test.ProductId, test.AccountId, test.ClientId, test.Features[0], test.Features[1])
			} else {
				expectedResult = fmt.Sprintf(test.JsonPattern, test.ProductId, test.AccountId, test.Features[0], test.Features[1])
			}
		default:
			if test.ClientId != "" {
				expectedResult = fmt.Sprintf(test.JsonPattern, test.ProductId, test.AccountId, test.ClientId)
			} else {
				expectedResult = fmt.Sprintf(test.JsonPattern, test.ProductId, test.AccountId)
			}
		}
		// Run test
		t.Run("Features: "+strings.Join(test.Features, ","), func(t *testing.T) {
			if data, err := CreateUsageData(test.ProductId, test.AccountId, test.ClientId, test.Features...); len(test.Features) > 0 {
				assert.NoError(t, err)
				body, err := json.Marshal([]ReportEcosystemUsageData{data})
				assert.NoError(t, err)
				assert.Equal(t, expectedResult, string(body))
			} else {
				assert.ErrorContains(t, err, "expected at least one feature to report usage on")
			}
		})
	}
}
