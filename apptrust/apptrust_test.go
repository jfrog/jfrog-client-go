package apptrust

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/apptrust/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	apptrustDetails := auth.NewApptrustDetails()
	apptrustDetails.SetUrl("http://localhost:8080/")
	apptrustDetails.SetUser("admin")
	apptrustDetails.SetPassword("password")

	config, err := config.NewConfigBuilder().
		SetServiceDetails(apptrustDetails).
		Build()
	assert.NoError(t, err)

	manager, err := New(config)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.Client())
}

func TestApptrustDetails(t *testing.T) {
	apptrustDetails := auth.NewApptrustDetails()
	assert.NotNil(t, apptrustDetails)

	// Test setting URL
	testUrl := "http://localhost:8080/"
	apptrustDetails.SetUrl(testUrl)
	assert.Equal(t, testUrl, apptrustDetails.GetUrl())

	// Test setting user
	testUser := "testuser"
	apptrustDetails.SetUser(testUser)
	assert.Equal(t, testUser, apptrustDetails.GetUser())
}
