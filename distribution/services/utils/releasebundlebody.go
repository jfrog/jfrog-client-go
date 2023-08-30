package utils

import "github.com/jfrog/jfrog-client-go/utils/distribution"

// REST body for create and update a release bundle
type ReleaseBundleBody struct {
	DryRun            bool          `json:"dry_run"`
	SignImmediately   *bool         `json:"sign_immediately,omitempty"`
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
	QueryName    string                     `json:"query_name,omitempty"`
	Aql          string                     `json:"aql,omitempty"`
	PathMappings []distribution.PathMapping `json:"mappings,omitempty"`
	AddedProps   []AddedProps               `json:"added_props,omitempty"`
}

type AddedProps struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}
