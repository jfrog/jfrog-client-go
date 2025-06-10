package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
)

type onemodelDetails struct {
	auth.CommonConfigFields
}

func NewOnemodelDetails() auth.ServiceDetails {
	return &onemodelDetails{}
}

func (rt *onemodelDetails) GetVersion() (string, error) {
	panic("Failed: Method is not implemented")
}
