package _go

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/version"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	register(&publishZipAndModApi{})
}

// Support for Artifactory 6.6.1 and above API
type publishZipAndModApi struct {
}

func (pwa *publishZipAndModApi) isCompatible(artifactoryVersion string) bool {
	propertiesApi := "6.6.1"
	if version.Compare(artifactoryVersion, propertiesApi) < 0 && artifactoryVersion != "development" {
		return false
	}
	return true
}

func (pwa *publishZipAndModApi) PublishPackage(params GoParams, client *rthttpclient.ArtifactoryHttpClient, ArtDetails auth.ArtifactoryDetails) error {
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
	resp, _, err := client.UploadFile(params.GetZipPath(), zipUrl, "", &clientDetails, GoUploadRetries, nil)
	if err != nil {
		return err
	}
	err = errorutils.CheckResponseStatus(resp, http.StatusCreated)
	if err != nil {
		return err
	}
	err = createUrlPath(moduleId[0], params.GetVersion(), params.GetProps(), ".mod", &url)
	if err != nil {
		return err
	}
	addGoVersion(params, &url)
	resp, _, err = client.UploadFile(params.GetModPath(), url, "", &clientDetails, GoUploadRetries, nil)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatus(resp, http.StatusCreated)
}

func addGoVersion(params GoParams, urlPath *string) {
	*urlPath += ";go.version=" + url.QueryEscape(params.GetVersion())
}
