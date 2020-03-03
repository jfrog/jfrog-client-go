package utils

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

func SetGpgPassphrase(gpgPassphrase string, headers *map[string]string) {
	if len(gpgPassphrase) > 0 {
		utils.AddHeader("X-GPG-PASSPHRASE", gpgPassphrase, headers)
	}
}
