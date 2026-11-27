package utils

import (
	"errors"
	"io/fs"
	"strings"

	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// RemoveDirsHoldingExcludedItems drops every folder that contains an item matched by params.Exclusions.
//
// A wildcard search with IncludeDirs returns folders next to files. The exclusions filter only the
// rows that match them, so a folder that holds an excluded file is still returned, and the later
// "top chain" reduce replaces all the files inside it with that folder. Artifactory then acts on the
// whole folder, excluded files included. This function removes such a folder.
// The other files of the folder could be lost, if the search did not return them as separate items.
// This is the case for a non-recursive search, or for a name pattern such as "repo/*-release".
// For these folders the content is searched again, with the exclusions, and the files are added.
//
// Cost: one extra AQL query (only when exclusions exist), plus one more for each removed folder.
// Returns candidates unchanged when there is nothing to remove.
func RemoveDirsHoldingExcludedItems(params *CommonParams, flags CommonConf, candidates *content.ContentReader) (result *content.ContentReader, err error) {
	if len(params.GetExclusions()) == 0 {
		return candidates, nil
	}
	if length, lengthErr := candidates.Length(); lengthErr != nil || length == 0 {
		return candidates, lengthErr
	}
	excludedParams := excludedItemsParams(params)
	excluded, err := SearchBySpecWithPattern(&excludedParams, flags, NONE)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, excluded.Close())
		if err != nil && result != nil && result != candidates {
			// The caller does not close a result that comes with an error.
			err = errors.Join(err, result.Close())
			result = nil
		}
	}()
	if length, lengthErr := excluded.Length(); lengthErr != nil || length == 0 {
		return candidates, lengthErr
	}
	kept, dirsToExpand, err := dropDirsHoldingItems(candidates, excluded)
	if err != nil || len(dirsToExpand) == 0 {
		return kept, err
	}
	if clientutils.RemovePlaceholderParentheses(params.Pattern, params.Target) != params.Pattern {
		// The destination of a file is built from the placeholders of the pattern, so a file that the pattern did not return cannot be placed.
		log.Warn("These folders hold excluded items. Their other files are not handled, because the pattern uses placeholders:", strings.Join(dirsToExpand, ", "))
		return kept, nil
	}
	if params.Offset > 0 || (params.Limit > 0 && len(params.SortBy) > 0) {
		// The first search cut its result by this offset, limit or order. The files of a removed folder cannot be placed in
		// that result exactly: they could be handled twice, more than the limit could be handled, or the order could be lost.
		// Not handling them is safer than handling too many.
		log.Warn("These folders hold excluded items. Their other files are not handled, because the spec uses an offset, or a limit with a sort order:", strings.Join(dirsToExpand, ", "))
		return kept, nil
	}
	result, err = addFilesOfDirs(params, flags, kept, dirsToExpand)
	if err = errors.Join(err, kept.Close()); err != nil {
		if result != nil {
			err = errors.Join(err, result.Close())
		}
		return nil, err
	}
	if params.Limit > 0 {
		// The searched files are not cut by the limit, so cut the final list once. The command never handles more than the limit.
		// Here the spec has no offset and no sort order, so the order by path is the order of the result.
		return limitItems(result, params.Limit)
	}
	return result, nil
}

// limitItems returns the first limit items of the reader. It closes the reader.
func limitItems(items *content.ContentReader, limit int) (result *content.ContentReader, err error) {
	defer func() {
		err = errors.Join(err, items.Close())
		if err != nil && result != nil {
			err = errors.Join(err, result.Close())
			result = nil
		}
	}()
	return copyItems(items, limit)
}

// copyItems writes the first max items of the reader to a new reader.
func copyItems(items *content.ContentReader, max int) (result *content.ContentReader, err error) {
	writer, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	closed := false
	defer func() {
		if err != nil {
			err = errors.Join(err, discardWriter(writer, closed))
		}
	}()
	count := 0
	for item := new(ResultItem); items.NextRecord(item) == nil; item = new(ResultItem) {
		// Read the reader to its end, so that its reading goroutine can finish.
		if count < max {
			writer.Write(*item)
			count++
		}
	}
	if err = items.GetError(); err != nil {
		return nil, err
	}
	closed = true
	if err = writer.Close(); err != nil {
		return nil, err
	}
	if writer.GetFilePath() == "" {
		return content.NewEmptyContentReader(content.DefaultKey), nil
	}
	return content.NewContentReader(writer.GetFilePath(), content.DefaultKey), nil
}

