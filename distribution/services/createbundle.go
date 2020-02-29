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

type ReleaseNotesSyntax string

const (
	Markdown  ReleaseNotesSyntax = "markdown"
	Asciidoc                     = "asciidoc"
	PlainText                    = "plain_text"
)

type CreateBundleService struct {
	client      *rthttpclient.ArtifactoryHttpClient
	DistDetails auth.CommonDetails
	DryRun      bool
}

func NewCreateBundleService(client *rthttpclient.ArtifactoryHttpClient) *CreateBundleService {
	return &CreateBundleService{client: client}
}

func (ps *CreateBundleService) GetDistDetails() auth.CommonDetails {
	return ps.DistDetails
}

func (cbs *CreateBundleService) CreateReleaseBundle(createBundleParams CreateBundleParams) error {
	var bundleQueries []BundleQuery
	// Create release bundle queries
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

	// Create release bundle struct
	releaseBundle := ReleaseBundleBody{
		Name:              createBundleParams.Name,
		Version:           createBundleParams.Version,
		DryRun:            cbs.DryRun,
		SignImmediately:   createBundleParams.SignImmediately,
		StoringRepository: createBundleParams.StoringRepository,
		Description:       createBundleParams.Description,
		BundleSpec: BundleSpec{
			Queries: bundleQueries,
		},
	}

	// Add relese notes if needed
	if createBundleParams.ReleaseNotes != "" {
		releaseBundle.ReleaseNotes = &ReleaseNotes{
			Syntax:  createBundleParams.ReleaseNotesSyntax,
			Content: createBundleParams.ReleaseNotes,
		}
	}

	return cbs.execCreateReleaseBundle(releaseBundle)
}

func (cbs *CreateBundleService) execCreateReleaseBundle(releaseBundle ReleaseBundleBody) error {
	httpClientsDetails := cbs.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(releaseBundle)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := cbs.DistDetails.GetUrl() + "api/v1/release_bundle"
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := cbs.client.SendPost(url, content, &httpClientsDetails)
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

type ReleaseBundleBody struct {
	Name              string        `json:"name"`
	Version           string        `json:"version"`
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

type CreateBundleParams struct {
	SpecFiles          []*artifactoryUtils.ArtifactoryCommonParams
	Name               string
	Version            string
	SignImmediately    bool
	StoringRepository  string
	Description        string
	ReleaseNotes       string
	ReleaseNotesSyntax ReleaseNotesSyntax
}

func NewCreateBundleParams(name, version string) CreateBundleParams {
	return CreateBundleParams{
		Name:    name,
		Version: version,
	}
}
