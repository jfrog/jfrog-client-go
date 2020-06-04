package tests

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/utils/tests"
)

const (
	JfrogTestsHome      = ".jfrogTest"
	JfrogHomeEnv        = "JFROG_CLI_HOME"
	CliIntegrationTests = "github.com/jfrog/jfrog-client-go/tests"
)

func TestMain(m *testing.M) {
	InitArtifactoryServiceManager()
	result := m.Run()
	cleanUp()
	os.Exit(result)
}

func InitArtifactoryServiceManager() {
	flag.Parse()
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
	// Create temp dir for integration tests
	err := fileutils.CreateReaderWriterTempDir()
	if err != nil {
		log.Error(("Creating temp folder failed: " + err.Error()))
		os.Exit(1)
	}
	createArtifactoryUploadManager()
	createArtifactorySearchManager()
	createArtifactoryDeleteManager()
	createArtifactoryDownloadManager()
	createArtifactorySecurityManager()
	createArtifactoryCreateLocalRepositoryManager()
	createArtifactoryUpdateLocalRepositoryManager()
	createArtifactoryCreateRemoteRepositoryManager()
	createArtifactoryUpdateRemoteRepositoryManager()
	createArtifactoryCreateVirtualRepositoryManager()
	createArtifactoryUpdateVirtualRepositoryManager()
	createArtifactoryDeleteRepositoryManager()
	createArtifactoryGetRepositoryManager()
	createArtifactoryReplicationCreateManager()
	createArtifactoryReplicationUpdateManager()
	createArtifactoryReplicationGetManager()
	createArtifactoryReplicationDeleteManager()
	if *DistUrl != "" {
		createDistributionManager()
	}
	createReposIfNeeded()
}

func cleanUp() {
	err := fileutils.CleanupReaderWriterTempFilesAndDirs()
	if err != nil {
		log.Error(("Deleting temp folder failed: " + err.Error()))
		os.Exit(1)
	}
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
