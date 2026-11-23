package services

import (
	"net/http"
	"net/url"
	"path"

	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const evidenceDeleteApi = "api/v1/evidence"

func (es *EvidenceService) DeleteEvidence(subjectRepoPath, evidenceName string) error {
	fullPath := path.Join(evidenceDeleteApi, subjectRepoPath, evidenceName)
	u := url.URL{Path: fullPath}

	requestFullUrl, err := clientutils.BuildUrl(es.GetEvidenceDetails().GetUrl(), u.String(), nil)
	if err != nil {
		return errorutils.CheckError(err)
	}

	httpClientDetails := es.GetEvidenceDetails().CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	log.Debug("Deleting evidence: ", requestFullUrl)

	resp, body, err := es.client.SendDelete(requestFullUrl, nil, &httpClientDetails)
	if err != nil {
		return err
	}

	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusNoContent)
}
