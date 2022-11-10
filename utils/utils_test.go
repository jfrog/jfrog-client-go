package utils

import (
	"github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestRemoveRepoFromPath(t *testing.T) {
	assertRemoveRepoFromPath("repo/abc/def", "/abc/def", t)
	assertRemoveRepoFromPath("repo/(*)", "/(*)", t)
	assertRemoveRepoFromPath("repo/", "/", t)
	assertRemoveRepoFromPath("/abc/def", "/abc/def", t)
	assertRemoveRepoFromPath("aaa", "aaa", t)
	assertRemoveRepoFromPath("", "", t)
}

func assertRemoveRepoFromPath(path, expected string, t *testing.T) {
	result := removeRepoFromPath(path)
	if expected != result {
		t.Error("Unexpected string built by removeRepoFromPath. Expected: `" + expected + "` Got `" + result + "`")
	}
}

func TestBuildTargetPath(t *testing.T) {
	assertBuildTargetPath("1(*)234", "1hello234", "{1}", "hello", true, t)
	assertBuildTargetPath("1234", "1hello234", "{1}", "{1}", true, t)
	assertBuildTargetPath("1(2*5)6", "123456", "{1}", "2345", true, t)
	assertBuildTargetPath("(*) something", "doing something", "{1} something else", "doing something else", true, t)
	assertBuildTargetPath("(switch) (this)", "switch this", "{2} {1}", "this switch", true, t)
	assertBuildTargetPath("before(*)middle(*)after", "before123middle456after", "{2}{1}{2}", "456123456", true, t)
	assertBuildTargetPath("foo/before(*)middle(*)after", "foo/before123middle456after", "{2}{1}{2}", "456123456", true, t)
	assertBuildTargetPath("foo/before(*)middle(*)after", "bar/before123middle456after", "{2}{1}{2}", "456123456", true, t)
	assertBuildTargetPath("foo/before(*)middle(*)after", "bar/before123middle456after", "{2}{1}{2}", "{2}{1}{2}", false, t)
	assertBuildTargetPath("foo/before(*)middle(*)", "bar/before123middle456after", "{2}{1}{2}", "456after123456after", true, t)
	assertBuildTargetPath("f(*)oo/before(*)after", "f123oo/before456after", "{2}{1}{2}", "456123456", true, t)
	assertBuildTargetPath("f(*)oo/before(*)after", "f123oo/before456after", "{2}{1}{2}", "456123456", false, t)
	assertBuildTargetPath("generic-(*)-(bar)", "generic-foo-bar/after/a.in", "{1}/{2}", "foo/bar", true, t)
	assertBuildTargetPath("generic-(*)-(bar)/(*)", "generic-foo-bar/after/a.in", "{1}/{2}/{3}", "foo/bar/after/a.in", true, t)
	assertBuildTargetPath("generic-(*)-(bar)", "generic-foo-bar/after/a.in", "{1}/{2}/after/a.in", "foo/bar/after/a.in", true, t)
	assertBuildTargetPath("", "nothing should change", "nothing should change", "nothing should change", true, t)
}

func assertBuildTargetPath(regexp, source, dest, expected string, ignoreRepo bool, t *testing.T) {
	result, _, err := BuildTargetPath(regexp, source, dest, ignoreRepo)
	if err != nil {
		t.Error(err.Error())
	}
	if expected != result {
		t.Error("Unexpected target string built. Expected: `" + expected + "` Got `" + result + "`")
	}
}

func TestSplitWithEscape(t *testing.T) {
	assertSplitWithEscape("", []string{""}, t)
	assertSplitWithEscape("a", []string{"a"}, t)
	assertSplitWithEscape("a/b", []string{"a", "b"}, t)
	assertSplitWithEscape("a/b/c", []string{"a", "b", "c"}, t)
	assertSplitWithEscape("a/b\\5/c", []string{"a", "b5", "c"}, t)
	assertSplitWithEscape("a/b\\\\5.2/c", []string{"a", "b\\5.2", "c"}, t)
	assertSplitWithEscape("a\\8/b\\5/c", []string{"a8", "b5", "c"}, t)
	assertSplitWithEscape("a\\\\8/b\\\\5.2/c", []string{"a\\8", "b\\5.2", "c"}, t)
	assertSplitWithEscape("a/b\\5/c\\0", []string{"a", "b5", "c0"}, t)
	assertSplitWithEscape("a/b\\\\5.2/c\\\\0", []string{"a", "b\\5.2", "c\\0"}, t)
}

