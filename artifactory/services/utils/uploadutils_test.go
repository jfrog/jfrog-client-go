package utils

import (
	"fmt"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"testing"
)

func TestGetArtifactToUpload(t *testing.T) {
	var debianPaths = []struct {
		rootFile        string
		target          string
		symlinkTarget   string
		isFlat          bool
		preserveSymlink bool
		expected        clientutils.Artifact
	}{
		{"/Location/on/fs/file1.in", "repo1/", "", true, false, clientutils.Artifact{"/Location/on/fs/file1.in", "repo1/file1.in", ""}},
		{"/Location/on/fs/file2.in", "repo1/", "", true, true, clientutils.Artifact{"/Location/on/fs/file2.in", "repo1/file2.in", ""}},
		{"Location/on/fs/file3.in", "repo1/", "", false, false, clientutils.Artifact{"Location/on/fs/file3.in", "repo1/Location/on/fs/file3.in", ""}},
		{"Location/on/fs/file12.in", "repo1/", "", false, true, clientutils.Artifact{"Location/on/fs/file12.in", "repo1/Location/on/fs/file12.in", ""}},
		{"folder1/on/fs/file4.in", "repo1/", "folder2/on/fs/file5", true, false, clientutils.Artifact{"folder2/on/fs/file5", "repo1/file4.in", "folder2/on/fs/file5"}},
		{"folder1/on/fs/file6.in", "repo1/", "folder2/on/fs/file7", true, true, clientutils.Artifact{"folder1/on/fs/file6.in", "repo1/file6.in", "folder2/on/fs/file7"}},
		{"folder1/on/fs/file8.in", "repo1/", "folder2/on/fs/file9", false, true, clientutils.Artifact{"folder1/on/fs/file8.in", "repo1/folder1/on/fs/file8.in", "folder2/on/fs/file9"}},
		{"folder1/on/fs/file10.in", "repo1/", "folder2/on/fs/file11", false, false, clientutils.Artifact{"folder2/on/fs/file11", "repo1/folder1/on/fs/file10.in", "folder2/on/fs/file11"}},
		{"/Location/on/fs/file13.in", "repo1/folder2/file14.in", "", true, false, clientutils.Artifact{"/Location/on/fs/file13.in", "repo1/folder2/file14.in", ""}},
		{"folder1/on/fs/file15.in", "repo1/folder-in-repo/file16.in", "folder2/on/fs/file17", true, false, clientutils.Artifact{"folder2/on/fs/file17", "repo1/folder-in-repo/file16.in", "folder2/on/fs/file17"}},
		{"folder1/on/fs/file18.in", "repo1/folder-in-repo/file19.in", "folder2/on/fs/file20", true, true, clientutils.Artifact{"folder1/on/fs/file18.in", "repo1/folder-in-repo/file19.in", "folder2/on/fs/file20"}},
		{"folder1/on/fs/file18.in", "repo1/folder-in-repo/file19.in", "folder2/on/fs/file20", false, true, clientutils.Artifact{"folder1/on/fs/file18.in", "repo1/folder-in-repo/file19.in", "folder2/on/fs/file20"}},
	}

	for _, v := range debianPaths {
		result := GetArtifactToUpload(v.rootFile, v.target, v.symlinkTarget, v.isFlat, v.preserveSymlink)
		if result != v.expected {
			t.Errorf("Expected:\n%s, received:\n%s", artifactToString(v.expected), artifactToString(result))
		}
	}
}

func artifactToString(artifact clientutils.Artifact) string {
	return fmt.Sprintf("{LocalPath: \"%s\", TargetPath: \"%s\", Symlink: \"%s\"}", artifact.LocalPath, artifact.TargetPath, artifact.Symlink)
}
