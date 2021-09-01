package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/url"
	"strings"
)

const (
	propsSeparator     = ";"
	keyValuesSeparator = "="
	valuesSeparator    = ","
)

type Properties struct {
	properties map[string][]string
}

type Property struct {
	Key   string
	Value string
}

func NewProperties() *Properties {
	return &Properties{properties: make(map[string][]string)}
}

// Parsing properties string to Properties struct.
func ParseProperties(propStr string) (*Properties, error) {
	props := NewProperties()
	err := props.ParseAndAddProperties(propStr)
	return props, err
}

func (props *Properties) ParseAndAddProperties(propStr string) error {
	propList := splitWhileIgnoringBackslashPrefixSeparators(propStr, propsSeparator)
	for _, prop := range propList {
		if prop == "" {
			continue
		}

		parts := splitWhileIgnoringBackslashPrefixSeparators(prop, keyValuesSeparator)
		if len(parts) != 2 {
			return errorutils.CheckError(errors.New(fmt.Sprintf("Invalid property format: %s - format should be key=val1,val2,...", prop)))
		}

		key := parts[0]
		values := splitWhileIgnoringBackslashPrefixSeparators(parts[1], valuesSeparator)
		for _, val := range values {
			props.properties[key] = append(props.properties[key], val)
		}
	}
	props.removeDuplicateValues()
	return nil
}

// Split slices s into all substrings separated by sep and returns a slice of the substrings between those separators,
// ignoring separators with a '\' prefix, which indicates that the separator is a part of the value.
func splitWhileIgnoringBackslashPrefixSeparators(str, separator string) (splitArray []string) {
	values := strings.Split(str, separator)
	for i, val := range values {
		// Let's use separator=',' for example:
		// If "\" is found, then it means that the original string contains the "\," which indicates this "," is part
		// of the value and not a separator.
		if strings.HasSuffix(val, "\\") && i+1 < len(values) {
			values[i+1] = val[:len(val)-1] + separator + values[i+1]
		} else {
			splitArray = append(splitArray, val)
		}
	}
	return
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

// Creates a string of the properties, ready to use in a URL.
// If concatValues is true, then the values of each property are concatenated together separated by a comma. For example: key=val1,val2,...
// Otherwise, each value of the property will be written with its key separately. For example: key=val1;key=val2;...
func (props *Properties) ToEncodedString(concatValues bool) string {
	encodedProps := ""
	for key, values := range props.properties {
		var jointProp string

		if concatValues {
			jointProp = fmt.Sprintf("%s=", url.QueryEscape(key))
		}
		for _, value := range values {
			if concatValues {
				propValue := strings.Replace(value, valuesSeparator, fmt.Sprintf("\\%s", valuesSeparator), -1)
				jointProp = fmt.Sprintf("%s%s%s", jointProp, url.QueryEscape(propValue), url.QueryEscape(valuesSeparator))
			} else {
				jointProp = fmt.Sprintf("%s%s=%s%s", jointProp, url.QueryEscape(key), url.QueryEscape(value), propsSeparator)
			}
		}
		// Trim the last comma/semicolon
		if concatValues {
			jointProp = strings.TrimSuffix(jointProp, url.QueryEscape(valuesSeparator))
		} else {
			jointProp = strings.TrimSuffix(jointProp, propsSeparator)
		}

		encodedProps = fmt.Sprintf("%s%s%s", encodedProps, propsSeparator, jointProp)
	}
	// Remove leading semicolon and return
	return strings.TrimPrefix(encodedProps, propsSeparator)
}

func (props *Properties) ToHeadersMap() map[string]string {
	headers := map[string]string{}
	for key, values := range props.properties {
		headers[key] = base64.StdEncoding.EncodeToString([]byte(strings.Join(values, valuesSeparator)))
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

// Merges multiple Properties structs into one and removes duplicate values
func MergeProperties(properties []*Properties) *Properties {
	mergedProps := NewProperties()
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
