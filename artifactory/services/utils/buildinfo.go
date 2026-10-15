package utils

import (
	"strings"
	"unicode"
)

const spaceEncoding = "%20"

var buildInfoRepositoryCharacters = map[rune]string{
	'/':  "%2F",
	'\\': "%5C",
	'|':  "%7C",
	'*':  "%2A",
	'?':  "%3F",
	'\'': "%22",
	':':  "%3A",
	';':  "%3B",
	'%':  "%25",
}

func EncodeForBuildInfoRepository(value string) string {
	results := ""
	for _, char := range value {
		if unicode.IsSpace(char) {
			char = ' '
		}
		if encoding, exist := buildInfoRepositoryCharacters[char]; exist {
			results += encoding
		} else {
			results += string(char)
		}
	}
	slashEncoding := buildInfoRepositoryCharacters['/']
	results = strings.ReplaceAll(results, slashEncoding+" ", slashEncoding+spaceEncoding)
	results = strings.ReplaceAll(results, " "+slashEncoding, spaceEncoding+slashEncoding)
	return results
}
