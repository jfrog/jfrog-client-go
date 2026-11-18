package services

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RTECO-1596: a download must not include "property" in its main AQL query, because Artifactory then joins the item
// properties table, which can be very slow on large instances. The symlink properties come from a separate query.
func TestDownloadFetchesSymlinkPropertiesInSeparateQuery(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks on Windows may require elevated privileges")
	}
	out := t.TempDir()
	var mutex sync.Mutex
	var queries []string
	var linkRequests atomic.Int32

	item := func(name string, props map[string]string) map[string]interface{} {
		result := map[string]interface{}{
			"repo": "repo", "path": "dir", "name": name, "type": "file", "size": 4,
			"actual_sha1": "a17c9aaa61e80a1bf71d0d850af4e5baa9800bbd", "actual_md5": "8d777f385d3dfec8815d20f7496026dc",
		}
		var properties []map[string]string
		for key, value := range props {
			properties = append(properties, map[string]string{"key": key, "value": value})
		}
		if properties != nil {
			result["properties"] = properties
		}
		return result
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "api/search/aql") {
			body, readErr := io.ReadAll(r.Body)
			assert.NoError(t, readErr)
			query := string(body)
			mutex.Lock()
			queries = append(queries, query)
			mutex.Unlock()
			results := []map[string]interface{}{item("target.txt", nil), item("link", nil)}
			if strings.Contains(query, "symlink.dest") {
				results = []map[string]interface{}{item("link", map[string]string{"symlink.dest": "target.txt"})}
			}
			response, marshalErr := json.Marshal(map[string]interface{}{"results": results})
			assert.NoError(t, marshalErr)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(response)
			return
		}
		if r.URL.Path == "/repo/dir/link" {
			linkRequests.Add(1)
		}
		_, _ = w.Write([]byte("data"))
	}))
	defer server.Close()

	details := &testServiceDetails{}
	details.SetUrl(server.URL + "/")
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	require.NoError(t, err)
	ds := NewDownloadService(details, client)
	ds.SetThreads(1)

	params := NewDownloadParams()
	params.Pattern = "repo/dir/*"
	params.Target = out + string(filepath.Separator)
	params.Symlink = true
	params.Flat = true

	summary, err := ds.DownloadFiles(params)
	require.NoError(t, err)
	assert.Equal(t, 2, summary.TotalSucceeded)

	require.Len(t, queries, 2)
	assert.NotContains(t, queries[0], `"property"`, "the main query must not include the item properties")
	assert.Contains(t, queries[1], `"@symlink.dest"`)
	assert.Contains(t, queries[1], `"@symlink.destsha1"`)

	// The link is created from the properties of the second query, not downloaded.
	dest, readlinkErr := os.Readlink(filepath.Join(out, "link"))
	require.NoError(t, readlinkErr)
	assert.Equal(t, "target.txt", dest)
	assert.Equal(t, int32(0), linkRequests.Load())
}
