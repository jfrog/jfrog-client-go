package utils

import (
	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
)

type ReleaseNotesSyntax string

const (
	Markdown  ReleaseNotesSyntax = "markdown"
	Asciidoc  ReleaseNotesSyntax = "asciidoc"
	PlainText ReleaseNotesSyntax = "plain_text"
)

type ReleaseBundleParams struct {
	SpecFiles          []*rtUtils.CommonParams
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

		// Create path mapping
		pathMappings := distribution.CreatePathMappings(specFile.Pattern, specFile.Target)

		// Create added properties
		addedProps := createAddedProps(specFile)

		// Append bundle query
		bundleQueries = append(bundleQueries, BundleQuery{Aql: aql, PathMappings: pathMappings, AddedProps: addedProps})
	}

	// Create release bundle struct
	releaseBundleBody := &ReleaseBundleBody{
		DryRun:            dryRun,
		SignImmediately:   &releaseBundleParams.SignImmediately,
		StoringRepository: releaseBundleParams.StoringRepository,
		Description:       releaseBundleParams.Description,
		BundleSpec: BundleSpec{
			Queries: bundleQueries,
		},
	}

	// Add release notes if needed
	if releaseBundleParams.ReleaseNotes != "" {
		releaseBundleBody.ReleaseNotes = &ReleaseNotes{
			Syntax:  releaseBundleParams.ReleaseNotesSyntax,
			Content: releaseBundleParams.ReleaseNotes,
		}
	}
	return releaseBundleBody, nil
}

// Create the AQL query from the input spec
func createAql(specFile *rtUtils.CommonParams) (string, error) {
	if specFile.GetSpecType() != rtUtils.AQL {
		query, err := rtUtils.CreateAqlBodyForSpecWithPattern(specFile)
		if err != nil {
			return "", err
		}
		specFile.Aql = rtUtils.Aql{ItemsFind: query}
	}
	return rtUtils.BuildQueryFromSpecFile(specFile, rtUtils.NONE), nil
}

// Create the AddedProps array from the input TargetProps string
func createAddedProps(specFile *rtUtils.CommonParams) []AddedProps {
	props := specFile.TargetProps

	var addedProps []AddedProps
	if props != nil {
		for key, values := range props.ToMap() {
			addedProps = append(addedProps, AddedProps{key, values})
		}
	}
	return addedProps
}

func AddGpgPassphraseHeader(gpgPassphrase string, headers *map[string]string) {
	if gpgPassphrase != "" {
		rtUtils.AddHeader("X-GPG-PASSPHRASE", gpgPassphrase, headers)
	}
}
