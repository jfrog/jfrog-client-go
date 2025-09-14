//go:build itest

package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessPing(t *testing.T) {
	initAccessTest(t)
	body, err := testsAccessPingService.Ping()
	assert.NoError(t, err)
	assert.Equal(t, "OK", string(body))
}
