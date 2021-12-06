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

const ArtifactoryMinSupportedVersion = "6.10.0"

// Support for Artifactory 6.10.0 and above API
type GoPublishCommand struct {
	artifactoryVersion string
	clientDetails      httputils.HttpClientDetails
	client             *jfroghttpclient.JfrogHttpClient
}

func (gpc *GoPublishCommand) verifyCompatibleVersion(artifactoryVersion string) error {
	propertiesApi := ArtifactoryMinSupportedVersion
	ver := version.NewVersion(artifactoryVersion)
	gpc.artifactoryVersion = artifactoryVersion
	if !ver.AtLeast(propertiesApi) {
		return errorutils.CheckErrorf("Unsupported version of Artifactory: %s\nSupports Artifactory version %s and above", artifactoryVersion, propertiesApi)
	}
	return nil
}

func (gpc *GoPublishCommand) PublishPackage(params GoParams, client *jfroghttpclient.JfrogHttpClient, ArtDetails auth.ServiceDetails) (*utils.OperationSummary, error) {
	rtUrl, err := utils.BuildArtifactoryUrl(ArtDetails.GetUrl(), "api/go/"+params.GetTargetRepo(), make(map[string]string))
	if err != nil {
		return nil, err
	}
	gpc.clientDetails = ArtDetails.CreateHttpClientDetails()
	gpc.client = client
	moduleId := strings.Split(params.GetModuleId(), ":")
	totalSucceed, totalFailed := 0, 0
	var filesDetails []clientutils.FileTransferDetails
	// Upload zip file
	success, failed, err := gpc.uploadFile(params, params.ZipPath, moduleId[0], ".zip", rtUrl, &filesDetails, gpc)
	if err != nil {
		return nil, err
	}
	totalSucceed, totalFailed = totalSucceed+success, totalFailed+failed
	// Upload mod file
	success, failed, err = gpc.uploadFile(params, params.ModPath, moduleId[0], ".mod", rtUrl, &filesDetails, gpc)
	if err != nil {
		return nil, err
	}
	totalSucceed, totalFailed = totalSucceed+success, totalFailed+failed
	if version.NewVersion(gpc.artifactoryVersion).AtLeast(ArtifactoryMinSupportedVersion) && params.GetInfoPath() != "" {
		// Upload info file. This is supported from Artifactory version 6.10.0 and above
		success, failed, err = gpc.uploadFile(params, params.InfoPath, moduleId[0], ".info", rtUrl, &filesDetails, gpc)
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

func (gpc *GoPublishCommand) uploadFile(params GoParams, filePath string, moduleId, ext, url string, filesDetails *[]clientutils.FileTransferDetails, pwa *GoPublishCommand) (success, failed int, err error) {
	success, failed = 0, 1
	details, err := pwa.upload(filePath, moduleId, params.GetVersion(), params.GetProps(), ext, url)
	if err != nil {
		return
	}
	success, failed = 1, 0
	*filesDetails = append(*filesDetails, *details)
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
func (gpc *GoPublishCommand) upload(localPath, moduleId, version, props, ext, urlPath string) (*clientutils.FileTransferDetails, error) {
	err := CreateUrlPath(moduleId, version, props, ext, &urlPath)
	if err != nil {
		return nil, err
	}
	addGoVersion(version, &urlPath)
	details, err := fileutils.GetFileDetails(localPath, true)
	if err != nil {
		return nil, err
	}
	utils.AddChecksumHeaders(gpc.clientDetails.Headers, details)
	resp, _, err := gpc.client.UploadFile(localPath, urlPath, "", &gpc.clientDetails, nil)
	if err != nil {
		return nil, err
	}
	sha256 := resp.Header.Get("X-Checksum-Sha256")
	if sha256 == "" {
		log.Info("Failed to extract file's sha256 from response body.\nFile: " + localPath)
	}
	// Remove urls properties suffix
	splitUrlPath := strings.Split(urlPath, ";")
	// Remove "api/go/" substring from url to get the actual file's path in Artifactory
	targetPath := strings.ReplaceAll(splitUrlPath[0], "api/go/", "")
	filesDetails := clientutils.FileTransferDetails{SourcePath: localPath, TargetPath: targetPath, Sha256: sha256}
	return &filesDetails, errorutils.CheckResponseStatus(resp, http.StatusCreated)
}
