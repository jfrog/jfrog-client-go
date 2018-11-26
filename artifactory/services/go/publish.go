package _go

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/httpclient"
)

var publishers []PublishGoPackage

type PublishGoPackage interface {
	isCompatible(artifactoryVersion string) bool
	PublishPackage(params GoParams, client *httpclient.HttpClient, ArtDetails auth.ArtifactoryDetails) error
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
