package services

type VisibilityMetric interface {
	Value() int
	MetricsName() string
}

type Metric struct {
	Value int    `json:"value"`
	Name  string `json:"metrics_name"`
}
