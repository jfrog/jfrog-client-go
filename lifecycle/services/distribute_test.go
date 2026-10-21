package services

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestCreateDistributeBody_IncludeEvidence verifies that the include_evidence
// flag is serialized into the distribute request body when set, and omitted
// (via omitempty) when left at its zero value.
func TestCreateDistributeBody_IncludeEvidence(t *testing.T) {
	tests := []struct {
		name            string
		includeEvidence bool
		wantSubstring   bool
	}{
		{name: "flag enabled is serialized", includeEvidence: true, wantSubstring: true},
		{name: "flag disabled is omitted", includeEvidence: false, wantSubstring: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dr := &DistributeReleaseBundleService{IncludeEvidence: tc.includeEvidence}
			body := dr.createDistributeBody()
			if body.IncludeEvidence != tc.includeEvidence {
				t.Fatalf("IncludeEvidence field = %v, want %v", body.IncludeEvidence, tc.includeEvidence)
			}
			marshalled, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("failed to marshal distribute body: %v", err)
			}
			gotSubstring := strings.Contains(string(marshalled), `"include_evidence":true`)
			if gotSubstring != tc.wantSubstring {
				t.Fatalf("body %s: contains include_evidence=%v, want %v", marshalled, gotSubstring, tc.wantSubstring)
			}
		})
	}
}
