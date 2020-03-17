package utils

type ReplicationParams struct {
	Username               string `json:"username"`
	Password               string `json:"password"`
	URL                    string `json:"url"`
	CronExp                string `json:"cronExp"`
	RepoKey                string `json:"repoKey"`
	EnableEventReplication bool   `json:"enableEventReplication"`
	SocketTimeoutMillis    int    `json:"socketTimeoutMillis"`
	Enabled                bool   `json:"enabled"`
	SyncDeletes            bool   `json:"syncDeletes"`
	SyncProperties         bool   `json:"syncProperties"`
	SyncStatistics         bool   `json:"syncStatistics"`
	PathPrefix             string `json:"pathPrefix"`
}
