package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePropSearchURI(t *testing.T) {
	tests := []struct {
		name   string
		uri    string
		want   SkillPropertySearchResult
		wantOK bool
	}{
		{
			name:   "valid URI",
			uri:    "https://example.jfrog.io/artifactory/api/storage/rafi-skills/4chan-reader/1.0.0/4chan-reader-1.0.0.zip",
			want:   SkillPropertySearchResult{Repo: "rafi-skills", Name: "4chan-reader", Version: "1.0.0", URI: "https://example.jfrog.io/artifactory/api/storage/rafi-skills/4chan-reader/1.0.0/4chan-reader-1.0.0.zip"},
			wantOK: true,
		},
		{
			name:   "valid URI with different version",
			uri:    "https://host.com/artifactory/api/storage/my-repo/my-skill/2.3.1/my-skill-2.3.1.zip",
			want:   SkillPropertySearchResult{Repo: "my-repo", Name: "my-skill", Version: "2.3.1", URI: "https://host.com/artifactory/api/storage/my-repo/my-skill/2.3.1/my-skill-2.3.1.zip"},
			wantOK: true,
		},
		{
			name:   "no api/storage segment",
			uri:    "https://host.com/artifactory/rafi-skills/4chan-reader/1.0.0/file.zip",
			want:   SkillPropertySearchResult{},
			wantOK: false,
		},
		{
			name:   "too few path segments",
			uri:    "https://host.com/artifactory/api/storage/repo/slug",
			want:   SkillPropertySearchResult{},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parsePropSearchURI(tt.uri)
			assert.Equal(t, tt.wantOK, ok)
			if ok {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// testServiceDetails is a minimal auth.ServiceDetails good enough to point a service
// at an httptest.Server. It can't import artifactory/auth's real implementation here:
// that package imports this one (artifactory/services), which would be a cycle.
type testServiceDetails struct {
	auth.CommonConfigFields
}

func (d *testServiceDetails) GetVersion() (string, error) { return "", nil }

// newMockSkillsServer wires a SkillsService to an httptest.Server driven by handler.
// The returned server must be closed by the caller.
func newMockSkillsServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *SkillsService) {
	t.Helper()
	testServer := httptest.NewServer(handler)

	details := &testServiceDetails{}
	details.SetUrl(testServer.URL + "/")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)

	svc := NewSkillsService(client)
	svc.ArtDetails = details
	return testServer, svc
}

// pagedVersionsHandler serves a fixed set of versions (newest first, matching real
// server order) with real limit/cursor pagination semantics: cursor is the version
// string to resume after, and nextCursor is omitted (not empty-stringed) once the
// last page has been served - mirroring the live behavior confirmed against
// agent-skills/readme-standard on bukgradlefix.jfrogdev.org.
func pagedVersionsHandler(t *testing.T, all []SkillVersion) (http.HandlerFunc, *[]*recordedRequest) {
	t.Helper()
	var requests []*recordedRequest
	return func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, &recordedRequest{RawQuery: r.URL.RawQuery})

		limit := 0
		if l := r.URL.Query().Get("limit"); l != "" {
			_, err := fmt.Sscanf(l, "%d", &limit)
			require.NoError(t, err)
		}
		cursor := r.URL.Query().Get("cursor")

		start := 0
		if cursor != "" {
			for i, v := range all {
				if v.Version == cursor {
					start = i + 1
					break
				}
			}
		}
		end := start + limit
		if limit <= 0 || end > len(all) {
			end = len(all)
		}
		page := all[start:end]

		resp := skillVersionsResponse{Items: page}
		if end < len(all) {
			resp.NextCursor = page[len(page)-1].Version
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		require.NoError(t, json.NewEncoder(w).Encode(resp))
	}, &requests
}

// recordedRequest captures just what these tests assert on.
type recordedRequest struct {
	RawQuery string
}

func TestSkillsService_ListVersions_DefaultsLimitWhenNonPositive(t *testing.T) {
	all := []SkillVersion{{Version: "1.0.1", CreatedAt: 1}}
	handler, requests := pagedVersionsHandler(t, all)
	server, svc := newMockSkillsServer(t, handler)
	defer server.Close()

	versions, nextCursor, err := svc.ListVersions("repo", "my-skill", 0, "")
	require.NoError(t, err)
	assert.Equal(t, all, versions)
	assert.Empty(t, nextCursor)
	require.Len(t, *requests, 1)
	assert.Contains(t, (*requests)[0].RawQuery, fmt.Sprintf("limit=%d", DefaultSkillVersionsLimit))
}

func TestSkillsService_ListVersions_TruncatesAndReturnsCursor(t *testing.T) {
	all := []SkillVersion{
		{Version: "1.0.3", CreatedAt: 3},
		{Version: "1.0.2", CreatedAt: 2},
		{Version: "1.0.1", CreatedAt: 1},
	}
	handler, _ := pagedVersionsHandler(t, all)
	server, svc := newMockSkillsServer(t, handler)
	defer server.Close()

	versions, nextCursor, err := svc.ListVersions("repo", "my-skill", 2, "")
	require.NoError(t, err)
	assert.Equal(t, all[:2], versions)
	assert.Equal(t, "1.0.2", nextCursor)
}

func TestSkillsService_ListVersions_CursorResumesAfterLastItem(t *testing.T) {
	all := []SkillVersion{
		{Version: "1.0.3", CreatedAt: 3},
		{Version: "1.0.2", CreatedAt: 2},
		{Version: "1.0.1", CreatedAt: 1},
	}
	handler, _ := pagedVersionsHandler(t, all)
	server, svc := newMockSkillsServer(t, handler)
	defer server.Close()

	versions, nextCursor, err := svc.ListVersions("repo", "my-skill", 2, "1.0.2")
	require.NoError(t, err)
	assert.Equal(t, all[2:], versions)
	assert.Empty(t, nextCursor)
}

func TestSkillsService_ListVersions_PaginatesUntilExhausted(t *testing.T) {
	// Simulates a skill with more versions than a single page: the caller must loop
	// on nextCursor to see all of them, and the loop must terminate (no infinite spin).
	all := make([]SkillVersion, 5)
	for i := range all {
		all[i] = SkillVersion{Version: fmt.Sprintf("1.0.%d", len(all)-i), CreatedAt: int64(len(all) - i)}
	}
	handler, requests := pagedVersionsHandler(t, all)
	server, svc := newMockSkillsServer(t, handler)
	defer server.Close()

	var collected []SkillVersion
	cursor := ""
	for i := 0; i < len(all)+1; i++ { // hard cap so a broken loop fails the test instead of hanging
		page, next, err := svc.ListVersions("repo", "my-skill", 2, cursor)
		require.NoError(t, err)
		collected = append(collected, page...)
		if next == "" {
			break
		}
		cursor = next
	}
	assert.Equal(t, all, collected)
	assert.Equal(t, 3, len(*requests)) // 2 + 2 + 1, minimum calls for 5 items at page size 2
}

func TestSkillsService_ListVersions_NotFound(t *testing.T) {
	server, svc := newMockSkillsServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	versions, nextCursor, err := svc.ListVersions("repo", "missing-skill", DefaultSkillVersionsLimit, "")
	require.Error(t, err)
	assert.Nil(t, versions)
	assert.Empty(t, nextCursor)
}

func TestSkillsService_VersionExists(t *testing.T) {
	all := []SkillVersion{{Version: "1.0.2"}, {Version: "1.0.1"}}
	handler, _ := pagedVersionsHandler(t, all)
	server, svc := newMockSkillsServer(t, handler)
	defer server.Close()

	exists, err := svc.VersionExists("repo", "my-skill", "1.0.1")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = svc.VersionExists("repo", "my-skill", "9.9.9")
	require.NoError(t, err)
	assert.False(t, exists)
}

// fixedPageSizeVersionsHandler ignores the requested limit and always serves exactly
// one item per call, regardless of DefaultSkillVersionsLimit. This simulates a server
// whose real version count exceeds a single page, so VersionExists is forced to
// follow nextCursor across multiple calls to find (or rule out) the target version.
func fixedPageSizeVersionsHandler(t *testing.T, all []SkillVersion) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		cursor := r.URL.Query().Get("cursor")
		start := 0
		if cursor != "" {
			for i, v := range all {
				if v.Version == cursor {
					start = i + 1
					break
				}
			}
		}
		resp := skillVersionsResponse{}
		if start < len(all) {
			resp.Items = all[start : start+1]
			if start+1 < len(all) {
				resp.NextCursor = resp.Items[0].Version
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		require.NoError(t, json.NewEncoder(w).Encode(resp))
	}
}

func TestSkillsService_VersionExists_SearchesBeyondFirstPage(t *testing.T) {
	all := []SkillVersion{{Version: "1.0.3"}, {Version: "1.0.2"}, {Version: "1.0.1"}}
	server, svc := newMockSkillsServer(t, fixedPageSizeVersionsHandler(t, all))
	defer server.Close()

	exists, err := svc.VersionExists("repo", "my-skill", "1.0.1")
	require.NoError(t, err, "must not stop after the first (single-item) page")
	assert.True(t, exists)

	exists, err = svc.VersionExists("repo", "my-skill", "9.9.9")
	require.NoError(t, err)
	assert.False(t, exists, "must exhaust all pages before concluding the version is absent")
}

func TestSkillsService_VersionExists_StopsOnNonAdvancingCursor(t *testing.T) {
	// A server bug returning the same cursor forever must not hang the caller.
	calls := 0
	server, svc := newMockSkillsServer(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls > 3 {
			t.Fatalf("VersionExists did not stop on a non-advancing cursor")
		}
		resp := skillVersionsResponse{
			Items:      []SkillVersion{{Version: "1.0.1"}},
			NextCursor: "stuck-cursor",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		require.NoError(t, json.NewEncoder(w).Encode(resp))
	})
	defer server.Close()

	exists, err := svc.VersionExists("repo", "my-skill", "9.9.9")
	require.NoError(t, err)
	assert.False(t, exists)
	assert.LessOrEqual(t, calls, 2, "must stop as soon as the cursor repeats")
}

func TestSkillsService_ListVersions_EscapesSlugInRequestPath(t *testing.T) {
	// r.URL.Path is decoded by net/http before the handler sees it, so the raw wire
	// form (r.RequestURI) is what actually proves the slug was escaped on the way out.
	var capturedRequestURI string
	server, svc := newMockSkillsServer(t, func(w http.ResponseWriter, r *http.Request) {
		capturedRequestURI = r.RequestURI
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		require.NoError(t, json.NewEncoder(w).Encode(skillVersionsResponse{}))
	})
	defer server.Close()

	_, _, err := svc.ListVersions("repo", "weird/slug name", 1, "")
	require.NoError(t, err)
	assert.Contains(t, capturedRequestURI, url.PathEscape("weird/slug name"))
	assert.NotContains(t, capturedRequestURI, "skills/weird/slug name/versions", "an unescaped slash would target a different path entirely")
}
