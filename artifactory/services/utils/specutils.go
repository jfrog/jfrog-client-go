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

type FileGetter interface {
	GetAql() Aql
	GetPattern() string
	SetPattern(pattern string)
	GetExclusions() []string
	GetTarget() string
	SetTarget(target string)
	IsExplode() bool
	GetProps() string
	GetSortOrder() string
	GetSortBy() []string
	GetOffset() int
	GetLimit() int
	GetBuild() string
	GetProject() string
	GetBundle() string
	GetSpecType() (specType SpecType)
	IsRecursive() bool
	IsIncludeDirs() bool
	GetArchiveEntries() string
	SetArchiveEntries(archiveEntries string)
	GetPatternType() clientutils.PatternType
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

func (params *CommonParams) IsExplode() bool {
	return params.Recursive
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
	switch {
	case params.Build != "" && params.Aql.ItemsFind == "" && (params.Pattern == "*" || params.Pattern == ""):
		specType = BUILD
	case params.Aql.ItemsFind != "":
		specType = AQL
	default:
		specType = WILDCARD
	}
	return specType
}
