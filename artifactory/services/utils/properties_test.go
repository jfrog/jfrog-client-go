package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToMap(t *testing.T) {
	properties := Properties{[]Property{
		{Key: "a", Value: "b"},
		{Key: "c", Value: "d"},
		{Key: "c", Value: "e"}},
	}
	propertiesMap := properties.ToMap()
	assert.Len(t, propertiesMap, 2)
	for key, values := range propertiesMap {
		switch key {
		case "a":
			assert.Equal(t, []string{"b"}, values)
		case "c":
			assert.ElementsMatch(t, []string{"d", "e"}, values)
		default:
			assert.Fail(t, "Unexpected key "+key)
		}
	}
}

func TestToEncodedString(t *testing.T) {
	tests := []struct {
		props    Properties
		expected string
	}{
		{Properties{[]Property{{Key: "a", Value: "b"}}}, "a=b"},
		{Properties{[]Property{{Key: "a;a", Value: "b;a"}}}, "a%3Ba=b%3Ba"},
		{Properties{[]Property{{Key: "a", Value: "b"}}}, "a=b"},
		{Properties{[]Property{{Key: ";a", Value: ";b"}}}, "%3Ba=%3Bb"},
		{Properties{[]Property{{Key: ";a", Value: ";b"}, {Key: ";a", Value: ";b"}, {Key: "aaa", Value: "bbb"}}}, "%3Ba=%3Bb;%3Ba=%3Bb;aaa=bbb"},
		{Properties{[]Property{{Key: "a;", Value: "b;"}}}, "a%3B=b%3B"},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			test.props.ToEncodedString()
			if test.expected != test.props.ToEncodedString() {
				t.Error("Failed to encode properties. The propertyes", test.props.ToEncodedString(), "expected to be encoded to", test.expected)
			}
		})
	}
}
