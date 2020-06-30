package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildUrlForGetStatus(t *testing.T) {
	service := NewDistributionStatusService(nil)

	// Get all release bundles
	url := service.BuildUrlForGetStatus("https://dummy-url/distribution/", "", "", "")
	assert.Equal(t, "https://dummy-url/distribution/api/v1/release_bundle/distribution", url)

	// Get release bundles by name
	url = service.BuildUrlForGetStatus("https://dummy-url/distribution/", "bundleName", "", "")
	assert.Equal(t, "https://dummy-url/distribution/api/v1/release_bundle/bundleName/distribution", url)

	// Get release bundle by name and version
	url = service.BuildUrlForGetStatus("https://dummy-url/distribution/", "bundleName", "22", "")
	assert.Equal(t, "https://dummy-url/distribution/api/v1/release_bundle/bundleName/22/distribution", url)

	// Get release bundle by name version and tracker ID
	url = service.BuildUrlForGetStatus("https://dummy-url/distribution/", "bundleName", "22", "123234")
	assert.Equal(t, "https://dummy-url/distribution/api/v1/release_bundle/bundleName/22/distribution/123234", url)
}
