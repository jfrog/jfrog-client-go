package _go

import (
	"encoding/base64"
	"github.com/Masterminds/semver"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/errors/httperrors"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"strings"
)

type GoService struct {
	client     *httpclient.HttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewGoService(client *httpclient.HttpClient) *GoService {
	return &GoService{client: client}
}

func (gs *GoService) GetJfrogHttpClient() *httpclient.HttpClient {
	return gs.client
}

func (gs *GoService) SetArtDetails(artDetails auth.ArtifactoryDetails) {
	gs.ArtDetails = artDetails
}

func (gs *GoService) PublishPackage(params GoParams) error {
	url, err := utils.BuildArtifactoryUrl(gs.ArtDetails.GetUrl(), "api/go/"+params.GetTargetRepo(), make(map[string]string))
	clientDetails := gs.ArtDetails.CreateHttpClientDetails()

	utils.AddHeader("X-GO-MODULE-VERSION", params.GetVersion(), &clientDetails.Headers)
	utils.AddHeader("X-GO-MODULE-CONTENT", base64.StdEncoding.EncodeToString(params.GetModContent()), &clientDetails.Headers)
	artifactoryVersion, err := gs.ArtDetails.GetVersion()
	if err != nil {
		return err
	}
	if !shouldUseHeaders(artifactoryVersion) {
		createUrlPath(params, &url)
	} else {
		addPropertiesHeaders(params.GetProps(), &clientDetails.Headers)
	}

	resp, body, err := gs.client.UploadFile(params.GetZipPath(), url, clientDetails, 0)
	if err != nil {
		return err
	}
	return httperrors.CheckResponseStatus(resp, body, 201)
}

// This is needed when using Artifactory older then 6.5.0
func addPropertiesHeaders(props string, headers *map[string]string) error {
	properties, err := utils.ParseProperties(props, utils.JoinCommas)
	if err != nil {
		return err
	}
	headersMap := properties.ToHeadersMap()
	for k, v := range headersMap {
		utils.AddHeader("X-ARTIFACTORY-PROPERTY-"+k, v, headers)
	}
	return nil
}

func createUrlPath(params GoParams, url *string) error {
	*url = strings.Join([]string{*url, params.GetModuleId(), "@v", params.GetVersion() + ".zip"}, "/")
	properties, err := utils.ParseProperties(params.GetProps(), utils.JoinCommas)
	if err != nil {
		return err
	}
	*url = strings.Join([]string{*url, properties.ToEncodedString()}, ";")
	if strings.HasSuffix(*url, ";") {
		tempUrl := *url
		tempUrl = tempUrl[:len(tempUrl)-1]
		*url = tempUrl
	}
	return nil
}

// Returns true if needed to use properties as header (Artifactory version between 6.2.0 and 6.5.0)
// or false if need to use matrix params (Artifactory version 6.5.0 and above).
func shouldUseHeaders(artifactoryVersion string) bool {
	propertiesApi := "6.5.0"
	return artifactoryVersion != "development" && semver.MustParse(artifactoryVersion).Compare(semver.MustParse(propertiesApi)) < 0
}

type GoParams struct {
	ZipPath    string
	ModContent []byte
	Version    string
	Props      string
	TargetRepo string
	ModuleId   string
}

func (gpi *GoParams) GetZipPath() string {
	return gpi.ZipPath
}

func (gpi *GoParams) GetModContent() []byte {
	return gpi.ModContent
}

func (gpi *GoParams) GetVersion() string {
	return gpi.Version
}

func (gpi *GoParams) GetProps() string {
	return gpi.Props
}

func (gpi *GoParams) GetTargetRepo() string {
	return gpi.TargetRepo
}

func (gpi *GoParams) GetModuleId() string {
	return gpi.ModuleId
}

func NewGoParams() GoParams {
	return GoParams{}
}
