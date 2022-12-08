package services

import "time"

type PipelineRunStatusResponse struct {
	TotalCount int         `json:"totalCount,omitempty"`
	Pipelines  []Pipelines `json:"pipelines,omitempty"`
}

type StaticPropertyBag struct {
	TriggeredByUserName    string `json:"triggeredByUserName,omitempty"`
	SignedPipelinesEnabled bool   `json:"signedPipelinesEnabled,omitempty"`
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
