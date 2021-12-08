package utils

import (
	"github.com/jfrog/jfrog-client-go/utils/tests"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/log"
)

func init() {
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
}

func getBaseTestDir(t *testing.T) string {
	pwd := tests.GetwdAndAssert(t)
	if pwd == "" {
		return ""
	}
	return filepath.Join(pwd, "tests", "testdata")
}
