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
		{"1.0.0-SNAPSHOT", "1.0.0-SNAPSHOT", 0},
		{"1.1.0-SNAPSHOT", "1.0.0-SNAPSHOT", 1},
		{"2.0.0-SNAPSHOT", "1.0.0-SNAPSHOT", 1},
		{"1.0", "1.0.0-SNAPSHOT", 1},
		{"1.1", "1.0.0-SNAPSHOT", 1},
		{"1.0.0-SNAPSHOT", "1.0", -1},
		{"1.0.0-SNAPSHOT", "1.1", -1},
		{"1", "2", -1},
		{"1.0", "2.0", -1},
		{"2.1", "2.0", 1},
		{"1.0", "1", 0},
		{"1.1", "1", 1},
		{"1", "1.1", -1},
		{"6.0-SNAPSHOT", "5.5.2", 1},
		{"6.0-SNAPSHOT", "6.5.0", -1},
		{"6.5.0-SNAPSHOT", "6.5.2", -1},
		{"7.0-SNAPSHOT", "6.0-SNAPSHOT", 1},
		{"6.1.0-SNAPSHOT", "6.2.0-SNAPSHOT", -1},
		{"v1.0.0", "1.0.0", 0},
		{"v1.0.1", "v1.0.0", 1},
		{"v5.10.0", "5.5.2", 1},
		{"5.5.2", "v5.15.2", -1},
		{"v5.6.2", "5.50.2", -1},
		{"5.5.2", "v5.0.2", 1},
		{"v15.5.2", "6.0.2", 1},
		{"51.5.2", "v6.0.2", 1},
		{"v5.0.3", "v5.0.20", -1},
		{"v5.0.20", "5.0.3", 1},
		{"v1.0.0", "v1.0.1", -1},
		{"v1.0.0-SNAPSHOT", "1.0.0-SNAPSHOT", 0},
		{"1.1.0-SNAPSHOT", "v1.0.0-SNAPSHOT", 1},
		{"v2.0.0-SNAPSHOT", "v1.0.0-SNAPSHOT", 1},
	}
	for _, test := range tests {
		t.Run(test.ver1+":"+test.ver2, func(t *testing.T) {
			result, err := compare(test.ver1, test.ver2)
			if err != nil {
				t.Error(err)
			}
			if result != test.expected {
				t.Error("ver1:", test.ver1, "ver2:", test.ver2, "Expecting:", test.expected, "got:", result)
			}
		})
	}
}
