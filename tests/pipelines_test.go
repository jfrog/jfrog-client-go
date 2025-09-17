//go:build itest

package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipelinesVersion(t *testing.T) {
	initPipelinesTest(t)
	version, err := GetPipelinesDetails().GetVersion()
	if !assert.NoError(t, err) {
		return
	}
	assert.NotEmpty(t, version)
}

func initPipelinesTest(t *testing.T) {
	if !*TestPipelines {
		t.Skip("Skipping pipelines test. To run pipelines test add the '-test.pipelines=true' option.")
	}
}
