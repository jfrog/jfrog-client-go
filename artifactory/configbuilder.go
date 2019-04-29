package artifactory

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func NewConfigBuilder() *artifactoryServicesConfigBuilder {
	configBuilder := &artifactoryServicesConfigBuilder{}
	configBuilder.threads = 3
	return configBuilder
}

type artifactoryServicesConfigBuilder struct {
	auth.ArtifactoryDetails
	certificatesPath string
	threads          int
	isDryRun         bool
	logger           log.Log
	insecureTls      bool
}

func (builder *artifactoryServicesConfigBuilder) SetArtDetails(artDetails auth.ArtifactoryDetails) *artifactoryServicesConfigBuilder {
	builder.ArtifactoryDetails = artDetails
	return builder
}

func (builder *artifactoryServicesConfigBuilder) SetCertificatesPath(certificatesPath string) *artifactoryServicesConfigBuilder {
	builder.certificatesPath = certificatesPath
	return builder
}

func (builder *artifactoryServicesConfigBuilder) SetThreads(threads int) *artifactoryServicesConfigBuilder {
	builder.threads = threads
	return builder
}

func (builder *artifactoryServicesConfigBuilder) SetDryRun(dryRun bool) *artifactoryServicesConfigBuilder {
	builder.isDryRun = dryRun
	return builder
}

func (builder *artifactoryServicesConfigBuilder) SetInsecureTls(insecureTls bool) *artifactoryServicesConfigBuilder {
	builder.insecureTls = insecureTls
	return builder
}

func (builder *artifactoryServicesConfigBuilder) Build() (Config, error) {
	c := &artifactoryServicesConfig{}
	c.ArtifactoryDetails = builder.ArtifactoryDetails
	c.threads = builder.threads
	c.logger = builder.logger
	c.certificatesPath = builder.certificatesPath
	c.dryRun = builder.isDryRun
	c.insecureTls = builder.insecureTls
	return c, nil
}

func (builder *artifactoryServicesConfigBuilder) SetLogger(logger log.Log) *artifactoryServicesConfigBuilder {
	builder.logger = logger
	return builder
}
