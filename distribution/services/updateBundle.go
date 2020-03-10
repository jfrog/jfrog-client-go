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

type ReleaseNotesSyntax string

const (
	Markdown  ReleaseNotesSyntax = "markdown"
	Asciidoc                     = "asciidoc"
	PlainText                    = "plain_text"
)

type UpdateReleaseBundleService struct {
	client        *rthttpclient.ArtifactoryHttpClient
	DistDetails   auth.CommonDetails
	DryRun        bool
	GpgPassphrase string
}

func NewUpdateReleseBundleService(client *rthttpclient.ArtifactoryHttpClient) *UpdateReleaseBundleService {
	return &UpdateReleaseBundleService{client: client}
}

func (ubs *UpdateReleaseBundleService) GetDistDetails() auth.CommonDetails {
	return ubs.DistDetails
}

func (ubs *UpdateReleaseBundleService) CreateReleaseBundle(createBundleParams CreateUpdateReleaseBundleParams) error {
	releaseBundleBody, err := CreateBundleBody(createBundleParams, ubs.DryRun)
	if err != nil {
		return err
	}
	return ubs.execUpdateReleaseBundle(createBundleParams.Name, createBundleParams.Version, releaseBundleBody)
}

func (ubs *UpdateReleaseBundleService) execUpdateReleaseBundle(name, version string, releaseBundle *ReleaseBundleBody) error {
	httpClientsDetails := ubs.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(releaseBundle)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := ubs.DistDetails.GetUrl() + "api/v1/release_bundle/" + name + "/" + version
	distrbutionServiceUtils.SetGpgPassphrase(ubs.GpgPassphrase, &httpClientsDetails.Headers)
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := ubs.client.SendPut(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Artifactory response: ", resp.Status)
	log.Output(utils.IndentJson(body))
	return errorutils.CheckError(err)
}

func CreateBundleBody(createBundleParams CreateUpdateReleaseBundleParams, dryRun bool) (*ReleaseBundleBody, error) {
	var bundleQueries []BundleQuery
	// Create release bundle queries
	for _, specFile := range createBundleParams.SpecFiles {
		if specFile.GetSpecType() != artifactoryUtils.AQL {
			query, err := artifactoryUtils.CreateAqlBodyForSpecWithPattern(specFile)
			if err != nil {
				return nil, err
			}
			specFile.Aql = artifactoryUtils.Aql{ItemsFind: query}
			aql := artifactoryUtils.BuildQueryFromSpecFile(specFile, artifactoryUtils.NONE)
			bundleQueries = append(bundleQueries, BundleQuery{Aql: aql})
		}
	}

	// Create release bundle struct
	releaseBundleBody := &ReleaseBundleBody{
		DryRun:            dryRun,
		SignImmediately:   createBundleParams.SignImmediately,
		StoringRepository: createBundleParams.StoringRepository,
		Description:       createBundleParams.Description,
		BundleSpec: BundleSpec{
			Queries: bundleQueries,
		},
	}

	// Add relese notes if needed
	if createBundleParams.ReleaseNotes != "" {
		releaseBundleBody.ReleaseNotes = &ReleaseNotes{
			Syntax:  createBundleParams.ReleaseNotesSyntax,
			Content: createBundleParams.ReleaseNotes,
		}
	}
	return releaseBundleBody, nil
}

type ReleaseBundleBody struct {
	DryRun            bool          `json:"dry_run"`
	SignImmediately   bool          `json:"sign_immediately,omitempty"`
	StoringRepository string        `json:"storing_repository,omitempty"`
	Description       string        `json:"description,omitempty"`
	ReleaseNotes      *ReleaseNotes `json:"release_notes,omitempty"`
	BundleSpec        BundleSpec    `json:"spec"`
}

type ReleaseNotes struct {
	Syntax  ReleaseNotesSyntax `json:"syntax,omitempty"`
	Content string             `json:"content,omitempty"`
}

type BundleSpec struct {
	Queries []BundleQuery `json:"queries"`
}

type BundleQuery struct {
	QueryName string `json:"query_name,omitempty"`
	Aql       string `json:"aql"`
}

type CreateUpdateReleaseBundleParams struct {
	SpecFiles          []*artifactoryUtils.ArtifactoryCommonParams
	Name               string
	Version            string
	SignImmediately    bool
	StoringRepository  string
	Description        string
	ReleaseNotes       string
	ReleaseNotesSyntax ReleaseNotesSyntax
}

func NewCreateUpdateBundleParams(name, version string) CreateUpdateReleaseBundleParams {
	return CreateUpdateReleaseBundleParams{
		Name:    name,
		Version: version,
	}
}
