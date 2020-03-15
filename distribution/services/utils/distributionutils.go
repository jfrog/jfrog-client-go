package utils

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

type ReleaseNotesSyntax string

const (
	Markdown  ReleaseNotesSyntax = "markdown"
	Asciidoc                     = "asciidoc"
	PlainText                    = "plain_text"
)

type ReleaseBundleParams struct {
	SpecFiles          []*utils.ArtifactoryCommonParams
	Name               string
	Version            string
	SignImmediately    bool
	StoringRepository  string
	Description        string
	ReleaseNotes       string
	ReleaseNotesSyntax ReleaseNotesSyntax
	GpgPassphrase      string
}

func NewReleaseBundleParams(name, version string) ReleaseBundleParams {
	return ReleaseBundleParams{
		Name:    name,
		Version: version,
	}
}

func CreateBundleBody(createBundleParams ReleaseBundleParams, dryRun bool) (*ReleaseBundleBody, error) {
	var bundleQueries []BundleQuery
	// Create release bundle queries
	for _, specFile := range createBundleParams.SpecFiles {
		if specFile.GetSpecType() != utils.AQL {
			query, err := utils.CreateAqlBodyForSpecWithPattern(specFile)
			if err != nil {
				return nil, err
			}
			specFile.Aql = utils.Aql{ItemsFind: query}
			aql := utils.BuildQueryFromSpecFile(specFile, utils.NONE)
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

func AddGpgPassphraseHeader(gpgPassphrase string, headers *map[string]string) {
	if gpgPassphrase != "" {
		utils.AddHeader("X-GPG-PASSPHRASE", gpgPassphrase, headers)
	}
}
