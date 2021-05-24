package services

import (
	"encoding/json"
	"errors"
	"net/http"

	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	distrbutionServiceUtils "github.com/jfrog/jfrog-client-go/distribution/services/utils"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type SignBundleService struct {
	client      *jfroghttpclient.JfrogHttpClient
	DistDetails auth.ServiceDetails
}

func NewSignBundleService(client *jfroghttpclient.JfrogHttpClient) *SignBundleService {
	return &SignBundleService{client: client}
}

func (sb *SignBundleService) GetDistDetails() auth.ServiceDetails {
	return sb.DistDetails
}

func (sb *SignBundleService) SignReleaseBundle(signBundleParams SignBundleParams) (*clientutils.Sha256Summary, error) {
	signBundleBody := &SignBundleBody{
		StoringRepository: signBundleParams.StoringRepository,
	}
	return sb.execSignReleaseBundle(signBundleParams.Name, signBundleParams.Version, signBundleParams.GpgPassphrase, signBundleBody)
}

func (sb *SignBundleService) execSignReleaseBundle(name, version, gpgPassphrase string, signBundleBody *SignBundleBody) (*clientutils.Sha256Summary, error) {
	summary := clientutils.NewSha256Summary()
	httpClientsDetails := sb.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(signBundleBody)
	if err != nil {
		return summary, errorutils.CheckError(err)
	}
	url := sb.DistDetails.GetUrl() + "api/v1/release_bundle/" + name + "/" + version + "/sign"
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	distrbutionServiceUtils.AddGpgPassphraseHeader(gpgPassphrase, &httpClientsDetails.Headers)
	resp, body, err := sb.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return summary, err
	}
	if resp.StatusCode != http.StatusOK {
		return summary, errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	summary.SetSucceeded(true)
	summary.SetSha256(resp.Header.Get("X-Checksum-Sha256"))

	log.Debug("Distribution response: ", resp.Status)
	log.Debug(clientutils.IndentJson(body))
	return summary, errorutils.CheckError(err)
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
