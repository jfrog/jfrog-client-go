package tests

import (
	clientutils "github.com/jfrog/jfrog-client-go/utils"
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
	t.Run("createSkipNonApplicablePolicy", createSkipNonApplicable)
}

func deletePolicy(t *testing.T, policyName string) {
	err := testsXrayPolicyService.Delete(policyName)
	assert.NoError(t, err)
}

func createMinSeverity(t *testing.T) {
	policyName := "create-min-severity" + getRunId()
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "min-severity" + getRunId(),
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.Low, false),
		Priority: 1,
	}
	createAndCheckPolicy(t, policyName, true, utils.Security, policyRule)
}

func createRangeSeverity(t *testing.T) {
	policyName := "create-range-severity" + getRunId()
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "range-severity" + getRunId(),
		Criteria: *utils.CreateCvssRangePolicyCriteria(3.4, 5.6),
		Priority: 1,
	}
	createAndCheckPolicy(t, policyName, true, utils.Security, policyRule)
}

func createLicenseAllowed(t *testing.T) {
	policyName := "create-allowed-licenses" + getRunId()
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "allowed-licenses" + getRunId(),
		Criteria: *utils.CreateLicensePolicyCriteria(true, true, true, "MIT", "Apache-2.0"),
		Priority: 1,
	}
	createAndCheckPolicy(t, policyName, true, utils.License, policyRule)
}

func createLicenseBanned(t *testing.T) {
	policyName := "create-banned-licenses" + getRunId()
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "banned-licenses" + getRunId(),
		Criteria: *utils.CreateLicensePolicyCriteria(false, true, true, "MIT", "Apache-2.0"),
		Priority: 1,
	}
	createAndCheckPolicy(t, policyName, true, utils.License, policyRule)
}

func create2Priorities(t *testing.T) {
	policyName := "create-2-priorties" + getRunId()
	defer deletePolicy(t, policyName)

	policyRule1 := utils.PolicyRule{
		Name:     "priority-1" + getRunId(),
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.Low, false),
		Priority: 1,
	}
	policyRule2 := utils.PolicyRule{
		Name:     "priority-2" + getRunId(),
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.Medium, false),
		Priority: 2,
	}
	createAndCheckPolicy(t, policyName, true, utils.Security, policyRule1, policyRule2)
}

func createPolicyActions(t *testing.T) {
	policyName := "create-policy-actions" + getRunId()
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "policy-actions" + getRunId(),
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.High, false),
		Priority: 1,
		Actions: &utils.PolicyAction{
			BlockDownload: utils.PolicyBlockDownload{
				Active:    clientutils.Pointer(true),
				Unscanned: clientutils.Pointer(true),
			},
			BlockReleaseBundleDistribution: clientutils.Pointer(true),
			FailBuild:                      clientutils.Pointer(true),
			NotifyDeployer:                 clientutils.Pointer(true),
			NotifyWatchRecipients:          clientutils.Pointer(true),
			CustomSeverity:                 utils.Information,
		},
	}
	createAndCheckPolicy(t, policyName, true, utils.Security, policyRule)
}

func createUpdatePolicy(t *testing.T) {
	policyName := "update-policy" + getRunId()
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "low-severity" + getRunId(),
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.Low, false),
		Priority: 1,
	}
	createAndCheckPolicy(t, policyName, true, utils.Security, policyRule)

	policyRule = utils.PolicyRule{
		Name:     "medium-severity" + getRunId(),
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.Medium, false),
		Priority: 1,
	}

	createAndCheckPolicy(t, policyName, false, utils.Security, policyRule)
}

func createSkipNonApplicable(t *testing.T) {
	policyName := "skip-non-applicable" + getRunId()
	defer deletePolicy(t, policyName)

	policyRule := utils.PolicyRule{
		Name:     "skip-non-applicable-rule" + getRunId(),
		Criteria: *utils.CreateSeverityPolicyCriteria(utils.Low, true),
		Priority: 1,
	}
	createAndCheckPolicy(t, policyName, true, utils.Security, policyRule)
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
	assert.True(t, policyRulesAreEqual(expected.Rules, actual.Rules))
}

// policyRulesAreEqual tells whether both PolicyRule slices contain the same elements, regardless of the order.
func policyRulesAreEqual(expectedRules, actualRules []utils.PolicyRule) bool {
	if len(expectedRules) != len(actualRules) {
		return false
	}
	for _, expectedRule := range expectedRules {
		for _, actualRule := range actualRules {
			if expectedRule.Name == actualRule.Name && expectedRule.Priority == actualRule.Priority && assert.ObjectsAreEqual(expectedRule.Criteria, actualRule.Criteria) {
				if expectedRule.Actions != nil {
					return assert.ObjectsAreEqual(expectedRule.Actions, actualRule.Actions)
				}
				return true
			}
		}
	}
	return false
}
