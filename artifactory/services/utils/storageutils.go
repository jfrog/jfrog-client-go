package utils

import "encoding/json"

type FolderInfo struct {
	Uri          string               `json:"uri,omitempty"`
	Repo         string               `json:"repo,omitempty"`
	Path         string               `json:"path,omitempty"`
	Created      string               `json:"created,omitempty"`
	CreatedBy    string               `json:"createdBy,omitempty"`
	LastModified string               `json:"lastModified,omitempty"`
	ModifiedBy   string               `json:"modifiedBy,omitempty"`
	LastUpdated  string               `json:"lastUpdated,omitempty"`
	Children     []FolderInfoChildren `json:"children,omitempty"`
}

type FolderInfoChildren struct {
	Uri    string `json:"uri,omitempty"`
	Folder bool   `json:"folder,omitempty"`
}

type FileListParams struct {
	Deep               bool
	Depth              int
	ListFolders        bool
	MetadataTimestamps bool
	IncludeRootPath    bool
}

func NewFileListParams() FileListParams {
	return FileListParams{}
}

type FileListResponse struct {
	Uri     string         `json:"uri,omitempty"`
	Created string         `json:"created,omitempty"`
	Files   []FileListFile `json:"files,omitempty"`
}

type FileListFile struct {
	Uri                string             `json:"uri,omitempty"`
	Size               json.Number        `json:"size,omitempty"`
	LastModified       string             `json:"lastModified,omitempty"`
	Folder             bool               `json:"folder,omitempty"`
	Sha1               string             `json:"sha1,omitempty"`
	Sha2               string             `json:"sha2,omitempty"`
	MetadataTimestamps MetadataTimestamps `json:"mdTimestamps,omitempty"`
}

type MetadataTimestamps struct {
	Properties string `json:"properties,omitempty"`
}

type StorageInfo struct {
	BinariesSummary         `json:"binariesSummary,omitempty"`
	RepositoriesSummaryList []RepositorySummary `json:"repositoriesSummaryList,omitempty"`
	FileStoreSummary        `json:"fileStoreSummary,omitempty"`
}

type BinariesSummary struct {
	BinariesCount  string `json:"binariesCount,omitempty"`
	BinariesSize   string `json:"binariesSize,omitempty"`
	ArtifactsSize  string `json:"artifactsSize,omitempty"`
	Optimization   string `json:"optimization,omitempty"`
	ItemsCount     string `json:"itemsCount,omitempty"`
	ArtifactsCount string `json:"artifactsCount,omitempty"`
}

type RepositorySummary struct {
	RepoKey          string      `json:"repoKey,omitempty"`
	RepoType         string      `json:"repoType,omitempty"`
	FoldersCount     json.Number `json:"foldersCount,omitempty"`
	FilesCount       json.Number `json:"filesCount,omitempty"`
	UsedSpace        string      `json:"usedSpace,omitempty"`
	UsedSpaceInBytes json.Number `json:"usedSpaceInBytes,omitempty"`
	ItemsCount       json.Number `json:"itemsCount,omitempty"`
	PackageType      string      `json:"packageType,omitempty"`
	ProjectKey       string      `json:"projectKey,omitempty"`
	Percentage       string      `json:"percentage,omitempty"`
}

type FileStoreSummary struct {
	StorageType      string `json:"storageType,omitempty"`
	StorageDirectory string `json:"storageDirectory,omitempty"`
	TotalSpace       string `json:"totalSpace,omitempty"`
	UsedSpace        string `json:"usedSpace,omitempty"`
	FreeSpace        string `json:"freeSpace,omitempty"`
}
