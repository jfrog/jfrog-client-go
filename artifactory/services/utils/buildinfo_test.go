package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var encodeForBuildInfoRepositoryProvider = []struct {
	value            string
	expectedEncoding string
}{
	// Shouldn't encode
	{"", ""},
	{"a", "a"},
	{"a b", "a b"},
	{"a.b", "a.b"},
	{"a&b", "a&b"},

	// Should encode
	{"a/b", "a%2Fb"},
	{"a\\b", "a%5Cb"},
	{"a:b", "a%3Ab"},
	{"a|b", "a%7Cb"},
	{"a*b", "a%2Ab"},
	{"a?b", "a%3Fb"},
	{"a  /  b", "a %20%2F%20 b"},

	// Should convert whitespace to space
	{"a\tb", "a b"},
	{"a\nb", "a b"},
}

func TestEncodeForBuildInfoRepository(t *testing.T) {
	for _, testCase := range encodeForBuildInfoRepositoryProvider {
		t.Run(testCase.value, func(t *testing.T) {
			assert.Equal(t, testCase.expectedEncoding, EncodeForBuildInfoRepository(testCase.value))
		})
	}
}
