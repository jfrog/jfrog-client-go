package services

import (
	"encoding/json"
	"errors"
	"net/http"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	distrbutionServiceUtils "github.com/jfrog/jfrog-client-go/distribution/services/utils"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type SignBundleService struct {
	client        *rthttpclient.ArtifactoryHttpClient
	DistDetails   auth.CommonDetails
	GpgPassphrase string
}

func NewSignBundleService(client *rthttpclient.ArtifactoryHttpClient) *SignBundleService {
	return &SignBundleService{client: client}
}

func (ps *SignBundleService) GetDistDetails() auth.CommonDetails {
	return ps.DistDetails
}

func (cbs *SignBundleService) SignReleaseBundle(signBundleParams SignBundleParams) error {
	signBundleBody := &SignBundleBody{
		StoringRepository: signBundleParams.StoringRepository,
	}
	return cbs.execSignReleaseBundle(signBundleParams.Name, signBundleParams.Version, signBundleBody)
}

func (cbs *SignBundleService) execSignReleaseBundle(name, version string, signBundleBody *SignBundleBody) error {
	httpClientsDetails := cbs.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(signBundleBody)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := cbs.DistDetails.GetUrl() + "api/v1/release_bundle/" + name + "/" + version + "/sign"
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	distrbutionServiceUtils.SetGpgPassphrase(cbs.GpgPassphrase, &httpClientsDetails.Headers)
	resp, body, err := cbs.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Distribution response: ", resp.Status)
	log.Output(utils.IndentJson(body))
	return errorutils.CheckError(err)
}

type SignBundleBody struct {
	StoringRepository string `json:"storing_repository,omitempty"`
}

type SignBundleParams struct {
	Name              string
	Version           string
	StoringRepository string
}

func NewSignBundleParams(name, version string) SignBundleParams {
	return SignBundleParams{
		Name:    name,
		Version: version,
	}
}
