package services

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

type DirectDownloadParams struct {
	*utils.CommonParams
	Pattern            string
	Target             string
	Flat               bool
	Recursive          bool
	Exclusions         []string
	SkipChecksum       bool
	MinSplitSizeMB     int64
	SplitCount         int
	Retries            int
	RetryWaitMilliSecs int
	SyncDeletesPath    string
	Quiet              bool
}

func (ddp *DirectDownloadParams) GetPattern() string {
	return ddp.Pattern
}

func (ddp *DirectDownloadParams) SetPattern(pattern string) {
	ddp.Pattern = pattern
}

func (ddp *DirectDownloadParams) GetTarget() string {
	return ddp.Target
}

func (ddp *DirectDownloadParams) SetTarget(target string) {
	ddp.Target = target
}

func (ddp *DirectDownloadParams) IsFlat() bool {
	return ddp.Flat
}

func (ddp *DirectDownloadParams) SetFlat(flat bool) {
	ddp.Flat = flat
}

func (ddp *DirectDownloadParams) IsRecursive() bool {
	return ddp.Recursive
}

func (ddp *DirectDownloadParams) SetRecursive(recursive bool) {
	ddp.Recursive = recursive
}

func (ddp *DirectDownloadParams) GetExclusions() []string {
	return ddp.Exclusions
}

func (ddp *DirectDownloadParams) SetExclusions(exclusions []string) {
	ddp.Exclusions = exclusions
}

func (ddp *DirectDownloadParams) IsSkipChecksum() bool {
	return ddp.SkipChecksum
}

func (ddp *DirectDownloadParams) SetSkipChecksum(skipChecksum bool) {
	ddp.SkipChecksum = skipChecksum
}

func (ddp *DirectDownloadParams) GetSyncDeletesPath() string {
	return ddp.SyncDeletesPath
}

func (ddp *DirectDownloadParams) SetSyncDeletesPath(syncDeletesPath string) {
	ddp.SyncDeletesPath = syncDeletesPath
}

func (ddp *DirectDownloadParams) GetRetries() int {
	return ddp.Retries
}

func (ddp *DirectDownloadParams) SetRetries(retries int) {
	ddp.Retries = retries
}

func (ddp *DirectDownloadParams) IsQuiet() bool {
	return ddp.Quiet
}

func (ddp *DirectDownloadParams) SetQuiet(quiet bool) {
	ddp.Quiet = quiet
}

func NewDirectDownloadParams() *DirectDownloadParams {
	return &DirectDownloadParams{
		CommonParams:       &utils.CommonParams{},
		MinSplitSizeMB:     5120,
		SplitCount:         3,
		Retries:            3,
		RetryWaitMilliSecs: 0,
		Recursive:          true,
	}
}
