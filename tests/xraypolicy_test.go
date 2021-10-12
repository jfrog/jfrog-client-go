package tests

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/xray/services/utils"
	"github.com/stretchr/testify/assert"
)

func TestXrayPolicy(t *testing.T) {
	initXrayTest(t)
	t.Run("createMinSeverity", createMinSeverity)
	t.Run("createRangeSeverity", createRangeSeverity)
	t.Run("createLicenseAllowed", createLicenseAllowed)
	t.Run("createLicenseBanned", createLicenseBanned)
	t.Run("create2Priorities", create2Priorities)
	t.Run("createPolicyActions", createPolicyActions)
	t.Run("createUpdatePolicy", createUpdatePolicy)
}

func deletePolicy(t *testing.T, policyName string) {
	err := testsXrayPolicyService.Delete(policyName)
	assert.NoError(t, err)
}

func createMinSeverity(t *testing.T) {
	policyName := "create-min-severity" + randomRunNumber
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "min-severity",
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.Low),
		Priority: 1,
	}
	createAndCheckPolicy(t, policyName, true, utils.Security, policyRule)
}

func createRangeSeverity(t *testing.T) {
	policyName := "create-range-severity" + randomRunNumber
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "range-severity",
		Criteria: *utils.CreateCvssRangePolicyCriteria(3.4, 5.6),
		Priority: 1,
	}
	createAndCheckPolicy(t, policyName, true, utils.Security, policyRule)
}

func createLicenseAllowed(t *testing.T) {
	policyName := "create-allowed-licenses" + randomRunNumber
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "allowed-licenses",
		Criteria: *utils.CreateLicensePolicyCriteria(true, true, true, "MIT", "Apache-2.0"),
		Priority: 1,
	}
	createAndCheckPolicy(t, policyName, true, utils.License, policyRule)
}

func createLicenseBanned(t *testing.T) {
	policyName := "create-banned-licenses" + randomRunNumber
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "banned-licenses",
		Criteria: *utils.CreateLicensePolicyCriteria(false, true, true, "MIT", "Apache-2.0"),
		Priority: 1,
	}
	createAndCheckPolicy(t, policyName, true, utils.License, policyRule)
}

func create2Priorities(t *testing.T) {
	policyName := "create-2-priorties" + randomRunNumber
	defer deletePolicy(t, policyName)

	policyRule1 := utils.PolicyRule{
		Name:     "priority-1",
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.Low),
		Priority: 1,
	}
	policyRule2 := utils.PolicyRule{
		Name:     "priority-2",
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.Medium),
		Priority: 2,
	}
	createAndCheckPolicy(t, policyName, true, utils.Security, policyRule1, policyRule2)
}

func createPolicyActions(t *testing.T) {
	policyName := "create-policy-actions" + randomRunNumber
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "policy-actions",
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.High),
		Priority: 1,
		Actions: &utils.PolicyAction{
			BlockDownload: utils.PolicyBlockDownload{
				Active:    true,
				Unscanned: true,
			},
			BlockReleaseBundleDistribution: true,
			FailBuild:                      true,
			NotifyDeployer:                 true,
			NotifyWatchRecipients:          true,
			CustomSeverity:                 utils.Information,
		},
	}
	createAndCheckPolicy(t, policyName, true, utils.Security, policyRule)
}

func createUpdatePolicy(t *testing.T) {
	policyName := "update-policy" + randomRunNumber
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "low-severity",
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.Low),
		Priority: 1,
	}
	createAndCheckPolicy(t, policyName, true, utils.Security, policyRule)

	policyRule = utils.PolicyRule{
		Name:     "medium-severity",
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.Medium),
		Priority: 1,
	}

	createAndCheckPolicy(t, policyName, false, utils.Security, policyRule)
}

func createPolicy(t *testing.T, policyName string, policyType utils.PolicyType, policyRules ...utils.PolicyRule) *utils.PolicyParams {
	policyParams := utils.PolicyParams{
		Name:        policyName,
		Type:        policyType,
		Description: "crate-policy-description",
		Rules:       policyRules,
	}
	err := testsXrayPolicyService.Create(policyParams)
	assert.NoError(t, err)
	return &policyParams
}

func updatePolicy(t *testing.T, policyName string, policyType utils.PolicyType, policyRules ...utils.PolicyRule) *utils.PolicyParams {
	policyParams := utils.PolicyParams{
		Name:        policyName,
		Type:        policyType,
		Description: "update-policy-description",
		Rules:       policyRules,
	}
	err := testsXrayPolicyService.Update(policyParams)
	assert.NoError(t, err)
	return &policyParams
}

func createAndCheckPolicy(t *testing.T, policyName string, create bool, policyType utils.PolicyType, policyRules ...utils.PolicyRule) {
	var expected *utils.PolicyParams
	if create {
		expected = createPolicy(t, policyName, policyType, policyRules...)
	} else {
		expected = updatePolicy(t, policyName, policyType, policyRules...)
	}

	// Get policy
	actual, err := testsXrayPolicyService.Get(policyName)
	assert.NoError(t, err)

	// Compare general policy details
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Description, actual.Description)

	// Compare rules
	assert.Len(t, actual.Rules, len(expected.Rules))
	for i, expectedRule := range expected.Rules {
		actualRule := actual.Rules[i]
		assert.Equal(t, expectedRule.Name, actualRule.Name)
		assert.Equal(t, expectedRule.Priority, actualRule.Priority)
		assert.Equal(t, expectedRule.Criteria, actualRule.Criteria)
		if expectedRule.Actions != nil {
			assert.Equal(t, expectedRule.Actions, actualRule.Actions)
		}
	}

}
