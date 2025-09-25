//go:build itest

package tests

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/utils/tests"
)

const (
	JfrogHomeEnv        = "JFROG_CLI_HOME"
	CliIntegrationTests = "github.com/jfrog/jfrog-client-go/tests"
)

func TestMain(m *testing.M) {
	setupIntegrationTests()
	os.Exit(m.Run())
}

func setupIntegrationTests() {
	flag.Parse()
	log.SetLogger(log.NewLogger(log.DEBUG, nil))

	checkFlags()

	if *TestUnit {
		return
	}

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
}

func TestUnitTests(t *testing.T) {
	initUnitTests(t)

	homePath := t.TempDir()

	setJfrogHome(t, homePath)

	packages := tests.GetTestPackages("./../...")
	packages = tests.ExcludeTestsPackage(packages, CliIntegrationTests)
	err := tests.RunTests(packages, false)
	assert.NoError(t, err)
}

func initUnitTests(t *testing.T) {
	if !*TestUnit {
		t.Skip("Skipping unit tests. To run unit tests add the '-test.unit=true' option.")
	}
}

func setJfrogHome(t *testing.T, homePath string) {
	err := os.Setenv(JfrogHomeEnv, homePath)
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.Unsetenv(JfrogHomeEnv)
		if err != nil {
			log.Warn(fmt.Sprintf("Failed to unset env %s: %s", JfrogHomeEnv, err.Error()))
		}
	})
}
