package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponentResolutionRequestVcs(t *testing.T) {
	tests := []struct {
		testName string
		vcs      *ComponentResolutionVcs
		wantVcs  any
	}{
		{
			testName: "omitted_outside_a_git_checkout",
		},
		{
			testName: "url_branch_and_revision",
			vcs:      &ComponentResolutionVcs{Url: "https://github.com/org/repo.git", Branch: "main", Revision: "abc123"},
			wantVcs:  map[string]any{"url": "https://github.com/org/repo.git", "branch": "main", "revision": "abc123"},
		},
		{
			testName: "empty_vcs_sends_no_empty_url",
			vcs:      &ComponentResolutionVcs{},
			wantVcs:  map[string]any{},
		},
		{
			testName: "detached_head_has_no_branch",
			vcs:      &ComponentResolutionVcs{Url: "https://github.com/org/repo.git", Revision: "abc123"},
			wantVcs:  map[string]any{"url": "https://github.com/org/repo.git", "revision": "abc123"},
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			body, err := json.Marshal(ComponentResolutionRequest{BuildTool: "npm", Repo: "npm-virtual", Lockfile: "{}", Vcs: test.vcs})
			require.NoError(t, err)
			var sent map[string]any
			require.NoError(t, json.Unmarshal(body, &sent))
			vcs, present := sent["vcs"]
			if test.wantVcs == nil {
				assert.False(t, present, "a request without vcs must not send the key")
				return
			}
			assert.Equal(t, test.wantVcs, vcs)
		})
	}
}
