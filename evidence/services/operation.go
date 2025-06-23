package services

import (
	"net/http"
	"net/url"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type EvidenceService struct {
	client          *jfroghttpclient.JfrogHttpClient
	evidenceDetails *auth.ServiceDetails
}

func NewEvidenceService(evidenceDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *EvidenceService {
	return &EvidenceService{evidenceDetails: &evidenceDetails, client: client}
}

func (es *EvidenceService) GetEvidenceDetails() auth.ServiceDetails {
	return *es.evidenceDetails
}

type EvidenceOperation interface {
	getOperationRestApi() string
	getRequestBody() []byte
	getProviderId() string
}

func (es *EvidenceService) doOperation(operation EvidenceOperation) ([]byte, error) {
	u := url.URL{Path: operation.getOperationRestApi()}
	queryParams := make(map[string]string)
	if operation.getProviderId() != "" {
		queryParams["providerId"] = operation.getProviderId()
	}

	requestFullUrl, err := clientutils.BuildUrl(es.GetEvidenceDetails().GetUrl(), u.String(), queryParams)
	if err != nil {
		return []byte{}, errorutils.CheckError(err)
	}

	httpClientDetails := es.GetEvidenceDetails().CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	log.Debug("Creating Evidence:")
	resp, body, err := es.client.SendPost(requestFullUrl, operation.getRequestBody(), &httpClientDetails)
	if err != nil {
		return []byte{}, err
	}

	return body, errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated)
}

type EvidenceDetails struct {
	SubjectUri  string `json:"subject_uri"`
	DSSEFileRaw []byte `json:"dsse_file_raw"`
	ProviderId  string `json:"provider_id"`
}
