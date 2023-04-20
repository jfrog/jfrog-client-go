package services

import (
	"fmt"
	"github.com/jfrog/gofrog/datastructures"
	xrayUtils "github.com/jfrog/jfrog-client-go/xray/services/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateScanGraphQueryParams(t *testing.T) {
	tests := []struct {
		testName      string
		projectKey    string
		repoPath      string
		watches       []string
		scanType      ScanType
		expectedQuery string
	}{
		{"with_project_key", "p1", "", nil, Binary,
			fmt.Sprintf("?%s%s&%s%s", projectQueryParam, "p1", scanTypeQueryParam, Binary)},

		{"with_repo_path", "", "r1", nil, Binary,
			fmt.Sprintf("?%s%s&%s%s", repoPathQueryParam, "r1", scanTypeQueryParam, Binary)},

		{"with_watches", "", "", []string{"w1", "w2"}, Binary,
			fmt.Sprintf("?%s%s&%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2", scanTypeQueryParam, Binary)},

		{"with_empty_watch_string", "", "", []string{""}, "",
			""},

		{"without_context", "", "", nil, Dependency,
			fmt.Sprintf("?%s%s", scanTypeQueryParam, Dependency)},

		{"without_scan_type", "", "", []string{"w1", "w2"}, "",
			fmt.Sprintf("?%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2")},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			params := XrayGraphScanParams{
				RepoPath:   test.repoPath,
				Watches:    test.watches,
				ProjectKey: test.projectKey,
				ScanType:   test.scanType,
			}
			actualQuery := createScanGraphQueryParams(params)
			if actualQuery != test.expectedQuery {
				t.Error(test.testName, "Expecting:", test.expectedQuery, "Got:", actualQuery)
			}
		})
	}
}

func TestFlattenGraph(t *testing.T) {
	nodeA := &xrayUtils.GraphNode{Id: "A"}
	nodeB := &xrayUtils.GraphNode{Id: "B"}
	nodeC := &xrayUtils.GraphNode{Id: "C"}
	nodeD := &xrayUtils.GraphNode{Id: "D"}
	nodeE := &xrayUtils.GraphNode{Id: "E"}
	nodeF := &xrayUtils.GraphNode{Id: "F"}

	// Set dependencies
	nodeA.Nodes = []*xrayUtils.GraphNode{nodeB, nodeC}
	nodeB.Nodes = []*xrayUtils.GraphNode{nodeC, nodeD}
	nodeC.Nodes = []*xrayUtils.GraphNode{nodeD}
	nodeD.Nodes = []*xrayUtils.GraphNode{nodeE, nodeF}
	nodeF.Nodes = []*xrayUtils.GraphNode{nodeA, nodeB, nodeC}

	// Create graph
	graph := []*xrayUtils.GraphNode{nodeA, nodeB, nodeC}
	flatGraph, err := FlattenGraph(graph)
	assert.NoError(t, err)

	// Check that the graph has been flattened correctly
	assert.Equal(t, len(flatGraph[0].Nodes), 6)
	set := datastructures.MakeSet[string]()
	for _, node := range flatGraph[0].Nodes {
		assert.Len(t, node.Nodes, 0)
		assert.False(t, set.Exists(node.Id))
		set.Add(node.Id)
	}
}
