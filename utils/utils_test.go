package utils

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
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

func equalSlicesIgnoreOrder(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	sort.Strings(s1)
	sort.Strings(s2)
	return reflect.DeepEqual(s1, s2)
}

func TestGetMaxPlaceholderIndex(t *testing.T) {
	type args struct {
		toReplace string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{"empty", args{""}, 0, nil},
		{"empty", args{"{}"}, 0, nil},
		{"basic", args{"{1}{5}{3}"}, 5, nil},
		{"basic", args{"}5{{3}"}, 3, nil},
		{"basic", args{"{1}}5}{3}"}, 3, nil},
		{"basic", args{"{1}5{}}{3}"}, 3, nil},
		{"special characters", args{"!@#$%^&*abc(){}}{{2}!@#$%^&*abc(){}}{{1}!@#$%^&*abc(){}}{"}, 2, nil},
		{"multiple digits", args{"{2}{100}fdsff{101}d#%{99}"}, 101, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getMaxPlaceholderIndex(tt.args.toReplace)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "getMaxPlaceholderIndex(%v)", tt.args.toReplace)
		})
	}
}

func TestReplacePlaceHolders(t *testing.T) {
	type args struct {
		groups    []string
		toReplace string
		isRegexp  bool
	}
	tests := []struct {
		name            string
		args            args
		expected        string
		expectedBoolean bool
	}{
		// First element in the group isn't relevant cause the matching loop starts from index 1.
		{"non regexp, empty group", args{[]string{}, "{1}-{2}-{3}", false}, "{1}-{2}-{3}", false},
		{"non regexp, empty group", args{[]string{""}, "{1}-{2}-{3}", false}, "{1}-{2}-{3}", false},
		{"regexp, empty group", args{[]string{}, "{1}-{2}-{3}", true}, "{1}-{2}-{3}", false},
		{"regexp, empty group", args{[]string{""}, "{1}-{2}-{3}", true}, "{1}-{2}-{3}", false},
		// Non regular expressions
		{"basic", args{[]string{"", "a", "b", "c"}, "{1}-{2}-{3}", false}, "a-b-c", true},
		{"opposite order", args{[]string{"", "a", "b", "c"}, "{3}-{2}-{1}-{4}", false}, "c-b-a-{4}", true},
		{"double", args{[]string{"", "a", "b"}, "{2}-{2}-{1}-{1}", false}, "b-b-a-a", true},
		{"skip placeholders indexes", args{[]string{"", "a", "b"}, "{4}-{1}", false}, "b-a", true},
		// Regular expressions
		{"basic", args{[]string{"", "a", "b", "c"}, "{1}-{2}-{3}", true}, "a-b-c", true},
		{"opposite order", args{[]string{"", "a", "b", "c"}, "{4}-{3}-{2}-{5}", true}, "{4}-c-b-{5}", true},
		{"double", args{[]string{"", "a", "b"}, "{2}-{2}-{1}-{1}", true}, "b-b-a-a", true},
		{"skip placeholders indexes", args{[]string{"", "a", "b"}, "{3}-{1}", true}, "{3}-a", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, replaceOccurred, err := ReplacePlaceHolders(tt.args.groups, tt.args.toReplace, tt.args.isRegexp)
			assert.NoError(t, err)
			assert.Equalf(t, tt.expected, result, "ReplacePlaceHolders(%v, %v, %v)", tt.args.groups, tt.args.toReplace, tt.args.isRegexp)
			assert.Equalf(t, tt.expectedBoolean, replaceOccurred, "ReplacePlaceHolders(%v, %v, %v)", tt.args.groups, tt.args.toReplace, tt.args.isRegexp)
		})
	}
}

func TestValidateMinimumVersion(t *testing.T) {
	minTestVersion := "6.9.0"
	tests := []struct {
		artifactoryVersion string
		expectedResult     bool
	}{
		{"6.5.0", false},
		{"6.2.0", false},
		{"5.9.0", false},
		{"6.0.0", false},
		{"6.6.0", false},
		{"6.9.0", true},
		{Development, true},
		{"6.10.2", true},
		{"6.15.2", true},
	}
	for _, test := range tests {
		t.Run(test.artifactoryVersion, func(t *testing.T) {
			err := ValidateMinimumVersion(Xray, test.artifactoryVersion, minTestVersion)
			if test.expectedResult {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, fmt.Sprintf(MinimumVersionMsg, Xray, test.artifactoryVersion, minTestVersion))
			}
		})
	}
}

