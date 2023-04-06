package fileutils

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jfrog/gofrog/datastructures"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/mholt/archiver/v3"
)

func IsSupportedArchive(filePath string) bool {
	iArchiver, err := archiver.ByExtension(filePath)
	if err != nil {
		return false
	}
	_, ok := iArchiver.(archiver.Unarchiver)
	return ok
}

// The 'archiver' dependency includes an API called 'Unarchive' to extract archive files. This API uses the archive file
// extension to determine the archive type.
// We therefore need to use the file name as it was in Artifactory, and not the file name which was downloaded. To achieve this,
// we added a new implementation of the 'Unarchive' func and use it instead of the default one.
// localArchivePath - The local file path to extract the archive
// originArchiveName - The archive file name
// destinationPath - The extraction destination directory
func Unarchive(localArchivePath, originArchiveName, destinationPath string) error {
	archive, err := byExtension(originArchiveName)
	if err != nil {
		return err
	}
	u, ok := archive.(archiver.Unarchiver)
	if !ok {
		return errorutils.CheckErrorf("format specified by source filename is not an archive format: " + originArchiveName)
	}
	if err = inspectArchive(archive, localArchivePath, destinationPath); err != nil {
		return err
	}
	return u.Unarchive(localArchivePath, destinationPath)
}

// Instead of using 'archiver.byExtension' that by default sets OverwriteExisting to false, we implement our own.
func byExtension(filename string) (interface{}, error) {
	var ec interface{}
	for _, c := range extCheckers {
		if err := c.CheckExt(filename); err == nil {
			ec = c
			break
		}
	}
	switch ec.(type) {
	case *archiver.Rar:
		archiveInstance := archiver.NewRar()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.Tar:
		archiveInstance := archiver.NewTar()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarBrotli:
		archiveInstance := archiver.NewTarBrotli()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarBz2:
		archiveInstance := archiver.NewTarBz2()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarGz:
		archiveInstance := archiver.NewTarGz()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarLz4:
		archiveInstance := archiver.NewTarLz4()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarSz:
		archiveInstance := archiver.NewTarSz()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarXz:
		archiveInstance := archiver.NewTarXz()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.TarZstd:
		archiveInstance := archiver.NewTarZstd()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.Zip:
		archiveInstance := archiver.NewZip()
		archiveInstance.OverwriteExisting = true
		return archiveInstance, nil
	case *archiver.Gz:
		return archiver.NewGz(), nil
	case *archiver.Bz2:
		return archiver.NewBz2(), nil
	case *archiver.Lz4:
		return archiver.NewLz4(), nil
	case *archiver.Snappy:
		return archiver.NewSnappy(), nil
	case *archiver.Xz:
		return archiver.NewXz(), nil
	case *archiver.Zstd:
		return archiver.NewZstd(), nil
	}
	return nil, errorutils.CheckErrorf("format unrecognized by filename: %s", filename)
}

var extCheckers = []archiver.ExtensionChecker{
	&archiver.TarBrotli{},
	&archiver.TarBz2{},
	&archiver.TarGz{},
	&archiver.TarLz4{},
	&archiver.TarSz{},
	&archiver.TarXz{},
	&archiver.TarZstd{},
	&archiver.Rar{},
	&archiver.Tar{},
	&archiver.Zip{},
	&archiver.Brotli{},
	&archiver.Gz{},
	&archiver.Bz2{},
	&archiver.Lz4{},
	&archiver.Snappy{},
	&archiver.Xz{},
	&archiver.Zstd{},
}

// Make sure the archive is free from Zip Slip and Zip symlinks attacks
func inspectArchive(archive interface{}, localArchivePath, destinationDir string) error {
	walker, ok := archive.(archiver.Walker)
	if !ok {
		return errorutils.CheckErrorf("couldn't inspect archive: " + localArchivePath)
	}

	uplinksValidator := newUplinksValidator()
	err := walker.Walk(localArchivePath, func(archiveEntry archiver.File) error {
		header, err := extractArchiveEntryHeader(archiveEntry)
		if err != nil {
			return err
		}
		pathInArchive := getPathInArchive(destinationDir, "", header.EntryPath)
		if !strings.HasPrefix(pathInArchive, destinationDir) {
			return errorutils.CheckErrorf(
				"illegal path in archive: '%s'. To prevent Zip Slip exploit, the path can't lead to an entry outside '%s'",
				header.EntryPath, destinationDir)
		}
		if (archiveEntry.Mode()&os.ModeSymlink) != 0 || len(header.TargetLink) > 0 {
			var targetLink string
			if targetLink, err = checkSymlinkEntry(header, archiveEntry, destinationDir); err != nil {
				return err
			}
			uplinksValidator.addTargetLink(pathInArchive, targetLink)
		}
		uplinksValidator.addEntryFile(pathInArchive, archiveEntry.IsDir())
		return err
	})
	if err != nil {
		return err
	}
	return uplinksValidator.ensureNoUplinkDirs()
}

