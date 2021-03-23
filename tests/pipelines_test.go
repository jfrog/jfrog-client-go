package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipelinesVersion(t *testing.T) {
	if *PipelinesUrl == "" {
		t.Skip("Pipelines is not being tested, skipping...")
	}

	version, err := GetPipelinesDetails().GetVersion()
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.NotEmpty(t, version)
}
