package services

import (
	"net/http"
	"testing"
)

// Regression test for the build-publish/build-scan async-indexing race
// (https://github.com/jfrog/jfrog-azure-devops-extension/issues/596).
// Before the fix, a 404 carrying the "not indexed" message in the "error" field was
// not recognized as transient, so the trigger failed instead of retrying.
func TestIsArtifactoryBuildNotFoundError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantRetry  bool // true => detected as transient not-found (non-nil error => triggerScan retries)
	}{
		{
			name:       "legacy_info_field_wasnt_found",
			statusCode: http.StatusNotFound,
			body:       `{"info":"Build my-build number 42 wasn't found in Artifactory"}`,
			wantRetry:  true,
		},
		{
			name:       "issue596_error_field_not_indexed",
			statusCode: http.StatusNotFound,
			body:       `{"error":"build doesn't exist or not indexed in Xray"}`,
			wantRetry:  true, // FAILS before the fix
		},
		{
			name:       "info_field_not_indexed_variant",
			statusCode: http.StatusNotFound,
			body:       `{"info":"the build is not indexed yet"}`,
			wantRetry:  true,
		},
		{
			name:       "not_found_unrelated_body_is_terminal",
			statusCode: http.StatusNotFound,
			body:       `{"error":"some unrelated 404"}`,
			wantRetry:  false,
		},
		{
			name:       "non_404_ignored",
			statusCode: http.StatusOK,
			body:       `{"info":"Build x number 1 wasn't found in Artifactory"}`,
			wantRetry:  false,
		},
		{
			name:       "unparseable_body_is_terminal",
			statusCode: http.StatusNotFound,
			body:       `not-json`,
			wantRetry:  false, // unparseable body => treat as a real 404, do not retry
		},
		{
			// Semantics: if EITHER field carries a transient phrase, retry (retry-biased / safe).
			name:       "both_fields_set_either_matches",
			statusCode: http.StatusNotFound,
			body:       `{"info":"some terminal text","error":"build not indexed"}`,
			wantRetry:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.statusCode}
			err := isArtifactoryBuildNotFoundError(resp, []byte(tt.body))
			if (err != nil) != tt.wantRetry {
				t.Fatalf("isArtifactoryBuildNotFoundError() error = %v, wantRetry = %v", err, tt.wantRetry)
			}
		})
	}
}
