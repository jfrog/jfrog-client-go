package utils

import (
	"strings"

	clientutils "github.com/jfrog/jfrog-client-go/utils"
)

const (
	WILDCARD SpecType = "wildcard"
	AQL      SpecType = "aql"
	BUILD    SpecType = "build"
)

type SpecType string

type Aql struct {
	ItemsFind string `json:"items.find"`
}

type PathMapping struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type CommonParams struct {
	Aql              Aql
	PathMapping      PathMapping
	Pattern          string
	Exclusions       []string
	Target           string
	Props            string
	TargetProps      *Properties
	ExcludeProps     string
	SortOrder        string
	SortBy           []string
	Offset           int
	Limit            int
	Build            string
	Project          string
	ExcludeArtifacts bool
	IncludeDeps      bool
	Bundle           string
	Recursive        bool
	IncludeDirs      bool
	Regexp           bool
	Ant              bool
	ArchiveEntries   string
	Transitive       bool
	Include          []string
}

func (params CommonParams) GetArchiveEntries() string {
	return params.ArchiveEntries
}

func (params *CommonParams) SetArchiveEntries(archiveEntries string) {
	params.ArchiveEntries = archiveEntries
}

func (params *CommonParams) GetPattern() string {
	return params.Pattern
}

func (params *CommonParams) SetPattern(pattern string) {
	params.Pattern = pattern
}

func (params *CommonParams) SetTarget(target string) {
	params.Target = target
}

func (params *CommonParams) GetTarget() string {
	return params.Target
}

func (params *CommonParams) GetProps() string {
	return params.Props
}

func (params *CommonParams) GetTargetProps() *Properties {
	return params.TargetProps
}

func (params *CommonParams) GetExcludeProps() string {
	return params.ExcludeProps
}

func (params *CommonParams) IsRecursive() bool {
	return params.Recursive
}

func (params *CommonParams) GetPatternType() clientutils.PatternType {
	return clientutils.GetPatternType(clientutils.PatternTypes{RegExp: params.Regexp, Ant: params.Ant})
}

func (params *CommonParams) GetAql() Aql {
	return params.Aql
}

func (params *CommonParams) GetBuild() string {
	return params.Build
}

func (params *CommonParams) GetProject() string {
	return params.Project
}

func (params *CommonParams) GetBundle() string {
	return params.Bundle
}

func (params CommonParams) IsIncludeDirs() bool {
	return params.IncludeDirs
}

func (params *CommonParams) SetProps(props string) {
	params.Props = props
}

func (params *CommonParams) SetTargetProps(targetProps *Properties) {
	params.TargetProps = targetProps
}

func (params *CommonParams) SetExcludeProps(excludeProps string) {
	params.ExcludeProps = excludeProps
}

func (params *CommonParams) GetSortBy() []string {
	return params.SortBy
}

func (params *CommonParams) GetSortOrder() string {
	return params.SortOrder
}

func (params *CommonParams) GetOffset() int {
	return params.Offset
}

func (params *CommonParams) GetLimit() int {
	return params.Limit
}

func (params *CommonParams) GetExclusions() []string {
	return params.Exclusions
}

func (aql *Aql) UnmarshalJSON(value []byte) error {
	str := string(value)
	first := strings.Index(str[strings.Index(str, "{")+1:], "{")
	last := strings.LastIndex(str, "}")

	aql.ItemsFind = str[first+1 : last]
	return nil
}

func (params CommonParams) GetSpecType() (specType SpecType) {
	hasNonTrivialPattern := params.Pattern != "" && params.Pattern != "*"
	// When a non-trivial pattern is combined with features that the BUILD path's AQL
	// does not support, fall through to WILDCARD which handles them correctly:
	//   - IncludeDeps: local deps (added from filesystem) need SHA1-based post-filtering.
	//   - Props/ExcludeProps: the BUILD path's AQL does not include property filters.
	patternRequiresWildcard := hasNonTrivialPattern && (params.IncludeDeps || params.Props != "" || params.ExcludeProps != "")
	switch {
	case params.Build != "" && params.Aql.ItemsFind == "" && !patternRequiresWildcard:
		// When a build is specified, use the BUILD path. The BUILD path uses the
		// dedicated build-artifacts API (fast, no AQL JOINs) and handles aggregated
		// builds. If a pattern is also specified, results are post-filtered by the
		// pattern in SearchBySpecWithBuild.
		specType = BUILD
	case params.Aql.ItemsFind != "":
		specType = AQL
	default:
		specType = WILDCARD
	}
	return specType
}
