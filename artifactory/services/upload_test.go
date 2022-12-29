package services

import (
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/stretchr/testify/assert"
)

func TestDebianProperties(t *testing.T) {
	var debianPaths = []struct {
		in       string
		expected string
	}{
		{"dist/comp/arch", ";deb.distribution=dist;deb.component=comp;deb.architecture=arch"},
		{"dist1,dist2/comp/arch", ";deb.distribution=dist1,dist2;deb.component=comp;deb.architecture=arch"},
		{"dist/comp1,comp2/arch", ";deb.distribution=dist;deb.component=comp1,comp2;deb.architecture=arch"},
		{"dist/comp/arch1,arch2", ";deb.distribution=dist;deb.component=comp;deb.architecture=arch1,arch2"},
		{"dist1,dist2/comp1,comp2/arch1,arch2", ";deb.distribution=dist1,dist2;deb.component=comp1,comp2;deb.architecture=arch1,arch2"},
	}

	for _, v := range debianPaths {
		result := getDebianProps(v.in)
		if result != v.expected {
			t.Errorf("getDebianProps(\"%s\") => '%s', want '%s'", v.in, result, v.expected)
		}
	}
}

func TestBuildUploadUrls(t *testing.T) {
	var testsParams = []struct {
		targetPath                  string
		targetProps                 string
		buildProps                  string
		expectedTargetPathWithProps string
	}{
		{"repo1/file1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file1;k1=v1;k2=v2"},
		{"repo1/file@1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file@1;k1=v1;k2=v2"},
		{"repo1/file;1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file%3B1;k1=v1;k2=v2"},
		{"repo1/file,1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file,1;k1=v1;k2=v2"},
		{"repo1/file^1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file%5E1;k1=v1;k2=v2"},
		{"repo1/file:1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file:1;k1=v1;k2=v2"},
		{"repo1/file1", "", "k2=v2", "http://localhost:8881/artifactory/repo1/file1;k2=v2"},
		{"repo1/file1", "k1=v1", "", "http://localhost:8881/artifactory/repo1/file1;k1=v1"},
		{"repo1/file1", "", "", "http://localhost:8881/artifactory/repo1/file1"},
	}

	for _, v := range testsParams {
		targetProps, e := utils.ParseProperties(v.targetProps)
		assert.NoError(t, e)
		actualTargetPathWithProps, e := buildUploadUrls("http://localhost:8881/artifactory/", v.targetPath, v.buildProps, "", targetProps)
		assert.NoError(t, e)
		assert.Equal(t, v.expectedTargetPathWithProps, actualTargetPathWithProps)
	}
}

func TestAddEscapingParenthesesWithTargetInArchive(t *testing.T) {
	type args struct {
		pattern         string
		target          string
		targetInArchive string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty parentheses", args{"()", "", "{2}"}, "\\(\\)"},
		{"empty parentheses", args{"()", "", "{}"}, "\\(\\)"},
		{"empty parentheses", args{"()", "", "{1}"}, "()"},
		{"empty parentheses", args{")(", "", "{1}"}, "\\)\\("},
		{"first parentheses", args{"(a)/(b)/(c)", "", "{2}/{3}"}, "\\(a\\)/(b)/(c)"},
		{"second parentheses", args{"(a)/(b)/(c)", "", "{1}/{3}"}, "(a)/\\(b\\)/(c)"},
		{"third parentheses", args{"(a)/(b)/(c)", "", "{1}/{2}"}, "(a)/(b)/\\(c\\)"},
		{"empty placeholders", args{"(a)/(b)/(c)", "", ""}, "\\(a\\)/\\(b\\)/\\(c\\)"},
		{"un-symmetric parentheses", args{")a)/(b)/(c(", "", ""}, "\\)a\\)/\\(b\\)/\\(c\\("},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, addEscapingParenthesesForUpload(tt.args.pattern, tt.args.target, tt.args.targetInArchive), "AddEscapingParentheses(%v, %v)", tt.args.pattern, tt.args.target)
		})
	}
}

func TestAddEscapingParenthesesWithTargetAndTargetInArchive(t *testing.T) {
	type args struct {
		pattern         string
		target          string
		targetInArchive string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty parentheses", args{"()", "{2}", "{3}"}, "\\(\\)"},
		{"empty parentheses", args{"()", "{}", "{}"}, "\\(\\)"},
		{"empty parentheses", args{"()()", "{2}", "{1}"}, "()()"},
		{"empty parentheses", args{"))(((", "{2}", "{1}"}, "\\)\\)\\(\\(\\("},
		{"first parentheses", args{"(a)/(b)/(c)/(d)", "{4}", "{2}/{3}"}, "\\(a\\)/(b)/(c)/(d)"},
		{"second parentheses", args{"(a)/(b)/(c)/(d)", "{1}/{4}", "{1}/{3}"}, "(a)/\\(b\\)/(c)/(d)"},
		{"last parentheses", args{"(a)/(b)/(c)/(d)", "{1}/{3}", "{2}/{3}"}, "(a)/(b)/(c)/\\(d\\)"},
		{"mixed parentheses", args{"(a)/(b)/(c)/(d)/(e)/(f)", "{5}", "{1}/{2}"}, "(a)/(b)/\\(c\\)/\\(d\\)/(e)/\\(f\\)"},
		{"out of range placeholders", args{"(a)/(b)/(c)", "{5}", "{4}"}, "\\(a\\)/\\(b\\)/\\(c\\)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, addEscapingParenthesesForUpload(tt.args.pattern, tt.args.target, tt.args.targetInArchive), "AddEscapingParentheses(%v, %v)", tt.args.pattern, tt.args.target)
		})
	}
}

func TestSkipDirUpload(t *testing.T) {
	data := []struct {
		targetFiles []string
		sourceDirs  []string
		targetDir   string
		sourceDir   string
		includeDirs bool
		result      bool
	}{
		{[]string{}, []string{}, "cli-rt1-1671381032/b", "testdata/a/b", true, false},
		{[]string{"dirdir/b/"}, []string{}, "dirdir", "testdata/a/b", true, true},
		{[]string{"cli-rt1-1671381032/b/"}, []string{}, "cli-rt1-1671381032/b", "testdata/a/b", true, true},
		{[]string{"cli-rt1-1671383851/c", "cli-rt1-1671383851/b3.in"}, []string{filepath.Join("testdata", "a", "b", "c")}, "cli-rt1-1671383851/b", filepath.Join("testdata", "a", "b"), true, true},
	}
	for _, d := range data {
		got := skipDirUpload(d.targetFiles, d.sourceDirs, d.targetDir, d.sourceDir, d.includeDirs)
		assert.Equal(t, d.result, got)
	}
}
