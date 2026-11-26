package utils

import (
	"errors"

	"github.com/jfrog/jfrog-client-go/utils/io/content"
)

// RemoveDirsHoldingExcludedItems drops every folder that contains an item matched by params.Exclusions.
//
// A wildcard search with IncludeDirs returns folders next to files. The exclusions filter only the
// rows that match them, so a folder that holds an excluded file is still returned, and the later
// "top chain" reduce replaces all the files inside it with that folder. Artifactory then acts on the
// whole folder, excluded files included. When such a folder is removed here, its remaining files
// stay in the list as separate items and are handled one by one, without the excluded file.
//
// Cost: one extra AQL query (only when exclusions exist) plus a linear merge of two sorted streams,
// the same method that is used for --exclude-props in delete. Returns candidates unchanged when
// there is nothing to remove.
func RemoveDirsHoldingExcludedItems(params *CommonParams, flags CommonConf, candidates *content.ContentReader) (result *content.ContentReader, err error) {
	if len(params.GetExclusions()) == 0 {
		return candidates, nil
	}
	if length, lengthErr := candidates.Length(); lengthErr != nil || length == 0 {
		return candidates, lengthErr
	}
	excludedParams := excludedItemsParams(params)
	var excluded *content.ContentReader
	excluded, err = SearchBySpecWithPattern(&excludedParams, flags, NONE)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, excluded.Close())
	}()
	if length, lengthErr := excluded.Length(); lengthErr != nil || length == 0 {
		return candidates, lengthErr
	}
	return dropDirsHoldingItems(candidates, excluded)
}

// dropDirsHoldingItems returns candidates without the folders that contain at least one of the items.
func dropDirsHoldingItems(candidates, items *content.ContentReader) (result *content.ContentReader, err error) {
	sortedItems, err := content.SortContentReader(ResultItem{}, items, true)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, sortedItems.Close())
	}()

	resultWriter, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	defer func() {
		// On failure the half-written output file must not stay on disk.
		if err != nil {
			err = errors.Join(err, resultWriter.RemoveOutputFilePath())
		}
	}()
	sortedDirs, err := FilterCandidateToBeDeleted(candidates, resultWriter, Folder)
	defer func() {
		for _, reader := range sortedDirs {
			err = errors.Join(err, reader.Close())
		}
	}()
	if err != nil {
		return nil, err
	}
	if len(sortedDirs) > 0 {
		if err = WriteCandidateDirsToBeDeleted(sortedDirs, sortedItems, resultWriter); err != nil {
			return nil, err
		}
	}
	if err = resultWriter.Close(); err != nil {
		return nil, err
	}
	return content.NewContentReader(resultWriter.GetFilePath(), content.DefaultKey), nil
}

// excludedItemsParams returns a copy of params for the search that finds the excluded items.
// The caller's spec must not change.
func excludedItemsParams(params *CommonParams) CommonParams {
	excludedParams := *params
	excludedParams.ExcludedOnly = true
	excludedParams.IncludeDirs = true
	// A folder returned by a non-recursive search still holds its whole content, so excluded
	// items below the first level must be found too.
	excludedParams.Recursive = true
	excludedParams.Props, excludedParams.ExcludeProps, excludedParams.Bundle, excludedParams.Build = "", "", "", ""
	// Every excluded item must be found, so no paging or sorting may cut the answer short.
	excludedParams.Limit, excludedParams.Offset, excludedParams.SortBy, excludedParams.SortOrder = 0, 0, nil, ""
	return excludedParams
}
