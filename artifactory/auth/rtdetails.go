package auth

import (
	"encoding/json"
	"errors"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strings"
	"sync"
)

func NewArtifactoryDetails() auth.CommonDetails {
	return &artifactoryDetails{}
}

var expiryHandleMutex sync.Mutex

type artifactoryDetails struct {
	auth.CommonConfigFields
}

func (rt *artifactoryDetails) GetVersion() (string, error) {
	var err error
	if rt.Version == "" {
		rt.Version, err = rt.getArtifactoryVersion()
		if err != nil {
			return "", err
		}
		log.Debug("The Artifactory version is:", rt.Version)
	}
	return rt.Version, nil
}

func (rt *artifactoryDetails) getArtifactoryVersion() (string, error) {
	cd := auth.CommonDetails(rt)
	client, err := rthttpclient.ArtifactoryClientBuilder().
		SetCommonDetails(&cd).
		Build()
	if err != nil {
		return "", err
	}
	httpDetails := rt.CreateHttpClientDetails()
	resp, body, _, err := client.SendGet(rt.GetUrl()+"api/system/version", true, &httpDetails)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	var version artifactoryVersion
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	return strings.TrimSpace(version.Version), nil
}

type artifactoryVersion struct {
	Version string `json:"version,omitempty"`
}