func assertSplitWithEscape(str string, expected []string, t *testing.T) {
	result := SplitWithEscape(str, '/')
	if !reflect.DeepEqual(result, expected) {
		t.Error("Unexpected string array built. Expected: `", expected, "` Got `", result, "`")
	}
}

func TestConvertLocalPatternToRegexp(t *testing.T) {
	var tests = []struct {
		localPath string
		expected  string
	}{
		{"./", "^.*$"},
		{".\\\\", "^.*$"},
		{".\\", "^.*$"},
		{"./abc", "abc"},
		{".\\\\abc", "abc"},
		{".\\abc", "abc"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, ConvertLocalPatternToRegexp(test.localPath, RegExp))
	}
}
func TestCleanPath(t *testing.T) {
	if io.IsWindows() {
		parameter := "\\\\foo\\\\baz\\\\..\\\\bar\\\\*"
		got := cleanPath(parameter)
		want := "\\\\foo\\\\bar\\\\*"
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
		parameter = "\\\\foo\\\\\\\\bar\\\\*"
		got = cleanPath(parameter)
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
		parameter = "\\\\foo\\\\.\\\\bar\\\\*"
		got = cleanPath(parameter)
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
		parameter = "\\\\foo\\\\.\\\\bar\\\\*\\\\"
		want = "\\\\foo\\\\bar\\\\*\\\\"
		got = cleanPath(parameter)
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
		parameter = "foo\\\\bar"
		got = cleanPath(parameter)
		want = "foo\\\\bar"
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
		parameter = ".\\\\foo\\\\bar\\\\"
		got = cleanPath(parameter)
		want = "foo\\\\bar\\\\"
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
	} else {
		parameter := "/foo/bar/"
		got := cleanPath(parameter)
		want := "/foo/bar/"
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
		parameter = "/foo/baz/../bar/*"
		got = cleanPath(parameter)
		want = "/foo/bar/*"
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
		parameter = "/foo//bar/*"
		got = cleanPath(parameter)
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
		parameter = "/foo/./bar/*"
		got = cleanPath(parameter)
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
		parameter = "/foo/./bar/*/"
		want = "/foo/bar/*/"
		got = cleanPath(parameter)
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
		parameter = "foo/bar"
		got = cleanPath(parameter)
		want = "foo/bar"
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
		parameter = "./foo/bar/"
		got = cleanPath(parameter)
		want = "foo/bar/"
		if got != want {
			t.Errorf("cleanPath(%s) == %s, want %s", parameter, got, want)
		}
	}
}
func TestIsWildcardParentheses(t *testing.T) {
	strA := "/tmp/cache/download/(github.com/)"
	strB := "/tmp/cache/download/(github.com/*)"
	parenthesesA := CreateParenthesesSlice(strA, "")
	parenthesesB := CreateParenthesesSlice(strA, "{1}")

	got := isWildcardParentheses(strA, parenthesesA)
	want := false
	if got != want {
		t.Errorf("TestIsWildcardParentheses() == %t, want %t", got, want)
	}

	got = isWildcardParentheses(strB, parenthesesB)
	want = true
	if got != want {
		t.Errorf("TestIsWildcardParentheses() == %t, want %t", got, want)
	}
}