// Make sure the extraction path of the symlink entry target is under the destination dir
func checkSymlinkEntry(header *archiveHeader, archiveEntry archiver.File, destinationDir string) (string, error) {
	targetLinkPath := header.TargetLink
	if targetLinkPath == "" {
		// The link destination path is not always in the archive header
		// In that case, we will look at the link content to get the link destination path
		content, err := io.ReadAll(archiveEntry.ReadCloser)
		if err != nil {
			return "", errorutils.CheckError(err)
		}
		targetLinkPath = string(content)
	}

	targetPathInArchive := getPathInArchive(destinationDir, filepath.Dir(header.EntryPath), targetLinkPath)
	if !strings.HasPrefix(targetPathInArchive, destinationDir) {
		return "", errorutils.CheckErrorf(
			"illegal link path in archive: '%s'. To prevent Zip Slip Symlink exploit, the path can't lead to an entry outside '%s'",
			targetLinkPath, destinationDir)
	}

	return targetPathInArchive, nil
}

// Get the path in archive of the entry or the target link
func getPathInArchive(destinationDir, entryDirInArchive, pathInArchive string) string {
	// If pathInArchive starts with '/' and we are on Windows, the path is illegal
	pathInArchive = strings.TrimSpace(pathInArchive)
	if os.IsPathSeparator('\\') && strings.HasPrefix(pathInArchive, "/") {
		return ""
	}

	pathInArchive = filepath.Clean(pathInArchive)
	if !filepath.IsAbs(pathInArchive) {
		// If path is relative, concatenate it to the destination dir
		pathInArchive = filepath.Join(destinationDir, entryDirInArchive, pathInArchive)
	}
	return pathInArchive
}

// Extract the header of the archive entry
func extractArchiveEntryHeader(f archiver.File) (*archiveHeader, error) {
	headerBytes, err := json.Marshal(f.Header)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	archiveHeader := &archiveHeader{}
	err = json.Unmarshal(headerBytes, archiveHeader)
	return archiveHeader, errorutils.CheckError(err)
}

type archiveHeader struct {
	EntryPath  string `json:"Name,omitempty"`
	TargetLink string `json:"Linkname,omitempty"`
}

// This validator blocks the option to extract an archive with a link to an ancestor directory.
// An ancestor directory is a directory located above the symlink in the hierarchy of the extraction dir, but not necessarily a direct ancestor.
// For example, a sibling of a parent is an ancestor directory.
// The purpose of the uplinksValidator is to prevent directories loop in the file system during extraction.
type uplinksValidator struct {
	entryFiles        *datastructures.Set[string]
	targetParentLinks map[string]string
}

func newUplinksValidator() *uplinksValidator {
	return &uplinksValidator{
		// Set of all entries that are not directories in the archive
		entryFiles: datastructures.MakeSet[string](),
		// Map of all links in the archive pointing to an ancestor entry
		targetParentLinks: make(map[string]string),
	}
}

func (lv *uplinksValidator) addTargetLink(pathInArchive, targetLink string) {
	if strings.Count(targetLink, string(filepath.Separator)) < strings.Count(pathInArchive, string(filepath.Separator)) {
		// Add the target link only if it is an ancestor
		lv.targetParentLinks[pathInArchive] = targetLink
	}
}

func (lv *uplinksValidator) addEntryFile(entryFile string, isDir bool) {
	if !isDir {
		// Add the entry only if it is not a directory
		lv.entryFiles.Add(entryFile)
	}
}

// Iterate over all links pointing to an ancestor directories and files.
// If a targetParentLink does not exist in the entryFiles list, it is a directory and therefore return an error.
func (lv *uplinksValidator) ensureNoUplinkDirs() error {
	for pathInArchive, targetLink := range lv.targetParentLinks {
		if lv.entryFiles.Exists(targetLink) {
			// Target link to a file
			continue
		}
		// Target link to a directory
		return errorutils.CheckErrorf(
			"illegal target link path in archive: '%s' -> '%s'. To prevent Zip Slip symlink exploit, a link can't lead to an ancestor directory",
			pathInArchive, targetLink)
	}
	return nil
}
