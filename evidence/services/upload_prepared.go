package services

import (
	"fmt"
	"net/url"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

// UploadPreparedSignedEvidence uploads a signed DSSE envelope to the root-relative
// post URL returned by PrepareEvidence.
func (es *EvidenceService) UploadPreparedSignedEvidence(postURL string, signedEnvelope []byte) ([]byte, error) {
	requestURL, err := es.resolvePreparedEvidencePostURL(postURL)
	if err != nil {
		return nil, err
	}
	return es.uploadEvidenceBody(requestURL, signedEnvelope)
}

func (es *EvidenceService) resolvePreparedEvidencePostURL(postURL string) (string, error) {
	post, err := url.Parse(postURL)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	if postURL == "" || post.IsAbs() || post.Host != "" || post.Path == "" || post.Path[0] != '/' {
		return "", fmt.Errorf("prepared Evidence post URL must be a non-empty root-relative URL")
	}

	base, err := url.Parse(es.GetEvidenceDetails().GetUrl())
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	if base.Scheme == "" || base.Host == "" {
		return "", fmt.Errorf("Evidence URL must include a scheme and host")
	}

	post.Scheme = base.Scheme
	post.Host = base.Host
	post.User = base.User
	return post.String(), nil
}
