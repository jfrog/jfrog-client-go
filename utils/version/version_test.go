package version

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		ver1     string
		ver2     string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"5.10.0", "5.5.2", 1},
		{"5.5.2", "5.15.2", -1},
		{"5.6.2", "5.50.2", -1},
		{"5.5.2", "5.0.2", 1},
		{"15.5.2", "6.0.2", 1},
		{"51.5.2", "6.0.2", 1},
		{"5.0.3", "5.0.20", -1},
		{"5.0.20", "5.0.3", 1},
		{"1.0.0", "1.0.1", -1},
		{"1.0.x-SNAPSHOT", "1.0.x-SNAPSHOT", 0},
		{"1.1.x-SNAPSHOT", "1.0.x-SNAPSHOT", 1},
		{"2.0.x-SNAPSHOT", "1.0.x-SNAPSHOT", 1},
		{"1.0", "1.0.x-SNAPSHOT", -1},
		{"1.1", "1.0.x-SNAPSHOT", 1},
		{"1.0.x-SNAPSHOT", "1.0", 1},
		{"1.0.x-SNAPSHOT", "1.1", -1},
		{"1", "2", -1},
		{"1.0", "2.0", -1},
		{"2.1", "2.0", 1},
		{"2.a", "2.b", -1},
		{"b", "a", 1},
		{"1.0", "1", 0},
		{"1.1", "1", 1},
		{"1", "1.1", -1},
		{"", "1", -1},
		{"1", "", 1},
		{"6.x-SNAPSHOT", "5.5.2", 1},
		{"6.x-SNAPSHOT", "6.5.0", 1},
		{"6.5.x-SNAPSHOT", "6.5.2", 1},
		{"7.x-SNAPSHOT", "6.x-SNAPSHOT", 1},
		{"6.1.x-SNAPSHOT", "6.2.x-SNAPSHOT", -1},
	}
	for _, test := range tests {
		t.Run(test.ver1+":"+test.ver2, func(t *testing.T) {
			assert.Equal(t, Compare(test.ver1, test.ver2), test.expected)
		})
	}
}

func TestIsMatch(t *testing.T) {
	tests := []struct {
		version    string
		constraint string
		expected   bool
	}{
		{"1.0.0", "1.0.0", true},
		{"2.0.0", "1.0.0", false},
		{"1.0.0", "*.0.0", true},
		{"1.0.0", "1.*.0", true},
		{"1.0.0", "1.0.*", true},
		{"1.2.0", "*.0.0", false},
		{"2.0.0", "1.*.0", false},
		{"2.0.0", "1.0.*", false},
		{"1.0.0-SNAPSHOT", "1.0.0-SNAPSHOT", true},
		{"1.0.0-SNAPSHOT", "1.0.0-*", true},
		{"1.0.0-SNAPSHOT", "1.*.0-SNAPSHOT", true},
		{"1.0.0-SNAPSHOT", "1.0.*-SNAPSHOT", true},
		{"2.0.0-SNAPSHOT", "1.0.0-SNAPSHOT", false},
		{"1.2.0-SNAPSHOT", "1.0.0-*", false},
		{"1.0.2-SNAPSHOT", "1.*.0-SNAPSHOT", false},
		{"1.0.0", "1.0.*-SNAPSHOT", false},
		{"1", "2", false},
		{"1", "*", true},
		{"1", "1.*", true},
		{"1.0.1", "1.*", true},
		{"1.0.1", "*", true},
		{"1.0", "2.0", false},
		{"1.0", "*", true},
		{"1.0", "1.*", true},
		{"1.0", "*.0", true},
		{"2.a", "2.b", false},
		{"2.a", "2.*", true},
		{"b", "*", true},
		{"", "1", false},
		{"1", "", true},
		{"", "*", true},
	}
	for _, test := range tests {
		t.Run(test.version+":"+test.constraint, func(t *testing.T) {
			assert.Equal(t, IsMatch(test.version, test.constraint), test.expected)
		})
	}
}
