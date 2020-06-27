package buildinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInclude(t *testing.T) {
	tests := []struct {
		description string
		config      Configuration
		input       map[string]string
		expected    map[string]string
		expectError bool
	}{
		{
			description: "empty input",
			config:      Configuration{},
			input:       map[string]string{},
			expected:    map[string]string{},
			expectError: false,
		},
		{
			description: "input with no pattern",
			config:      Configuration{},
			input: map[string]string{
				"USER":     "jfrog",
				"PASSWORD": "password",
			},
			expected:    map[string]string{},
			expectError: false,
		},
		{
			description: "input with pattern",
			config:      Configuration{EnvInclude: "*user*"},
			input: map[string]string{
				"USER":     "jfrog",
				"PASSWORD": "password",
			},
			expected: map[string]string{
				"USER": "jfrog",
			},
			expectError: false,
		},
		{
			description: "input with bad pattern",
			config:      Configuration{EnvInclude: "use[*"},
			input: map[string]string{
				"USER": "jfrog",
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out, err := tc.config.IncludeFilter()(tc.input)
			if tc.expectError {
				assert.NotNil(t, err)
			}

			assert.Equal(t, tc.expected, out)
		})
	}
}

func TestExclude(t *testing.T) {
	tests := []struct {
		description string
		config      Configuration
		input       map[string]string
		expected    map[string]string
		expectError bool
	}{
		{
			description: "empty input",
			config:      Configuration{},
			input:       map[string]string{},
			expected:    map[string]string{},
			expectError: false,
		},
		{
			description: "input with no pattern",
			config:      Configuration{},
			input: map[string]string{
				"USER":     "jfrog",
				"PASSWORD": "password",
			},
			expected: map[string]string{
				"USER":     "jfrog",
				"PASSWORD": "password",
			},
			expectError: false,
		},
		{
			description: "input with pattern",
			config:      Configuration{EnvExclude: "*pass*"},
			input: map[string]string{
				"USER":     "jfrog",
				"PASSWORD": "password",
			},
			expected: map[string]string{
				"USER": "jfrog",
			},
			expectError: false,
		},
		{
			description: "input with bad non-matching pattern",
			config:      Configuration{EnvExclude: "pas[*"},
			input: map[string]string{
				"USER": "jfrog",
			},
			expected: map[string]string{
				"USER": "jfrog",
			},
			expectError: false,
		},
		{
			description: "input with bad matching pattern",
			config:      Configuration{EnvExclude: "pas[*"},
			input: map[string]string{
				"USER":     "jfrog",
				"PASSWORD": "password",
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out, err := tc.config.ExcludeFilter()(tc.input)
			if tc.expectError {
				assert.NotNil(t, err)
			}

			assert.Equal(t, tc.expected, out)
		})
	}
}
