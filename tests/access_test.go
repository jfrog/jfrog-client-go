package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccessVersion(t *testing.T) {
	initAccessTest(t)
	_, err := GetAccessDetails().GetVersion()
	assert.Error(t, err)
}

func initAccessTest(t *testing.T) {
	if !*TestAccess {
		t.Skip("Skipping access test. To run access test add the '-test.access=true' option.")
	}
}
