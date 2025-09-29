package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
)

func NewJPDDetails() auth.ServiceDetails {
	return &jpdDetails{}
}

type jpdDetails struct {
	auth.CommonConfigFields
}

func (rt *jpdDetails) GetVersion() (string, error) {
	panic("Failed: Method is not implemented")
}
