package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubServiceDetails struct {
	auth.CommonConfigFields
}

func (d *stubServiceDetails) GetVersion() (string, error) { return "", nil }

type stubCommonConf struct {
	details auth.ServiceDetails
	client  *jfroghttpclient.JfrogHttpClient
}

func (c *stubCommonConf) GetArtifactoryDetails() auth.ServiceDetails           { return c.details }
func (c *stubCommonConf) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient { return c.client }

func aqlResultItem(name string, props map[string]string) map[string]interface{} {
	item := map[string]interface{}{"repo": "repo", "path": "dir", "name": name, "type": "file", "size": 4, "actual_sha1": "sha1-" + name}
	var properties []map[string]string
	for key, value := range props {
		properties = append(properties, map[string]string{"key": key, "value": value})
	}
	if properties != nil {
		item["properties"] = properties
	}
	return item
}

// RTECO-1596: a download (SYMLINK) gets the symlink properties from a separate AQL query, for plain, --limit,
// --sort-by and --transitive downloads. The second query is skipped when the main query finds nothing, and its empty
// result keeps the main result as is.
func TestSearchBySpecWithAqlFetchesSymlinkProps(t *testing.T) {
	symlinkProps := map[string]string{ArtifactorySymlink: "target.txt", SymlinkSha1: "abc"}
	tests := []struct {
		name            string
		params          CommonParams
		mainResults     []map[string]interface{}
		symlinkResults  []map[string]interface{}
		symlinkStatus   int
		expectedQueries int
		expectedErr     bool
		expectedProps   map[string]map[string]string
		checkQueries    func(t *testing.T, queries []string)
	}{
		{
			name:            "plain download",
			mainResults:     []map[string]interface{}{aqlResultItem("target.txt", nil), aqlResultItem("link", nil)},
			symlinkResults:  []map[string]interface{}{aqlResultItem("link", symlinkProps)},
			expectedQueries: 2,
			expectedProps:   map[string]map[string]string{"target.txt": nil, "link": symlinkProps},
			checkQueries: func(t *testing.T, queries []string) {
				assert.NotContains(t, queries[0], `"property"`)
				assert.NotContains(t, queries[1], ".transitive()")
			},
		},
		{
			name:            "main query finds nothing",
			mainResults:     []map[string]interface{}{},
			expectedQueries: 1,
			expectedProps:   map[string]map[string]string{},
		},
		{
			name:            "symlink query finds nothing",
			mainResults:     []map[string]interface{}{aqlResultItem("a.txt", nil), aqlResultItem("b.txt", nil)},
			symlinkResults:  []map[string]interface{}{},
			expectedQueries: 2,
			expectedProps:   map[string]map[string]string{"a.txt": nil, "b.txt": nil},
		},
		{
			name:            "symlink query fails",
			mainResults:     []map[string]interface{}{aqlResultItem("link", nil)},
			symlinkStatus:   http.StatusInternalServerError,
			expectedQueries: 2,
			expectedErr:     true,
		},
		{
			name:            "limit",
			params:          CommonParams{Limit: 10},
			mainResults:     []map[string]interface{}{aqlResultItem("target.txt", nil), aqlResultItem("link", nil)},
			symlinkResults:  []map[string]interface{}{aqlResultItem("link", symlinkProps)},
			expectedQueries: 2,
			// Before RTECO-1596 this path filtered by symlink.dest only, so symlink.destsha1 could be missing.
			expectedProps: map[string]map[string]string{"target.txt": nil, "link": symlinkProps},
			checkQueries: func(t *testing.T, queries []string) {
				assert.Contains(t, queries[0], ".limit(10)")
				assert.NotContains(t, queries[0], `"property"`)
			},
		},
		{
			name:            "sort by",
			params:          CommonParams{SortBy: []string{"name"}},
			mainResults:     []map[string]interface{}{aqlResultItem("link", nil), aqlResultItem("target.txt", nil)},
			symlinkResults:  []map[string]interface{}{aqlResultItem("link", symlinkProps)},
			expectedQueries: 2,
			expectedProps:   map[string]map[string]string{"target.txt": nil, "link": symlinkProps},
			checkQueries: func(t *testing.T, queries []string) {
				assert.Contains(t, queries[0], ".sort(")
			},
		},
		{
			name:            "transitive",
			params:          CommonParams{Transitive: true},
			mainResults:     []map[string]interface{}{aqlResultItem("target.txt", nil), aqlResultItem("link", nil)},
			symlinkResults:  []map[string]interface{}{aqlResultItem("link", symlinkProps)},
			expectedQueries: 2,
			expectedProps:   map[string]map[string]string{"target.txt": nil, "link": symlinkProps},
			checkQueries: func(t *testing.T, queries []string) {
				assert.True(t, strings.HasSuffix(queries[0], ".transitive()"))
				assert.True(t, strings.HasSuffix(queries[1], ".transitive()"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mutex sync.Mutex
			var queries []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				query := string(body)
				mutex.Lock()
				queries = append(queries, query)
				mutex.Unlock()
				results := test.mainResults
				if strings.Contains(query, "@"+ArtifactorySymlink) {
					if test.symlinkStatus != 0 {
						w.WriteHeader(test.symlinkStatus)
						return
					}
					results = test.symlinkResults
				}
				response, err := json.Marshal(map[string]interface{}{"results": results})
				assert.NoError(t, err)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(response)
			}))
			defer server.Close()

			details := &stubServiceDetails{}
			details.SetUrl(server.URL + "/")
			client, err := jfroghttpclient.JfrogClientBuilder().Build()
			require.NoError(t, err)

			params := test.params
			params.Aql = Aql{ItemsFind: `{"repo":"repo","path":"dir","name":{"$match":"*"}}`}
			reader, err := SearchBySpecWithAql(&params, &stubCommonConf{details: details, client: client}, SYMLINK)
			if reader != nil {
				defer readerCloseAndAssert(t, reader)
			}
			require.Len(t, queries, test.expectedQueries)
			if test.checkQueries != nil {
				test.checkQueries(t, queries)
			}
			if test.expectedErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			actualProps := map[string]map[string]string{}
			for item := new(ResultItem); reader.NextRecord(item) == nil; item = new(ResultItem) {
				var props map[string]string
				for _, prop := range item.Properties {
					if props == nil {
						props = map[string]string{}
					}
					props[prop.Key] = prop.Value
				}
				actualProps[item.Name] = props
			}
			require.NoError(t, reader.GetError())
			assert.Equal(t, test.expectedProps, actualProps)
		})
	}
}
