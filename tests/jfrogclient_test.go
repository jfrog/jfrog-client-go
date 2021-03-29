package tests

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/utils/tests"
)

const (
	JfrogTestsHome      = ".jfrogTest"
	JfrogHomeEnv        = "JFROG_CLI_HOME"
	CliIntegrationTests = "github.com/jfrog/jfrog-client-go/tests"
)

func TestMain(m *testing.M) {
	InitServiceManagers()
	result := m.Run()
	os.Exit(result)
}

func InitServiceManagers() {
	flag.Parse()
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
	log.Info(*RtUrl)
	log.Info(*DistUrl)
	log.Info(*XrayUrl)
	log.Info(*PipelinesUrl)
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
	err := createReposIfNeeded()
	if err != nil {
		log.Error(err.Error())
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
