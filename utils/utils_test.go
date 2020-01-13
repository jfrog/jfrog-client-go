package utils

import (
	"reflect"
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
	assertBuildTargetPath("(*) somthing", "doing somthing", "{1} somthing else", "doing somthing else", true, t)
	assertBuildTargetPath("(switch) (this)", "switch this", "{2} {1}", "this switch", true, t)
	assertBuildTargetPath("before(*)middle(*)after", "before123middle456after", "{2}{1}{2}", "456123456", true, t)
	assertBuildTargetPath("foo/before(*)middle(*)after", "foo/before123middle456after", "{2}{1}{2}", "456123456", true, t)
	assertBuildTargetPath("foo/before(*)middle(*)after", "bar/before123middle456after", "{2}{1}{2}", "456123456", true, t)
	assertBuildTargetPath("foo/before(*)middle(*)after", "bar/before123middle456after", "{2}{1}{2}", "{2}{1}{2}", false, t)
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
