package utils

import (
	"strings"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"

	buildinfo "github.com/jfrog/build-info-go/entities"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
)

type Result struct {
	SuccessCount []int
	TotalCount   []int
}

func NewResult(threads int) *Result {
	return &Result{SuccessCount: make([]int, threads),
		TotalCount: make([]int, threads)}
}

func SumIntArray(arr []int) int {
	sum := 0
	for _, i := range arr {
		sum += i
	}
	return sum
}

type OperationSummary struct {
	// A ContentReader of FileTransferDetails structs
	TransferDetailsReader *content.ContentReader
	// A ContentReader of ArtifactDetails structs
	ArtifactsDetailsReader *content.ContentReader
	TotalSucceeded         int
	TotalFailed            int
}

type ArtifactDetails struct {
	// Path of the artifact in Artifactory
	ArtifactoryPath string             `json:"artifactoryPath,omitempty"`
	Checksums       buildinfo.Checksum `json:"checksums,omitempty"`
}

func (cs *OperationSummary) Close() error {
	err := cs.TransferDetailsReader.Close()
	if err != nil {
		return err
	}
	return cs.ArtifactsDetailsReader.Close()
}

func (ad *ArtifactDetails) ToBuildInfoArtifact() (buildinfo.Artifact, error) {
	artifact := buildinfo.Artifact{Checksum: buildinfo.Checksum{}}
	artifact.Sha1 = ad.Checksums.Sha1
	artifact.Md5 = ad.Checksums.Md5
	artifact.Sha256 = ad.Checksums.Sha256
	// Artifact name in build info as the name in artifactory
	filename, _ := fileutils.GetFileAndDirFromPath(ad.ArtifactoryPath)
	artifact.Name = filename
	if i := strings.LastIndex(filename, "."); i != -1 {
		artifact.Type = filename[i+1:]
	}
	// The 'path' property in the build-info should not include the repository. We therefore remove the repository from the path.
	if i := strings.Index(ad.ArtifactoryPath, "/"); i != -1 {
		artifact.Path = ad.ArtifactoryPath[i+1:]
	} else {
		return artifact, errorutils.CheckErrorf("artifact path:' " + ad.ArtifactoryPath + "' lacks repository name")
	}
	return artifact, nil
}

func (ad *ArtifactDetails) ToBuildInfoDependency() buildinfo.Dependency {
	dependency := buildinfo.Dependency{Checksum: buildinfo.Checksum{}}
	dependency.Sha1 = ad.Checksums.Sha1
	dependency.Md5 = ad.Checksums.Md5
	// Artifact name in build info as the name in artifactory
	filename, _ := fileutils.GetFileAndDirFromPath(ad.ArtifactoryPath)
	dependency.Id = filename
	return dependency
}

func ConvertArtifactsDetailsToBuildInfoArtifacts(artifactsDetailsReader *content.ContentReader) ([]buildinfo.Artifact, error) {
	var buildArtifacts []buildinfo.Artifact
	for artifactDetails := new(ArtifactDetails); artifactsDetailsReader.NextRecord(artifactDetails) == nil; artifactDetails = new(ArtifactDetails) {
		artifact, err := artifactDetails.ToBuildInfoArtifact()
		if err != nil {
			return nil, err
		}
		buildArtifacts = append(buildArtifacts, artifact)
	}
	return buildArtifacts, artifactsDetailsReader.GetError()
}

func ConvertArtifactsDetailsToBuildInfoDependencies(artifactsDetailsReader *content.ContentReader) ([]buildinfo.Dependency, error) {
	var buildDependencies []buildinfo.Dependency
	for artifactDetails := new(ArtifactDetails); artifactsDetailsReader.NextRecord(artifactDetails) == nil; artifactDetails = new(ArtifactDetails) {
		buildDependencies = append(buildDependencies, artifactDetails.ToBuildInfoDependency())
	}
	return buildDependencies, artifactsDetailsReader.GetError()
}
