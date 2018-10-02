package utils

import (
	"errors"
	"github.com/jfrog/jfrog-cli-go/jfrog-client/utils/errorutils"
	"net/url"
	"strings"
	"fmt"
	"encoding/base64"
)

type PropertyParseOptions int

const (
	// Parsing properties
	SplitCommas PropertyParseOptions = iota
	JoinCommas
)

type Properties struct {
	Properties []Property
}

type Property struct {
	Key   string
	Value string
}

// Parsing properties string to Properties struct.
func ParseProperties(propStr string, option PropertyParseOptions) (*Properties, error) {
	props := &Properties{}
	propList := strings.Split(propStr, ";")
	for _, prop := range propList {
		if prop == "" {
			continue
		}

		key, values, err := splitProp(prop)
		if err != nil {
			return props, err
		}

		switch option {
		case SplitCommas:
			for _, val := range strings.Split(values, ",") {
				props.Properties = append(props.Properties, Property{key, val})
			}
		case JoinCommas:
			props.Properties = append(props.Properties, Property{key, values})
		}
	}
	return props, nil
}

func (props *Properties) ToEncodedString() string {
	encodedProps := ""
	for _, v := range props.Properties {
		jointProp := fmt.Sprintf("%s=%s", url.QueryEscape(v.Key), url.QueryEscape(v.Value))
		encodedProps = fmt.Sprintf("%s;%s", encodedProps, jointProp)
	}
	// Remove leading semicolon
	if strings.HasPrefix(encodedProps, ";") {
		return encodedProps[1:]
	}
	return encodedProps
}

func (props *Properties) ToHeadersMap() map[string]string {
	headers := map[string]string{}
	for _, v := range props.Properties {
		headers[v.Key] = base64.StdEncoding.EncodeToString([]byte(v.Value))
	}
	return headers
}

// Split properties string of format key=value to key value strings
func splitProp(prop string) (string, string, error) {
	splitIndex := strings.Index(prop, "=")
	if splitIndex < 1 || len(prop[splitIndex+1:]) < 1 {
		err := errorutils.CheckError(errors.New("Invalid property: " + prop))
		return "", "", err
	}
	return prop[:splitIndex], prop[splitIndex+1:], nil
}
