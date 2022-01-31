package utils

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type createPathFilePairsTest struct {
	pattern         string
	recursive       bool
	expectedTriples []RepoPathFile
}

var pathFilesDataProvider = []createPathFilePairsTest{
	{"a", true,
		[]RepoPathFile{{"r", ".", "a"}}},
	{"a/*", true,
		[]RepoPathFile{{"r", "a", "*"}, {"r", "a/*", "*"}}},
	{"a/a*b", true,
		[]RepoPathFile{{"r", "a", "a*b"}, {"r", "a/a*", "*b"}}},
	{"a/a*b*", true,
		[]RepoPathFile{{"r", "a/a*", "*b*"}, {"r", "a/a*", "*b*"}, {"r", "a/a*b*", "*"}}},
	{"a/a*b*/a/b", true,
		[]RepoPathFile{{"r", "a/a*b*/a", "b"}}},
	{"*/a*/*b*a*", true,
		[]RepoPathFile{{"r", "*/a*", "*b*a*"}, {"r", "*/a*/*", "*b*a*"}, {"r", "*/a*/*b*", "*a*"}, {"r", "*/a*/*b*a*", "*"}}},
	{"*", true,
		[]RepoPathFile{{"r", "*", "*"}}},
	{"*/*", true,
		[]RepoPathFile{{"r", "*", "*"}, {"r", "*/*", "*"}}},
	{"*/a.z", true,
		[]RepoPathFile{{"r", "*", "a.z"}}},
	{"a", false,
		[]RepoPathFile{{"r", ".", "a"}}},
	{"/*", false,
		[]RepoPathFile{{"r", "", "*"}}},
	{"/a*b", false,
		[]RepoPathFile{{"r", "", "a*b"}}},
	{"a*b*", true,
		[]RepoPathFile{{"r", "a*", "*b*"}, {"r", "a*b*", "*"}, {"r", ".", "a*b*"}}},
	{"*b*", true,
		[]RepoPathFile{{"r", "*b*", "*"}, {"r", "*", "*b*"}}},
}

type createRepoPathFileTriplesTest struct {
	pattern            string
	recursive          bool
	expectedTriples    []RepoPathFile
	expectedSingleRepo bool
}

var repoPathFilesDataProvider = []createRepoPathFileTriplesTest{
	{"a/*", true,
		[]RepoPathFile{{"a", "*", "*"}}, true},
	{"a/a*b", true,
		[]RepoPathFile{{"a", "a*", "*b"}, {"a", ".", "a*b"}}, true},
	{"a/a*b*", true,
		[]RepoPathFile{{"a", "a*b*", "*"}, {"a", "a*", "*b*"}, {"a", ".", "a*b*"}}, true},
	{"a/a*b*/a/b", true,
		[]RepoPathFile{{"a", "a*b*/a", "b"}}, true},
	{"*a/b*/*c*d*", true,
		[]RepoPathFile{{"*", "*a/b*/*c*d*", "*"}, {"*", "*a/b*/*c*", "*d*"}, {"*", "*a/b*/*", "*c*d*"}, {"*", "*a/b*", "*c*d*"},
			{"*a", "b*", "*c*d*"}, {"*a", "b*/*c*", "*d*"}, {"*a", "b*/*", "*c*d*"}, {"*a", "b*/*c*d*", "*"}}, false},
	{"*aa/b*/*c*d*", true,
		[]RepoPathFile{{"*", "*aa/b*/*c*d*", "*"}, {"*", "*aa/b*/*c*", "*d*"}, {"*", "*aa/b*/*", "*c*d*"}, {"*", "*aa/b*", "*c*d*"},
			{"*aa", "b*", "*c*d*"}, {"*aa", "b*/*c*", "*d*"}, {"*aa", "b*/*", "*c*d*"}, {"*aa", "b*/*c*d*", "*"}}, false},
	{"*/a*/*b*a*", true,
		[]RepoPathFile{{"*", "*a*/*b*a*", "*"}, {"*", "*a*", "*b*a*"}, {"*", "*a*/*b*", "*a*"}, {"*", "*a*/*", "*b*a*"}}, false},
	{"*", true,
		[]RepoPathFile{{"*", "*", "*"}}, false},
	{"*/*", true,
		[]RepoPathFile{{"*", "*", "*"}}, false},
	{"*/a.z", true,
		[]RepoPathFile{{"*", "*", "a.z"}}, false},
	{"a/b", true,
		[]RepoPathFile{{"a", ".", "b"}}, true},
	{"a/b", false,
		[]RepoPathFile{{"a", ".", "b"}}, true},
	{"a//*", false,
		[]RepoPathFile{{"a", "", "*"}}, true},
	{"r//a*b", false,
		[]RepoPathFile{{"r", "", "a*b"}}, true},
	{"a*b", true,
		[]RepoPathFile{{"a*", "*", "*b"}, {"a*b", "*", "*"}}, false},
	{"a*b*", true,
		[]RepoPathFile{{"a*", "*b*", "*"}, {"a*", "*", "*b*"}, {"a*b*", "*", "*"}}, false},
}

func TestCreatePathFilePairs(t *testing.T) {
	for _, sample := range pathFilesDataProvider {
		t.Run(sample.pattern+"_recursive_"+strconv.FormatBool(sample.recursive), func(t *testing.T) {
			validateRepoPathFile(createPathFilePairs("r", sample.pattern, sample.recursive), sample.expectedTriples, sample.pattern, t)
		})
	}
}

func TestCreateRepoPathFileTriples(t *testing.T) {
	for _, sample := range repoPathFilesDataProvider {
		t.Run(sample.pattern+"_recursive_"+strconv.FormatBool(sample.recursive), func(t *testing.T) {
			repoPathFileTriples, singleRepo, err := createRepoPathFileTriples(sample.pattern, sample.recursive)
			assert.NoError(t, err)
			assert.Equal(t, sample.expectedSingleRepo, singleRepo)
			validateRepoPathFile(repoPathFileTriples, sample.expectedTriples, sample.pattern, t)
		})
	}
}

func validateRepoPathFile(actual, expected []RepoPathFile, pattern string, t *testing.T) {
	if len(actual) != len(expected) {
		t.Errorf("Wrong triple.\nPattern:  %v\nExpected: %v\nActual:   %v", pattern, expected, actual)
	}
	for _, triple := range expected {
		found := false
		for _, actualTriple := range actual {
			if triple.repo == actualTriple.repo && triple.path == actualTriple.path && triple.file == actualTriple.file {
				found = true
			}
		}
		if found == false {
			t.Errorf("Wrong triple for pattern: '%s'. Missing %v between %v", pattern, triple, actual)
		}
	}
}
