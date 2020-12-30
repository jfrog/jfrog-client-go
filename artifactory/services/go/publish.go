package _go

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

var publishers []PublishGoPackage

type PublishGoPackage interface {
	isCompatible(artifactoryVersion string) bool
	PublishPackage(params GoParams, client *jfroghttpclient.JfrogHttpClient, ArtDetails auth.ServiceDetails) error
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
