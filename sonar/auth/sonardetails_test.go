package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSonarDetails(t *testing.T) {
	details := NewSonarDetails()
	assert.NotNil(t, details)
	assert.IsType(t, &sonarDetails{}, details)
}

func TestSonarDetails_GetVersion(t *testing.T) {
	details := NewSonarDetails()

	assert.Panics(t, func() {
		details.GetVersion()
	}, "GetVersion should panic as it's not implemented")
}
