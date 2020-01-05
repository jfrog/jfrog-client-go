package tests

import (
	"flag"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/utils/tests"
	"os"
	"path/filepath"
	"testing"
)

const (
	JfrogTestsHome      = ".jfrogTest"
	JfrogHomeEnv        = "JFROG_CLI_HOME"
	CliIntegrationTests = "github.com/jfrog/jfrog-client-go/tests"
)

func TestMain(m *testing.M) {
	InitArtifactoryServiceManager()
	result := m.Run()
	os.Exit(result)
}

func InitArtifactoryServiceManager() {
	flag.Parse()
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
	createArtifactoryUploadManager()
	createArtifactorySearchManager()
	createArtifactoryDeleteManager()
	createArtifactoryDownloadManager()
	createArtifactorySecurityManager()
	createReposIfNeeded()
}

func TestUnitTests(t *testing.T) {
	homePath, err := filepath.Abs(JfrogTestsHome)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	setJfrogHome(homePath)
	packages := tests.GetTestPackages("./../...")
	packages = tests.ExcludeTestsPackage(packages, CliIntegrationTests)
	tests.RunTests(packages, false)
	cleanUnitTestsJfrogHome(homePath)
}

func setJfrogHome(homePath string) {
	if err := os.Setenv(JfrogHomeEnv, homePath); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func cleanUnitTestsJfrogHome(homePath string) {
	os.RemoveAll(homePath)
	if err := os.Unsetenv(JfrogHomeEnv); err != nil {
		os.Exit(1)
	}
}
