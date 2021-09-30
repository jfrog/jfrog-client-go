package auth

import (
	"errors"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

func NewAccessDetails() auth.ServiceDetails {
	return &accessDetails{}
}

type accessDetails struct {
	auth.CommonConfigFields
}

func (rt *accessDetails) GetVersion() (string, error) {
	return "", errorutils.CheckError(errors.New("failed: Method is not implemented. Access has no separate version API"))
}
