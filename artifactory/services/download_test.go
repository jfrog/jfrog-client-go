package services

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetResultItemBuild(t *testing.T) {
	props := []utils.Property{
		{Key: "build.name", Value: "MyBuild"},
		{Key: "build.number", Value: "1"},
		{Key: "build.timestamp", Value: "654321"},
		{Key: "other-key", Value: "other-value"},
	}
	assert.Equal(t, getResultItemBuild(props), "MyBuild/1/654321")
}
