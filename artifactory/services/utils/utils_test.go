package utils

import (
	"os"
	"path/filepath"

	"github.com/jfrog/jfrog-client-go/utils/log"
)

func init() {
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
}

func getBaseTestDir() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(pwd, "tests", "testdata"), nil
}
