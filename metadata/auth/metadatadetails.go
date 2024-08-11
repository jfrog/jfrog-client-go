package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
)

type metadataDetails struct {
	auth.CommonConfigFields
}

func NewMetadataDetails() auth.ServiceDetails {
	return &metadataDetails{}
}

func (rt *metadataDetails) GetVersion() (string, error) {
	panic("Failed: Method is not implemented")
}
