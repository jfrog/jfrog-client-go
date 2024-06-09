package services

import (
	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
	"net/url"
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
}

func (es *EvidenceService) doOperation(operation EvidenceOperation) ([]byte, error) {
	u := url.URL{Path: operation.getOperationRestApi()}
	requestFullUrl, err := url.Parse(es.GetEvidenceDetails().GetUrl() + u.String())
	if err != nil {
		return []byte{}, errorutils.CheckError(err)
	}

	httpClientDetails := es.GetEvidenceDetails().CreateHttpClientDetails()
	rtUtils.SetContentType("application/json", &httpClientDetails.Headers)

	resp, body, err := es.client.SendPost(requestFullUrl.String(), operation.getRequestBody(), &httpClientDetails)
	if err != nil {
		return []byte{}, err
	}

	return body, errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated)
}

type EvidenceDetails struct {
	SubjectUri  string `json:"subject_uri"`
	DSSEFileRaw []byte `json:"dsse_file_raw"`
}
