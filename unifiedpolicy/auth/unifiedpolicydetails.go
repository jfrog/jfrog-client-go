package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
)

func NewUnifiedPolicyDetails() auth.ServiceDetails {
	return &unifiedPolicyDetails{}
}

type unifiedPolicyDetails struct {
	auth.CommonConfigFields
}

func (up *unifiedPolicyDetails) GetVersion() (string, error) {
	return "1.0.0", nil
}
