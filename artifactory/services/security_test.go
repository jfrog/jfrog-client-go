package services

import "testing"

// Test mapping ExpiresIn to expires_in request value and default handling
func TestBuildCreateTokenUrlValuesExpiresIn(t *testing.T) {
	tests := []struct {
		testName string
		input    int
		output   string
	}{
		{"never expires", 0, "0"},
		{"expires", 1800, "1800"},
		{"default", -1, ""},
	}
	for _, test := range tests {
		values := buildCreateTokenUrlValues(CreateTokenParams{
			ExpiresIn: test.input,
		})
		if values.Get("expires_in") != test.output {
			t.Errorf("Test name: %s: Expected: %s, Got: %s", test.testName, test.output, values.Get("expires_in"))
		}
	}
}

// Test default value -1 in NewCreateTokenParams
func TestNewCreateTokenParams(t *testing.T) {
	values := buildCreateTokenUrlValues(NewCreateTokenParams())
	if values.Get("expires_in") != "" {
		t.Errorf("default expires_in should be empty")
	}
}