// dropDirsHoldingItems returns candidates without the folders that contain at least one of the items,
// and without the folders that are items themselves.
// It also returns the full path (repo/path/name) of the removed folders that must be searched again.
// A removed folder that is inside another removed folder is not returned: the search of that folder includes it.
func dropDirsHoldingItems(candidates, items *content.ContentReader) (result *content.ContentReader, dirsToExpand []string, err error) {
	holders, itemDirs, err := collectItemDirs(items)
	if err != nil {
		return nil, nil, err
	}
	resultWriter, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, nil, err
	}
	closed := false
	defer func() {
		if err != nil {
			err = errors.Join(err, discardWriter(resultWriter, closed))
		}
	}()
	droppedDirs := make(map[string]bool)
	for candidate := new(ResultItem); candidates.NextRecord(candidate) == nil; candidate = new(ResultItem) {
		if candidate.Type != string(Folder) || candidate.Name == "." {
			resultWriter.Write(*candidate)
			continue
		}
		dirPath := fullPath(candidate)
		switch {
		case itemDirs[dirPath]:
			droppedDirs[dirPath] = false
		case holders[dirPath]:
			droppedDirs[dirPath] = true
		default:
			resultWriter.Write(*candidate)
		}
	}
	if err = candidates.GetError(); err != nil {
		return nil, nil, err
	}
	candidates.Reset()
	if len(droppedDirs) == 0 {
		// Nothing is removed, so the copy is not needed.
		closed = true
		return candidates, nil, discardWriter(resultWriter, false)
	}
	closed = true
	if err = resultWriter.Close(); err != nil {
		return nil, nil, err
	}
	for dirPath, holdsItems := range droppedDirs {
		if holdsItems && !hasAncestorIn(dirPath, droppedDirs) {
			dirsToExpand = append(dirsToExpand, dirPath)
		}
	}
	return content.NewContentReader(resultWriter.GetFilePath(), content.DefaultKey), dirsToExpand, nil
}

// collectItemDirs returns two sets of folder paths:
// holders - every folder that contains at least one of the items, at any depth.
// itemDirs - the items that are folders.
func collectItemDirs(items *content.ContentReader) (holders, itemDirs map[string]bool, err error) {
	holders, itemDirs = make(map[string]bool), make(map[string]bool)
	for item := new(ResultItem); items.NextRecord(item) == nil; item = new(ResultItem) {
		itemPath := fullPath(item)
		if item.Type == string(Folder) {
			itemDirs[itemPath] = true
		}
		for dirPath := parentPath(itemPath); strings.Contains(dirPath, "/") && !holders[dirPath]; dirPath = parentPath(dirPath) {
			holders[dirPath] = true
		}
	}
	items.Reset()
	return holders, itemDirs, items.GetError()
}

// addFilesOfDirs returns kept plus all the files in the folders, found with the exclusions of params.
// A file that is found more than once is returned once.
func addFilesOfDirs(params *CommonParams, flags CommonConf, kept *content.ContentReader, dirs []string) (result *content.ContentReader, err error) {
	log.Debug("Searching the content of the folders that hold excluded items:", strings.Join(dirs, ", "))
	combinedWriter, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	closed := false
	defer func() {
		if err != nil {
			err = errors.Join(err, discardWriter(combinedWriter, closed))
		}
	}()
	for item := new(ResultItem); kept.NextRecord(item) == nil; item = new(ResultItem) {
		combinedWriter.Write(*item)
	}
	if err = kept.GetError(); err != nil {
		return nil, err
	}
	for _, dir := range dirs {
		if err = writeFilesOfDir(params, flags, dir, combinedWriter); err != nil {
			return nil, err
		}
	}
	closed = true
	if err = combinedWriter.Close(); err != nil {
		return nil, err
	}
	if combinedWriter.GetFilePath() == "" {
		return content.NewEmptyContentReader(content.DefaultKey), nil
	}
	combined := content.NewContentReader(combinedWriter.GetFilePath(), content.DefaultKey)
	defer func() {
		if err = errors.Join(err, combined.Close()); err != nil && result != nil {
			err = errors.Join(err, result.Close())
			result = nil
		}
	}()
	return sortWithoutDuplicates(combined)
}

