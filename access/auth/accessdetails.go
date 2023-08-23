package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
)

func NewAccessDetails() auth.ServiceDetails {
	return &accessDetails{}
}

type accessDetails struct {
	auth.CommonConfigFields
}

func (rt *accessDetails) GetXscUrl() string {
	panic("Failed: Method is not implemented")
}

func (rt *accessDetails) SetXscUrl(url string) {
	panic("Failed: Method is not implemented")
}

func (rt *accessDetails) GetVersion() (string, error) {
	panic("Failed: Method is not implemented")
}

func (rt *accessDetails) GetXscVersion() string {
	panic("Failed: Method is not implemented")
}
