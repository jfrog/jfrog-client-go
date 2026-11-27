package utils

import (
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/stretchr/testify/assert"
)

// exItem builds a ResultItem in the test repository with the given type, path and name.
func exItem(itemType ResultItemType, path, name string) ResultItem {
	return ResultItem{Repo: "repo", Path: path, Name: name, Type: string(itemType)}
}

func exReader(t *testing.T, items ...ResultItem) *content.ContentReader {
	writer, err := content.NewContentWriter(content.DefaultKey, true, false)
	assert.NoError(t, err)
	for _, item := range items {
		writer.Write(item)
	}
	assert.NoError(t, writer.Close())
	return content.NewContentReader(writer.GetFilePath(), content.DefaultKey)
}

func exNames(t *testing.T, reader *content.ContentReader) []string {
	var names []string
	for item := new(ResultItem); reader.NextRecord(item) == nil; item = new(ResultItem) {
		names = append(names, item.Path+"/"+item.Name)
	}
	assert.NoError(t, reader.GetError())
	sort.Strings(names)
	return names
}

func TestDropDirsHoldingItems(t *testing.T) {
	sub := exItem(Folder, "a", "sub")
	top := exItem(Folder, ".", "a")
	topLong := exItem(Folder, ".", "ab")
	topDash := exItem(Folder, ".", "a-b")
	sub2 := exItem(Folder, "a", "sub2")
	other := exItem(Folder, "a", "other")
	subSub := exItem(Folder, "a/sub", "sub")
	subB := exItem(Folder, "a/sub", "b")
	one := exItem(File, "a", "one.txt")
	four := exItem(File, "a/sub", "four.txt")
	five := exItem(File, "a/sub", "five.txt")
	six := exItem(File, "a/sub/b", "six.txt")
	deepFour := exItem(File, "a/sub/sub", "four.txt")
	sub2File := exItem(File, "a/sub2", "x.txt")
	tests := []struct {
		name       string
		candidates []ResultItem
		excluded   []ResultItem
		expected   []string
	}{
		{"folder with an excluded file is dropped, other files stay",
			[]ResultItem{one, sub, five}, []ResultItem{four}, []string{"a/one.txt", "a/sub/five.txt"}},
		{"unrelated folder is kept",
			[]ResultItem{sub, other, one}, []ResultItem{four}, []string{"a/one.txt", "a/other"}},
		{"sibling folder with the same name prefix is kept",
			[]ResultItem{sub, sub2, sub2File}, []ResultItem{four}, []string{"a/sub2", "a/sub2/x.txt"}},
		{"nested folders: only folders on the excluded path are dropped",
			[]ResultItem{sub, subSub, subB, five, six}, []ResultItem{deepFour}, []string{"a/sub/b", "a/sub/b/six.txt", "a/sub/five.txt"}},
		{"excluded folder drops its parent folder",
			[]ResultItem{top, one, other}, []ResultItem{sub}, []string{"a/one.txt", "a/other"}},
		{"top level folder is kept when only a folder with a longer name holds the excluded item",
			[]ResultItem{top, topLong, one}, []ResultItem{exItem(File, "ab", "four.txt")}, []string{"./a", "a/one.txt"}},
		{"top level folder with a dash sibling",
			[]ResultItem{top, topDash, one}, []ResultItem{exItem(File, "a-b", "four.txt")}, []string{"./a", "a/one.txt"}},
		{"nothing excluded keeps everything",
			[]ResultItem{sub, one}, nil, []string{"a/one.txt", "a/sub"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidates, excluded := exReader(t, test.candidates...), exReader(t, test.excluded...)
			defer func() {
				assert.NoError(t, candidates.Close())
				assert.NoError(t, excluded.Close())
			}()
			result, _, err := dropDirsHoldingItems(candidates, excluded)
			assert.NoError(t, err)
			defer func() { assert.NoError(t, result.Close()) }()
			assert.Equal(t, test.expected, exNames(t, result))
		})
	}
}

func TestBuildExcludedOnlyQuery(t *testing.T) {
	params := CommonParams{Pattern: "repo/a/*", Recursive: true, Exclusions: []string{"*/*four*"}, ExcludedOnly: true}
	aql, err := CreateAqlBodyForSpecWithPattern(&params)
	assert.NoError(t, err)
	assert.Contains(t, aql, `"$and":[{"$or":[{"$and":[{"path":{"$match":"*"}},{"name":{"$match":"*four*"}}]}`)
	assert.NotContains(t, aql, `$nmatch`)

	params.ExcludedOnly = false
	aql, err = CreateAqlBodyForSpecWithPattern(&params)
	assert.NoError(t, err)
	assert.Contains(t, aql, `$nmatch`)
	assert.NotContains(t, aql, `"$match":"*four*"`)
}

