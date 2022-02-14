package utils

import "encoding/json"

type ReplicationBody struct {
	Username               string `json:"username"`
	Password               string `json:"password"`
	URL                    string `json:"url"`
	CronExp                string `json:"cronExp"`
	RepoKey                string `json:"repoKey"`
	Proxy                  string `json:"proxy"`
	EnableEventReplication bool   `json:"enableEventReplication"`
	SocketTimeoutMillis    int    `json:"socketTimeoutMillis"`
	Enabled                bool   `json:"enabled"`
	SyncDeletes            bool   `json:"syncDeletes"`
	SyncProperties         bool   `json:"syncProperties"`
	SyncStatistics         bool   `json:"syncStatistics"`
	PathPrefix             string `json:"pathPrefix"`
}

type ReplicationParams struct {
	Username string
	Password string
	Url      string
	CronExp  string
	// Source replication repository.
	RepoKey                  string
	Proxy                    string
	EnableEventReplication   bool
	SocketTimeoutMillis      int
	Enabled                  bool
	SyncDeletes              bool
	SyncProperties           bool
	SyncStatistics           bool
	PathPrefix               string
	IncludePathPrefixPattern string
}

// UnmarshalJSON overrides the default JSON unmarshal function because the POST request to create a replication
// has a field named `proxy` but the GET request returns a JSON with a field named `proxyRef`
func (rp *ReplicationParams) UnmarshalJSON(data []byte) error {
	type Alias ReplicationParams

	aux := &struct {
		ProxyRef string `json:"proxyRef"`
		*Alias
	}{
		Alias: (*Alias)(rp),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	rp.Proxy = aux.ProxyRef

	return nil
}

func CreateReplicationBody(params ReplicationParams) *ReplicationBody {
	return &ReplicationBody{
		Username:               params.Username,
		Password:               params.Password,
		URL:                    params.Url,
		CronExp:                params.CronExp,
		RepoKey:                params.RepoKey,
		Proxy:                  params.Proxy,
		EnableEventReplication: params.EnableEventReplication,
		SocketTimeoutMillis:    params.SocketTimeoutMillis,
		Enabled:                params.Enabled,
		SyncDeletes:            params.SyncDeletes,
		SyncProperties:         params.SyncProperties,
		SyncStatistics:         params.SyncStatistics,
		PathPrefix:             params.PathPrefix,
	}
}
