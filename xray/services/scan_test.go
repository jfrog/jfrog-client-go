package services

import (
	"fmt"
	"github.com/jfrog/gofrog/datastructures"
	"github.com/jfrog/jfrog-client-go/utils/log"
	xrayUtils "github.com/jfrog/jfrog-client-go/xray/services/utils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
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
	// Create random trees with the following 8 IDs
	timestamp := time.Now().Unix()
	log.Info("Timestamp for test:", timestamp)
	//#nosec G404
	seed := rand.New(rand.NewSource(timestamp))
	depIds := []string{"dep1", "dep2", "dep3", "dep4", "dep5", "dep6", "dep7", "dep8"}
	tree1 := generateTreeWithIDs(depIds, seed)
	tree2 := generateTreeWithIDs(depIds, seed)
	tree3 := generateTreeWithIDs(depIds, seed)

	// Create graph
	flatGraph, err := FlattenGraph([]*xrayUtils.GraphNode{tree1, tree2, tree3})
	assert.NoError(t, err)

	// Check that the graph has been flattened correctly
	assert.Len(t, flatGraph.Nodes, 8)
	set := datastructures.MakeSet[string]()
	for _, node := range flatGraph.Nodes {
		assert.Len(t, node.Nodes, 0)
		assert.False(t, set.Exists(node.Id))
		set.Add(node.Id)
	}
}

func generateTreeWithIDs(remainingIDs []string, seed *rand.Rand) *xrayUtils.GraphNode {
	if len(remainingIDs) == 0 {
		return nil
	}
	// Shuffle IDs
	seed.Shuffle(len(remainingIDs), func(i, j int) {
		remainingIDs[i], remainingIDs[j] = remainingIDs[j], remainingIDs[i]
	})

	nodeID, remainingIDs := remainingIDs[0], remainingIDs[1:]
	node := &xrayUtils.GraphNode{Id: nodeID}

	numChildren := seed.Intn(5) + 1
	for i := 0; i < numChildren; i++ {
		child := generateTreeWithIDs(remainingIDs, seed)
		if child != nil {
			node.Nodes = append(node.Nodes, child)
		}
	}

	return node
}
