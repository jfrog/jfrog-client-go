package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
)

func NewJfConnectDetails() auth.ServiceDetails {
	return &jfConnectDetails{}
}

type jfConnectDetails struct {
	auth.CommonConfigFields
}

func (jc *jfConnectDetails) GetVersion() (string, error) {
	panic("Failed: Method is not implemented")
}
