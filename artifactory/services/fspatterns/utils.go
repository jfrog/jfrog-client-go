package fspatterns

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jfrog/gofrog/crypto"
	"github.com/jfrog/jfrog-client-go/utils/log"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
)

// Return all the existing paths of the provided root path
func ListFiles(rootPath string, isRecursive, includeDirs, excludeWithRelativePath, preserveSymlink bool, excludePathPattern string) ([]string, error) {
	return ListFilesFilterPatternAndSize(rootPath, isRecursive, includeDirs, excludeWithRelativePath, preserveSymlink, excludePathPattern, nil)
}

// Return all the existing paths of the provided root path
func ListFilesFilterPatternAndSize(rootPath string, isRecursive, includeDirs, excludeWithRelativePath, preserveSymlink bool, excludePathPattern string, sizeThreshold *SizeThreshold) ([]string, error) {
	filterFunc := filterFilesFunc(rootPath, includeDirs, excludeWithRelativePath, preserveSymlink, excludePathPattern, sizeThreshold)
	return fileutils.ListFilesWithFilterFunc(rootPath, isRecursive, !preserveSymlink, filterFunc)
}

// Transform to regexp and prepare Exclude patterns to be used, exclusion patterns must be absolute paths.
func PrepareExcludePathPattern(exclusions []string, patternType utils.PatternType, isRecursive bool) string {
	excludePathPattern := ""

	for _, singleExclusion := range exclusions {
		if len(singleExclusion) > 0 {
			singleExclusion = utils.ReplaceTildeWithUserHome(singleExclusion)
			singleExclusion = utils.ConvertLocalPatternToRegexp(singleExclusion, patternType)
			if isRecursive && strings.HasSuffix(singleExclusion, fileutils.GetFileSeparator()) {
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

// Returns a function that filters files according to the provided parameters
func filterFilesFunc(rootPath string, includeDirs, excludeWithRelativePath, preserveSymlink bool, excludePathPattern string, sizeThreshold *SizeThreshold) func(filePath string) (included bool, err error) {
	return func(path string) (included bool, err error) {
		if path == "." {
			return false, nil
		}
		if !includeDirs {
			isDir, err := fileutils.IsDirExists(path, preserveSymlink)
			if err != nil || isDir {
				return false, err
			}
		}
		var isExcludedByPattern bool
		isExcludedByPattern, err = isPathExcluded(path, excludePathPattern, rootPath, excludeWithRelativePath)
		if err != nil {
			return false, err
		}
		if isExcludedByPattern {
			log.Debug(fmt.Sprintf("The path '%s' is excluded", path))
			return false, nil
		}

		if sizeThreshold != nil {
			fileInfo, err := fileutils.GetFileInfo(path, preserveSymlink)
			if err != nil {
				return false, errorutils.CheckError(err)
			}
			// Check if the file size is within the limits
			if !fileInfo.IsDir() && !sizeThreshold.IsSizeWithinThreshold(fileInfo.Size()) {
				log.Debug(fmt.Sprintf("The path '%s' is excluded", path))
				return false, nil
			}
		}
		return true, nil
	}
}

// Return the actual sub-paths that match the regex provided.
// Excluded sub-paths are not returned
func SearchPatterns(path string, preserveSymlinks, includeDirs bool, regexp *regexp.Regexp) (matches []string, isDir bool, err error) {
	isDir, err = fileutils.IsDirExists(path, false)
	if err != nil {
		return
	}
	isSymlinkFlow := preserveSymlinks && fileutils.IsPathSymlink(path)
	if isDir && !includeDirs && !isSymlinkFlow {
		return
	}
	// Upload directory. We ignore IsDir in a symlink flow since we want to create a dummy file instead that holds the symlink property.
	// Properties cannot be assigned to repositories in Artifactory.
	if isSymlinkFlow {
		isDir = false
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

func isPathExcluded(path, excludePathPattern, rootPath string, excludeWithRelativePath bool) (excludedPath bool, err error) {
	if len(excludePathPattern) > 0 {
		if excludeWithRelativePath {
			path = strings.TrimPrefix(path, rootPath)
		}
		excludedPath, err = regexp.MatchString(excludePathPattern, path)
		err = errorutils.CheckError(err)
	}
	return
}

// If filePath is path to a symlink we should return the link content e.g where the link points
func GetFileSymlinkPath(filePath string) (string, error) {
	if filePath == "" {
		return "", nil
	}
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

// Find parentheses in 'target' and 'archive-target', merge the results to one slice with no duplication.
func getPlaceholderParentheses(pattern, target, archiveTarget string) utils.ParenthesesSlice {
	targetParentheses := utils.CreateParenthesesSlice(pattern, target)
	archiveTargetParentheses := utils.CreateParenthesesSlice(pattern, archiveTarget)
	parenthesesMap := make(map[utils.Parentheses]bool)
	var parenthesesSlice []utils.Parentheses
	// Target parentheses
	for _, v := range targetParentheses.Parentheses {
		parenthesesSlice = append(parenthesesSlice, v)
		parenthesesMap[v] = true
	}
	// Archive target parentheses
	for _, v := range archiveTargetParentheses.Parentheses {
		if parenthesesMap[v] {
			continue
		}
		parenthesesSlice = append(parenthesesSlice, v)
		parenthesesMap[v] = true
	}
	return utils.NewParenthesesSlice(parenthesesSlice)
}

// Get the local root path, from which to start collecting artifacts to be uploaded to Artifactory.
// If path does not exist error will be returned.
func GetRootPath(pattern, target, archiveTarget string, patternType utils.PatternType, preserveSymLink bool) (string, error) {
	placeholderParentheses := getPlaceholderParentheses(pattern, target, archiveTarget)
	rootPath := utils.GetRootPath(pattern, patternType, placeholderParentheses)
	if !fileutils.IsPathExists(rootPath, preserveSymLink) {
		return "", errorutils.CheckErrorf("path does not exist: %s", rootPath)
	}

	return rootPath, nil
}

// When handling symlink we want to simulate the creation of empty file
func CreateSymlinkFileDetails() (*fileutils.FileDetails, error) {
	checksums, err := crypto.CalcChecksums(bytes.NewBuffer([]byte(fileutils.SymlinkFileContent)))
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	details := new(fileutils.FileDetails)
	details.Checksum.Md5 = checksums[crypto.MD5]
	details.Checksum.Sha1 = checksums[crypto.SHA1]
	details.Checksum.Sha256 = checksums[crypto.SHA256]
	details.Size = int64(0)
	return details, nil
}
