package _go

import (
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/version"
)

func init() {
	register(&publishZipAndModApi{})
}

const ArtifactoryMinSupportedVersionForInfoFile = "6.10.0"

// Support for Artifactory 6.6.1 and above API
type publishZipAndModApi struct {
	artifactoryVersion string
	clientDetails      httputils.HttpClientDetails
	client             *jfroghttpclient.JfrogHttpClient
}

func (pwa *publishZipAndModApi) isCompatible(artifactoryVersion string) bool {
	propertiesApi := "6.6.1"
	version := version.NewVersion(artifactoryVersion)
	pwa.artifactoryVersion = artifactoryVersion
	return version.AtLeast(propertiesApi)
}

func (pwa *publishZipAndModApi) PublishPackage(params GoParams, client *jfroghttpclient.JfrogHttpClient, ArtDetails auth.ServiceDetails) (*utils.OperationSummary, error) {
	url, err := utils.BuildArtifactoryUrl(ArtDetails.GetUrl(), "api/go/"+params.GetTargetRepo(), make(map[string]string))
	if err != nil {
		return nil, err
	}
	pwa.clientDetails = ArtDetails.CreateHttpClientDetails()
	pwa.client = client
	moduleId := strings.Split(params.GetModuleId(), ":")
	totalSucceed, totalFailed := 0, 0
	var filesDetails []clientutils.FileTransferDetails
	// Upload zip file
	success, failed, err := uploadArchive(params, params.ZipPath, moduleId[0], ".zip", url, &filesDetails, pwa)
	if err != nil {
		return nil, err
	}
	totalSucceed, totalFailed = totalSucceed+success, totalFailed+failed
	// Upload mod file
	success, failed, err = uploadArchive(params, params.ModPath, moduleId[0], ".mod", url, &filesDetails, pwa)
	if err != nil {
		return nil, err
	}
	totalSucceed, totalFailed = totalSucceed+success, totalFailed+failed
	if version.NewVersion(pwa.artifactoryVersion).AtLeast(ArtifactoryMinSupportedVersionForInfoFile) && params.GetInfoPath() != "" {
		// Upload info file. This supported from Artifactory version 6.10.0 and above
		success, failed, err = uploadArchive(params, params.InfoPath, moduleId[0], ".info", url, &filesDetails, pwa)
		totalSucceed, totalFailed = totalSucceed+success, totalFailed+failed
		if err != nil {
			return nil, err
		}
	}
	tempFile, err := clientutils.SaveFileTransferDetailsInTempFile(&filesDetails)
	if err != nil {
		return nil, err
	}
	return &utils.OperationSummary{TotalSucceeded: totalSucceed, TotalFailed: totalFailed, TransferDetailsReader: content.NewContentReader(tempFile, "files")}, nil
}

func uploadArchive(params GoParams, archivePath string, moduleId, ext, url string, filesDetails *[]clientutils.FileTransferDetails, pwa *publishZipAndModApi) (success, failed int, err error) {
	success, failed = 0, 1
	details, err := pwa.upload(archivePath, moduleId, params.GetVersion(), params.GetProps(), ext, url)
	if err != nil {
		return
	}
	if details != nil {
		success, failed = 1, 0
		*filesDetails = append(*filesDetails, *details)
	}
	return
}

func addGoVersion(version string, urlPath *string) {
	*urlPath += ";go.version=" + url.QueryEscape(version)
}

// localPath - The location of the file on the file system.
// moduleId - The name of the module for example github.com/jfrog/jfrog-client-go.
// version - The version of the project that being uploaded.
// props - The properties to be assigned for each artifact
// ext - The extension of the file: zip, mod, info. This extension will be joined with the version for the path. For example v1.2.3.info or v1.2.3.zip
// urlPath - The url including the repository. For example: http://127.0.0.1/artifactory/api/go/go-local
func (pwa *publishZipAndModApi) upload(localPath, moduleId, version, props, ext, urlPath string) (*clientutils.FileTransferDetails, error) {
	err := CreateUrlPath(moduleId, version, props, ext, &urlPath)
	if err != nil {
		return nil, err
	}
	addGoVersion(version, &urlPath)
	details, err := fileutils.GetFileDetails(localPath)
	if err != nil {
		return nil, err
	}
	utils.AddChecksumHeaders(pwa.clientDetails.Headers, details)
	resp, _, err := pwa.client.UploadFile(localPath, urlPath, "", &pwa.clientDetails, GoUploadRetries, nil)
	if err != nil {
		return nil, err
	}
	sha256 := resp.Header.Get("X-Checksum-Sha256")
	if err != nil {
		log.Error("Failed to extract file's sha256 from response body.\nFile: " + localPath + "\nError message:" + err.Error())
	}
	// Remove urls properties suffix
	splitUrlPath := strings.Split(urlPath, ";")
	filesDetails := clientutils.FileTransferDetails{SourcePath: localPath, TargetPath: splitUrlPath[0], Sha256: sha256}
	return &filesDetails, errorutils.CheckResponseStatus(resp, http.StatusCreated)
}
