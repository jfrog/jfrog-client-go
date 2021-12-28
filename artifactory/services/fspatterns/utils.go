package fspatterns

import (
	"bytes"
	"fmt"
	biutils "github.com/jfrog/build-info-go/utils"
	"os"
	"regexp"
	"strings"

	serviceutils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
)

// Return all the existing paths of the provided root path
func GetPaths(rootPath string, isRecursive, includeDirs, isSymlink bool) ([]string, error) {
	var paths []string
	var err error
	if isRecursive {
		paths, err = fileutils.ListFilesRecursiveWalkIntoDirSymlink(rootPath, !isSymlink)
	} else {
		paths, err = fileutils.ListFiles(rootPath, includeDirs)
	}
	if err != nil {
		return paths, err
	}
	return paths, nil
}

// Transform to regexp and prepare Exclude patterns to be used
func PrepareExcludePathPattern(params serviceutils.FileGetter) string {
	excludePathPattern := ""
	for _, singleExclusion := range params.GetExclusions() {
		if len(singleExclusion) > 0 {
			singleExclusion = utils.ReplaceTildeWithUserHome(singleExclusion)
			singleExclusion = utils.ConvertLocalPatternToRegexp(singleExclusion, params.GetPatternType())
			if params.IsRecursive() && strings.HasSuffix(singleExclusion, fileutils.GetFileSeparator()) {
				singleExclusion += "*"
			}
			excludePathPattern += fmt.Sprintf(`(%s)|`, singleExclusion)
		}
	}
	if len(excludePathPattern) > 0 {
		excludePathPattern = excludePathPattern[:len(excludePathPattern)-1]
	}
	return excludePathPattern
}

// Return only subpaths of the provided by the user path that matched to the provided regexp.
// Subpaths that matched to an exclude pattern won't returned
func PrepareAndFilterPaths(path, excludePathPattern string, preserveSymlinks, includeDirs bool, regexp *regexp.Regexp) (matches []string, isDir, isSymlinkFlow bool, err error) {
	isDir, err = fileutils.IsDirExists(path, false)
	if err != nil {
		return
	}

	excludedPath, err := IsPathExcluded(path, excludePathPattern)
	if err != nil {
		return
	}

	if excludedPath {
		return
	}
	isSymlinkFlow = preserveSymlinks && fileutils.IsPathSymlink(path)

	if isDir && !includeDirs && !isSymlinkFlow {
		return
	}
	matches = regexp.FindStringSubmatch(path)
	return
}

func GetSingleFileToUpload(rootPath, targetPath string, flat bool) (utils.Artifact, error) {
	symlinkPath, err := GetFileSymlinkPath(rootPath)
	if err != nil {
		return utils.Artifact{}, err
	}

	var uploadPath string
	if !strings.HasSuffix(targetPath, "/") {
		uploadPath = targetPath
	} else {
		localPath := rootPath

		if flat {
			uploadPath, _ = fileutils.GetFileAndDirFromPath(localPath)
			uploadPath = targetPath + uploadPath
		} else {
			uploadPath = targetPath + localPath
			uploadPath = utils.TrimPath(uploadPath)
		}
	}

	return utils.Artifact{LocalPath: rootPath, TargetPath: uploadPath, SymlinkTargetPath: symlinkPath}, nil
}

func IsPathExcluded(path string, excludePathPattern string) (excludedPath bool, err error) {
	if len(excludePathPattern) > 0 {
		excludedPath, err = regexp.MatchString(excludePathPattern, path)
	}
	return
}

// If filePath is path to a symlink we should return the link content e.g where the link points
func GetFileSymlinkPath(filePath string) (string, error) {
	fileInfo, e := os.Lstat(filePath)
	if errorutils.CheckError(e) != nil {
		return "", e
	}
	var symlinkPath = ""
	if fileutils.IsFileSymlink(fileInfo) {
		symlinkPath, e = os.Readlink(filePath)
		if errorutils.CheckError(e) != nil {
			return "", e
		}
	}
	return symlinkPath, nil
}

// Get the local root path, from which to start collecting artifacts to be uploaded to Artifactory.
// If path dose not exist error will be returned.
func GetRootPath(pattern, target string, patternType clientutils.PatternType, preserveSymLink bool) (string, error) {
	placeholderParentheses := clientutils.NewParenthesesSlice(pattern, target)
	rootPath := utils.GetRootPath(pattern, patternType, placeholderParentheses)
	if !fileutils.IsPathExists(rootPath, preserveSymLink) {
		return "", errorutils.CheckErrorf("Path does not exist: " + rootPath)
	}

	return rootPath, nil
}

// When handling symlink we want to simulate the creation of empty file
func CreateSymlinkFileDetails() (*fileutils.FileDetails, error) {
	checksumInfo, err := biutils.CalcChecksums(bytes.NewBuffer([]byte(fileutils.SymlinkFileContent)))
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	details := new(fileutils.FileDetails)
	details.Checksum.Md5 = checksumInfo[biutils.MD5]
	details.Checksum.Sha1 = checksumInfo[biutils.SHA1]
	details.Checksum.Sha256 = checksumInfo[biutils.SHA256]
	details.Size = int64(0)
	return details, nil
}
