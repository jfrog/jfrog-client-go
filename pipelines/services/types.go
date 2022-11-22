package services

import "time"

type PipResponse struct {
	TotalCount int         `json:"totalCount"`
	Pipelines  []Pipelines `json:"pipelines"`
}
type StaticPropertyBag struct {
	TriggeredByUserName    string `json:"triggeredByUserName"`
	SignedPipelinesEnabled bool   `json:"signedPipelinesEnabled"`
}
type Run struct {
	ID                int               `json:"id"`
	RunNumber         int               `json:"runNumber"`
	CreatedAt         time.Time         `json:"createdAt"`
	StartedAt         time.Time         `json:"startedAt"`
	EndedAt           time.Time         `json:"endedAt"`
	StatusCode        int               `json:"statusCode"`
	StaticPropertyBag StaticPropertyBag `json:"staticPropertyBag"`
	DurationSeconds   int               `json:"durationSeconds"`
}
type Pipelines struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	LatestCompletedRunID int    `json:"latestCompletedRunId"`
	PipelineSourceBranch string `json:"pipelineSourceBranch"`
	ProjectID            int    `json:"projectId"`
	LatestRunID          int    `json:"latestRunId"`
	PipelineSourceID     int    `json:"pipelineSourceId"`
	Run                  Run    `json:"run"`
}
