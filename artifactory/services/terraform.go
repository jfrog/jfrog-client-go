package services

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/utils/version"
	"net/http"
	"strings"
)

const ArtifactoryMinSupportedVersion = "6.10.0" // change version

// Support for Artifactory 6.10.0 and above API
type TerraformPublishCommand struct {
	artifactoryVersion string
	clientDetails      httputils.HttpClientDetails
	client             *jfroghttpclient.JfrogHttpClient
}

type TerraformService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewTerraformService(client *jfroghttpclient.JfrogHttpClient) *TerraformService {
	return &TerraformService{client: client}
}

func (gs *TerraformService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return gs.client
}

func (gs *TerraformService) SetServiceDetails(artDetails auth.ServiceDetails) {
	gs.ArtDetails = artDetails
}

//func (gs *TerraformService) PublishModule(params TerraformParams) (*utils.OperationSummary, error) {
//	artifactoryVersion, err := gs.ArtDetails.GetVersion()
//	if err != nil {
//		return nil, err
//	}
//	publisher := &TerraformPublishCommand{}
//	// PublishPackage supports Artifactory version "6.10.0" and above.
//	err = publisher.verifyCompatibleVersion(artifactoryVersion)
//	if err != nil {
//		return nil, err
//	}
//	return publisher.PublishPackage(params, gs.client, gs.ArtDetails)
//}

type TerraformParams struct {
	ZipPath string
	//ModPath    string
	//InfoPath   string
	//ModContent []byte
	Namespace  string
	ModuleName string
	Provider   string
	Tag        string
	TargetRepo string
	//ModuleId   string
}

func (gp *TerraformParams) GetZipPath() string {
	return gp.ZipPath
}

func (gp *TerraformParams) GetNamespace() string {
	return gp.Namespace
}

func (gp *TerraformParams) GetModuleName() string {
	return gp.ModuleName
}

func (gp *TerraformParams) GetProvider() string {
	return gp.Provider
}

func (gp *TerraformParams) GetTag() string {
	return gp.Tag
}

func (gp *TerraformParams) GetTargetRepo() string {
	return gp.TargetRepo
}

func NewTerraformParams() TerraformParams {
	return TerraformParams{}
}

func (gpc *TerraformPublishCommand) verifyCompatibleVersion(artifactoryVersion string) error {
	propertiesApi := ArtifactoryMinSupportedVersion
	version := version.NewVersion(artifactoryVersion)
	gpc.artifactoryVersion = artifactoryVersion
	if !version.AtLeast(propertiesApi) {
		return errorutils.CheckErrorf("Unsupported version of Artifactory: %s\nSupports Artifactory version %s and above", artifactoryVersion, propertiesApi)
	}
	return nil
}

//func (gpc *TerraformPublishCommand) PublishPackage(params TerraformParams, client *jfroghttpclient.JfrogHttpClient, ArtDetails auth.ServiceDetails) (*utils.OperationSummary, error) {
//	url, err := utils.BuildArtifactoryUrl(ArtDetails.GetUrl(), "api/terraform/"+params.GetTargetRepo(), make(map[string]string))
//	if err != nil {
//		return nil, err
//	}
//	gpc.clientDetails = ArtDetails.CreateHttpClientDetails()
//	gpc.client = client
//	moduleId := strings.Split(params.GetModuleId(), ":")
//	totalSucceed, totalFailed := 0, 0
//	var filesDetails []clientutils.FileTransferDetails
//	// Upload zip file
//	success, failed, err := gpc.uploadFile(params, params.ZipPath, moduleId[0], ".zip", url, &filesDetails, gpc)
//	if err != nil {
//		return nil, err
//	}
//	totalSucceed, totalFailed = totalSucceed+success, totalFailed+failed
//	// Upload mod file
//	success, failed, err = gpc.uploadFile(params, params.ModPath, moduleId[0], ".mod", url, &filesDetails, gpc)
//	if err != nil {
//		return nil, err
//	}
//	totalSucceed, totalFailed = totalSucceed+success, totalFailed+failed
//	if version.NewVersion(gpc.artifactoryVersion).AtLeast(ArtifactoryMinSupportedVersion) && params.GetInfoPath() != "" {
//		// Upload info file. This is supported from Artifactory version 6.10.0 and above
//		success, failed, err = gpc.uploadFile(params, params.InfoPath, moduleId[0], ".info", url, &filesDetails, gpc)
//		totalSucceed, totalFailed = totalSucceed+success, totalFailed+failed
//		if err != nil {
//			return nil, err
//		}
//	}
//	tempFile, err := clientutils.SaveFileTransferDetailsInTempFile(&filesDetails)
//	if err != nil {
//		return nil, err
//	}
//	return &utils.OperationSummary{TotalSucceeded: totalSucceed, TotalFailed: totalFailed, TransferDetailsReader: content.NewContentReader(tempFile, "files")}, nil
//}
//
//func (gpc *TerraformPublishCommand) uploadFile(params TerraformParams, filePath, namespace, moduleName, provider, tag, url string, filesDetails *[]clientutils.FileTransferDetails, pwa *TerraformPublishCommand) (success, failed int, err error) {
//	success, failed = 0, 1
//	details, err := pwa.upload(filePath, moduleId, params.GetVersion(), params.GetProps(), ext, url)
//	if err != nil {
//		return
//	}
//	success, failed = 1, 0
//	*filesDetails = append(*filesDetails, *details)
//	return
//}

// localPath - The location of the file on the file system.
// moduleId - The name of the module for example github.com/jfrog/jfrog-client-go.
// version - The version of the project that being uploaded.
// props - The properties to be assigned for each artifact
// ext - The extension of the file: zip, mod, info. This extension will be joined with the version for the path. For example v1.2.3.info or v1.2.3.zip
// urlPath - The url including the repository. For example: http://127.0.0.1/artifactory/api/go/go-local
func (gpc *TerraformPublishCommand) upload(localPath, moduleId, version, props, ext, urlPath string) (*clientutils.FileTransferDetails, error) {
	err := CreateUrlPath(moduleId, version, props, ext, &urlPath)
	if err != nil {
		return nil, err
	}
	//addGoVersion(version, &urlPath)
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

func CreateUrlPath(namespace, moduleName, provider, tag string, url *string) error {
	*url = strings.Join([]string{*url, namespace, moduleName, provider, tag + ".zip"}, "/")
	//properties, err := utils.ParseProperties(props)
	//if err != nil {
	//	return err
	//}

	//*url = strings.Join([]string{*url, properties.ToEncodedString(true)}, ";")
	//if strings.HasSuffix(*url, ";") {
	//	tempUrl := *url
	//	tempUrl = tempUrl[:len(tempUrl)-1]
	//	*url = tempUrl
	//}
	return nil
}
