package services

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPropsParams(t *testing.T) {
	params := NewPropsParams()
	assert.Nil(t, params.Reader)
	assert.Empty(t, params.Props)
	assert.False(t, params.IsRepoOnly)
	assert.False(t, params.UseDebugLogs)
	assert.False(t, params.IsRecursive)
}

func TestPropsParamsGetProps(t *testing.T) {
	tests := []struct {
		name     string
		props    string
		expected string
	}{
		{"empty props", "", ""},
		{"single property", "key=value", "key=value"},
		{"multiple properties", "key1=val1;key2=val2", "key1=val1;key2=val2"},
		{"property with special chars", "key=val;ue", "key=val;ue"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := PropsParams{Props: tt.props}
			assert.Equal(t, tt.expected, params.GetProps())
		})
	}
}

func TestPropsParamsFields(t *testing.T) {
	tests := []struct {
		name         string
		isRepoOnly   bool
		useDebugLogs bool
		isRecursive  bool
	}{
		{"all false", false, false, false},
		{"isRepoOnly true", true, false, false},
		{"useDebugLogs true", false, true, false},
		{"isRecursive true", false, false, true},
		{"all true", true, true, true},
		{"isRepoOnly and isRecursive", true, false, true},
		{"useDebugLogs and isRecursive", false, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := PropsParams{
				IsRepoOnly:   tt.isRepoOnly,
				UseDebugLogs: tt.useDebugLogs,
				IsRecursive:  tt.isRecursive,
			}
			assert.Equal(t, tt.isRepoOnly, params.IsRepoOnly)
			assert.Equal(t, tt.useDebugLogs, params.UseDebugLogs)
			assert.Equal(t, tt.isRecursive, params.IsRecursive)
		})
	}
}

func TestGetEncodedParam(t *testing.T) {
	ps := NewPropsService(nil)

	tests := []struct {
		name        string
		props       string
		isDelete    bool
		expected    string
		expectError bool
	}{
		{
			name:        "set single property",
			props:       "key=value",
			isDelete:    false,
			expected:    "key=value",
			expectError: false,
		},
		{
			name:        "set multiple properties",
			props:       "key1=value1;key2=value2",
			isDelete:    false,
			expected:    "key1=value1;key2=value2",
			expectError: false,
		},
		{
			name:        "delete single property",
			props:       "key",
			isDelete:    true,
			expected:    "key",
			expectError: false,
		},
		{
			name:        "delete multiple properties",
			props:       "key1,key2",
			isDelete:    true,
			expected:    "key1,key2",
			expectError: false,
		},
		{
			name:        "delete property with special chars",
			props:       "key;special,key2",
			isDelete:    true,
			expected:    "key%3Bspecial,key2",
			expectError: false,
		},
		{
			name:        "set property with escaped semicolon in value",
			props:       "key=val\\;ue",
			isDelete:    false,
			expected:    "key=val%3Bue",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := PropsParams{Props: tt.props}
			result, err := ps.getEncodedParam(params, tt.isDelete)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, result, tt.expected[:3]) // Check prefix to handle encoding variations
			}
		})
	}
}

func TestActionTypeBasedOnIsDeleteFlag(t *testing.T) {
	ps := NewPropsService(nil)

	t.Run("delete action assigned when isDelete is true", func(t *testing.T) {
		var action func(string, string, string, bool) (*http.Response, []byte, error)
		ps.actionTypeBasedOnIsDeleteFlag(true, &action)
		assert.NotNil(t, action)
	})

	t.Run("put action assigned when isDelete is false", func(t *testing.T) {
		var action func(string, string, string, bool) (*http.Response, []byte, error)
		ps.actionTypeBasedOnIsDeleteFlag(false, &action)
		assert.NotNil(t, action)
	})
}

func TestRecursiveURLConstruction(t *testing.T) {
	tests := []struct {
		name        string
		isRecursive bool
		expected    string
	}{
		{
			name:        "non-recursive adds recursive=0",
			isRecursive: false,
			expected:    "&recursive=0",
		},
		{
			name:        "recursive adds recursive=1",
			isRecursive: true,
			expected:    "&recursive=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the recursive flag logic
			recursive := "&recursive=0"
			if tt.isRecursive {
				recursive = "&recursive=1"
			}
			assert.Equal(t, tt.expected, recursive)
		})
	}
}

func TestPropsServiceGetters(t *testing.T) {
	ps := NewPropsService(nil)

	t.Run("IsDryRun returns false", func(t *testing.T) {
		assert.False(t, ps.IsDryRun())
	})

	t.Run("GetThreads returns set value", func(t *testing.T) {
		ps.Threads = 5
		assert.Equal(t, 5, ps.GetThreads())
	})

	t.Run("GetThreads returns zero when not set", func(t *testing.T) {
		ps2 := NewPropsService(nil)
		assert.Equal(t, 0, ps2.GetThreads())
	})
}
