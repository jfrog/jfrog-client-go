package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccessVersion(t *testing.T) {
	initAccessTest(t)
	version, err := GetAccessDetails().GetVersion()
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.NotEmpty(t, version)
}

func initAccessTest(t *testing.T) {
	if !*TestAccess {
		t.Skip("Skipping access test. To run access test add the '-test.access=true' option.")
	}
}
