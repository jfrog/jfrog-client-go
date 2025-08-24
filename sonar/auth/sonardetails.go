package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
)

type sonarDetails struct {
	auth.CommonConfigFields
}

func NewSonarDetails() auth.ServiceDetails {
	return &sonarDetails{}
}

func (sd *sonarDetails) GetVersion() (string, error) {
	panic("Failed: Method is not implemented")
}
