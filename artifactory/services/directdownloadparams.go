package services

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

type DirectDownloadParams struct {
	*utils.CommonParams
	Flat                    bool
	Explode                 bool
	Symlink                 bool
	ValidateSymlink         bool
	BypassArchiveInspection bool
	// Min split size in Kilobytes
	MinSplitSize int64
	SplitCount   int
	SkipChecksum bool

	// Optional fields (Sha256,Size) to avoid AQL request:
	Sha256 string
	// Size in bytes
	Size *int64
}

func (ddp *DirectDownloadParams) IsFlat() bool {
	return ddp.Flat
}

func (ddp *DirectDownloadParams) IsBypassArchiveInspection() bool {
	return ddp.BypassArchiveInspection
}

func (ddp *DirectDownloadParams) IsSymlink() bool {
	return ddp.Symlink
}

func (ddp *DirectDownloadParams) ValidateSymlinks() bool {
	return ddp.ValidateSymlink
}

func (ddp *DirectDownloadParams) IsExplode() bool {
	return ddp.Explode
}

func (ddp *DirectDownloadParams) GetFile() *utils.CommonParams {
	return ddp.CommonParams
}

func (ddp *DirectDownloadParams) IsSkipChecksum() bool {
	return ddp.SkipChecksum
}

func (ddp *DirectDownloadParams) IsExcludeArtifacts() bool {
	return ddp.ExcludeArtifacts
}

func (ddp *DirectDownloadParams) IsIncludeDeps() bool {
	return ddp.IncludeDeps
}

func NewDirectDownloadParams() DirectDownloadParams {
	return DirectDownloadParams{CommonParams: &utils.CommonParams{}, MinSplitSize: 5120, SplitCount: 3}
}
