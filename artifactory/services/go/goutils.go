package _go

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"strings"
)

func CreateUrlPath(pathInArtifactory, props string, goApiUrl *string) error {
	*goApiUrl += pathInArtifactory
	properties, err := utils.ParseProperties(props)
	if err != nil {
		return err
	}

	*goApiUrl = strings.Join([]string{*goApiUrl, properties.ToEncodedString(true)}, ";")
	if strings.HasSuffix(*goApiUrl, ";") {
		tempUrl := *goApiUrl
		tempUrl = tempUrl[:len(tempUrl)-1]
		*goApiUrl = tempUrl
	}
	return nil
}
