package version

import (
	"github.com/jfrog/jfrog-client-go/utils"
	"regexp"
	"strconv"
	"strings"
)

// If ver1 == ver2 returns 0
// If ver1 > ver2 returns 1
// If ver1 < ver2 returns -1
func Compare(ver1, ver2 string) int {
	if ver1 == ver2 {
		return 0
	}

	ver1Tokens := strings.Split(ver1, ".")
	ver2Tokens := strings.Split(ver2, ".")

	maxIndex := len(ver1Tokens)
	if len(ver2Tokens) > maxIndex {
		maxIndex = len(ver2Tokens)
	}

	for tokenIndex := 0; tokenIndex < maxIndex; tokenIndex++ {
		ver1Token := "0"
		if len(ver1Tokens) >= tokenIndex+1 {
			ver1Token = strings.TrimSpace(ver1Tokens[tokenIndex])
		}
		ver2Token := "0"
		if len(ver2Tokens) >= tokenIndex+1 {
			ver2Token = strings.TrimSpace(ver2Tokens[tokenIndex])
		}
		compare := compareTokens(ver1Token, ver2Token)
		if compare != 0 {
			return compare
		}
	}

	return 0
}

func compareTokens(ver1Token, ver2Token string) int {
	// Ignoring error because we strip all the non numeric values in advance.
	ver1TokenInt, _ := strconv.Atoi(getFirstNumeral(ver1Token))
	ver2TokenInt, _ := strconv.Atoi(getFirstNumeral(ver2Token))

	switch {
	case ver1TokenInt > ver2TokenInt:
		return 1
	case ver1TokenInt < ver2TokenInt:
		return -1
	default:
		return strings.Compare(ver1Token, ver2Token)
	}
}

// Return true if version matches the constraint
func IsMatch(version, constraint string) bool {
	versionTokens := strings.Split(version, ".")
	constraintTokens := strings.Split(constraint, ".")

	maxIndex := len(versionTokens)
	if len(constraintTokens) > maxIndex {
		maxIndex = len(constraintTokens)
	}

	// If the last token in constraint was "*", we want that the following tokens be also "*"
	isAsterisk := false
	for tokenIndex := 0; tokenIndex < maxIndex; tokenIndex++ {
		versionToken := "0"
		if len(versionTokens) >= tokenIndex+1 {
			versionToken = strings.TrimSpace(versionTokens[tokenIndex])
		}
		constraintToken := "0"
		if isAsterisk {
			constraintToken = "*"
		}
		if len(constraintTokens) >= tokenIndex+1 {
			constraintToken = strings.TrimSpace(constraintTokens[tokenIndex])
		}
		isAsterisk = constraintToken == "*"
		if !isVersionMatch(versionToken, constraintToken) {
			return false
		}
	}
	return true
}

// Return true if version token matches constraint token
func isVersionMatch(version string, constraint string) bool {
	regexConstraint := utils.WildcardToRegex(constraint)
	isMatched, _ := regexp.MatchString(regexConstraint, version)
	return isMatched
}

func getFirstNumeral(token string) string {
	numeric := ""
	for i := 0; i < len(token); i++ {
		n := token[i : i+1]
		if _, err := strconv.Atoi(n); err != nil {
			return "999999"
		}
		numeric += n
	}
	return numeric
}
