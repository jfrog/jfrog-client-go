package evidence

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	artifactoryAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	evidenceServices "github.com/jfrog/jfrog-client-go/evidence/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareEvidence(t *testing.T) {
	tests := []struct {
		name        string
		includePAE  bool
		expectedPAE string
	}{
		{name: "without PAE"},
		{name: "with PAE", includePAE: true, expectedPAE: "DSSEv1 28 application/vnd.in-toto+json 7 payload"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/evidence/api/v1/evidence/prepare", r.URL.Path)
				assert.Equal(t, test.includePAE, r.URL.Query().Get("include_pae") == "true")

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.JSONEq(t, `{
					"predicate": {"result": "passed"},
					"predicate_type": "https://example.com/predicate/v1",
					"provider_id": "ci",
					"subject": {
						"subject_type": "entity",
						"entity_type": "gitCommit",
						"entity_id": "abc123"
					},
					"project_key": "proj",
					"attachments": [{
						"repository": "reports-local",
						"path": "reports/result.json"
					}]
				}`, string(body))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				response := map[string]any{
					"post_url":          "/evidence/api/v1/entity/gitCommit/abc123?project=proj",
					"dsse_payload":      "cGF5bG9hZA==",
					"dsse_payload_type": "application/vnd.in-toto+json",
					"attachments": []map[string]string{{
						"repository": "reports-local",
						"path":       "reports/result.json",
						"sha256":     "abc",
					}},
				}
				if test.expectedPAE != "" {
					response["pre_authentication_encoding"] = test.expectedPAE
				}
				require.NoError(t, json.NewEncoder(w).Encode(response))
			}))
			defer testServer.Close()

			serviceDetails := artifactoryAuth.NewArtifactoryDetails()
			serviceDetails.SetUrl(testServer.URL + "/evidence/")
			client, err := jfroghttpclient.JfrogClientBuilder().Build()
			require.NoError(t, err)
			service := evidenceServices.NewEvidenceService(serviceDetails, client)

			response, err := service.PrepareEvidence(evidenceServices.PrepareEvidenceRequest{
				Predicate:     json.RawMessage(`{"result":"passed"}`),
				PredicateType: "https://example.com/predicate/v1",
				ProviderID:    "ci",
				Subject: evidenceServices.PrepareEvidenceSubject{
					SubjectType: evidenceServices.SubjectTypeEntity,
					EntityType:  "gitCommit",
					EntityID:    "abc123",
				},
				ProjectKey: "proj",
				Attachments: []evidenceServices.PrepareEvidenceAttachment{{
					Repository: "reports-local",
					Path:       "reports/result.json",
				}},
			}, test.includePAE)
			require.NoError(t, err)
			assert.Equal(t, "/evidence/api/v1/entity/gitCommit/abc123?project=proj", response.PostURL)
			assert.Equal(t, "cGF5bG9hZA==", response.DSSEPayload)
			assert.Equal(t, "application/vnd.in-toto+json", response.DSSEPayloadType)
			assert.Equal(t, test.expectedPAE, response.PreAuthenticationEncoding)
			require.Len(t, response.Attachments, 1)
			assert.Equal(t, "abc", response.Attachments[0].SHA256)
		})
	}
}

func TestPrepareEvidence_ErrorResponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"errors":[{"message":"subject not found"}]}`))
		require.NoError(t, err)
	}))
	defer testServer.Close()

	serviceDetails := artifactoryAuth.NewArtifactoryDetails()
	serviceDetails.SetUrl(testServer.URL + "/evidence/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)
	service := evidenceServices.NewEvidenceService(serviceDetails, client)

	_, err = service.PrepareEvidence(evidenceServices.PrepareEvidenceRequest{
		Predicate:     json.RawMessage(`{}`),
		PredicateType: "https://example.com/predicate/v1",
		Subject: evidenceServices.PrepareEvidenceSubject{
			SubjectType: evidenceServices.SubjectTypeArtifact,
			RepoPath:    "repo/path/file",
		},
	}, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
	assert.Contains(t, err.Error(), "subject not found")
}

func TestPrepareEvidence_InvalidPredicateJSON(t *testing.T) {
	service := evidenceServices.NewEvidenceService(artifactoryAuth.NewArtifactoryDetails(), nil)

	_, err := service.PrepareEvidence(evidenceServices.PrepareEvidenceRequest{
		Predicate:     json.RawMessage(`{invalid`),
		PredicateType: "https://example.com/predicate/v1",
		Subject: evidenceServices.PrepareEvidenceSubject{
			SubjectType: evidenceServices.SubjectTypeEntity,
			EntityType:  "gitCommit",
			EntityID:    "abc123",
		},
	}, false)
	require.Error(t, err)
}
