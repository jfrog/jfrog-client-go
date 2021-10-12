package tests

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/utils/tests"
)

const (
	JfrogTestsHome      = ".jfrogTest"
	JfrogHomeEnv        = "JFROG_CLI_HOME"
	CliIntegrationTests = "github.com/jfrog/jfrog-client-go/tests"
)

func TestMain(m *testing.M) {
	exitCode := setupIntegrationTests()
	if exitCode != 0 {
		os.Exit(exitCode)
	}
	result := m.Run()
	exitCode = teardownIntegrationTests()
	if result == 0 {
		os.Exit(exitCode)
	}
	os.Exit(result)
}

func setupIntegrationTests() int {
	flag.Parse()
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
	if *TestArtifactory || *TestDistribution || *TestXray {
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
		createArtifactoryCreateFederatedRepositoryManager()
		createArtifactoryUpdateFederatedRepositoryManager()
		createArtifactoryDeleteRepositoryManager()
		createArtifactoryGetRepositoryManager()
		createArtifactoryReplicationCreateManager()
		createArtifactoryReplicationUpdateManager()
		createArtifactoryReplicationGetManager()
		createArtifactoryReplicationDeleteManager()
		createArtifactoryPermissionTargetManager()
		createArtifactoryUserManager()
		createArtifactoryGroupManager()
		createArtifactoryBuildInfoManager()
		createArtifactoryFederationManager()
	}

	if *TestDistribution {
		createDistributionManager()
	}
	if *TestXray {
		createXrayWatchManager()
		createXrayPolicyManager()
		createXrayBinMgrManager()
	}
	if *TestPipelines {
		createPipelinesIntegrationsManager()
		createPipelinesSourcesManager()
	}
	if *TestAccess {
		createAccessProjectManager()
	}
	if err := createRepo(); err != nil {
		log.Error(err.Error())
		return 1
	}
	return 0
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
	err = tests.RunTests(packages, false)
	assert.NoError(t, err)
	cleanUnitTestsJfrogHome(homePath)
}

func setJfrogHome(homePath string) {
	if err := os.Setenv(JfrogHomeEnv, homePath); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func cleanUnitTestsJfrogHome(homePath string) {
	err := os.RemoveAll(homePath)
	if err != nil {
		log.Error(err.Error())
	}
	if err := os.Unsetenv(JfrogHomeEnv); err != nil {
		os.Exit(1)
	}
}
