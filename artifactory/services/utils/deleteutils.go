package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
)

func WildcardToDirsPath(deletePattern, searchResult string) (string, error) {
	if !strings.HasSuffix(deletePattern, "/") {
		return "", errors.New("Delete pattern must end with \"/\"")
	}

	regexpPattern := "^" + strings.Replace(deletePattern, "*", "([^/]*|.*)", -1)
	r, err := regexp.Compile(regexpPattern)
	errorutils.CheckError(err)
	if err != nil {
		return "", err
	}

	groups := r.FindStringSubmatch(searchResult)
	if len(groups) > 0 {
		return groups[0], nil
	}
	return "", nil
}

// Write all the dir results in 'bufferFiles' to the 'resultWriter'.
// However, skip dirs with artifact(s) that should not be deleted.
func WriteCandidateDirsToBeDeleted(bufferFiles []*content.ContentReader, artifactNotToBeDeleteReader *content.ContentReader, resultWriter *content.ContentWriter) (err error) {
	dirsToBeDeletedReader, err := MergeSortedFiles(bufferFiles)
	if err != nil {
		return
	}
	var candidateDirToBeDeletedPath string
	var artifactNotToBeDeletePath string
	var candidateDirToBeDeleted, artifactNotToBeDeleted *ResultItem
	for {
		// Fetch the next 'candidateDirToBeDeleted'.
		if candidateDirToBeDeleted == nil {
			candidateDirToBeDeleted = new(ResultItem)
			if err = dirsToBeDeletedReader.NextRecord(candidateDirToBeDeleted); err != nil {
				break
			}
			candidateDirToBeDeletedPath = strings.ToLower(candidateDirToBeDeleted.GetItemRelativePath())
		}
		// Fetch the next 'artifactNotToBeDelete'.
		if artifactNotToBeDeleted == nil {
			artifactNotToBeDeleted = new(ResultItem)
			if err = artifactNotToBeDeleteReader.NextRecord(artifactNotToBeDeleted); err != nil {
				// No artifacts left, write remaining dirs to be deleted to result file.
				resultWriter.Write(*candidateDirToBeDeleted)
				writeRemainCandidate(resultWriter, dirsToBeDeletedReader)
				break
			}
			artifactNotToBeDeletePath = strings.ToLower(artifactNotToBeDeleted.GetItemRelativePath())
		}
		// Found an 'artifact not to be deleted' in 'dir to be deleted', therefore skip writing the dir to the result file.
		if strings.HasPrefix(artifactNotToBeDeletePath, candidateDirToBeDeletedPath) {
			candidateDirToBeDeleted = nil
			continue
		}
		// 'artifactNotToBeDeletePath' & 'candidateDirToBeDeletedPath' are sorted, if 'candidateDirToBeDeleted'. As a result 'candidateDirToBeDeleted' cant be a prefix for any of the remaining artifacts.
		if artifactNotToBeDeletePath > candidateDirToBeDeletedPath {
			resultWriter.Write(*candidateDirToBeDeleted)
			candidateDirToBeDeleted = nil
			continue
		}
		artifactNotToBeDeleted = nil
	}
	err = artifactNotToBeDeleteReader.GetError()
	return
}

func writeRemainCandidate(cw *content.ContentWriter, mergeResult *content.ContentReader) {
	for toBeDeleted := new(ResultItem); mergeResult.NextRecord(toBeDeleted) == nil; toBeDeleted = new(ResultItem) {
		cw.Write(*toBeDeleted)
	}
}

func FilterCandidateToBeDeleted(deleteCandidates *content.ContentReader, resultWriter *content.ContentWriter) ([]*content.ContentReader, error) {
	paths := make(map[string]ResultItem)
	pathsKeys := make([]string, 0, MAX_BUFFER_SIZE)
	dirsToBeDeleted := []*content.ContentReader{}
	for candidate := new(ResultItem); deleteCandidates.NextRecord(candidate) == nil; candidate = new(ResultItem) {
		// Save all dirs candidate in a diffrent temp file.
		if candidate.Type == "folder" {
			pathsKeys = append(pathsKeys, candidate.GetItemRelativePath())
			paths[candidate.GetItemRelativePath()] = *candidate
			if len(pathsKeys) == MAX_BUFFER_SIZE {
				sortedCandidateDirsFile, err := SortAndSaveBufferToFile(paths, pathsKeys, true)
				if err != nil {
					return nil, err
				}
				dirsToBeDeleted = append(dirsToBeDeleted, sortedCandidateDirsFile)
				// Init buffer.
				paths = make(map[string]ResultItem)
				pathsKeys = make([]string, 0, MAX_BUFFER_SIZE)
			}
		} else {
			// Write none dir results.
			resultWriter.Write(candidate)
		}
	}
	if err := deleteCandidates.GetError(); err != nil {
		return nil, err
	}
	if len(pathsKeys) > 0 {
		sortedFile, err := SortAndSaveBufferToFile(paths, pathsKeys, true)
		if err != nil {
			return nil, err
		}
		dirsToBeDeleted = append(dirsToBeDeleted, sortedFile)
	}
	return dirsToBeDeleted, nil
}
