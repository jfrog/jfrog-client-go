package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var separator = string(os.PathSeparator)

var paths = getFileSystemsPathsForTestingAntPattern(separator)

var testAntPathToRegExpDataProvider = []struct {
	description           string
	antPattern            string
	paths                 []string
	expectedMatchingPaths []string
}{
	{"check '?' in file's name", filepath.Join("dev", "a", "b?.txt"), paths, []string{filepath.Join("dev", "a", "bb.txt"), filepath.Join("dev", "a", "bc.txt")}},
	{"check '?' in directory's name", filepath.Join("dev", "a?", "b.txt"), paths, []string{filepath.Join("dev", "aa", "b.txt")}},
	{"check '*' in file's name", filepath.Join("dev", "a", "b*.txt"), paths, []string{filepath.Join("dev", "a", "b.txt"), filepath.Join("dev", "a", "bb.txt"), filepath.Join("dev", "a", "bc.txt")}},
	{"check '*' in directory's name", filepath.Join("dev", "*", "b.txt"), paths, []string{filepath.Join("dev", "a", "b.txt"), filepath.Join("dev", "aa", "b.txt")}},
	{"check '*' in directory's name", filepath.Join("dev", "*", "a", "b.txt"), paths, nil},
	{"check '**' in directory path", filepath.Join("**", "b.txt"), paths, []string{filepath.Join("dev", "a", "b.txt"), filepath.Join("dev", "aa", "b.txt"), filepath.Join("dev", "a1", "a2", "a3", "b.txt"), filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("test", "a", "b.txt"), filepath.Join("test", "aa", "b.txt")}},
	{"check '**' in the beginning and the end of path", filepath.Join("**", "a2", "**"), paths, []string{filepath.Join("dev", "a1", "a2", "a3", "b.txt"), filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "bc.txt"), "a2"}},
	{"check '**' in the beginning and the end of path", "**a2**", paths, []string{filepath.Join("dev", "a1", "a2", "a3", "b.txt"), filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "bc.txt"), "a2"}},
	{"check double '**'", filepath.Join("**", "a2", "**", "**"), paths, []string{filepath.Join("dev", "a1", "a2", "a3", "b.txt"), filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "bc.txt")}},
	{"check '**' in the beginning and the end of file", filepath.Join("**", "b.zip", "**"), paths, []string{filepath.Join("dev", "aa", "b.zip"), filepath.Join("test", "aa", "b.zip"), "b.zip", filepath.Join("test2", "b.zip")}},
	{"combine '**' and '*'", filepath.Join("**", "a2", "*"), paths, []string{filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "bc.txt")}},
	{"combine '**' and '*'", filepath.Join("**", "a2", "*", "**"), paths, []string{filepath.Join("dev", "a1", "a2", "a3", "b.txt"), filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "bc.txt")}},
	{"combine all signs", filepath.Join("**", "b?.*"), paths, []string{filepath.Join("dev", "a", "bb.txt"), filepath.Join("dev", "a", "bc.txt"), filepath.Join("dev", "aa", "bb.txt"), filepath.Join("dev", "aa", "bc.txt"), filepath.Join("dev", "aa", "bc.zip"), filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "bc.txt"), filepath.Join("test", "a", "bb.txt"), filepath.Join("test", "a", "bc.txt"), filepath.Join("test", "aa", "bb.txt"), filepath.Join("test", "aa", "bc.txt"), filepath.Join("test", "aa", "bc.zip")}},
	{"'**' all files", "**", paths, paths},
	{"test2/**/b/**", filepath.Join("test2", "**", "b", "**"), paths, []string{filepath.Join("test2", "a", "b", "c.zip")}},
	{"*/b.zip", filepath.Join("*", "b.zip"), paths, []string{filepath.Join("test2", "b.zip")}},
	{"**/dev/**/a3/*c*", filepath.Join("dev", "**", "a3", "*c*"), paths, []string{filepath.Join("dev", "a1", "a2", "a3", "bc.txt")}},
	{"**/dev/**/a3/**", filepath.Join("dev", "**", "a3", "**"), paths, []string{filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "a3", "b.txt")}},
	{"exclude 'temp/foo5/a'", filepath.Join("**", "foo", "**"), paths, []string{filepath.Join("tmp", "foo", "a"), filepath.Join("tmp", "foo")}},
	{"include dirs", filepath.Join("tmp", "*", "**"), paths, []string{"tmp" + separator, filepath.Join("tmp", "foo", "a"), filepath.Join("tmp", "foo5", "a"), filepath.Join("tmp", "foo"), filepath.Join("tmp", "foo5")}},
	{"include dirs", filepath.Join("tmp", "**"), paths, []string{"tmp" + separator, filepath.Join("tmp", "foo", "a"), filepath.Join("tmp", "foo5", "a"), filepath.Join("tmp", "foo"), filepath.Join("tmp", "foo5")}},
	{"double and single wildcard", filepath.Join("**", "tmp*", "**"), paths, []string{"tmp" + separator, filepath.Join("tmp", "foo", "a"), filepath.Join("tmp", "foo5", "a"), filepath.Join("tmp", "foo"), filepath.Join("tmp", "foo5"), filepath.Join("Wrapper", "tmp", "boo"), filepath.Join("Wrapper", "tmp12", "boo")}},
	{"exclude only sub dir", filepath.Join("**", "loo", "**", "bar", "**"), paths, []string{filepath.Join("kmp", "loo", "bar"), filepath.Join("kmp", "loo", "bar", "b"), filepath.Join("kmp", "loo", "bar", "a")}},
	{"**/", "**" + separator, paths, paths},
	{"xxx/x*", filepath.Join("tmp", "f*"), paths, []string{filepath.Join("tmp", "foo"), filepath.Join("tmp", "foo5")}},
	{"xxx/x*x", filepath.Join("tmp", "f*5"), paths, []string{filepath.Join("tmp", "foo5")}},
	{"xxx/x*", filepath.Join("dev", "a1", "a2", "b*"), paths, []string{filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "bc.txt")}},
	{"xxx/*x*", filepath.Join("dev", "a1", "a2", "*c*"), paths, []string{filepath.Join("dev", "a1", "a2", "bc.txt")}},
	{"*", "*", paths, []string{"b.zip", "a2"}},
	{"*", "*", []string{"a", "a" + separator, filepath.Join("a", "b")}, []string{"a"}},
}

