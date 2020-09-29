package utils

import (
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"strings"
)

// Returns an Artifact struct for uploading.
// rootPath - Candidate file to upload.
// targetPath - Path in Artifactory uploading to.
// symlinkTarget - If candidate file is a symlink, this is the symlink target.
// isFlat - Value of 'flat' flag.
// preserveSymlink - Value of 'symlinks' flag.
func GetArtifactToUpload(rootPath, targetPath, symlinkTarget string, isFlat, preserveSymlink bool) clientutils.Artifact {
	uploadTarget := getUploadTarget(rootPath, targetPath, isFlat)
	if preserveSymlink || symlinkTarget == "" {
		// If preserving symlinks or symlink target is empty, use root path name for upload (symlink itself / regular file).
		return clientutils.Artifact{LocalPath: rootPath, TargetPath: uploadTarget, Symlink: symlinkTarget}
	}
	// Actual file to upload is the symlink path.
	return clientutils.Artifact{LocalPath: symlinkTarget, TargetPath: uploadTarget, Symlink: symlinkTarget}
}

func getUploadTarget(fileToUploadPath, targetInArtifactory string, flat bool) string {
	var uploadTarget string
	if !strings.HasSuffix(targetInArtifactory, "/") {
		uploadTarget = targetInArtifactory
	} else {
		if flat {
			targetFileName, _ := fileutils.GetFileAndDirFromPath(fileToUploadPath)
			uploadTarget = targetInArtifactory + targetFileName
		} else {
			uploadTarget = targetInArtifactory + clientutils.TrimPath(fileToUploadPath)
		}
	}
	return uploadTarget
}
