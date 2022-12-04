package utils

import (
	"github.com/stretchr/testify/assert"
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

	pattern = "(a)养只/(b)"
	target = "{1}/{2}"
	got = findParentheses(pattern, target)
	want = []Parentheses{{0, 2}, {10, 12}}
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

func TestAddEscapingParentheses(t *testing.T) {
	type args struct {
		pattern string
		target  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty parentheses", args{"()", "{2}"}, "\\(\\)"},
		{"empty parentheses", args{"()", "{}"}, "\\(\\)"},
		{"empty parentheses", args{"()", "{1}"}, "()"},
		{"empty parentheses", args{")(", "{1}"}, "\\)\\("},
		{"first parentheses", args{"(a)/(b)/(c)", "{2}/{3}"}, "\\(a\\)/(b)/(c)"},
		{"second parentheses", args{"(a)/(b)/(c)", "{1}/{3}"}, "(a)/\\(b\\)/(c)"},
		{"third parentheses", args{"(a)/(b)/(c)", "{1}/{2}"}, "(a)/(b)/\\(c\\)"},
		{"empty placeholders", args{"(a)/(b)/(c)", ""}, "\\(a\\)/\\(b\\)/\\(c\\)"},
		{"un-symmetric parentheses", args{")a)/(b)/(c(", ""}, "\\)a\\)/\\(b\\)/\\(c\\("},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, AddEscapingParentheses(tt.args.pattern, tt.args.target), "AddEscapingParentheses(%v, %v)", tt.args.pattern, tt.args.target)
		})
	}
}
