package _go

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/jfrog/gofrog/version"
	"github.com/jfrog/jfrog-client-go/utils/log"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
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
	goApiUrl, err := utils.BuildArtifactoryUrl(ArtDetails.GetUrl(), "api/go/", make(map[string]string))
	if err != nil {
		return nil, err
	}
	gpc.clientDetails = ArtDetails.CreateHttpClientDetails()
	gpc.client = client
	moduleId := strings.Split(params.GetModuleId(), ":")
	totalSucceed, totalFailed := 0, 0
	var filesDetails []clientutils.FileTransferDetails
	// Upload zip file
	success, failed, err := gpc.uploadFile(params, params.ZipPath, moduleId[0], ".zip", goApiUrl, &filesDetails, gpc)
	if err != nil {
		return nil, err
	}
	totalSucceed, totalFailed = totalSucceed+success, totalFailed+failed
	// Upload mod file
	success, failed, err = gpc.uploadFile(params, params.ModPath, moduleId[0], ".mod", goApiUrl, &filesDetails, gpc)
	if err != nil {
		return nil, err
	}
	totalSucceed, totalFailed = totalSucceed+success, totalFailed+failed
	if version.NewVersion(gpc.artifactoryVersion).AtLeast(ArtifactoryMinSupportedVersion) && params.GetInfoPath() != "" {
		// Upload info file. This is supported from Artifactory version 6.10.0 and above
		success, failed, err = gpc.uploadFile(params, params.InfoPath, moduleId[0], ".info", goApiUrl, &filesDetails, gpc)
		totalSucceed, totalFailed = totalSucceed+success, totalFailed+failed
		if err != nil {
			return nil, err
		}
	}
	fileTransferDetailsTempFile, err := clientutils.SaveFileTransferDetailsInTempFile(&filesDetails)
	if err != nil {
		return nil, err
	}

	return &utils.OperationSummary{TotalSucceeded: totalSucceed, TotalFailed: totalFailed, TransferDetailsReader: content.NewContentReader(fileTransferDetailsTempFile, "files")}, nil
}

func (gpc *GoPublishCommand) uploadFile(params GoParams, filePath string, moduleId, ext, goApiUrl string, filesDetails *[]clientutils.FileTransferDetails, pwa *GoPublishCommand) (success, failed int, err error) {
	success, failed = 0, 1
	pathInArtifactory := strings.Join([]string{params.GetTargetRepo(), moduleId, "@v", params.GetVersion() + ext}, "/")
	details, err := pwa.upload(filePath, pathInArtifactory, params.GetVersion(), params.GetProps(), goApiUrl)
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
// pathInArtifactory - The path of the file in Artifactory for example: go-repo/github.com/jfrog/jfrog-client-go/@v/v1.1.1.zip
// version - The version of the project that being uploaded.
// props - The properties to be assigned for each artifact
// ext - The extension of the file: zip, mod, info. This extension will be joined with the version for the path. For example: v1.2.3.info or v1.2.3.zip
// goApiUrl - The URL of the Go API in Artifactory. For example: http://127.0.0.1/artifactory/api/go/
func (gpc *GoPublishCommand) upload(localPath, pathInArtifactory, version, props, goApiUrl string) (*clientutils.FileTransferDetails, error) {
	rtUrl := strings.ReplaceAll(goApiUrl, "api/go/", "")
	err := CreateUrlPath(pathInArtifactory, props, &goApiUrl)
	if err != nil {
		return nil, err
	}
	addGoVersion(version, &goApiUrl)
	details, err := fileutils.GetFileDetails(localPath, true)
	if err != nil {
		return nil, err
	}
	utils.AddChecksumHeaders(gpc.clientDetails.Headers, details)
	resp, body, err := gpc.client.UploadFile(localPath, goApiUrl, "", &gpc.clientDetails, nil)
	if err != nil {
		return nil, err
	}
	sha256 := resp.Header.Get("X-Checksum-Sha256")
	if sha256 == "" {
		log.Info("Failed to extract file's sha256 from response body.\nFile: " + localPath)
	}
	filesDetails := clientutils.FileTransferDetails{SourcePath: localPath, TargetPath: pathInArtifactory, RtUrl: rtUrl, Sha256: sha256}
	return &filesDetails, errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated)
}
