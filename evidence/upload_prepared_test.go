package evidence

import (
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

func TestUploadPreparedSignedEvidence(t *testing.T) {
	signedEnvelope := []byte(`{"payload":"cGF5bG9hZA==","payloadType":"application/vnd.in-toto+json","signatures":[{"keyid":"key","sig":"signature"}]}`)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/evidence/api/v1/entity/gitCommit/abc123", r.URL.Path)
		assert.Equal(t, "proj", r.URL.Query().Get("project"))
		assert.Equal(t, "ci provider", r.URL.Query().Get("providerId"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, signedEnvelope, body)

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(`{"verified":true}`))
		require.NoError(t, err)
	}))
	defer testServer.Close()

	serviceDetails := artifactoryAuth.NewArtifactoryDetails()
	serviceDetails.SetUrl(testServer.URL + "/evidence/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)
	service := evidenceServices.NewEvidenceService(serviceDetails, client)

	body, err := service.UploadPreparedSignedEvidence(
		"/evidence/api/v1/entity/gitCommit/abc123?project=proj&providerId=ci+provider",
		signedEnvelope,
	)
	require.NoError(t, err)
	assert.JSONEq(t, `{"verified":true}`, string(body))
}

func TestUploadPreparedSignedEvidence_ErrorResponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"errors":[{"message":"invalid signature"}]}`))
		require.NoError(t, err)
	}))
	defer testServer.Close()

	serviceDetails := artifactoryAuth.NewArtifactoryDetails()
	serviceDetails.SetUrl(testServer.URL + "/evidence/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)
	service := evidenceServices.NewEvidenceService(serviceDetails, client)

	_, err = service.UploadPreparedSignedEvidence("/evidence/api/v1/subject/repo/file", []byte(`{}`))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "400")
	assert.Contains(t, err.Error(), "invalid signature")
}

func TestUploadPreparedSignedEvidence_RejectsNonRootRelativeURL(t *testing.T) {
	serviceDetails := artifactoryAuth.NewArtifactoryDetails()
	serviceDetails.SetUrl("https://example.jfrog.io/evidence/")
	service := evidenceServices.NewEvidenceService(serviceDetails, nil)

	for _, postURL := range []string{
		"",
		"api/v1/entity/gitCommit/abc123",
		"https://other.example/evidence/api/v1/entity/gitCommit/abc123",
		"//other.example/evidence/api/v1/entity/gitCommit/abc123",
	} {
		t.Run(postURL, func(t *testing.T) {
			_, err := service.UploadPreparedSignedEvidence(postURL, []byte(`{}`))
			require.Error(t, err)
			assert.Contains(t, err.Error(), "root-relative")
		})
	}
}