func TestExcludedItemsParamsSearchesNestedItemsWhenNotRecursive(t *testing.T) {
	params := CommonParams{Pattern: "repo/*", Recursive: false, Exclusions: []string{"repo/a/keep.txt"}, Props: "p=1", Limit: 5}
	excludedParams := excludedItemsParams(&params)
	assert.True(t, excludedParams.Recursive)
	assert.True(t, excludedParams.ExcludedOnly)
	assert.Empty(t, excludedParams.Props)
	assert.Zero(t, excludedParams.Limit)
	// The caller's spec is not changed.
	assert.False(t, params.Recursive)
	assert.False(t, params.ExcludedOnly)
	assert.Equal(t, "p=1", params.Props)

	// The nested exclusion must be reachable: the query is not limited to the first level.
	body, err := CreateAqlBodyForSpecWithPattern(&excludedParams)
	assert.NoError(t, err)
	assert.Contains(t, body, `{"path":{"$match":"a"}}`)
	assert.Contains(t, body, `{"name":{"$match":"keep.txt"}}`)
}

func TestExcludedItemsParamsIgnoresNameOfSourcePattern(t *testing.T) {
	params := CommonParams{Pattern: "repo/*-release", Exclusions: []string{"repo/a-release/keep.txt"}}
	excludedParams := excludedItemsParams(&params)
	assert.Equal(t, "repo/*", excludedParams.Pattern)
	// The caller's spec is not changed.
	assert.Equal(t, "repo/*-release", params.Pattern)

	body, err := CreateAqlBodyForSpecWithPattern(&excludedParams)
	assert.NoError(t, err)
	// The source part of the query matches every item of the repository.
	assert.Contains(t, body, `"path":{"$match":"*"},"name":{"$match":"*"}`)
	assert.Contains(t, body, `{"name":{"$match":"keep.txt"}}`)
}

func TestExcludedOnlyQueryRejectsBundle(t *testing.T) {
	params := CommonParams{Pattern: "repo/*", Exclusions: []string{"*/*four*"}, ExcludedOnly: true, Bundle: "name/1"}
	_, err := CreateAqlBodyForSpecWithPattern(&params)
	assert.Error(t, err)
}

func TestDropDirsHoldingItemsReturnsDirsToExpand(t *testing.T) {
	// repo/a/sub holds an excluded file, and so does repo/a/sub/sub, a folder inside it.
	// repo/a/gone is an excluded folder, and repo/a/other holds nothing excluded.
	candidates := exReader(t, exItem(Folder, "a", "sub"), exItem(Folder, "a/sub", "sub"), exItem(Folder, "a", "gone"),
		exItem(Folder, "a", "other"), exItem(Folder, ".", "a"))
	excluded := exReader(t, exItem(File, "a/sub", "four.txt"), exItem(File, "a/sub/sub", "four.txt"), exItem(Folder, "a", "gone"))
	defer func() {
		assert.NoError(t, candidates.Close())
		assert.NoError(t, excluded.Close())
	}()
	result, dirs, err := dropDirsHoldingItems(candidates, excluded)
	assert.NoError(t, err)
	defer func() { assert.NoError(t, result.Close()) }()
	assert.Equal(t, []string{"a/other"}, exNames(t, result))
	// repo/a is a holder too, and it contains all the others, so only it is searched again.
	assert.Equal(t, []string{"repo/a"}, dirs)
}

func TestDropDirsHoldingItemsKeepsRootFolderItem(t *testing.T) {
	root := ResultItem{Repo: "repo", Path: ".", Name: ".", Type: string(Folder)}
	candidates, excluded := exReader(t, root, exItem(Folder, "a", "sub")), exReader(t, exItem(File, "a/sub", "four.txt"))
	defer func() {
		assert.NoError(t, candidates.Close())
		assert.NoError(t, excluded.Close())
	}()
	result, _, err := dropDirsHoldingItems(candidates, excluded)
	assert.NoError(t, err)
	defer func() { assert.NoError(t, result.Close()) }()
	assert.Equal(t, []string{"./."}, exNames(t, result))
}

func TestCollectItemDirs(t *testing.T) {
	items := exReader(t, exItem(File, "a/sub", "four.txt"), exItem(Folder, "b", "gone"), exItem(File, ".", "top.txt"))
	defer func() { assert.NoError(t, items.Close()) }()
	holders, itemDirs, err := collectItemDirs(items)
	assert.NoError(t, err)
	assert.Equal(t, map[string]bool{"repo/a/sub": true, "repo/a": true, "repo/b": true}, holders)
	assert.Equal(t, map[string]bool{"repo/b/gone": true}, itemDirs)
}

