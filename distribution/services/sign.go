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
	client      *rthttpclient.ArtifactoryHttpClient
	DistDetails auth.ServiceDetails
}

func NewSignBundleService(client *rthttpclient.ArtifactoryHttpClient) *SignBundleService {
	return &SignBundleService{client: client}
}

func (sb *SignBundleService) GetDistDetails() auth.ServiceDetails {
	return sb.DistDetails
}

func (sb *SignBundleService) SignReleaseBundle(signBundleParams SignBundleParams) error {
	signBundleBody := &SignBundleBody{
		StoringRepository: signBundleParams.StoringRepository,
	}
	return sb.execSignReleaseBundle(signBundleParams.Name, signBundleParams.Version, signBundleParams.GpgPassphrase, signBundleBody)
}

func (sb *SignBundleService) execSignReleaseBundle(name, version, gpgPassphrase string, signBundleBody *SignBundleBody) error {
	httpClientsDetails := sb.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(signBundleBody)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := sb.DistDetails.GetUrl() + "api/v1/release_bundle/" + name + "/" + version + "/sign"
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	distrbutionServiceUtils.AddGpgPassphraseHeader(gpgPassphrase, &httpClientsDetails.Headers)
	resp, body, err := sb.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Distribution response: ", resp.Status)
	log.Debug(utils.IndentJson(body))
	return errorutils.CheckError(err)
}

type SignBundleBody struct {
	StoringRepository string `json:"storing_repository,omitempty"`
}

type SignBundleParams struct {
	Name              string
	Version           string
	StoringRepository string
	GpgPassphrase     string
}

func NewSignBundleParams(name, version string) SignBundleParams {
	return SignBundleParams{
		Name:    name,
		Version: version,
	}
}
