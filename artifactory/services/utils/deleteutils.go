package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
)

func WildcardToDirsPath(deletePattern, searchResult string) (string, error) {
	if !strings.HasSuffix(deletePattern, "/") {
		return "", errors.New("delete pattern must end with \"/\"")
	}

	regexpPattern := "^" + strings.Replace(deletePattern, "*", "([^/]*|.*)", -1)
	r, err := regexp.Compile(regexpPattern)
	if err != nil {
		return "", errorutils.CheckError(err)
	}

	groups := r.FindStringSubmatch(searchResult)
	if len(groups) > 0 {
		return groups[0], nil
	}
	return "", nil
}

// Write all the dirs to be deleted into 'resultWriter'.
// However, skip dirs with files(s) that should not be deleted.
// In order to accomplish this, we check if the dirs are a prefix of any artifact, witch means the folder contains the artifact and should not be deleted.
// Optimization: In order not to scan for each dir the entire artifact reader and see if it is a prefix or not, we rely on the fact that the dirs and artifacts are sorted.
// We have two sorted readers in ascending order, we will start scanning from the beginning of the lists and compare whether the folder is a prefix of the current artifact,
// in case this is true the dir should not be deleted and we can move on to the next dir, otherwise we have to continue to the next dir or artifact.
// To know this, we will choose to move on with the lexicographic largest between the two.
//
// candidateDirsReaders - Sorted list of dirs to be deleted.
// filesNotToBeDeleteReader - Sorted files that should not be deleted.
// resultWriter - The filtered list of dirs to be deleted.
func WriteCandidateDirsToBeDeleted(candidateDirsReaders []*content.ContentReader, filesNotToBeDeleteReader *content.ContentReader, resultWriter *content.ContentWriter) (err error) {
	dirsToBeDeletedReader, err := content.MergeSortedReaders(ResultItem{}, candidateDirsReaders, true)
	if err != nil {
		return
	}
	defer func() {
		e := dirsToBeDeletedReader.Close()
		if err == nil {
			err = e
		}
	}()
	var candidateDirToBeDeletedPath string
	var itemNotToBeDeletedLocation string
	var candidateDirToBeDeleted, artifactNotToBeDeleted *ResultItem
	for {
		// Fetch the next 'candidateDirToBeDeleted'.
		if candidateDirToBeDeleted == nil {
			candidateDirToBeDeleted = new(ResultItem)
			if err = dirsToBeDeletedReader.NextRecord(candidateDirToBeDeleted); err != nil {
				break
			}
			if candidateDirToBeDeleted.Name == "." {
				continue
			}
			candidateDirToBeDeletedPath = candidateDirToBeDeleted.GetItemRelativePath()
		}
		// Fetch the next 'artifactNotToBeDelete'.
		if artifactNotToBeDeleted == nil {
			artifactNotToBeDeleted = new(ResultItem)
			if err = filesNotToBeDeleteReader.NextRecord(artifactNotToBeDeleted); err != nil {
				// No artifacts left, write remaining dirs to be deleted to result file.
				resultWriter.Write(*candidateDirToBeDeleted)
				writeRemainCandidate(resultWriter, dirsToBeDeletedReader)
				break
			}
			itemNotToBeDeletedLocation = artifactNotToBeDeleted.GetItemRelativeLocation()
		}
		// Found an 'artifact not to be deleted' in 'dir to be deleted', therefore skip writing the dir to the result file.
		if strings.HasPrefix(itemNotToBeDeletedLocation, candidateDirToBeDeletedPath) {
			candidateDirToBeDeleted = nil
			continue
		}
		// 'artifactNotToBeDeletePath' & 'candidateDirToBeDeletedPath' are both sorted. As a result 'candidateDirToBeDeleted' cant be a prefix for any of the remaining artifacts.
		if itemNotToBeDeletedLocation > candidateDirToBeDeletedPath {
			resultWriter.Write(*candidateDirToBeDeleted)
			candidateDirToBeDeleted = nil
			continue
		}
		artifactNotToBeDeleted = nil
	}
	err = filesNotToBeDeleteReader.GetError()
	filesNotToBeDeleteReader.Reset()
	return
}

func writeRemainCandidate(cw *content.ContentWriter, mergeResult *content.ContentReader) {
	for toBeDeleted := new(ResultItem); mergeResult.NextRecord(toBeDeleted) == nil; toBeDeleted = new(ResultItem) {
		cw.Write(*toBeDeleted)
	}
}

func FilterCandidateToBeDeleted(deleteCandidates *content.ContentReader, resultWriter *content.ContentWriter, candidateType string) ([]*content.ContentReader, error) {
	paths := make(map[string]content.SortableContentItem)
	pathsKeys := make([]string, 0, utils.MaxBufferSize)
	toBeDeleted := []*content.ContentReader{}
	for candidate := new(ResultItem); deleteCandidates.NextRecord(candidate) == nil; candidate = new(ResultItem) {
		// Save all candidates, of the requested type, to a different temp file.
		if candidate.Type == candidateType {
			if candidateType == "folder" && candidate.Name == "." {
				continue
			}
			pathsKeys = append(pathsKeys, candidate.GetItemRelativePath())
			paths[candidate.GetItemRelativePath()] = *candidate
			if len(pathsKeys) == utils.MaxBufferSize {
				sortedCandidateDirsFile, err := content.SortAndSaveBufferToFile(paths, pathsKeys, true)
				if err != nil {
					return nil, err
				}
				toBeDeleted = append(toBeDeleted, sortedCandidateDirsFile)
				// Init buffer.
				paths = make(map[string]content.SortableContentItem)
				pathsKeys = make([]string, 0, utils.MaxBufferSize)
			}
		} else {
			// Write none results of the requested type.
			resultWriter.Write(*candidate)
		}
	}
	if err := deleteCandidates.GetError(); err != nil {
		return nil, err
	}
	deleteCandidates.Reset()
	if len(pathsKeys) > 0 {
		sortedFile, err := content.SortAndSaveBufferToFile(paths, pathsKeys, true)
		if err != nil {
			return nil, err
		}
		toBeDeleted = append(toBeDeleted, sortedFile)
	}
	return toBeDeleted, nil
}
