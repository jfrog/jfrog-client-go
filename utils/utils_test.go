package utils

import (
	"github.com/magiconair/properties/assert"
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
	assert.Equal(t, result, expected)
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
	assertBuildTargetPath("", "nothing should change", "nothing should change", "nothing should change", true, t)

}

func assertBuildTargetPath(regexp, source, dest, expected string, ignoreRepo bool, t *testing.T) {
	result, err := BuildTargetPath(regexp, source, dest, ignoreRepo)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, result, expected)
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
	assert.Equal(t, result, expected)
}

func TestPathToRegExp(t *testing.T) {
	// Unix - Absolute
	assertTestPathToRegExp("/just/a/simple/path", "^/just/a/simple/path$", t)
	assertTestPathToRegExp("/just/a/wild*card/path", "^/just/a/wild.*card/path$", t)
	assertTestPathToRegExp("/just/a/wild?card/path", "^/just/a/wild.?card/path$", t)
	assertTestPathToRegExp("/just/a/directory", "^/just/a/directory$", t)
	assertTestPathToRegExp("/directory/", "^/directory/.*$", t)
	assertTestPathToRegExp("/s.p^e$c+ial/characters", "^/s\\.p\\^e\\$c\\+ial/characters$", t)
	assertTestPathToRegExp("/al*l/toge*?th$er/", "^/al.*l/toge.*.?th\\$er/.*$", t)

	// Unix - Relative
	assertTestPathToRegExp("just/a/simple/path", "^just/a/simple/path$", t)
	assertTestPathToRegExp("just/a/wild*card/path", "^just/a/wild.*card/path$", t)
	assertTestPathToRegExp("just/a/wild?card/path", "^just/a/wild.?card/path$", t)
	assertTestPathToRegExp("just/a/directory", "^just/a/directory$", t)
	assertTestPathToRegExp("directory/", "^directory/.*$", t)
	assertTestPathToRegExp("s.p^e$c+ial/characters", "^s\\.p\\^e\\$c\\+ial/characters$", t)
	assertTestPathToRegExp("al*l/toge*?th$er/", "^al.*l/toge.*.?th\\$er/.*$", t)

	// Windows - Absolute
	assertTestPathToRegExp("C:\\just\\a\\simple\\path", "^C:\\just\\a\\simple\\path$", t)
	assertTestPathToRegExp("C:\\just\\a\\wild*card\\path", "^C:\\just\\a\\wild.*card\\path$", t)
	assertTestPathToRegExp("C:\\just\\a\\wild?card\\path", "^C:\\just\\a\\wild.?card\\path$", t)
	assertTestPathToRegExp("C:\\just\\a\\directory", "^C:\\just\\a\\directory$", t)
	assertTestPathToRegExp("C:\\directory\\", "^C:\\directory\\.*$", t)
	assertTestPathToRegExp("C:\\s.p^e$c+ial\\characters", "^C:\\s\\.p\\^e\\$c\\+ial\\characters$", t)
	assertTestPathToRegExp("C:\\al*l\\toge*?th$er\\", "^C:\\al.*l\\toge.*.?th\\$er\\.*$", t)

	// Windows - Relative
	assertTestPathToRegExp("just\\a\\simple\\path", "^just\\a\\simple\\path$", t)
	assertTestPathToRegExp("just\\a\\wild*card\\path", "^just\\a\\wild.*card\\path$", t)
	assertTestPathToRegExp("just\\a\\wild?card\\path", "^just\\a\\wild.?card\\path$", t)
	assertTestPathToRegExp("just\\a\\directory", "^just\\a\\directory$", t)
	assertTestPathToRegExp("directory\\", "^directory\\.*$", t)
	assertTestPathToRegExp("s.p^e$c+ial\\characters", "^s\\.p\\^e\\$c\\+ial\\characters$", t)
	assertTestPathToRegExp("al*l\\toge*?th$er\\", "^al.*l\\toge.*.?th\\$er\\.*$", t)
}

func assertTestPathToRegExp(path, expected string, t *testing.T) {
	assert.Equal(t, pathToRegExp(path), expected)
}
