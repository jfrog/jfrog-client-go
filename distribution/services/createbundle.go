package services

import (
	"encoding/json"
	"errors"
	"net/http"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type ReleaseBundle struct {
	Name              string       `json:"name"`
	Version           string       `json:"version"`
	DryRun            bool         `json:"dry_run,omitempty"`
	SignImmediately   bool         `json:"sign_immediately,omitempty"`
	StoringRepository string       `json:"storing_repository,omitempty"`
	Description       string       `json:"description"`
	ReleaseNotes      ReleaseNotes `json:"release_notes,omitempty"`
	BundleSpec        BundleSpec   `json:"spec"`
}

type ReleaseNotes struct {
	Syntax  ReleaseNotesSyntax `json:"syntax,omitempty"`
	Content string             `json:"content"`
}

type BundleSpec struct {
	Queries []BundleQuery `json:"queries"`
}

type BundleQuery struct {
	QueryName string `json:"query_name,omitempty"`
	Aql       string `json:"aql"`
}

type ReleaseNotesSyntax string

const (
	Markdown  ReleaseNotesSyntax = "markdown"
	Asciidoc                     = "asciidoc"
	PlainText                    = "plain_text"
)

type CreateBundleService struct {
	client             *rthttpclient.ArtifactoryHttpClient
	DistDetails        auth.CommonDetails
	DryRun             bool
	Name               string
	Version            string
	SignImmediately    bool
	StoringRepository  string
	Description        string
	ReleaseNotesPath   string
	ReleaseNotesSyntax ReleaseNotesSyntax
}

func NewCreateBundleService(client *rthttpclient.ArtifactoryHttpClient) *CreateBundleService {
	return &CreateBundleService{client: client}
}

func (ps *CreateBundleService) GetDistDetails() auth.CommonDetails {
	return ps.DistDetails
}

func (cbs *CreateBundleService) CreateReleaseBundle(createBundleParams CreateBundleParams) error {
	var bundleQueries []BundleQuery
	for _, specFile := range createBundleParams.SpecFiles {
		if specFile.GetSpecType() != artifactoryUtils.AQL {
			query, err := artifactoryUtils.CreateAqlBodyForSpecWithPattern(specFile)
			if err != nil {
				return err
			}
			specFile.Aql = artifactoryUtils.Aql{ItemsFind: query}
			aql := artifactoryUtils.BuildQueryFromSpecFile(specFile, artifactoryUtils.NONE)
			bundleQueries = append(bundleQueries, BundleQuery{Aql: aql})
		}
	}
	releaseBundle := ReleaseBundle{
		Name:              cbs.Name,
		Version:           cbs.Version,
		DryRun:            cbs.DryRun,
		SignImmediately:   cbs.SignImmediately,
		StoringRepository: cbs.StoringRepository,
		Description:       cbs.Description,
		ReleaseNotes: ReleaseNotes{
			Syntax: cbs.ReleaseNotesSyntax,
			// Content: cbs.ReleaseNotesPath, // TODO
		},
		BundleSpec: BundleSpec{
			Queries: bundleQueries,
		},
	}

	return cbs.execCreateReleaseBundle(releaseBundle)
}

func (cbs *CreateBundleService) execCreateReleaseBundle(releaseBundle ReleaseBundle) error {
	httpClientsDetails := cbs.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(releaseBundle)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := cbs.DistDetails.GetUrl() + "api/v1/release_bundle"
	resp, body, err := cbs.client.SendPost(url, content, &httpClientsDetails)
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Artifactory response: ", resp.Status)
	return errorutils.CheckError(err)
}

type CreateBundleParams struct {
	SpecFiles          []*artifactoryUtils.ArtifactoryCommonParams
	Name               string
	Version            string
	SignImmediately    bool
	StoringRepository  string
	Description        string
	ReleaseNotesPath   string
	ReleaseNotesSyntax ReleaseNotesSyntax
}

func NewCreateBundleParams(name, version string) CreateBundleParams {
	return CreateBundleParams{
		Name:    name,
		Version: version,
	}
}
