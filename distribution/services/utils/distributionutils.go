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

func CreateBundleBody(releaseBundleParams ReleaseBundleParams, dryRun bool) (*ReleaseBundleBody, error) {
	var bundleQueries []BundleQuery
	// Create release bundle queries
	for _, specFile := range releaseBundleParams.SpecFiles {
		// Create AQL
		aql, err := createAql(specFile)
		if err != nil {
			return nil, err
		}

		// Create added properties
		addProps, err := createAddProps(specFile)
		if err != nil {
			return nil, err
		}

		// Append bundle query
		bundleQueries = append(bundleQueries, BundleQuery{Aql: aql, AddProps: addProps})
	}

	// Create release bundle struct
	releaseBundleBody := &ReleaseBundleBody{
		DryRun:            dryRun,
		SignImmediately:   releaseBundleParams.SignImmediately,
		StoringRepository: releaseBundleParams.StoringRepository,
		Description:       releaseBundleParams.Description,
		BundleSpec: BundleSpec{
			Queries: bundleQueries,
		},
	}

	// Add relese notes if needed
	if releaseBundleParams.ReleaseNotes != "" {
		releaseBundleBody.ReleaseNotes = &ReleaseNotes{
			Syntax:  releaseBundleParams.ReleaseNotesSyntax,
			Content: releaseBundleParams.ReleaseNotes,
		}
	}
	return releaseBundleBody, nil
}

// Create the AQL query from the input spec
func createAql(specFile *utils.ArtifactoryCommonParams) (string, error) {
	if specFile.GetSpecType() != utils.AQL {
		query, err := utils.CreateAqlBodyForSpecWithPattern(specFile)
		if err != nil {
			return "", err
		}
		specFile.Aql = utils.Aql{ItemsFind: query}
	}
	return utils.BuildQueryFromSpecFile(specFile, utils.NONE), nil
}

// Create the AddProps array from the input AddProps string
func createAddProps(specFile *utils.ArtifactoryCommonParams) ([]AddProps, error) {
	props, err := utils.ParseProperties(specFile.AddProps, utils.SplitCommas)
	if err != nil {
		return nil, err
	}

	var addProps []AddProps
	for key, values := range props.ToMap() {
		addProps = append(addProps, AddProps{key, values})
	}
	return addProps, nil
}

func AddGpgPassphraseHeader(gpgPassphrase string, headers *map[string]string) {
	if gpgPassphrase != "" {
		utils.AddHeader("X-GPG-PASSPHRASE", gpgPassphrase, headers)
	}
}
