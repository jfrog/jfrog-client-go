package services

import (
	"fmt"
	"github.com/jfrog/gofrog/datastructures"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateScanGraphQueryParams(t *testing.T) {
	tests := []struct {
		testName      string
		params        XrayGraphScanParams
		expectedQuery string
	}{
		{
			testName:      "with_project_key",
			params:        XrayGraphScanParams{ProjectKey: "p1", ScanType: Binary},
			expectedQuery: fmt.Sprintf("?%s%s&%s%s", projectQueryParam, "p1", scanTypeQueryParam, Binary),
		},
		{
			testName:      "with_repo_path",
			params:        XrayGraphScanParams{RepoPath: "r1", ScanType: Binary},
			expectedQuery: fmt.Sprintf("?%s%s&%s%s", repoPathQueryParam, "r1", scanTypeQueryParam, Binary),
		},
		{
			testName:      "with_watches",
			params:        XrayGraphScanParams{Watches: []string{"w1", "w2"}, ScanType: Binary},
			expectedQuery: fmt.Sprintf("?%s%s&%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2", scanTypeQueryParam, Binary),
		},
		{
			testName:      "with_empty_watch_string",
			params:        XrayGraphScanParams{Watches: []string{""}},
			expectedQuery: "",
		},
		{
			testName:      "without_context",
			params:        XrayGraphScanParams{ScanType: Dependency, XrayVersion: MinXrayVersionGitRepoKey},
			expectedQuery: fmt.Sprintf("?%s%s", scanTypeQueryParam, Dependency),
		},
		{
			testName:      "without_scan_type",
			params:        XrayGraphScanParams{Watches: []string{"w1", "w2"}},
			expectedQuery: fmt.Sprintf("?%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2"),
		},
		{
			testName:      "with_git_repo_url",
			params:        XrayGraphScanParams{GitRepoHttpsCloneUrl: "http://some-url", ScanType: Dependency, XrayVersion: MinXrayVersionGitRepoKey},
			expectedQuery: fmt.Sprintf("?%s%s&%s%s", scanTypeQueryParam, Dependency, gitRepoKeyQueryParam, "some-url.git"),
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actualQuery := createScanGraphQueryParams(test.params)
			if actualQuery != test.expectedQuery {
				t.Error(test.testName, "Expecting:", test.expectedQuery, "Got:", actualQuery)
			}
		})
	}
}

func TestFlattenGraph(t *testing.T) {
	nodeA := &GraphNode{Id: "A"}
	nodeB := &GraphNode{Id: "B"}
	nodeC := &GraphNode{Id: "C"}
	nodeD := &GraphNode{Id: "D"}
	nodeE := &GraphNode{Id: "E"}
	nodeF := &GraphNode{Id: "F"}

	// Set dependencies
	nodeA.Nodes = []*GraphNode{nodeB, nodeC}
	nodeB.Nodes = []*GraphNode{nodeC, nodeD}
	nodeC.Nodes = []*GraphNode{nodeD}
	nodeD.Nodes = []*GraphNode{nodeE, nodeF}
	nodeF.Nodes = []*GraphNode{nodeA, nodeB, nodeC}

	// Create graph
	graph := []*GraphNode{nodeA, nodeB, nodeC}
	flatGraph := FlattenGraph(graph)

	// Check that the graph has been flattened correctly
	assert.Equal(t, len(flatGraph[0].Nodes), 6)
	set := datastructures.MakeSet[string]()
	for _, node := range flatGraph[0].Nodes {
		assert.Len(t, node.Nodes, 0)
		assert.False(t, set.Exists(node.Id))
		set.Add(node.Id)
	}
}
