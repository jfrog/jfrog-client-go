package distribution

import (
	"github.com/jfrog/gofrog/stringutils"
	"regexp"
)

var fileSpecCaptureGroup = regexp.MustCompile(`({\d})`)

// Create the path mapping from the input spec
func CreatePathMappings(pattern, target string) []PathMapping {
	if len(target) == 0 {
		return []PathMapping{}
	}

	// Convert the file spec pattern and target to match the path mapping input and output specifications, respectfully.
	return []PathMapping{{
		// The file spec pattern is wildcard based. Convert it to Regex:
		Input: stringutils.WildcardPatternToRegExp(pattern),
		// The file spec target contain placeholders-style matching groups, like {1}.
		// Convert it to REST-APIs matching groups style, like $1.
		Output: fileSpecCaptureGroup.ReplaceAllStringFunc(target, func(s string) string {
			// Remove curly parenthesis and prepend $
			return "$" + s[1:2]
		}),
	}}
}

type PathMapping struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}