// In each case, we take an array of paths, simulating a filesystem hierarchy, and an ANT pattern expression and
// check if the conversion to regular expression worked.
func TestAntPathToRegExp(t *testing.T) {
	for _, test := range testAntPathToRegExpDataProvider {
		t.Run(test.description, func(t *testing.T) {
			regExpStr := AntToRegex(cleanPath(test.antPattern))
			var matches []string
			for _, checkedPath := range test.paths {
				match, _ := regexp.MatchString(regExpStr, checkedPath)
				if match {
					matches = append(matches, checkedPath)
				}
			}
			if !equalSlicesIgnoreOrder(matches, test.expectedMatchingPaths) {
				t.Error("Unmatched! : ant pattern `" + test.antPattern + "` matches paths:\n[" + strings.Join(test.expectedMatchingPaths, ",") + "]\nbut got:\n[" + strings.Join(matches, ",") + "]")
			}
		})
	}
}

var testAntToRegexProvider = []struct {
	ant           string
	expectedRegex string
}{
	{"a.zip", "^a\\.zip$"},
	{"ab", "^ab$"},
	{"**", "^(.*)$"},
	{"**/", "^(.*/)*(.*)$"},
	{"**/*", "^(.*/)*([^/]*)$"},
	{"/**", "^(/.*)*$"},
	{"*/**", "^([^/]*)(/.*)*$"},
	{"/**/ab", "^/(.*/)*ab$"},
	{"/**/ab*", "^/(.*/)*ab([^/]*)$"},
	{"/**/ab/", "^/(.*/)*ab(/.*)*$"},
	{"/**/ab/*", "^/(.*/)*ab/([^/]*)$"},
	{"/**/ab*/", "^/(.*/)*ab([^/]*)(/.*)*$"},
	{"ab/**/", "^ab/(.*/)*(.*)$"},
	{"*ab/**/", "^([^/]*)ab/(.*/)*(.*)$"},
	{"/ab/**/", "^/ab/(.*/)*(.*)$"},
	{"/ab*/**/", "^/ab([^/]*)/(.*/)*(.*)$"},
	{"/**/ab/**/", "^/(.*/)*ab/(.*/)*(.*)$"},
	{"/**/a*b/**/", "^/(.*/)*a([^/]*)b/(.*/)*(.*)$"},
	{"/**/ab/**/cd/**/ef/", "^/(.*/)*ab/(.*/)*cd/(.*/)*ef(/.*)*$"},
}

