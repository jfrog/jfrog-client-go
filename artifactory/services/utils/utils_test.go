package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/log"
)

func init() {
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
}

func getBaseTestDir(t *testing.T) string {
	pwd, err := os.Getwd()
	assert.NoError(t, err, "Failed to get current dir")
	if err != nil {
		return ""
	}
	return filepath.Join(pwd, "tests", "testdata")
}
