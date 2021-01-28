package utils

import (
	"regexp"

	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils"
)

type ReleaseNotesSyntax string

const (
	Markdown  ReleaseNotesSyntax = "markdown"
	Asciidoc                     = "asciidoc"
	PlainText                    = "plain_text"
)

var fileSpecCaptureGroup = regexp.MustCompile("({\\d})")

type ReleaseBundleParams struct {
	SpecFiles          []*rtUtils.ArtifactoryCommonParams
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
		pathMappings := createPathMappings(specFile)

		// Create added properties
		addedProps, err := createAddedProps(specFile)
		if err != nil {
			return nil, err
		}

		// Append bundle query
		bundleQueries = append(bundleQueries, BundleQuery{Aql: aql, PathMappings: pathMappings, AddedProps: addedProps})
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
func createAql(specFile *rtUtils.ArtifactoryCommonParams) (string, error) {
	if specFile.GetSpecType() != rtUtils.AQL {
		query, err := rtUtils.CreateAqlBodyForSpecWithPattern(specFile)
		if err != nil {
			return "", err
		}
		specFile.Aql = rtUtils.Aql{ItemsFind: query}
	}
	return rtUtils.BuildQueryFromSpecFile(specFile, rtUtils.NONE), nil
}

// Creat the path mapping from the input spec
func createPathMappings(specFile *rtUtils.ArtifactoryCommonParams) []PathMapping {
	if len(specFile.Target) == 0 {
		return []PathMapping{}
	}

	// Convert the file spec pattern and target to match the path mapping input and output specifications, respectfully.
	return []PathMapping{{
		// The file spec pattern is wildcard based. Convert it to Regex:
		Input: utils.PathToRegExp(specFile.Pattern),
		// The file spec target contain placeholders-style matching groups, like {1}.
		// Convert it to REST API's matching groups style, like $1.
		Output: fileSpecCaptureGroup.ReplaceAllStringFunc(specFile.Target, func(s string) string {
			// Remove curly parenthesis and prepend $
			return "$" + s[1:2]
		}),
	}}
}

// Create the AddedProps array from the input TargetProps string
func createAddedProps(specFile *rtUtils.ArtifactoryCommonParams) ([]AddedProps, error) {
	props, err := rtUtils.ParseProperties(specFile.TargetProps, rtUtils.SplitCommas)
	if err != nil {
		return nil, err
	}

	var addedProps []AddedProps
	for key, values := range props.ToMap() {
		addedProps = append(addedProps, AddedProps{key, values})
	}
	return addedProps, nil
}

func AddGpgPassphraseHeader(gpgPassphrase string, headers *map[string]string) {
	if gpgPassphrase != "" {
		rtUtils.AddHeader("X-GPG-PASSPHRASE", gpgPassphrase, headers)
	}
}
