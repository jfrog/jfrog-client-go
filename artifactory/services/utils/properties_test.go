package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestParseProperties(t *testing.T) {
	tests := []struct {
		propsString string
		option      PropertyParseOptions
		expected    Properties
	}{
		{"y=a,b", SplitCommas, Properties{[]Property{{Key: "y", Value: "a"}, {Key: "y", Value: "b"}}}},
		{"y=a\\,b", SplitCommas, Properties{[]Property{{Key: "y", Value: "a,b"}}}},
		{"y=a,b\\", SplitCommas, Properties{[]Property{{Key: "y", Value: "a"}, {Key: "y", Value: "b\\"}}}},
		{"y=a,b\\,", SplitCommas, Properties{[]Property{{Key: "y", Value: "a"}, {Key: "y", Value: "b,"}}}},
		{"y=a,b\\,c,d", SplitCommas, Properties{[]Property{{Key: "y", Value: "a"}, {Key: "y", Value: "b,c"}, {Key: "y", Value: "d"}}}},
		{"y=a,b\\,c\\,d", SplitCommas, Properties{[]Property{{Key: "y", Value: "a"}, {Key: "y", Value: "b,c,d"}}}},
		{"y=a,b\\,c\\,d\\,e", SplitCommas, Properties{[]Property{{Key: "y", Value: "a"}, {Key: "y", Value: "b,c,d,e"}}}},
		{"y=\\,a b", SplitCommas, Properties{[]Property{{Key: "y", Value: ",a b"}}}},
	}
	for _, test := range tests {
		t.Run(test.propsString, func(t *testing.T) {
			props, err := ParseProperties(test.propsString, test.option)
			if err != nil {
				t.Error("Failed to parse property string.", err)
			}
			if !reflect.DeepEqual(test.expected.Properties, props.Properties) {
				t.Error("Failed to parse property string.", props, "expected to be parsed to", test.expected)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	tests := []struct {
		props Properties
	}{
		{Properties{[]Property{}}},
		{Properties{[]Property{{Key: "a", Value: "b"}}}},
		{Properties{[]Property{{Key: "a", Value: "b"}, {Key: "c", Value: "d"}}}},
		{Properties{[]Property{{Key: "a", Value: "b"}, {Key: "c", Value: "d"}, {Key: "c", Value: "e"}}}},
	}
	for _, test := range tests {
		t.Run(test.props.ToEncodedString(), func(t *testing.T) {
			expectedProps := test.props.Properties
			propsMap := test.props.ToMap()
			assert.LessOrEqual(t, len(propsMap), len(expectedProps))
			for _, expectedProp := range expectedProps {
				values := propsMap[expectedProp.Key]
				assert.NotNil(t, values, "Key "+expectedProp.Key+" not found in map")
				assert.Contains(t, values, expectedProp.Value)
			}
		})
	}
}