func TestAntToRegex(t *testing.T) {
	for _, test := range testAntToRegexProvider {
		t.Run("'"+test.ant+"'", func(t *testing.T) {
			regExpStr := AntToRegex(cleanPath(strings.ReplaceAll(test.ant, "/", separator)))
			expectedRegExpStr := strings.ReplaceAll(test.expectedRegex, "/", getFileSeparatorForAntToRegex())
			if regExpStr != expectedRegExpStr {
				t.Error("Unmatched! : ant pattern `" + test.ant + "` translated to:\n" + regExpStr + "\nbut expect it to be:\n" + expectedRegExpStr + "")
			}
		})
	}
}

func getFileSystemsPathsForTestingAntPattern(separator string) []string {
	return []string{
		filepath.Join("dev", "a", "b.txt"),
		filepath.Join("dev", "a", "bb.txt"),
		filepath.Join("dev", "a", "bc.txt"),
		filepath.Join("dev", "aa", "b.txt"),
		filepath.Join("dev", "aa", "bb.txt"),
		filepath.Join("dev", "aa", "bc.txt"),
		filepath.Join("dev", "aa", "b.zip"),
		filepath.Join("dev", "aa", "bc.zip"),
		filepath.Join("dev", "a1", "a2", "a3", "b.txt"),
		filepath.Join("dev", "a1", "a2", "b.txt"),
		filepath.Join("dev", "a1", "a2", "a3", "bc.txt"),
		filepath.Join("dev", "a1", "a2", "bc.txt"),
		"a2",
		filepath.Join("test", "a", "b.txt"),
		filepath.Join("test", "a", "bb.txt"),
		filepath.Join("test", "a", "bc.txt"),
		filepath.Join("test", "aa", "b.txt"),
		filepath.Join("test", "aa", "bb.txt"),
		filepath.Join("test", "aa", "bc.txt"),
		filepath.Join("test", "aa", "b.zip"),
		filepath.Join("test", "aa", "bc.zip"),

		filepath.Join("test2", "a", "b", "c.zip"),
		filepath.Join("test2", "a", "bb", "c.zip"),
		filepath.Join("test2", "b.zip"),
		"b.zip",
		"tmp" + separator,
		filepath.Join("tmp", "foo"),
		filepath.Join("Wrapper", "tmp", "boo"),
		filepath.Join("Wrapper", "tmp12", "boo"),
		filepath.Join("tmp", "foo", "a"),
		filepath.Join("tmp", "foo5"),
		filepath.Join("tmp", "foo5", "a"),

		filepath.Join("kmp", "loo"),
		filepath.Join("kmp", "loo", "bar", "a"),
		filepath.Join("kmp", "loo", "bar", "b"),
		filepath.Join("kmp", "loo", "bar"),
		filepath.Join("kmp", "loo", "lar"),
		filepath.Join("kmp", "loo", "lar", "a"),
		filepath.Join("kmp", "loo", "lar"),
		filepath.Join("kmp", "loo", "kar", "a"),
		filepath.Join("kmp", "loo", "kar"),
	}
}
