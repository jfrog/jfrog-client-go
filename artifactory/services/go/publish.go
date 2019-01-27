package _go

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/utils/httpclient"
)

var publishers []PublishGoPackage

type PublishGoPackage interface {
	isCompatible(artifactoryVersion string) bool
	PublishPackage(params GoParams, client *rthttpclient.ArtifactoryHttpClient, ArtDetails auth.ArtifactoryDetails) error
}

func register(publishApi PublishGoPackage) {
	publishers = append(publishers, publishApi)
}

// Returns the compatible publisher to Artifactory
func GetCompatiblePublisher(artifactoryVersion string) PublishGoPackage {
	for _, publisher := range publishers {
		if publisher.isCompatible(artifactoryVersion) {
			return publisher
		}
	}
	return nil
}