// writeFilesOfDir writes the files below the folder, without the excluded ones, to the writer.
func writeFilesOfDir(params *CommonParams, flags CommonConf, dir string, writer *content.ContentWriter) (err error) {
	if strings.ContainsAny(dir, "\"\\*?") {
		// The name cannot be used safely in a search pattern. Not handling the files is safer than handling other files.
		log.Warn("The files of the folder are not handled, because its name has one of the characters \" \\ * ?:", dir)
		return nil
	}
	filesParams := *params
	filesParams.Pattern = dir + "/*"
	filesParams.Recursive = true
	filesParams.IncludeDirs = false
	// The files of a removed folder must not be cut by paging.
	filesParams.Limit, filesParams.Offset, filesParams.SortBy, filesParams.SortOrder = 0, 0, nil, ""
	files, err := SearchBySpecWithPattern(&filesParams, flags, NONE)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, files.Close())
	}()
	for file := new(ResultItem); files.NextRecord(file) == nil; file = new(ResultItem) {
		// A wildcard character in the folder name could make the pattern match more than the folder. Never act on such an item.
		if strings.HasPrefix(fullPath(file), dir+"/") {
			writer.Write(*file)
		}
	}
	return files.GetError()
}

// sortWithoutDuplicates returns the items of the reader sorted, with each item once.
func sortWithoutDuplicates(items *content.ContentReader) (result *content.ContentReader, err error) {
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
	closed := false
	defer func() {
		if err != nil {
			err = errors.Join(err, discardWriter(resultWriter, closed))
		}
	}()
	previousKey := ""
	for item := new(ResultItem); sortedItems.NextRecord(item) == nil; item = new(ResultItem) {
		if key := item.GetSortKey(); key != previousKey {
			resultWriter.Write(*item)
			previousKey = key
		}
	}
	if err = sortedItems.GetError(); err != nil {
		return nil, err
	}
	closed = true
	if err = resultWriter.Close(); err != nil {
		return nil, err
	}
	if resultWriter.GetFilePath() == "" {
		return content.NewEmptyContentReader(content.DefaultKey), nil
	}
	return content.NewContentReader(resultWriter.GetFilePath(), content.DefaultKey), nil
}

// discardWriter stops the writer and removes its output file, if it has one.
// Close is not safe to call twice, so closed tells if the caller already called it.
func discardWriter(writer *content.ContentWriter, closed bool) (err error) {
	if !closed {
		err = writer.Close()
	}
	// The output file is created with the first write.
	if writer.GetFilePath() != "" {
		// The file can be gone already, because a reader on it was closed first.
		if removeErr := writer.RemoveOutputFilePath(); !errors.Is(removeErr, fs.ErrNotExist) {
			err = errors.Join(err, removeErr)
		}
	}
	return err
}

// fullPath returns repo/path/name of an item.
func fullPath(item *ResultItem) string {
	return strings.TrimSuffix(item.GetItemRelativePath(), "/")
}

// parentPath returns the path without its last element.
func parentPath(path string) string {
	return path[:max(strings.LastIndex(path, "/"), 0)]
}

// hasAncestorIn returns true when one of the folders above dirPath is in dirs.
func hasAncestorIn(dirPath string, dirs map[string]bool) bool {
	for ancestor := parentPath(dirPath); strings.Contains(ancestor, "/"); ancestor = parentPath(ancestor) {
		if _, found := dirs[ancestor]; found {
			return true
		}
	}
	return false
}

// excludedItemsParams returns a copy of params for the search that finds the excluded items.
// The caller's spec must not change.
func excludedItemsParams(params *CommonParams) CommonParams {
	excludedParams := *params
	excludedParams.ExcludedOnly = true
	excludedParams.IncludeDirs = true
	// The name part of the source pattern (for example "*-release") must not hide an excluded item
	// that is deeper in the folder, so the lookup covers the whole repository. The exclusions narrow it.
	excludedParams.Pattern = repositoryPattern(params.Pattern)
	// A folder returned by a non-recursive search still holds its whole content, so excluded
	// items below the first level must be found too.
	excludedParams.Recursive = true
	excludedParams.Props, excludedParams.ExcludeProps, excludedParams.Bundle, excludedParams.Build = "", "", "", ""
	excludedParams.ArchiveEntries = ""
	// Every excluded item must be found, so no paging or sorting may cut the answer short.
	excludedParams.Limit, excludedParams.Offset, excludedParams.SortBy, excludedParams.SortOrder = 0, 0, nil, ""
	return excludedParams
}

// repositoryPattern returns a pattern that matches everything in the repository part of the source pattern.
func repositoryPattern(pattern string) string {
	// Parentheses mark a placeholder. They are not part of the repository name.
	repository := strings.NewReplacer("(", "", ")", "").Replace(strings.SplitN(strings.TrimPrefix(pattern, "/"), "/", 2)[0])
	if repository == "" {
		repository = "*"
	}
	return repository + "/*"
}
