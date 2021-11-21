package buildinfo

import (
	"github.com/jfrog/jfrog-client-go/auth"
)

type Configuration struct {
	ArtDetails auth.ServiceDetails
	BuildUrl   string
	DryRun     bool
	EnvInclude string
	EnvExclude string
}

func (config *Configuration) GetArtifactoryDetails() auth.ServiceDetails {
	return config.ArtDetails
}

func (config *Configuration) SetArtifactoryDetails(artDetails auth.ServiceDetails) {
	config.ArtDetails = artDetails
}

func (config *Configuration) IsDryRun() bool {
	return config.DryRun
}
