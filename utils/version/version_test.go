package version

import "testing"

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
		{"go1.13", "go1.13.0", 0},
		{"go1.13", "go1.14.1", -1},
		{"go1.13", "go1.12.100", 1},
	}
	for _, test := range tests {
		t.Run(test.ver1+":"+test.ver2, func(t *testing.T) {
			version := Version{version: test.ver2}
			result := version.Compare(test.ver1)
			if result != test.expected {
				t.Error("ver1:", test.ver1, "ver2:", test.ver2, "Expecting:", test.expected, "got:", result)
			}
		})
	}
}

func TestAtLeast(t *testing.T) {
	tests := []struct {
		ver1     string
		ver2     string
		expected bool
	}{
		{"1.0.0", "1.0.0", true},
		{"1.0.1", "1.0.0", true},
		{"5.10.0", "5.5.2", true},
		{"1.0.x-SNAPSHOT", "1.0.x-SNAPSHOT", true},
		{"1.1.x-SNAPSHOT", "1.0.x-SNAPSHOT", true},
		{"2.0.x-SNAPSHOT", "1.0.x-SNAPSHOT", true},
		{"development", "5.5", true},
		{"6.2.0", "6.5.0", false},
		{"6.6.0", "6.8.0", false},
	}
	for _, test := range tests {
		t.Run(test.ver1+":"+test.ver2, func(t *testing.T) {
			version := Version{version: test.ver1}
			result := version.AtLeast(test.ver2)
			if result != test.expected {
				t.Error("ver1:", test.ver1, "ver2:", test.ver2, "Expecting:", test.expected, "got:", result)
			}
		})
	}
}
