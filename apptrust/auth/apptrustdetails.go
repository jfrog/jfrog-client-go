package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
)

func NewApptrustDetails() auth.ServiceDetails {
	return &apptrustDetails{}
}

type apptrustDetails struct {
	auth.CommonConfigFields
}

func (rt *apptrustDetails) GetVersion() (string, error) {
	panic("Failed: Method is not implemented")
}
