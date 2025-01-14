package distribution

import (
	"github.com/jfrog/gofrog/stringutils"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"regexp"
)

var fileSpecCaptureGroup = regexp.MustCompile(`({\d})`)

// Create the path mapping from the input spec which includes pattern and target
func CreatePathMappingsFromPatternAndTarget(pattern, target string) []utils.PathMapping {
	if len(target) == 0 {
		return []utils.PathMapping{}
	}

	// Convert the file spec pattern and target to match the path mapping input and output specifications, respectfully.
	return []utils.PathMapping{{
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

// Create the path mapping from the input spec
func CreatePathMappings(input, output string) []utils.PathMapping {
	if len(input) == 0 || len(output) == 0 {
		return []utils.PathMapping{}
	}

	return []utils.PathMapping{{
		Input:  input,
		Output: output,
	}}
}

func GetProjectQueryParam(projectKey string) map[string]string {
	queryParams := make(map[string]string)
	if projectKey != "" {
		queryParams["project"] = projectKey
	}
	return queryParams
}