func TestSetEnvWithResetCallback(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		args   args
		init   func()
		finish func()
	}{
		{
			name: "existing environment variable",
			args: args{key: "TEST_KEY", value: "test_value"},
			init: func() {
				assert.NoError(t, os.Setenv("TEST_KEY", "test-init-value"))
			},
			finish: func() {
				assert.Equal(t, os.Getenv("TEST_KEY"), "test-init-value")
			},
		},
		{
			name: "non-existing environment variable",
			args: args{key: "NEW_TEST_KEY", value: "test_value"},
			init: func() {

			},
			finish: func() {
				_, exist := os.LookupEnv("NEW_TEST_KEY")
				assert.False(t, exist)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.init()
			resetCallback, err := SetEnvWithResetCallback(tt.args.key, tt.args.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.value, os.Getenv(tt.args.key))
			assert.NoError(t, resetCallback())
			tt.finish()
		})
	}
}

func calculateExpectedBounds(attempt int, initialDelay, maxDelay time.Duration) (minExpected, maxExpected time.Duration) {
	expDelayFloat := float64(initialDelay) * math.Pow(2, float64(attempt))
	cappedDelayFloat := math.Min(expDelayFloat, float64(maxDelay))
	minJitterFactor := 0.8
	maxJitterFactor := 1.2
	minExpected = time.Duration(cappedDelayFloat * minJitterFactor)
	maxExpected = time.Duration(cappedDelayFloat * maxJitterFactor)
	if minExpected < 0 {
		minExpected = 0
	}
	return
}

func TestCalculateBackoff(t *testing.T) {
	testCases := []struct {
		name         string
		attempt      int
		initialDelay time.Duration
		maxDelay     time.Duration
	}{
		{
			name:         "Attempt 0 - No cap",
			attempt:      0,
			initialDelay: 10 * time.Millisecond,
			maxDelay:     1 * time.Second,
		},
		{
			name:         "Attempt 1 - No cap",
			attempt:      1,
			initialDelay: 10 * time.Millisecond,
			maxDelay:     1 * time.Second,
		},
		{
			name:         "Attempt 2 - No cap",
			attempt:      2,
			initialDelay: 10 * time.Millisecond,
			maxDelay:     1 * time.Second,
		},
		{
			name:         "Attempt 5 - No cap",
			attempt:      5,
			initialDelay: 10 * time.Millisecond,
			maxDelay:     1 * time.Second,
		},
		{
			name:         "Attempt 0 - With cap (initial delay is capped)",
			attempt:      0,
			initialDelay: 50 * time.Millisecond,
			maxDelay:     30 * time.Millisecond,
		},
		{
			name:         "Attempt 3 - With cap (exponential delay capped)",
			attempt:      3,
			initialDelay: 100 * time.Millisecond,
			maxDelay:     500 * time.Millisecond,
		},
		{
			name:         "Attempt 10 - Max delay reached",
			attempt:      10,
			initialDelay: 1 * time.Millisecond,
			maxDelay:     200 * time.Millisecond,
		},
		{
			name:         "Zero initial delay",
			attempt:      2,
			initialDelay: 0 * time.Millisecond,
			maxDelay:     1 * time.Second,
		},
		{
			name:         "Negative attempt (should still work due to float64 conversion)",
			attempt:      -1,
			initialDelay: 10 * time.Millisecond,
			maxDelay:     1 * time.Second,
		},
	}
	const numIterations = 1000
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			minExpected, maxExpected := calculateExpectedBounds(tc.attempt, tc.initialDelay, tc.maxDelay)

			for i := 0; i < numIterations; i++ {
				actualDelay := CalculateBackoff(tc.attempt, tc.initialDelay, tc.maxDelay)
				assert.Truef(t, actualDelay >= minExpected,
					"Iteration %d: Actual delay %v is less than minimum expected %v (Test: %s)",
					i, actualDelay, minExpected, tc.name)
				assert.Truef(t, actualDelay < maxExpected || (actualDelay == maxExpected && maxExpected == 0),
					"Iteration %d: Actual delay %v is greater than or equal to maximum expected %v (Test: %s)",
					i, actualDelay, maxExpected, tc.name)
				assert.Truef(t, actualDelay >= 0,
					"Iteration %d: Actual delay %v is negative (Test: %s)",
					i, actualDelay, tc.name)
			}
		})
	}
}
