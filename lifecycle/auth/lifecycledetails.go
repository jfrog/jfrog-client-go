package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
)

func NewLifecycleDetails() auth.ServiceDetails {
	return &lifecycleDetails{}
}

type lifecycleDetails struct {
	auth.CommonConfigFields
}

func (rt *lifecycleDetails) GetVersion() (string, error) {
	panic("Failed: Method is not implemented")
}
