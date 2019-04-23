package bintray

import (
	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func NewConfigBuilder() *bintrayServicesConfigBuilder {
	configBuilder := &bintrayServicesConfigBuilder{}
	configBuilder.threads = 3
	return configBuilder
}

type bintrayServicesConfigBuilder struct {
	auth.BintrayDetails
	threads  int
	isDryRun bool
	logger   log.Log
}

func (builder *bintrayServicesConfigBuilder) SetBintrayDetails(artDetails auth.BintrayDetails) *bintrayServicesConfigBuilder {
	builder.BintrayDetails = artDetails
	return builder
}

func (builder *bintrayServicesConfigBuilder) SetThreads(threads int) *bintrayServicesConfigBuilder {
	builder.threads = threads
	return builder
}

func (builder *bintrayServicesConfigBuilder) SetDryRun(dryRun bool) *bintrayServicesConfigBuilder {
	builder.isDryRun = dryRun
	return builder
}

func (builder *bintrayServicesConfigBuilder) Build() Config {
	c := &bintrayServicesConfig{}
	c.BintrayDetails = builder.BintrayDetails
	c.threads = builder.threads
	c.logger = builder.logger
	return c
}

func (builder *bintrayServicesConfigBuilder) SetLogger(logger log.Log) *bintrayServicesConfigBuilder {
	builder.logger = logger
	return builder
}
