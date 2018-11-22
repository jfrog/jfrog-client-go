package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var userAgent = getDefaultUserAgent()

func getVersion() string {
	return "0.1.0"
}

func GetUserAgent() string {
	return userAgent
}

func SetUserAgent(newUserAgent string) {
	userAgent = newUserAgent
}

func getDefaultUserAgent() string {
	return fmt.Sprintf("jfrog-client-go/%s", getVersion())
}

// Get the local root path, from which to start collecting artifacts to be used for:
// 1. Uploaded to Artifactory,
// 2. Adding to the local build-info, to be later published to Artifactory.
func GetRootPath(path string, useRegExp bool) string {
	// The first step is to split the local path pattern into sections, by the file seperator.
	seperator := "/"
	sections := strings.Split(path, seperator)
	if len(sections) == 1 {
		seperator = "\\"
		sections = strings.Split(path, seperator)
	}

	// Now we start building the root path, making sure to leave out the sub-directory that includes the pattern.
	rootPath := ""
	for _, section := range sections {
		if section == "" {
			continue
		}
		if useRegExp {
			if strings.Index(section, "(") != -1 {
				break
			}
		} else {
			if strings.Index(section, "*") != -1 {
				break
			}
		}
		if rootPath != "" {
			rootPath += seperator
		}
		if section == "~" {
			rootPath += GetUserHomeDir()
		} else {
			rootPath += section
		}
	}
	if len(sections) > 0 && sections[0] == "" {
		rootPath = seperator + rootPath
	}
	if rootPath == "" {
		return "."
	}
	return rootPath
}

func StringToBool(boolVal string, defaultValue bool) (bool, error) {
	if len(boolVal) > 0 {
		result, err := strconv.ParseBool(boolVal)
		errorutils.CheckError(err)
		return result, err
	}
	return defaultValue, nil
}

func AddTrailingSlashIfNeeded(url string) string {
	if url != "" && !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url
}

func IndentJson(jsonStr []byte) string {
	var content bytes.Buffer
	err := json.Indent(&content, jsonStr, "", "  ")
	if err == nil {
		return content.String()
	}
	return string(jsonStr)
}

func MergeMaps(src map[string]string, dst map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}

func CopyMap(src map[string]string) (dst map[string]string) {
	if dst == nil {
		dst = make(map[string]string)
	}
	for k, v := range src {
		dst[k] = v
	}
	return
}

func PrepareLocalPathForUpload(localPath string, useRegExp bool) string {
	if localPath == "./" || localPath == ".\\" {
		return "^.*$"
	}
	if strings.HasPrefix(localPath, "./") {
		localPath = localPath[2:]
	} else if strings.HasPrefix(localPath, ".\\") {
		localPath = localPath[3:]
	}
	if !useRegExp {
		localPath = pathToRegExp(localPath)
	}
	return localPath
}

func pathToRegExp(localPath string) string {
	var SPECIAL_CHARS = []string{".", "^", "$", "+"}
	for _, char := range SPECIAL_CHARS {
		localPath = strings.Replace(localPath, char, "\\"+char, -1)
	}
	var wildcard = ".*"
	localPath = strings.Replace(localPath, "*", wildcard, -1)
	if strings.HasSuffix(localPath, "/") || strings.HasSuffix(localPath, "\\") {
		localPath += wildcard
	}
	localPath = "^" + localPath + "$"
	return localPath
}

// Replaces matched regular expression from path to corresponding {i} at target.
// Example 1:
//      pattern = "repoA/1(.*)234" ; path = "repoA/1hello234" ; target = "{1}" ; ignoreRepo = false
//      returns "hello"
// Example 2:
//      pattern = "repoA/1(.*)234" ; path = "repoB/1hello234" ; target = "{1}" ; ignoreRepo = true
//      returns "hello"
func BuildTargetPath(pattern, path, target string, ignoreRepo bool) (string, error) {
	if ignoreRepo {
		pattern = removeRepoFromPath(pattern)
		path = removeRepoFromPath(path)
	}
	pattern = pathToRegExp(pattern)
	r, err := regexp.Compile(pattern)
	err = errorutils.CheckError(err)
	if err != nil {
		return "", err
	}

	groups := r.FindStringSubmatch(path)
	size := len(groups)
	if size > 0 {
		for i := 1; i < size; i++ {
			group := strings.Replace(groups[i], "\\", "/", -1)
			target = strings.Replace(target, "{"+strconv.Itoa(i)+"}", group, -1)
		}
	}
	return target, nil
}

func GetLogMsgPrefix(threadId int, dryRun bool) string {
	var strDryRun string
	if dryRun {
		strDryRun = "[Dry run] "
	}
	return "[Thread " + strconv.Itoa(threadId) + "] " + strDryRun
}

func TrimPath(path string) string {
	path = strings.Replace(path, "\\", "/", -1)
	path = strings.Replace(path, "//", "/", -1)
	path = strings.Replace(path, "../", "", -1)
	path = strings.Replace(path, "./", "", -1)
	return path
}

func Bool2Int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ReplaceTildeWithUserHome(path string) string {
	if len(path) > 1 && path[0:1] == "~" {
		return GetUserHomeDir() + path[1:]
	}
	return path
}

func GetUserHomeDir() string {
	if IsWindows() {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return strings.Replace(home, "\\", "\\\\", -1)
	}
	return os.Getenv("HOME")
}

func GetBoolEnvValue(flagName string, defValue bool) (bool, error) {
	envVarValue := os.Getenv(flagName)
	if envVarValue == "" {
		return defValue, nil
	}
	val, err := strconv.ParseBool(envVarValue)
	err = CheckErrorWithMessage(err, "can't parse environment variable "+flagName)
	return val, err
}

func CheckErrorWithMessage(err error, message string) error {
	if err != nil {
		log.Error(message)
		err = errorutils.CheckError(err)
	}
	return err
}

func ConvertSliceToMap(slice []string) map[string]bool {
	mapFromSlice := make(map[string]bool)
	for _, value := range slice {
		mapFromSlice[value] = true
	}
	return mapFromSlice
}

func removeRepoFromPath(path string) string {
	if idx := strings.Index(path, "/"); idx != -1 {
		return path[idx:]
	}
	return path
}

// Split str by the provided separator, escaping the separator if it is prefixed by a back-slash.
func SplitWithEscape(str string, separator rune) []string {
	var parts []string
	var current bytes.Buffer
	escaped := false
	for _, char := range str {
		if char == '\\' {
			if escaped {
				current.WriteRune(char)
			}
			escaped = true
		} else if char == separator && !escaped {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			escaped = false
			current.WriteRune(char)
		}
	}
	parts = append(parts, current.String())
	return parts
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

type Artifact struct {
	LocalPath  string
	TargetPath string
	Symlink    string
}
