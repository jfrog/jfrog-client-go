package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/jfrog/jfrog-client-go/xray/services/utils"
	"github.com/stretchr/testify/assert"
)

func TestXrayIgnoreRule(t *testing.T) {
	initXrayTest(t)
	t.Run("createCveIgnoreRule", createCveIgnoreRule)
	t.Run("createVulnerabilitesAndLicensesIgnoreRule", createVulnerabilitesAndLicensesIgnoreRule)
	t.Run("createIgnoreRuleOnWatch", createIgnoreRuleOnWatch)
}

func deleteIgnoreRule(t *testing.T, ignoreRuleId string) {
	err := testsXrayIgnoreRuleService.Delete(ignoreRuleId)
	assert.NoError(t, err)
}

func createCveIgnoreRule(t *testing.T) {
	var ignoreRuleId string
	defer func() {
		deleteIgnoreRule(t, ignoreRuleId)
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

	ignoreRuleId = createIgnoreRule(t, ignoreRuleFilter)
	assert.NotEmpty(t, ignoreRuleId)
}

func createVulnerabilitesAndLicensesIgnoreRule(t *testing.T) {
	var ignoreRuleId string
	defer func() {
		deleteIgnoreRule(t, ignoreRuleId)
	}()

	vulnerabilities := []string{"any"}
	licenses := []string{"any"}
	releaseBundle := utils.IgnoreFilterNameVersion{
		Name: "testRB",
	}
	releaseBundles := []utils.IgnoreFilterNameVersion{releaseBundle}
	ignoreRuleFilter := utils.IgnoreFilters{
		Vulnerabilities: vulnerabilities,
		Licenses:        licenses,
		ReleaseBundles:  releaseBundles,
	}

	ignoreRuleId = createIgnoreRule(t, ignoreRuleFilter)
	assert.NotEmpty(t, ignoreRuleId)
}

func createIgnoreRuleOnWatch(t *testing.T) {
	cve := []string{"CVE-2022-31197"}
	policyName := fmt.Sprintf("%s-%s", "test-policy-for-dummy-watch", getRunId())
	watchName := fmt.Sprintf("%s-%s", "test-watch-for-ignore-rule", getRunId())
	err := createDummyWatch(policyName, watchName)
	defer func() {
		assert.NoError(t, testsXrayWatchService.Delete(watchName))
		assert.NoError(t, testsXrayPolicyService.Delete(policyName))
	}()
	assert.NoError(t, err)
	watches := []string{watchName}

	var ignoreRuleId string
	defer func() {
		deleteIgnoreRule(t, ignoreRuleId)
	}()

	ignoreRuleFilter := utils.IgnoreFilters{
		CVEs:    cve,
		Watches: watches,
	}

	ignoreRuleId = createIgnoreRule(t, ignoreRuleFilter)
	assert.NotEmpty(t, ignoreRuleId)
}

func createIgnoreRule(t *testing.T, ignoreRuleFilter utils.IgnoreFilters) (ignoreRuleId string) {
	ignoreRuleParams := utils.IgnoreRuleParams{
		Notes:         "Create new ignore rule" + getRunId(),
		ExpiresAt:     time.Now().AddDate(0, 0, 1),
		IgnoreFilters: ignoreRuleFilter,
	}

	ignoreRuleId, err := testsXrayIgnoreRuleService.Create(ignoreRuleParams)
	assert.NoError(t, err)
	return ignoreRuleId
}

func createDummyWatch(policyName string, watchName string) error {
	if err := createDummyPolicy(policyName); err != nil {
		return err
	}
	params := utils.WatchParams{
		Name:   watchName,
		Active: true,
		Repositories: utils.WatchRepositoriesParams{
			Type: utils.WatchRepositoriesAll,
		},
		Policies: []utils.AssignedPolicy{
			{
				Name: policyName,
				Type: "security",
			},
		},
	}
	return testsXrayWatchService.Create(params)
}
