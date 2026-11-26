package utils

import (
	"sort"
	"testing"

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
			result, err := dropDirsHoldingItems(candidates, excluded)
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
