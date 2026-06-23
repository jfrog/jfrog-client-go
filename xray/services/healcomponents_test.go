package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyXrayServiceDetails struct {
	auth.CommonConfigFields
}

func (d *dummyXrayServiceDetails) GetVersion() (string, error) {
	return "", nil
}

func TestComponentsHealService_Heal_NpmNoChanges(t *testing.T) {
	input := `{"lockfileVersion":3}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, err := json.Marshal(map[string]string{"lockfile": input})
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	}))
	defer server.Close()

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)

	svc := NewComponentsHealService(client)
	svc.XrayDetails = &dummyXrayServiceDetails{}
	svc.XrayDetails.SetUrl(server.URL + "/")

	resp, disabled, err := svc.Heal(ComponentResolutionRequest{
		BuildTool: "npm",
		Repo:      "npm-virtual",
		Lockfile:  input,
	})
	require.NoError(t, err)
	assert.False(t, disabled)
	assert.Equal(t, input, resp.Lockfile)
	assert.Empty(t, resp.Changes)
}

func TestComponentsHealService_Heal_SelfHealDisabled(t *testing.T) {
	input := `{"lockfileVersion":3}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "self-heal is disabled", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)

	svc := NewComponentsHealService(client)
	svc.XrayDetails = &dummyXrayServiceDetails{}
	svc.XrayDetails.SetUrl(server.URL + "/")

	resp, disabled, err := svc.Heal(ComponentResolutionRequest{
		BuildTool: "npm",
		Repo:      "npm-virtual",
		Lockfile:  input,
	})
	require.NoError(t, err)
	assert.True(t, disabled)
	assert.Equal(t, input, resp.Lockfile)
	assert.Empty(t, resp.Changes)
}
