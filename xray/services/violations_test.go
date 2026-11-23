package services

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const includeIgnoredKey = "include_ignored"

func newTestViolationsService(t *testing.T, serverUrl string) *ViolationsService {
	details := &testXrayDetails{}
	details.SetUrl(serverUrl + "/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)
	service := NewViolationsService(client)
	service.XrayDetails = details
	return service
}

// Returns the "filters" object of the marshalled request, decoded as a generic map.
func marshalFilters(t *testing.T, request utils.ViolationsRequest) (raw []byte, filters map[string]any) {
	raw, err := json.Marshal(request)
	require.NoError(t, err)
	var decoded map[string]any
	require.NoError(t, json.Unmarshal(raw, &decoded))
	filters, ok := decoded["filters"].(map[string]any)
	require.True(t, ok, "request body has no filters object: %s", raw)
	return raw, filters
}

func TestViolationsRequest_DefaultHasNoIncludeIgnored(t *testing.T) {
	raw, filters := marshalFilters(t, utils.NewViolationsRequest())

	assert.NotContains(t, string(raw), includeIgnoredKey)
	assert.NotContains(t, filters, includeIgnoredKey)
	// Byte-level: the default body must stay identical to what it was before include_ignored was added.
	assert.Equal(t, `{"filters":{"resources":{}},"pagination":{"limit":9223372036854775807,"order_by":"created","direction":"asc"}}`, string(raw))
}

func TestViolationsRequest_IncludeIgnoredTrue(t *testing.T) {
	raw, filters := marshalFilters(t, utils.NewViolationsRequest().IncludeIgnored(true))

	assert.Contains(t, string(raw), `"include_ignored":true`)
	assert.Equal(t, true, filters[includeIgnoredKey])
}

func TestViolationsRequest_IncludeIgnoredFalseAfterTrueRemovesKey(t *testing.T) {
	raw, filters := marshalFilters(t, utils.NewViolationsRequest().IncludeIgnored(true).IncludeIgnored(false))

	assert.NotContains(t, string(raw), includeIgnoredKey)
	assert.NotContains(t, filters, includeIgnoredKey)
}

func TestGetViolations_SendsIncludeIgnoredToViolationsEndpoint(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.NoError(t, json.Unmarshal(body, &gotBody))
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(`{"total_violations":0,"violations":[]}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	response, err := newTestViolationsService(t, server.URL).GetViolations(utils.NewViolationsRequest().IncludeDetails(true).IncludeIgnored(true))

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/v1/violations", gotPath)
	filters, ok := gotBody["filters"].(map[string]any)
	require.True(t, ok, "server received no filters object: %v", gotBody)
	assert.Equal(t, true, filters[includeIgnoredKey])
	assert.Equal(t, true, filters["include_details"])
}

func TestGetViolations_ParsesIsIgnoredFlags(t *testing.T) {
	response := getViolationsFromServer(t, `{
		"total_violations": 2,
		"violations": [
			{
				"issue_id": "XRAY-1", "is_ignored": true,
				"matched_policies": [
					{"policy": "p1", "rule": "r1", "is_blocking": false, "is_ignored": true},
					{"policy": "p2", "rule": "r2", "is_blocking": true, "is_ignored": false}
				]
			},
			{
				"issue_id": "XRAY-2", "is_ignored": false,
				"matched_policies": [{"policy": "p3", "rule": "r3", "is_blocking": true}]
			}
		]
	}`)

	require.Len(t, response.Violations, 2)
	ignored, notIgnored := response.Violations[0], response.Violations[1]

	require.NotNil(t, ignored.IsIgnored)
	assert.True(t, *ignored.IsIgnored)
	require.Len(t, ignored.Policies, 2)
	assert.True(t, ignored.Policies[0].IsIgnored)
	assert.False(t, ignored.Policies[1].IsIgnored)

	require.NotNil(t, notIgnored.IsIgnored, "an explicit is_ignored:false must be distinguishable from a missing field")
	assert.False(t, *notIgnored.IsIgnored)
	require.Len(t, notIgnored.Policies, 1)
	assert.False(t, notIgnored.Policies[0].IsIgnored)
}

func TestGetViolations_MissingIsIgnoredIsNil(t *testing.T) {
	response := getViolationsFromServer(t, `{
		"total_violations": 1,
		"violations": [{"issue_id": "XRAY-1", "matched_policies": [{"policy": "p1", "rule": "r1", "is_blocking": true}]}]
	}`)

	require.Len(t, response.Violations, 1)
	assert.Nil(t, response.Violations[0].IsIgnored)
	require.Len(t, response.Violations[0].Policies, 1)
	assert.False(t, response.Violations[0].Policies[0].IsIgnored)
}

func TestGetViolations_ParsesIgnoreRuleInfo(t *testing.T) {
	response := getViolationsFromServer(t, `{
		"total_violations": 1,
		"violations": [{
			"issue_id": "XRAY-1",
			"is_ignored": true,
			"ignore_rule_info": {
				"id": "rule-id-1",
				"ignore_rule_type": "ignore_rule_on_violation",
				"author": "user@example.com",
				"created": "2026-11-01T10:00:00Z",
				"notes": "accepted risk",
				"is_expired": false,
				"expires_at": "2027-01-01T00:00:00Z"
			}
		}]
	}`)

	require.Len(t, response.Violations, 1)
	info := response.Violations[0].IgnoreInfo
	require.NotNil(t, info)
	assert.Equal(t, "rule-id-1", info.Id)
	assert.Equal(t, "ignore_rule_on_violation", info.Type)
	assert.Equal(t, "user@example.com", info.Author)
	assert.Equal(t, "accepted risk", info.Notes)
	assert.False(t, info.IsExpired)
	assert.Equal(t, "2027-01-01T00:00:00Z", info.ExpirationDate)
}

func TestGetViolations_UnexpectedStatusPreservesHttpResponseError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error":"bad request body"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	response, err := newTestViolationsService(t, server.URL).GetViolations(utils.NewViolationsRequest())

	require.Error(t, err)
	assert.Nil(t, response)
	var httpErr *errorutils.HttpResponseError
	require.True(t, errors.As(err, &httpErr), "expected *errorutils.HttpResponseError in the error chain, got: %v", err)
	assert.Equal(t, http.StatusBadRequest, httpErr.StatusCode)
	// The message text must stay unchanged.
	assert.True(t, strings.HasPrefix(err.Error(), "got unexpected server response while attempting to get violations:\nserver response: 400 Bad Request"), "unexpected error message: %s", err.Error())
}

// Serves the given JSON as the violations response and returns the parsed result.
func getViolationsFromServer(t *testing.T, responseBody string) *ViolationsResponse {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(responseBody))
		assert.NoError(t, err)
	}))
	defer server.Close()

	response, err := newTestViolationsService(t, server.URL).GetViolations(utils.NewViolationsRequest())
	require.NoError(t, err)
	require.NotNil(t, response)
	return response
}
