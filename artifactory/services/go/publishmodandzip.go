package _go

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/errors/httperrors"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/version"
	"net/url"
	"strings"
)

func init() {
	register(&publishModAndZipWithApi{})
}

// Support for Artifactory 6.6.0 and above API
type publishModAndZipWithApi struct {
}

func (pwa *publishModAndZipWithApi) isCompatible(artifactoryVersion string) bool {
	propertiesApi := "6.6.1"
	if version.Compare(artifactoryVersion, propertiesApi) < 0 && artifactoryVersion != "development" {
		return false
	}
	return true
}

func (pwa *publishModAndZipWithApi) PublishPackage(params GoParams, client *httpclient.HttpClient, ArtDetails auth.ArtifactoryDetails) error {
	url, err := utils.BuildArtifactoryUrl(ArtDetails.GetUrl(), "api/go/"+params.GetTargetRepo(), make(map[string]string))
	if err != nil {
		return err
	}
	zipUrl := url
	moduleId := strings.Split(params.GetModuleId(), ":")
	err = createUrlPath(moduleId[0], params.GetVersion(), params.GetProps(), ".zip", &zipUrl)
	if err != nil {
		return err
	}
	clientDetails := ArtDetails.CreateHttpClientDetails()

	addGoVersion(params, &zipUrl)
	resp, body, err := client.UploadFile(params.GetZipPath(), zipUrl, clientDetails, 0)
	if err != nil {
		return err
	}
	err = httperrors.CheckResponseStatus(resp, body, 201)
	if err != nil {
		return err
	}
	err = createUrlPath(moduleId[0], params.GetVersion(), params.GetProps(), ".mod", &url)
	if err != nil {
		return err
	}
	addGoVersion(params, &url)
	resp, body, err = client.UploadFile(params.GetModPath(), url, clientDetails, 0)
	if err != nil {
		return err
	}
	return httperrors.CheckResponseStatus(resp, body, 201)
}

func addGoVersion(params GoParams, urlPath *string) {
	*urlPath += ";go.version=" + url.QueryEscape(params.GetVersion())
}
