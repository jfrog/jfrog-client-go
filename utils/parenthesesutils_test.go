package utils

import (
	"reflect"
	"testing"
)

func TestSortNoDuplicates(t *testing.T) {
	got := []int{3, 2, 1, 4, 3, 2, 1, 4, 1}
	beforSsortNoDuplicates := got
	sortNoDuplicates(&got)
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("sortNoDuplicates(%v) == %v, want %v", beforSsortNoDuplicates, got, want)
	}
}

func TestFindParentheses(t *testing.T) {
	pattern := "(a/(b)"
	target := "{1}"
	got := findParentheses(pattern, target)
	want := []Parentheses{{3, 5}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("sortNoDuplicates(%s, %s) == %v, want %v", pattern, target, got, want)
	}

	pattern = "(a/(b"
	target = "{1}"
	got = findParentheses(pattern, target)
	if len(got) != 0 {
		t.Errorf("sortNoDuplicates(%s, %s) == %v, want []]", pattern, target, got)
	}

	pattern = "(a/(b)"
	target = "{1"
	got = findParentheses(pattern, target)
	if len(got) != 0 {
		t.Errorf("sortNoDuplicates(%s, %s) == %v, want []]", pattern, target, got)
	}

	pattern = "(a)/(b)"
	target = "{1}/{2}"
	got = findParentheses(pattern, target)
	want = []Parentheses{{0, 2}, {4, 6}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("sortNoDuplicates(%s, %s) == %v, want %v", pattern, target, got, want)
	}
}

func TestGetPlaceHoldersValues(t *testing.T) {
	target := "{1} {2}/{3}[4]{"
	got := getPlaceHoldersValues(target)
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("sortNoDuplicates(%v) == %v, want %v", target, got, want)
	}
}
