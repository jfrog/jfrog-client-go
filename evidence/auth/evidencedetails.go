package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
)

func NewEvidenceDetails() auth.ServiceDetails {
	return &evidenceDetails{}
}

type evidenceDetails struct {
	auth.CommonConfigFields
}

func (rt *evidenceDetails) GetVersion() (string, error) {
	panic("Failed: Method is not implemented")
}
