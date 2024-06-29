package tests

import (
	"testing"
	"time"

	"github.com/jfrog/jfrog-client-go/xray/services/utils"
	"github.com/stretchr/testify/assert"
)

func TestXrayIgnoreRule(t *testing.T) {
	initXrayTest(t)
	t.Run("createCveIgnoreRule", createCveIgnoreRule)
}

func deleteIgnoreRule(t *testing.T, ignoreRuleId string) {
	err := testsXrayIgnoreRuleService.Delete(ignoreRuleId)
	assert.NoError(t, err)
}

func createCveIgnoreRule(t *testing.T) {
	var ignoreRuleId string
	defer func() {
		if ignoreRuleId != "" {
			deleteIgnoreRule(t, ignoreRuleId)
		}
	}()

	component := utils.IgnoreFilterNameVersion{
		Name:    "gav://org.postgresql:postgresql",
		Version: "42.2.3.jre7",
	}
	components := []utils.IgnoreFilterNameVersion{component}

	cve := []string{"CVE-2022-31197"}
	ignoreRuleFilter := utils.IgnoreFilters{
		CVEs:       cve,
		Components: components,
	}

	createIgnoreRule(t, &ignoreRuleId, ignoreRuleFilter)
}

func createIgnoreRule(t *testing.T, ignoreRuleId *string, ignoreRuleFilter utils.IgnoreFilters) *utils.IgnoreRuleParams {
	ignoreRuleParams := utils.IgnoreRuleParams{
		Notes:         "Create new ignore rule" + getRunId(),
		ExpiresAt:     time.Date(2025, time.June, 28, 14, 30, 0, 0, time.UTC),
		IgnoreFilters: ignoreRuleFilter,
	}

	err := testsXrayIgnoreRuleService.Create(ignoreRuleParams, ignoreRuleId)
	assert.NoError(t, err)
	return &ignoreRuleParams
}
