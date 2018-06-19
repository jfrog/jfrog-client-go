package bintray

import (
	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type Config interface {
	GetUrl() string
	GetKey() string
	GetThreads() int
	GetMinSplitSize() int64
	GetSplitCount() int
	GetMinChecksumDeploy() int64
	IsDryRun() bool
	GetBintrayDetails() auth.BintrayDetails
	GetLogger() log.Log
}
