package distribution

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDistributeV1BodyWithPriority(t *testing.T) {
	rules := []*DistributionCommonParams{
		{SiteName: "edge1", Priority: "medium"},
		{SiteName: "edge2", Priority: "low"},
	}

	body := CreateDistributeV1BodyWithPriority(rules, false, true, "high")
	assert.Equal(t, "high", body.Priority)
	assert.True(t, body.AutoCreateRepo)
	assert.False(t, body.DryRun)
	require.Len(t, body.DistributionRules, 2)
	assert.Equal(t, "edge1", body.DistributionRules[0].SiteName)
	assert.Equal(t, "medium", body.DistributionRules[0].Priority)
	assert.Equal(t, "edge2", body.DistributionRules[1].SiteName)
	assert.Equal(t, "low", body.DistributionRules[1].Priority)

	raw, err := json.Marshal(body)
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"priority":"high"`)
	assert.Contains(t, string(raw), `"site_name":"edge1"`)
	assert.Contains(t, string(raw), `"priority":"medium"`)
	assert.Contains(t, string(raw), `"priority":"low"`)
}

func TestCreateDistributeV1BodyOmitsEmptyPriority(t *testing.T) {
	body := CreateDistributeV1Body([]*DistributionCommonParams{{SiteName: "edge1"}}, true, false)
	assert.Empty(t, body.Priority)
	require.Len(t, body.DistributionRules, 1)
	assert.Empty(t, body.DistributionRules[0].Priority)

	raw, err := json.Marshal(body)
	require.NoError(t, err)
	assert.NotContains(t, string(raw), `"priority"`)
	assert.Contains(t, string(raw), `"dry_run":true`)
}

func TestDistributionCommonParamsPriorityAccessors(t *testing.T) {
	params := &DistributionCommonParams{SiteName: "edge1"}
	assert.Empty(t, params.GetPriority())
	params.SetPriority("high")
	assert.Equal(t, "high", params.GetPriority())
}
