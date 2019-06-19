package utils

import (
	"strconv"
	"testing"
)

type CreateRepoPathFileTest struct {
	pattern   string
	recursive bool
	expected  []RepoPathFile
}

var pathFilesDataProvider = []CreateRepoPathFileTest{
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

var repoPathFilesDataProvider = []CreateRepoPathFileTest{
	{"a/*", true,
		[]RepoPathFile{{"a", "*", "*"}}},
	{"a/a*b", true,
		[]RepoPathFile{{"a", "a*", "*b"}, {"a", ".", "a*b"}}},
	{"a/a*b*", true,
		[]RepoPathFile{{"a", "a*b*", "*"}, {"a", "a*", "*b*"}, {"a", ".", "a*b*"}}},
	{"a/a*b*/a/b", true,
		[]RepoPathFile{{"a", "a*b*/a", "b"}}},
	{"*a/b*/*c*d*", true,
		[]RepoPathFile{{"*", "*a/b*/*c*d*", "*"}, {"*", "*a/b*/*c*", "*d*"}, {"*", "*a/b*/*", "*c*d*"}, {"*", "*a/b*", "*c*d*"},
			{"*a", "b*", "*c*d*"}, {"*a", "b*/*c*", "*d*"}, {"*a", "b*/*", "*c*d*"}, {"*a", "b*/*c*d*", "*"}}},
	{"*aa/b*/*c*d*", true,
		[]RepoPathFile{{"*", "*aa/b*/*c*d*", "*"}, {"*", "*aa/b*/*c*", "*d*"}, {"*", "*aa/b*/*", "*c*d*"}, {"*", "*aa/b*", "*c*d*"},
			{"*aa", "b*", "*c*d*"}, {"*aa", "b*/*c*", "*d*"}, {"*aa", "b*/*", "*c*d*"}, {"*aa", "b*/*c*d*", "*"}}},
	{"*/a*/*b*a*", true,
		[]RepoPathFile{{"*", "a*/*b*a*", "*"}, {"*", "a*", "*b*a*"}, {"*", "a*/*b*", "*a*"}, {"*", "a*/*", "*b*a*"}}},
	{"*", true,
		[]RepoPathFile{{"*", "*", "*"}}},
	{"*/*", true,
		[]RepoPathFile{{"*", "*", "*"}}},
	{"*/a.z", true,
		[]RepoPathFile{{"*", ".", "a.z"}}},
	{"a/b", true,
		[]RepoPathFile{{"a", ".", "b"}}},
	{"a/b", false,
		[]RepoPathFile{{"a", ".", "b"}}},
	{"a//*", false,
		[]RepoPathFile{{"a", "", "*"}}},
	{"r//a*b", false,
		[]RepoPathFile{{"r", "", "a*b"}}},
	{"a*b", true,
		[]RepoPathFile{{"a*", "*", "*b"}, {"a*b", "*", "*"}}},
	{"a*b*", true,
		[]RepoPathFile{{"a*", "*b*", "*"}, {"a*", "*", "*b*"}, {"a*b*", "*", "*"}}},
}

func TestCreatePathFilePairs(t *testing.T) {
	for _, sample := range pathFilesDataProvider {
		t.Run(sample.pattern+"_recursive_"+strconv.FormatBool(sample.recursive), func(t *testing.T) {
			validateRepoPathFile(createPathFilePairs("r", sample.pattern, sample.recursive), sample.expected, sample.pattern, t)
		})
	}
}

func TestCreateRepoPathFileTriples(t *testing.T) {
	for _, sample := range repoPathFilesDataProvider {
		t.Run(sample.pattern+"_recursive_"+strconv.FormatBool(sample.recursive), func(t *testing.T) {
			validateRepoPathFile(createRepoPathFileTriples(sample.pattern, sample.recursive), sample.expected, sample.pattern, t)
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