func TestSortWithoutDuplicates(t *testing.T) {
	items := exReader(t, exItem(File, "a", "two.txt"), exItem(File, "a", "one.txt"), exItem(File, "a", "two.txt"))
	defer func() { assert.NoError(t, items.Close()) }()
	result, err := sortWithoutDuplicates(items)
	assert.NoError(t, err)
	defer func() { assert.NoError(t, result.Close()) }()
	assert.Equal(t, []string{"a/one.txt", "a/two.txt"}, exNames(t, result))
}

func TestRepositoryPattern(t *testing.T) {
	tests := map[string]string{
		"repo/a/*":     "repo/*",
		"repo/*-rel":   "repo/*",
		"(repo)/(*)/b": "repo/*",
		"/repo/a":      "repo/*",
		"":             "*/*",
		"repo":         "repo/*",
	}
	for pattern, expected := range tests {
		assert.Equal(t, expected, repositoryPattern(pattern), pattern)
	}
}

func TestExcludedItemsParamsClearsArchiveEntries(t *testing.T) {
	params := CommonParams{Pattern: "repo/a", ArchiveEntries: "*.class"}
	assert.Empty(t, excludedItemsParams(&params).ArchiveEntries)
	assert.Equal(t, "*.class", params.ArchiveEntries)
}

func TestPlaceholderDetection(t *testing.T) {
	// The check used by RemoveDirsHoldingExcludedItems: only a parenthesis with a matching {n} in the target counts.
	hasPlaceholders := func(pattern, target string) bool {
		return clientutils.RemovePlaceholderParentheses(pattern, target) != pattern
	}
	assert.False(t, hasPlaceholders("repo/a/*", "dst/"))
	assert.False(t, hasPlaceholders("repo/a/*", ""))
	assert.False(t, hasPlaceholders("repo/a/(*)", "dst/"))
	assert.True(t, hasPlaceholders("repo/a/(*)", "dst/{1}"))
}

func TestLimitItems(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		expected []string
	}{
		{"limit below length", 2, []string{"a/1.txt", "a/2.txt"}},
		{"limit equal to length", 3, []string{"a/1.txt", "a/2.txt", "a/3.txt"}},
		{"limit above length", 10, []string{"a/1.txt", "a/2.txt", "a/3.txt"}},
		{"limit one", 1, []string{"a/1.txt"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := exReader(t, exItem(File, "a", "1.txt"), exItem(File, "a", "2.txt"), exItem(File, "a", "3.txt"))
			result, err := limitItems(reader, test.limit)
			assert.NoError(t, err)
			defer func() { assert.NoError(t, result.Close()) }()
			assert.Equal(t, test.expected, exNames(t, result))
		})
	}
}

// A reader with more records than the channel buffer must be read to its end by copyItems.
// Otherwise the reading goroutine blocks, and Reset (which waits for it) never returns.
func TestCopyItemsLetsTheReaderFinishOnLargeInput(t *testing.T) {
	items := make([]ResultItem, 0, clientutils.MaxBufferSize+100)
	for i := 0; i < clientutils.MaxBufferSize+100; i++ {
		items = append(items, exItem(File, "a", fmt.Sprintf("%06d.txt", i)))
	}
	reader := exReader(t, items...)
	defer func() { assert.NoError(t, reader.Close()) }()
	result, err := copyItems(reader, 5)
	assert.NoError(t, err)
	defer func() { assert.NoError(t, result.Close()) }()
	assert.Len(t, exNames(t, result), 5)

	resetDone := make(chan struct{})
	go func() {
		reader.Reset()
		close(resetDone)
	}()
	select {
	case <-resetDone:
	case <-time.After(30 * time.Second):
		t.Fatal("the reading goroutine of the input reader did not finish")
	}
}

// A reader that was closed first removes the file of the writer. Discarding the writer must not report this as an error.
func TestDiscardWriterIgnoresFileThatIsAlreadyRemoved(t *testing.T) {
	writer, err := content.NewContentWriter(content.DefaultKey, true, false)
	assert.NoError(t, err)
	writer.Write(exItem(File, "a", "1.txt"))
	assert.NoError(t, writer.Close())
	reader := content.NewContentReader(writer.GetFilePath(), content.DefaultKey)
	assert.NoError(t, reader.Close())
	assert.NoError(t, discardWriter(writer, true))
}

func TestDiscardWriterRemovesFile(t *testing.T) {
	writer, err := content.NewContentWriter(content.DefaultKey, true, false)
	assert.NoError(t, err)
	writer.Write(exItem(File, "a", "1.txt"))
	assert.NoError(t, discardWriter(writer, false))
	_, statErr := os.Stat(writer.GetFilePath())
	assert.True(t, os.IsNotExist(statErr))
}
