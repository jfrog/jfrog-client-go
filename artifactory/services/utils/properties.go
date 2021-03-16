package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/url"
	"strings"
)

const PropsSeparator = ";"
const ValuesSeparator = ","

type Properties struct {
	properties map[string][]string
}

type Property struct {
	Key   string
	Value string
}

// Parsing properties string to Properties struct.
func ParseProperties(propStr string) (*Properties, error) {
	props := &Properties{properties: make(map[string][]string)}
	err := props.ParseAndAddProperties(propStr)
	return props, err
}

func (props *Properties) ParseAndAddProperties(propStr string) error {
	propList := strings.Split(propStr, PropsSeparator)
	for _, prop := range propList {
		if prop == "" {
			continue
		}

		parts := strings.Split(prop, "=")
		if len(parts) != 2 {
			return errorutils.CheckError(errors.New("Invalid property: " + prop))
		}

		key := parts[0]
		values := strings.Split(parts[1], ValuesSeparator)
		for i, val := range values {
			// If "\" is found, then it means that the original string contains the "\," which indicate this "," is part of the value
			// and not a separator
			if strings.HasSuffix(val, "\\") && i+1 < len(values) {
				values[i+1] = val[:len(val)-1] + ValuesSeparator + values[i+1]
			} else {
				props.properties[key] = append(props.properties[key], val)
			}
		}
	}
	props.removeDuplicateValues()
	return nil
}

func (props *Properties) AddProperty(key, value string) {
	if _, exist := props.properties[key]; exist {
		for _, existValue := range props.properties[key] {
			if existValue == value {
				return
			}
		}
	}
	props.properties[key] = append(props.properties[key], value)
}

func (props *Properties) ToEncodedString(concatValues bool) string {
	encodedProps := ""
	for key, values := range props.properties {
		var jointProp string

		if concatValues {
			jointProp = fmt.Sprintf("%s=", url.QueryEscape(key))
		}
		for _, value := range values {
			if concatValues {
				propValue := strings.Replace(value, ValuesSeparator, fmt.Sprintf("\\%s", ValuesSeparator), -1)
				jointProp = fmt.Sprintf("%s%s%s", jointProp, url.QueryEscape(propValue), url.QueryEscape(ValuesSeparator))
			} else {
				jointProp = fmt.Sprintf("%s%s=%s%s", jointProp, url.QueryEscape(key), url.QueryEscape(value), PropsSeparator)
			}
		}
		// Trim the last comma/semicolon
		if concatValues {
			jointProp = strings.TrimSuffix(jointProp, url.QueryEscape(ValuesSeparator))
		} else {
			jointProp = strings.TrimSuffix(jointProp, PropsSeparator)
		}

		encodedProps = fmt.Sprintf("%s%s%s", encodedProps, PropsSeparator, jointProp)
	}
	// Remove leading semicolon and return
	return strings.TrimPrefix(encodedProps, PropsSeparator)
}

func (props *Properties) ToHeadersMap() map[string]string {
	headers := map[string]string{}
	for key, values := range props.properties {
		headers[key] = base64.StdEncoding.EncodeToString([]byte(strings.Join(values, ValuesSeparator)))
	}
	return headers
}

// Convert properties from Slice to map
func (props *Properties) ToMap() map[string][]string {
	return props.properties
}

func (props *Properties) KeysLen() int {
	return len(props.properties)
}

func (props *Properties) removeDuplicateValues() {
	for key, values := range props.properties {
		props.properties[key] = removeDuplicates(values)
	}
}

func removeDuplicates(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, item := range stringSlice {
		if _, exist := keys[item]; !exist {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func MergeProperties(properties []*Properties) *Properties {
	mergedProps := &Properties{properties: make(map[string][]string)}
	for _, propsStruct := range properties {
		if propsStruct != nil {
			for key, values := range propsStruct.properties {
				mergedProps.properties[key] = append(mergedProps.properties[key], values...)
			}
		}
	}
	mergedProps.removeDuplicateValues()
	return mergedProps
}
