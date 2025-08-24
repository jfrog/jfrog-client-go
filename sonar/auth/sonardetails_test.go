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
