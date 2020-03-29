package auth

import (
	"encoding/json"
	"errors"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"net/http"
	"strings"
	"sync"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func NewDistributionDetails() *distributionDetails {
	return &distributionDetails{}
}

var expiryHandleMutex sync.Mutex

type distributionDetails struct {
	auth.CommonConfigFields
}

func (ds *distributionDetails) GetVersion() (string, error) {
	var err error
	if ds.Version == "" {
		ds.Version, err = ds.getDistributionVersion()
		if err != nil {
			return "", err
		}
		log.Debug("The Distribution version is:", ds.Version)
	}
	return ds.Version, nil
}

func (ds *distributionDetails) getDistributionVersion() (string, error) {
	cd := auth.CommonDetails(ds)
	client, err := rthttpclient.ArtifactoryClientBuilder().
		SetCommonDetails(&cd).
		Build()
	if err != nil {
		return "", err
	}
	httpDetails := ds.CreateHttpClientDetails()
	resp, body, _, err := client.SendGet(ds.GetUrl()+"api/v1/system/info", true, &httpDetails)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	var version distributionVersion
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	return strings.TrimSpace(version.Version), nil
}

type distributionVersion struct {
	Version string `json:"version,omitempty"`
}
