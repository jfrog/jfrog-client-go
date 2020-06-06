package utils

type MultipleReplicationBody struct {
	RepoKey                string            `json:"repoKey"`
	EnableEventReplication bool              `json:"enableEventReplication"`
	CronExp                string            `json:"cronExp"`
	Replications           []ReplicationBody `json:"replications"`
}

type MultipleReplicationParams struct {
	RepoKey                string
	CronExp                string
	EnableEventReplication bool
	Replications           []ReplicationBody
}

func CreateMultipleReplicationBody(params MultipleReplicationParams) *MultipleReplicationBody {
	return &MultipleReplicationBody{
		CronExp:                params.CronExp,
		RepoKey:                params.RepoKey,
		EnableEventReplication: params.EnableEventReplication,
		Replications:           params.Replications,
	}
}
