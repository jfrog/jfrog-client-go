package services

import "testing"

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
