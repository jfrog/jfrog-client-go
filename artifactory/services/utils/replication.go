package utils

type replicationBody struct {
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
	// Deprecated
	PathPrefix               string `json:"pathPrefix"`
	IncludePathPrefixPattern string `json:"includePathPrefixPattern"`
}

type GetReplicationBody struct {
	replicationBody
	ProxyRef string `json:"proxyRef"`
}

type UpdateReplicationBody struct {
	replicationBody
	Proxy string `json:"proxy"`
}

type ReplicationParams struct {
	Username string
	Password string
	Url      string
	CronExp  string
	// Source replication repository.
	RepoKey                string
	Proxy                  string
	EnableEventReplication bool
	SocketTimeoutMillis    int
	Enabled                bool
	SyncDeletes            bool
	SyncProperties         bool
	SyncStatistics         bool
	// Deprecated
	PathPrefix               string
	IncludePathPrefixPattern string
}

func CreateUpdateReplicationBody(params ReplicationParams) *UpdateReplicationBody {
	return &UpdateReplicationBody{
		replicationBody: replicationBody{
			Username:                 params.Username,
			Password:                 params.Password,
			URL:                      params.Url,
			CronExp:                  params.CronExp,
			RepoKey:                  params.RepoKey,
			EnableEventReplication:   params.EnableEventReplication,
			SocketTimeoutMillis:      params.SocketTimeoutMillis,
			Enabled:                  params.Enabled,
			SyncDeletes:              params.SyncDeletes,
			SyncProperties:           params.SyncProperties,
			SyncStatistics:           params.SyncStatistics,
			PathPrefix:               params.PathPrefix,
			IncludePathPrefixPattern: params.IncludePathPrefixPattern,
		},
		Proxy: params.Proxy,
	}
}

func CreateReplicationParams(body GetReplicationBody) *ReplicationParams {
	return &ReplicationParams{
		Username:                 body.Username,
		Password:                 body.Password,
		Url:                      body.URL,
		CronExp:                  body.CronExp,
		RepoKey:                  body.RepoKey,
		Proxy:                    body.ProxyRef,
		EnableEventReplication:   body.EnableEventReplication,
		SocketTimeoutMillis:      body.SocketTimeoutMillis,
		Enabled:                  body.Enabled,
		SyncDeletes:              body.SyncDeletes,
		SyncProperties:           body.SyncProperties,
		SyncStatistics:           body.SyncStatistics,
		PathPrefix:               body.PathPrefix,
		IncludePathPrefixPattern: body.IncludePathPrefixPattern,
	}
}
