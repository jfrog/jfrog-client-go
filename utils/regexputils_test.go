package utils

import (
	"fmt"
	gofrogio "github.com/jfrog/gofrog/io"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"strings"
	"testing"
)

func TestRemoveCredentialsFromLine(t *testing.T) {
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
	regExpProtocol, err := GetRegExp(CredentialsInUrlRegexp)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name         string
		regex        gofrogio.CmdOutputPattern
		expectedLine string
		matched      bool
	}{
		{"http", gofrogio.CmdOutputPattern{RegExp: regExpProtocol, Line: "This is an example line http://user:password@127.0.0.1:8081/artifactory/path/to/repo"}, "This is an example line http://***.***@127.0.0.1:8081/artifactory/path/to/repo", true},
		{"https", gofrogio.CmdOutputPattern{RegExp: regExpProtocol, Line: "This is an example line https://user:password@127.0.0.1:8081/artifactory/path/to/repo"}, "This is an example line https://***.***@127.0.0.1:8081/artifactory/path/to/repo", true},
		{"Special characters 1", gofrogio.CmdOutputPattern{RegExp: regExpProtocol, Line: "This is an example line https://u-s!<e>_r:!p-a&%%s%sword@127.0.0.1:8081/artifactory/path/to/repo"}, "This is an example line https://***.***@127.0.0.1:8081/artifactory/path/to/repo", true},
		{"Special characters 2", gofrogio.CmdOutputPattern{RegExp: regExpProtocol, Line: "This is an example line https://!user:[p]a(s)sword@127.0.0.1:8081/artifactory/path/to/repo"}, "This is an example line https://***.***@127.0.0.1:8081/artifactory/path/to/repo", true},
		{"No credentials", gofrogio.CmdOutputPattern{RegExp: regExpProtocol, Line: "This is an example line https://127.0.0.1:8081/artifactory/path/to/repo"}, "This is an example line https://127.0.0.1:8081/artifactory/path/to/repo", false},
		{"No http", gofrogio.CmdOutputPattern{RegExp: regExpProtocol, Line: "This is an example line"}, "This is an example line", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.regex.MatchedResults = test.regex.RegExp.FindStringSubmatch(test.regex.Line)
			if test.matched && len(test.regex.MatchedResults) > 3 {
				t.Error(fmt.Sprintf("Expected to find 3 results, however, found %d.", len(test.regex.MatchedResults)))
			}
			if test.matched && test.regex.MatchedResults[0] == "" {
				t.Error("Expected to find a match.")
			}
			if test.matched {
				actual := MaskCredentials(test.regex.Line, test.regex.MatchedResults[0])
				if !strings.EqualFold(actual, test.expectedLine) {
					t.Errorf("Expected: %s, The Regex found %s and the masked line: %s", test.expectedLine, test.regex.MatchedResults[0], actual)
				}
			}
			if !test.matched && len(test.regex.MatchedResults) != 0 {
				t.Error("Expected to find zero match, found:", test.regex.MatchedResults[0])
			}
		})
	}
}

func TestReturnErrorOnNotFound(t *testing.T) {
	regExpProtocol, err := GetRegExp(`^go: ([^\/\r\n]+\/[^\r\n\s:]*).*(404( Not Found)?[\s]?)$`)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name  string
		regex gofrogio.CmdOutputPattern
		error bool
	}{
		{"Without Error", gofrogio.CmdOutputPattern{RegExp: regExpProtocol, Line: "go: github.com/jfrog/jfrog-client-go@v0.2.1: This is an example line http://user:password@127.0.0.1:8081/artifactory/path/to/repo"}, false},
		{"With Error No Response Message", gofrogio.CmdOutputPattern{RegExp: regExpProtocol, Line: "go: github.com/jfrog/jfrog-client-go@v0.2.1: This is an example line http://user:password@127.0.0.1:8081/artifactory/path/to/repo: 404"}, true},
		{"With Error With response message", gofrogio.CmdOutputPattern{RegExp: regExpProtocol, Line: "go: github.com/jfrog/jfrog-client-go@v0.2.1: This is an example line http://user:password@127.0.0.1:8081/artifactory/path/to/repo: 404 Not Found"}, true},
		{"On Different Message", gofrogio.CmdOutputPattern{RegExp: regExpProtocol, Line: "go: finding github.com/elazarl/go-bindata-assetfs v0.0.0-20151224045452-57eb5e1fc594"}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.regex.MatchedResults = test.regex.RegExp.FindStringSubmatch(test.regex.Line)
			if test.error && len(test.regex.MatchedResults) < 3 {
				t.Error(fmt.Sprintf("Expected to find at least 3 results, however, found %d.", len(test.regex.MatchedResults)))
			}
			if test.error && test.regex.MatchedResults[0] == "" {
				t.Error("Expected to find 404 not found, found nothing.")
			}
			if !test.error && len(test.regex.MatchedResults) != 0 {
				t.Error("Expected regex to return empty result. Got:", test.regex.MatchedResults[0])
			}
		})
	}
}
