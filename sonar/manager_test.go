package sonar

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/config"
	sonarauth "github.com/jfrog/jfrog-client-go/sonar/auth"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	serviceDetails := sonarauth.NewSonarDetails()
	serviceDetails.SetUrl("https://sonarcloud.io")
	serviceDetails.SetAccessToken("test-token")

	cfg, err := config.NewConfigBuilder().SetServiceDetails(serviceDetails).Build()
	assert.NoError(t, err)

	manager, err := NewManager(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.IsType(t, &sonarManager{}, manager)
}

func TestNewManager_WithInvalidConfig(t *testing.T) {
	// Test with nil config
	manager, err := NewManager(nil)
	assert.Error(t, err)
	assert.Nil(t, manager)
}

func TestSonarManager_Client(t *testing.T) {
	serviceDetails := sonarauth.NewSonarDetails()
	serviceDetails.SetUrl("https://sonarcloud.io")
	serviceDetails.SetAccessToken("test-token")

	cfg, err := config.NewConfigBuilder().SetServiceDetails(serviceDetails).Build()
	assert.NoError(t, err)

	manager, err := NewManager(cfg)
	assert.NoError(t, err)

	sm := manager.(*sonarManager)
	client := sm.Client()
	assert.NotNil(t, client)
}

func TestSonarManager_InterfaceMethods(t *testing.T) {
	serviceDetails := sonarauth.NewSonarDetails()
	serviceDetails.SetUrl("https://sonarcloud.io")
	serviceDetails.SetAccessToken("test-token")

	cfg, err := config.NewConfigBuilder().SetServiceDetails(serviceDetails).Build()
	assert.NoError(t, err)

	manager, err := NewManager(cfg)
	assert.NoError(t, err)

	// Call interface methods; network outcome is not asserted
	_, _ = manager.GetQualityGateAnalysis("test-analysis-id")
	_, _ = manager.GetTaskDetails("test-task-id")
	_, _ = manager.GetSonarIntotoStatementRaw("test-task-id")
}
