package utils

import "testing"
import "reflect"

func TestCleanRepoFromPath(t *testing.T) {
	assertCleanRepoFromPath("repo/abc/def", "/abc/def", t)
	assertCleanRepoFromPath("repo/(*)", "/(*)", t)
	assertCleanRepoFromPath("repo/", "/", t)
	assertCleanRepoFromPath("/abc/def", "/abc/def", t)
	assertCleanRepoFromPath("aaa", "aaa", t)
	assertCleanRepoFromPath("", "", t)
}

func assertCleanRepoFromPath(path, expected string, t *testing.T) {
	result := CleanRepoFromPath(path)
	if expected != result {
		t.Error("Unexpected string built by CleanRepoFromPath. Expected: `" + expected + "` Got `" + result + "`")
	}
}

func TestReformatTargetByPaths(t *testing.T) {
	assertReformatTargetByPaths("repo/1(*)234", "repo/1hello234", "{1}", "hello", t)
	assertReformatTargetByPaths("/1234", "/1hello234", "{1}", "{1}", t)
	assertReformatTargetByPaths("1(2*5)6", "123456", "{1}", "2345", t)
	assertReformatTargetByPaths("repo/(*) somthing", "repo/doing somthing", "{1} somthing else", "doing somthing else", t)
	assertReformatTargetByPaths("(switch) (this)", "switch this", "{2} {1}", "this switch", t)
	assertReformatTargetByPaths("before(*)middle(*)after", "before123middle456after", "{2}{1}{2}", "456123456", t)
	assertReformatTargetByPaths("", "nothing should change", "nothing should change", "nothing should change", t)
}

func assertReformatTargetByPaths(regexp, source, dest, expected string, t *testing.T) {
	result, err := ReformatTargetByPaths(regexp, source, dest)
	if err != nil {
		t.Error(err.Error())
	}
	if expected != result {
		t.Error("Unexpected string built. Expected: `" + expected + "` Got `" + result + "`")
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
