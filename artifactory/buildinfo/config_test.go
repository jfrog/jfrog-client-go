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
		prefix      string
		expected    map[string]string
		expectError bool
	}{
		{
			description: "empty input",
			config:      Configuration{},
			input:       map[string]string{},
			prefix:      "",
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
			prefix:      "",
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
			prefix: "",
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
			prefix:      "",
			expected:    nil,
			expectError: true,
		},
		{
			description: "input with prefix",
			config:      Configuration{EnvInclude: "*user*"},
			input: map[string]string{
				"buildInfo.env.USER":     "jfrog",
				"buildInfo.env.PASSWORD": "password",
			},
			prefix: "buildInfo.env.",
			expected: map[string]string{
				"buildInfo.env.USER": "jfrog",
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out, err := tc.config.IncludeFilter(tc.prefix)(tc.input)
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
		prefix      string
		expected    map[string]string
		expectError bool
	}{
		{
			description: "empty input",
			config:      Configuration{},
			input:       map[string]string{},
			prefix:      "",
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
			prefix: "",
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
			prefix: "",
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
			prefix: "",
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
			prefix:      "",
			expected:    nil,
			expectError: true,
		},
		{
			description: "input with prefix",
			config:      Configuration{EnvExclude: "*pass*"},
			input: map[string]string{
				"buildInfo.env.USER":     "jfrog",
				"buildInfo.env.PASSWORD": "password",
			},
			prefix: "buildInfo.env.",
			expected: map[string]string{
				"buildInfo.env.USER": "jfrog",
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out, err := tc.config.ExcludeFilter(tc.prefix)(tc.input)
			if tc.expectError {
				assert.NotNil(t, err)
			}

			assert.Equal(t, tc.expected, out)
		})
	}
}