// TestAntPathToRegExp check the functionality of antPatternToRegExp function.
// Each time we take an array of paths, simulating files hierarchy on a filesystem, and an ANT pattern expression -
//
//	and see if the conversion to regular expression worked as expected.
func TestAntPathToRegExp(t *testing.T) {
	separator := getFileSeparator()
	var paths = getFileSystemsPathsForTestingAntPattern(separator)
	tests := []struct {
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
		{"check '**' in the beginning and the end of path", filepath.Join("**", "a2", "**"), paths, []string{filepath.Join("dev", "a1", "a2", "a3", "b.txt"), filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "bc.txt")}},
		{"check '**' in the beginning and the end of path", filepath.Join("**a2**"), paths, []string{filepath.Join("dev", "a1", "a2", "a3", "b.txt"), filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "bc.txt")}},
		{"check double '**'", filepath.Join("**", "a2", "**", "**"), paths, []string{filepath.Join("dev", "a1", "a2", "a3", "b.txt"), filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "bc.txt")}},
		{"check '**' in the beginning and the end of file", filepath.Join("**", "b.zip", "**"), paths, []string{filepath.Join("dev", "aa", "b.zip"), filepath.Join("test", "aa", "b.zip"), filepath.Join("b.zip"), filepath.Join("test2", "b.zip")}},
		{"combine '**' and '*'", filepath.Join("**", "a2", "*"), paths, []string{filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "bc.txt")}},
		{"combine '**' and '*'", filepath.Join("**", "a2", "*", "**"), paths, []string{filepath.Join("dev", "a1", "a2", "a3", "b.txt"), filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "bc.txt")}},
		{"combine all signs", filepath.Join("**", "b?.*"), paths, []string{filepath.Join("dev", "a", "bb.txt"), filepath.Join("dev", "a", "bc.txt"), filepath.Join("dev", "aa", "bb.txt"), filepath.Join("dev", "aa", "bc.txt"), filepath.Join("dev", "aa", "bc.zip"), filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "bc.txt"), filepath.Join("test", "a", "bb.txt"), filepath.Join("test", "a", "bc.txt"), filepath.Join("test", "aa", "bb.txt"), filepath.Join("test", "aa", "bc.txt"), filepath.Join("test", "aa", "bc.zip")}},
		{"'**' all files", filepath.Join("**"), paths, paths},
		{"test2/**/b/**", filepath.Join("test2", "**", "b", "**"), paths, []string{filepath.Join("test2", "a", "b", "c.zip")}},
		{"*/b.zip", filepath.Join("*", "b.zip"), paths, []string{filepath.Join("test2", "b.zip")}},
		{"**/dev/**/a3/*c*", filepath.Join("dev", "**", "a3", "*c*"), paths, []string{filepath.Join("dev", "a1", "a2", "a3", "bc.txt")}},
		{"**/dev/**/a3/**", filepath.Join("dev", "**", "a3", "**"), paths, []string{filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "a3", "b.txt")}},
		{"exclude 'temp/foo5/a'", filepath.Join("**", "foo", "**"), paths, []string{filepath.Join("tmp", "foo", "a"), filepath.Join("tmp", "foo")}},
		{"include dirs", filepath.Join("tmp", "*", "**"), paths, []string{filepath.Join("tmp", "foo", "a"), filepath.Join("tmp", "foo5", "a"), filepath.Join("tmp", "foo"), filepath.Join("tmp", "foo5")}},
		{"include dirs", filepath.Join("tmp", "**"), paths, []string{filepath.Join("tmp", "foo", "a"), filepath.Join("tmp", "foo5", "a"), filepath.Join("tmp", "foo"), filepath.Join("tmp", "foo5"), "tmp" + separator}},
		{"**/", "**" + separator, paths, paths},
		{"xxx/x*", filepath.Join("tmp", "f*"), paths, []string{filepath.Join("tmp", "foo"), filepath.Join("tmp", "foo5")}},
		{"xxx/x*x", filepath.Join("tmp", "f*5"), paths, []string{filepath.Join("tmp", "foo5")}},
		{"xxx/x*", filepath.Join("dev", "a1", "a2", "b*"), paths, []string{filepath.Join("dev", "a1", "a2", "b.txt"), filepath.Join("dev", "a1", "a2", "bc.txt")}},
		{"xxx/*x*", filepath.Join("dev", "a1", "a2", "*c*"), paths, []string{filepath.Join("dev", "a1", "a2", "bc.txt")}},
		{"*", filepath.Join("*"), []string{"a", "a" + separator, filepath.Join("a", "b")}, []string{"a"}},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			regExpStr := antPatternToRegExp(cleanPath(test.antPattern))
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
		filepath.Join("b.zip"),
		"tmp" + separator,
		filepath.Join("tmp", "foo"),
		filepath.Join("tmp", "foo", "a"),
		filepath.Join("tmp", "foo5"),
		filepath.Join("tmp", "foo5", "a"),
	}
}

func equalSlicesIgnoreOrder(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	sort.Strings(s1)
	sort.Strings(s2)
	return reflect.DeepEqual(s1, s2)
}
