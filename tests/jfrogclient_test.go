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
	setupIntegrationTests()
	result := m.Run()
	teardownIntegrationTests()
	os.Exit(result)
}

func setupIntegrationTests() {
	flag.Parse()
	log.SetLogger(log.NewLogger(log.DEBUG, nil))

	if *TestArtifactory || *TestDistribution || *TestXray || *TestRepositories || *TestMultipartUpload {
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
		createArtifactoryGetPackageManager()
		createArtifactoryReplicationCreateManager()
		createArtifactoryReplicationUpdateManager()
		createArtifactoryReplicationGetManager()
		createArtifactoryReplicationDeleteManager()
		createArtifactoryPermissionTargetManager()
		createArtifactoryUserManager()
		createArtifactoryGroupManager()
		createArtifactoryBuildInfoManager()
		createArtifactoryFederationManager()
		createArtifactorySystemManager()
		createArtifactoryStorageManager()
		createArtifactoryAqlManager()
	}

	if *TestDistribution {
		createDistributionManager()
	}
	if *TestXray {
		createXrayWatchManager()
		createXrayPolicyManager()
		createXrayBinMgrManager()
		createXrayIgnoreRuleManager()
	}
	if *TestPipelines {
		createPipelinesIntegrationsManager()
		createPipelinesSourcesManager()
		createPipelinesRunManager()
		createPipelinesSyncManager()
		createPipelinesSyncStatusManager()
	}
	if *TestAccess {
		createAccessPingManager()
		createAccessProjectManager()
		createAccessInviteManager()
		createAccessTokensManager()
	}
	if err := createRepo(); err != nil {
		log.Error(err.Error())
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
	err = tests.RunTests(packages, false)
	assert.NoError(t, err)
	cleanUnitTestsJfrogHome(t, homePath)
}

func setJfrogHome(homePath string) {
	if err := os.Setenv(JfrogHomeEnv, homePath); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func cleanUnitTestsJfrogHome(t *testing.T, homePath string) {
	tests.RemoveAllAndAssert(t, homePath)
	if err := os.Unsetenv(JfrogHomeEnv); err != nil {
		os.Exit(1)
	}
}
