package services

import "time"

type PipelineRunStatusResponse struct {
	TotalCount int         `json:"totalCount,omitempty"`
	Pipelines  []Pipelines `json:"pipelines,omitempty"`
}

type StaticPropertyBag struct {
	TriggeredByUserName    string `json:"triggeredByUserName,omitempty"`
	SignedPipelinesEnabled *bool  `json:"signedPipelinesEnabled,omitempty"`
}

type Run struct {
	ID                int               `json:"id,omitempty"`
	RunNumber         int               `json:"runNumber,omitempty"`
	CreatedAt         time.Time         `json:"createdAt,omitempty"`
	StartedAt         time.Time         `json:"startedAt,omitempty"`
	EndedAt           time.Time         `json:"endedAt,omitempty"`
	StatusCode        int               `json:"statusCode,omitempty"`
	StaticPropertyBag StaticPropertyBag `json:"staticPropertyBag,omitempty"`
	DurationSeconds   int               `json:"durationSeconds,omitempty"`
}

type Pipelines struct {
	ID                   int    `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	LatestCompletedRunID int    `json:"latestCompletedRunId,omitempty"`
	PipelineSourceBranch string `json:"pipelineSourceBranch,omitempty"`
	ProjectID            int    `json:"projectId,omitempty"`
	LatestRunID          int    `json:"latestRunId,omitempty"`
	PipelineSourceID     int    `json:"pipelineSourceId,omitempty"`
	Run                  Run    `json:"run,omitempty"`
}

type PipelineSyncStatus struct {
	CommitData                   CommitData `json:"commitData,omitempty"`
	ID                           int        `json:"id,omitempty"`
	ProjectID                    int        `json:"projectId,omitempty"`
	PipelineSourceID             int        `json:"pipelineSourceId,omitempty"`
	PipelineSourceBranch         string     `json:"pipelineSourceBranch,omitempty"`
	IsSyncing                    *bool      `json:"isSyncing,omitempty"`
	LastSyncStatusCode           int        `json:"lastSyncStatusCode,omitempty"`
	LastSyncStartedAt            time.Time  `json:"lastSyncStartedAt,omitempty"`
	LastSyncEndedAt              time.Time  `json:"lastSyncEndedAt,omitempty"`
	LastSyncLogs                 string     `json:"lastSyncLogs,omitempty"`
	IsMissingConfig              *bool      `json:"isMissingConfig,omitempty"`
	TriggeredByResourceVersionID int        `json:"triggeredByResourceVersionId,omitempty"`
	CreatedAt                    time.Time  `json:"createdAt,omitempty"`
	UpdatedAt                    time.Time  `json:"updatedAt,omitempty"`
}

type CommitData struct {
	CommitSha string `json:"commitSha,omitempty"`
	Committer string `json:"committer,omitempty"`
	CommitMsg string `json:"commitMsg,omitempty"`
	Source    string `json:"source,omitempty"`
}

type PipelineResources struct {
	ValuesYmlPropertyBag     interface{} `json:"valuesYmlPropertyBag,omitempty"`
	PipelinesYmlPropertyBag  interface{} `json:"pipelinesYmlPropertyBag,omitempty"`
	ID                       int         `json:"id,omitempty"`
	Name                     interface{} `json:"name,omitempty"`
	ProjectID                int         `json:"projectId,omitempty"`
	ProjectIntegrationID     int         `json:"projectIntegrationId,omitempty"`
	RepositoryFullName       string      `json:"repositoryFullName,omitempty"`
	ConfigFolder             string      `json:"configFolder,omitempty"`
	IsMultiBranch            *bool       `json:"isMultiBranch,omitempty"`
	IsInternal               interface{} `json:"isInternal,omitempty"`
	Branch                   interface{} `json:"branch,omitempty"`
	BranchExcludePattern     string      `json:"branchExcludePattern,omitempty"`
	BranchIncludePattern     string      `json:"branchIncludePattern,omitempty"`
	FileFilter               string      `json:"fileFilter,omitempty"`
	Environments             interface{} `json:"environments,omitempty"`
	IsSyncing                *bool       `json:"isSyncing,omitempty"`
	LastSyncStatusCode       int         `json:"lastSyncStatusCode,omitempty"`
	LastSyncStartedAt        time.Time   `json:"lastSyncStartedAt,omitempty"`
	LastSyncEndedAt          time.Time   `json:"lastSyncEndedAt,omitempty"`
	LastSyncLogs             string      `json:"lastSyncLogs,omitempty"`
	ResourceID               int         `json:"resourceId,omitempty"`
	CreatedBy                int         `json:"createdBy,omitempty"`
	UpdatedBy                int         `json:"updatedBy,omitempty"`
	TemplateID               interface{} `json:"templateId,omitempty"`
	AccessResourceNamePrefix string      `json:"accessResourceNamePrefix,omitempty"`
	CreatedAt                time.Time   `json:"createdAt,omitempty"`
	UpdatedAt                time.Time   `json:"updatedAt,omitempty"`
}

type WorkspacesResponse struct {
	ValuesYmlPropertyBag    interface{}             `json:"valuesYmlPropertyBag,omitempty"`
	PipelinesYmlPropertyBag PipelinesYmlPropertyBag `json:"pipelinesYmlPropertyBag,omitempty"`
	ID                      int                     `json:"id,omitempty"`
	Name                    string                  `json:"name,omitempty"`
	ProjectID               int                     `json:"projectId,omitempty"`
	IsSyncing               *bool                   `json:"isSyncing,omitempty"`
	LastSyncStatusCode      int                     `json:"lastSyncStatusCode,omitempty"`
	LastSyncStartedAt       time.Time               `json:"lastSyncStartedAt,omitempty"`
	LastSyncEndedAt         interface{}             `json:"lastSyncEndedAt,omitempty"`
	LastSyncLogs            string                  `json:"lastSyncLogs,omitempty"`
	CreatedBy               int                     `json:"createdBy,omitempty"`
}

type PipelinesYmlPropertyBag struct {
	Resources []Resources `json:"resources,omitempty"`
	Pipelines []Pipelines `json:"pipelines,omitempty"`
}

type Resources struct {
	Configuration Configuration `json:"configuration,omitempty"`
	Name          string        `json:"name,omitempty"`
	Type          string        `json:"type,omitempty"`
}

type Configuration struct {
	Branches    Branches `json:"branches,omitempty"`
	GitProvider string   `json:"gitProvider,omitempty"`
	Path        string   `json:"path,omitempty"`
}

type Branches struct {
	Include string `json:"include,omitempty"`
}

type PipelinesRunID struct {
	LatestRunID   int    `json:"latestRunId,omitempty"`
	Name          string `json:"name,omitempty"`
	SyntaxVersion string `json:"syntaxVersion,omitempty"`
}

type Console struct {
	ConsoleID        string      `json:"consoleId,omitempty"`
	IsSuccess        *bool       `json:"isSuccess,omitempty"`
	IsShown          interface{} `json:"isShown,omitempty"`
	ParentConsoleID  string      `json:"parentConsoleId,omitempty"`
	StepletID        int         `json:"stepletId,omitempty"`
	PipelineID       int         `json:"pipelineId,omitempty"`
	Timestamp        int64       `json:"timestamp,omitempty"`
	TimestampEndedAt interface{} `json:"timestampEndedAt,omitempty"`
	Type             string      `json:"type,omitempty"`
	Message          string      `json:"message,omitempty"`
	CreatedAt        time.Time   `json:"createdAt,omitempty"`
	UpdatedAt        time.Time   `json:"updatedAt,omitempty"`
}

// Validation Types

type ValidationResponse struct {
	IsValid *bool              `json:"isValid,omitempty"`
	Errors  []ValidationErrors `json:"errors,omitempty"`
}

type ValidationErrors struct {
	Text       string `json:"text,omitempty"`
	LineNumber int    `json:"lineNumber,omitempty"`
}
