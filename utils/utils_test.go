package utils

import (
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
	result, err := BuildTargetPath(regexp, source, dest, ignoreRepo)
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

func TestCleanPath(t *testing.T) {
	if IsWindows() {
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
	parenthesesA := NewParenthesesSlice(strA, "")
	parenthesesB := NewParenthesesSlice(strA, "{1}")

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

func TestAntPathToRegExp(t *testing.T) {
	var fileSystemPaths []string = []string{
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
	}
	s := getFileSeparator()
	tests := []struct {
		name               string
		antPattern         string
		expectedRegExp     string
		allFileSystemPaths []string
		matchedPaths       []string
	}{
		{"check '?' in file's name", filepath.Join("dev", "a", "b?.txt"), addRegExpPrefixAndSuffix("dev" + s + "a" + s + "b.{1}\\.txt"), fileSystemPaths, []string{filepath.Join("dev", "a", "bb.txt"), filepath.Join("dev", "a", "bc.txt")}},
		{"check '?' in directory's name", filepath.Join("dev", "a?", "b.txt"), addRegExpPrefixAndSuffix("dev" + s + "a.{1}" + s + "b\\.txt"), fileSystemPaths, []string{filepath.Join("dev", "aa", "b.txt")}},
		{"check '*' in file's name", filepath.Join("dev", "a", "b*.txt"), addRegExpPrefixAndSuffix("dev" + s + "a" + s + "b([^" + s + "]*)\\.txt"), fileSystemPaths, []string{filepath.Join("dev", "a", "b.txt"), filepath.Join("dev", "a", "bb.txt"), filepath.Join("dev", "a", "bc.txt")}},
		{"check '*' in directory's name", filepath.Join("dev", "*", "b.txt"), addRegExpPrefixAndSuffix("dev" + s + "([^" + s + "]*)" + s + "b\\.txt"), fileSystemPaths, []string{filepath.Join("dev", "a", "b.txt"), filepath.Join("dev", "aa", "b.txt")}},
		{"check '**' in directory path", filepath.Join("**", "b.txt"), addRegExpPrefixAndSuffix("(.*" + s + ")?b\\.txt"), fileSystemPaths, []string{filepath.Join("dev", "a", "b.txt"), filepath.Join("dev", "aa", "b.txt"), filepath.Join("test", "a", "b.txt"), filepath.Join("test", "aa", "b.txt"), filepath.Join("dev", "a1", "a2", "a3", "b.txt"), filepath.Join("dev", "a1", "a2", "b.txt")}},
		{"combine all signs", filepath.Join("**", "b?.*"), addRegExpPrefixAndSuffix("(.*" + s + ")?b.{1}\\.([^" + s + "]*)"), fileSystemPaths, []string{filepath.Join("dev", "a", "bb.txt"), filepath.Join("dev", "a", "bc.txt"), filepath.Join("dev", "aa", "bb.txt"), filepath.Join("dev", "aa", "bc.txt"), filepath.Join("dev", "aa", "bc.zip"), filepath.Join("dev", "a1", "a2", "a3", "bc.txt"), filepath.Join("dev", "a1", "a2", "bc.txt"), filepath.Join("test", "a", "bb.txt"), filepath.Join("test", "a", "bc.txt"), filepath.Join("test", "aa", "bb.txt"), filepath.Join("test", "aa", "bc.txt"), filepath.Join("test", "aa", "bc.zip")}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			regExpStr := antPatternToRegExp(test.antPattern)
			if test.expectedRegExp != regExpStr {
				t.Error("Unmatched! for ant pattern `" + test.antPattern + "\n` : Expected `" + test.expectedRegExp + "` Got `" + regExpStr + "`")
			}
			var matches []string
			for _, checkedPath := range fileSystemPaths {
				match, _ := regexp.MatchString(regExpStr, checkedPath)
				if match {
					matches = append(matches, checkedPath)
				}
			}
			if !equalSlicesIgnoreOrder(matches, test.matchedPaths) {
				t.Error("Unmatched! : ant pattern `" + test.antPattern + "` matches paths:\n[" + strings.Join(test.matchedPaths, ",") + "]\nbut got:\n[" + strings.Join(matches, ",") + "]")
			}
		})
	}
}

func addRegExpPrefixAndSuffix(str string) string {
	return "^" + str + "$"
}

func equalSlicesIgnoreOrder(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	sort.Strings(s1)
	sort.Strings(s2)
	return reflect.DeepEqual(s1, s2)
}
