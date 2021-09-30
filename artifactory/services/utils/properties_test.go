package utils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func TestToEncodedString(t *testing.T) {
	tests := []struct {
		props    Properties
		expected string
	}{
		{Properties{properties: map[string][]string{"a": {"b"}}}, "a=b"},
		{Properties{properties: map[string][]string{"a;a": {"b;a"}}}, "a%3Ba=b%3Ba"},
		{Properties{properties: map[string][]string{"a": {"b"}}}, "a=b"},
		{Properties{properties: map[string][]string{";a": {";b"}}}, "%3Ba=%3Bb"},
		{Properties{properties: map[string][]string{";a": {";b"}, "aaa": {"bbb"}}}, "%3Ba=%3Bb;aaa=bbb"},
		{Properties{properties: map[string][]string{"a;": {"b;"}}}, "a%3B=b%3B"},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			actual := test.props.ToEncodedString(true)
			assert.ElementsMatch(t, strings.Split(test.expected, ";"), strings.Split(actual, ";"), "failed to encode properties")
		})
	}
}

func TestParseProperties(t *testing.T) {
	tests := []struct {
		propsString string
		expected    Properties
	}{
		{"y=a,b", Properties{properties: map[string][]string{"y": {"a", "b"}}}},
		{"y=a\\,b", Properties{properties: map[string][]string{"y": {"a,b"}}}},
		{"y=a,b\\", Properties{properties: map[string][]string{"y": {"a", "b\\"}}}},
		{"y=a,b\\,", Properties{properties: map[string][]string{"y": {"a", "b,"}}}},
		{"y=a,b\\,c,d", Properties{properties: map[string][]string{"y": {"a", "b,c", "d"}}}},
		{"y=a,b\\,c\\,d", Properties{properties: map[string][]string{"y": {"a", "b,c,d"}}}},
		{"y=a,b\\,c\\,d\\,e", Properties{properties: map[string][]string{"y": {"a", "b,c,d,e"}}}},
		{"y=\\,a b", Properties{properties: map[string][]string{"y": {",a b"}}}},
		{"y=a,b;x=i;y=a;x=j", Properties{properties: map[string][]string{"y": {"a", "b"}, "x": {"i", "j"}}}},
		{"y=a=a,b;x=i;y=a=a;x=j=", Properties{properties: map[string][]string{"y": {"a=a", "b"}, "x": {"i", "j="}}}},
		{"y=a,b;x=i\\;;y=a;x=j\\;\\;", Properties{properties: map[string][]string{"y": {"a", "b"}, "x": {"i;", "j;;"}}}},
		{"y=a,b;x=\\i;y=a\\;;x=j\\;=\\,", Properties{properties: map[string][]string{"y": {"a", "b", "a;"}, "x": {"\\i", "j;=,"}}}},
	}
	for _, test := range tests {
		t.Run(test.propsString, func(t *testing.T) {
			props, err := ParseProperties(test.propsString)
			if err != nil {
				t.Error("Failed to parse property string.", err)
			}
			if !reflect.DeepEqual(test.expected.ToMap(), props.ToMap()) {
				t.Error("Failed to parse property string.", props, "expected to be parsed to", test.expected)
			}
		})
	}
}

func TestMergeProperties(t *testing.T) {
	propsA := &Properties{properties: map[string][]string{"x": {"a", "b"}, "y": {"c"}}}
	propsB := &Properties{properties: map[string][]string{"x": {"a", "d"}}}
	mergedProps := MergeProperties([]*Properties{propsA, propsB})
	expected := &Properties{properties: map[string][]string{"x": {"a", "b", "d"}, "y": {"c"}}}
	if !reflect.DeepEqual(expected, mergedProps) {
		t.Error("Failed to marge Properties. expected:", expected, "actual:", mergedProps)
	}
}
