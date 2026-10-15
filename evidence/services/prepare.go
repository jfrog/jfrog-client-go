package services

import (
	"encoding/json"
	"net/http"

	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const prepareEvidenceAPI = "api/v1/evidence/prepare"

type SubjectType string

const (
	SubjectTypeArtifact           SubjectType = "artifact"
	SubjectTypeBuild              SubjectType = "build"
	SubjectTypePackage            SubjectType = "package"
	SubjectTypeReleaseBundle      SubjectType = "release_bundle"
	SubjectTypeApplicationVersion SubjectType = "application_version"
	SubjectTypeEntity             SubjectType = "entity"
)

// PrepareEvidenceRequest contains the data used by Evidence to generate an in-toto statement for signing.
// For entity subjects, at most one of EntityRepo, ProjectKey, and ApplicationKey may be set.
type PrepareEvidenceRequest struct {
	Predicate      json.RawMessage             `json:"predicate"`
	PredicateType  string                      `json:"predicate_type"`
	Markdown       string                      `json:"markdown,omitempty"`
	ProviderID     string                      `json:"provider_id,omitempty"`
	Subject        PrepareEvidenceSubject      `json:"subject"`
	EntityRepo     string                      `json:"entity_repo,omitempty"`
	ProjectKey     string                      `json:"project_key,omitempty"`
	ApplicationKey string                      `json:"application_key,omitempty"`
	Attachments    []PrepareEvidenceAttachment `json:"attachments,omitempty"`
}

// PrepareEvidenceSubject identifies the subject for which Evidence prepares the statement.
// The fields required by Evidence depend on SubjectType.
type PrepareEvidenceSubject struct {
	SubjectType   SubjectType `json:"subject_type"`
	SubjectSHA256 string      `json:"subject_sha256,omitempty"`

	RepoPath string `json:"repo_path,omitempty"`

	EntityType string `json:"entity_type,omitempty"`
	EntityID   string `json:"entity_id,omitempty"`

	BuildName      string `json:"build_name,omitempty"`
	BuildNumber    string `json:"build_number,omitempty"`
	BuildTimestamp string `json:"build_timestamp,omitempty"`

	ReleaseBundleName    string `json:"release_bundle_name,omitempty"`
	ReleaseBundleVersion string `json:"release_bundle_version,omitempty"`

	ApplicationKey     string `json:"application_key,omitempty"`
	ApplicationVersion string `json:"application_version,omitempty"`

	PackageRepo    string `json:"package_repo,omitempty"`
	PackageName    string `json:"package_name,omitempty"`
	PackageVersion string `json:"package_version,omitempty"`
}

// PrepareEvidenceAttachment identifies an attachment in Artifactory.
// SHA256 is optional in a prepare request and is resolved by Evidence.
type PrepareEvidenceAttachment struct {
	Repository string `json:"repository"`
	Path       string `json:"path"`
	SHA256     string `json:"sha256,omitempty"`
}

type PrepareEvidenceResponse struct {
	PostURL                   string                      `json:"post_url"`
	DSSEPayload               string                      `json:"dsse_payload"`
	DSSEPayloadType           string                      `json:"dsse_payload_type"`
	PreAuthenticationEncoding string                      `json:"pre_authentication_encoding,omitempty"`
	Attachments               []PrepareEvidenceAttachment `json:"attachments,omitempty"`
}

// PrepareEvidence asks Evidence to generate an in-toto statement for external signing.
// When includePAE is true, the response also contains the DSSE pre-authentication encoding.
func (es *EvidenceService) PrepareEvidence(request PrepareEvidenceRequest, includePAE bool) (*PrepareEvidenceResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	queryParams := make(map[string]string)
	if includePAE {
		queryParams["include_pae"] = "true"
	}
	requestFullURL, err := clientutils.BuildUrl(es.GetEvidenceDetails().GetUrl(), prepareEvidenceAPI, queryParams)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	httpClientDetails := es.GetEvidenceDetails().CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	log.Debug("Preparing Evidence for signing")
	resp, body, err := es.client.SendPost(requestFullURL, requestBody, &httpClientDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}

	var response PrepareEvidenceResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, errorutils.CheckError(err)
	}
	return &response, nil
}
