package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jfrog/gofrog/stringutils"
	"github.com/jfrog/jfrog-client-go/utils/io"
)

var (
	// Replace ** with a special string.
	doubleStartSpecialString = "__JFROG_DOUBLE_STAR__"

	// Match **/ ('__JFROG_DOUBLE_STAR__\/...')
	prefixDoubleStarRegex = regexp.MustCompile(fmt.Sprintf("%s%s", doubleStartSpecialString, getFileSeparatorForDoubleStart()))

	// match /** ('\/__JFROG_DOUBLE_STAR__...')
	postfixDoubleStarRegex = regexp.MustCompile(fmt.Sprintf("%s%s", getFileSeparatorForDoubleStart(), doubleStartSpecialString))

	// match **  ('...__JFROG_DOUBLE_STAR__...')
	middleDoubleStarNoSeparateRegex = regexp.MustCompile(doubleStartSpecialString)
)

func getFileSeparatorForDoubleStart() string {
	if io.IsWindows() {
		return `\\\\`
	}
	return `\/`
}

func AntToRegex(antPattern string) string {
	antPattern = stringutils.EscapeSpecialChars(antPattern)
	antPattern = antQuestionMarkToRegex(antPattern)
	return "^" + antStarsToRegex(antPattern) + "$"
}

func getFileSeparatorForAntToRegex() string {
	if io.IsWindows() {
		return `\\`
	}
	return `/`
}

func antStarsToRegex(antPattern string) string {
	separator := getFileSeparatorForAntToRegex()
	antPattern = addMissingShorthand(antPattern)

	// Replace ** with a special string, so it doesn't get mixed up with single *
	antPattern = strings.ReplaceAll(antPattern, "**", doubleStartSpecialString)

	// ant `*` => regexp `([^/]*)` : `*` matches zero or more characters except from `/`.
	antPattern = strings.ReplaceAll(antPattern, `*`, "([^"+separator+"]*)")

	// ant `**/` => regexp `(.*/)*` : Matches zero or more 'directories' at the beginning of the path.
	antPattern = prefixDoubleStarRegex.ReplaceAllString(antPattern, "(.*"+separator+")*")

	// ant `/**` => regexp `(/.*)*` : Matches zero or more 'directories' at the end of the path.
	antPattern = postfixDoubleStarRegex.ReplaceAllString(antPattern, "("+separator+".*)*")

	// ant `**` => regexp `(.*)*` : Matches zero or more 'directories'.
	return  middleDoubleStarNoSeparateRegex.ReplaceAllString(antPattern, "(.*)")
}

func antQuestionMarkToRegex(antPattern string) string {
	return strings.ReplaceAll(antPattern, "?", ".")
}

func addMissingShorthand(antPattern string) string {
	// There is one "shorthand": if a pattern ends with / or \, then ** is appended. For example, mypackage/test/ is interpreted as if it were mypackage/test/**.
	if strings.HasSuffix(antPattern, string(os.PathSeparator)) {
		return antPattern + "**"
	}
	return antPattern
}
