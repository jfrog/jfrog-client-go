package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jfrog/jfrog-client-go/utils/io"

	"github.com/jfrog/build-info-go/entities"
	"github.com/jfrog/gofrog/stringutils"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	Development = "development"
	Agent       = "jfrog-client-go"
	Version     = "1.18.0"
)

// In order to limit the number of items loaded from a reader into the memory, we use a buffers with this size limit.
var MaxBufferSize = 50000

var userAgent = getDefaultUserAgent()

func getVersion() string {
	return Version
}

func GetUserAgent() string {
	return userAgent
}

func SetUserAgent(newUserAgent string) {
	userAgent = newUserAgent
}

func getDefaultUserAgent() string {
	return fmt.Sprintf("%s/%s", Agent, getVersion())
}

// Get the local root path, from which to start collecting artifacts to be used for:
// 1. Uploaded to Artifactory,
// 2. Adding to the local build-info, to be later published to Artifactory.
func GetRootPath(path string, patternType PatternType, parentheses ParenthesesSlice) string {
	// The first step is to split the local path pattern into sections, by the file separator.
	separator := "/"
	sections := strings.Split(path, separator)
	if len(sections) == 1 {
		separator = "\\"
		sections = strings.Split(path, separator)
	}

	// Now we start building the root path, making sure to leave out the sub-directory that includes the pattern.
	rootPath := ""
	for _, section := range sections {
		if section == "" {
			continue
		}
		if patternType == RegExp {
			if strings.Contains(section, "(") {
				break
			}
		} else {
			if strings.Contains(section, "*") {
				break
			}
			if strings.Contains(section, "(") {
				temp := rootPath + section
				if isWildcardParentheses(temp, parentheses) {
					break
				}
			}
			if patternType == AntPattern {
				if strings.Contains(section, "?") {
					break
				}
			}
		}
		if rootPath != "" {
			rootPath += separator
		}
		if section == "~" {
			rootPath += GetUserHomeDir()
		} else {
			rootPath += section
		}
	}
	if len(sections) > 0 && sections[0] == "" {
		rootPath = separator + rootPath
	}
	if rootPath == "" {
		return "."
	}
	return rootPath
}

// Return true if the ‘str’ argument contains open parentasis, that is related to a placeholder.
// The ‘parentheses’ argument contains all the indexes of placeholder parentheses.
func isWildcardParentheses(str string, parentheses ParenthesesSlice) bool {
	toFind := "("
	currStart := 0
	for {
		idx := strings.Index(str, toFind)
		if idx == -1 {
			break
		}
		if parentheses.IsPresent(idx) {
			return true
		}
		currStart += idx + len(toFind)
		str = str[idx+len(toFind):]
	}
	return false
}

func StringToBool(boolVal string, defaultValue bool) (bool, error) {
	if len(boolVal) > 0 {
		result, err := strconv.ParseBool(boolVal)
		return result, errorutils.CheckError(err)
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
	return doIndentJson(jsonStr, "", "  ")
}

func IndentJsonArray(jsonStr []byte) string {
	return doIndentJson(jsonStr, "  ", "  ")
}

func doIndentJson(jsonStr []byte, prefix, indent string) string {
	var content bytes.Buffer
	err := json.Indent(&content, jsonStr, prefix, indent)
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
	dst = make(map[string]string)
	for k, v := range src {
		dst[k] = v
	}
	return
}

func ConvertLocalPatternToRegexp(localPath string, patternType PatternType) string {
	if localPath == "./" || localPath == ".\\" {
		return "^.*$"
	}
	if strings.HasPrefix(localPath, "./") {
		localPath = localPath[2:]
	} else if strings.HasPrefix(localPath, ".\\") {
		localPath = localPath[3:]
	}
	if patternType == AntPattern {
		localPath = antPatternToRegExp(cleanPath(localPath))
	} else if patternType == WildCardPattern {
		localPath = stringutils.WildcardPatternToRegExp(cleanPath(localPath))
	}

	return localPath
}

// Clean /../ | /./ using filepath.Clean.
func cleanPath(path string) string {
	temp := path[len(path)-1:]
	path = filepath.Clean(path)
	if temp == `\` || temp == "/" {
		path += temp
	}
	// Since filepath.Clean replaces \\ with \, we revert this action.
	path = strings.Replace(path, `\`, `\\`, -1)
	return path
}

func antPatternToRegExp(localPath string) string {
	localPath = stringutils.EscapeSpecialChars(localPath)
	separator := getFileSeparator()
	// 'xxx/' => 'xxx/**'
	if strings.HasSuffix(localPath, separator) {
		localPath += "**"
	}
	var wildcard = ".*"
	// ant `*` ~ regexp `([^/]*)` : `*` matches zero or more characters except from `/`.
	var regAsterisk = "([^" + separator + "]*)"
	// ant `\*` ~ regexp `([^/]+)` : `\*` matches one or more characters (except from `/`) with a `/` prefix.
	var regAsteriskWithSeparatorPrefix = "([^" + separator + "]+)"
	// ant `**` ~ regexp `(.*)?` : `**` matches zero or more 'directories' in a path.
	var doubleRegAsterisk = "(" + wildcard + ")?"
	var doubleRegAsteriskWithSeparatorPrefix = "(" + wildcard + separator + ")?"
	var doubleRegAsteriskWithSeparatorSuffix = "(" + separator + wildcard + ")?"

	// `?` => `.{1}` : `?` matches one character.
	localPath = strings.Replace(localPath, `?`, ".{1}", -1)
	// `*` => `([^/]*)`
	localPath = strings.Replace(localPath, `*`, regAsterisk, -1)
	// `**` => `(.*)?`
	localPath = strings.Replace(localPath, regAsterisk+regAsterisk, doubleRegAsterisk, -1)

	// `\([^/]*)` => `\([^/]+)` : there are 2 cases with '*':
	//		1. xxx/x* : * will represent 0 or more characters.
	//		2. xxx/* : * will represent 1 or more characters.
	// This "replace" handles the second option.
	localPath = strings.Replace(localPath, separator+regAsterisk, separator+regAsteriskWithSeparatorPrefix, -1)
	// `(.*)?/` => `(.*/)?`
	localPath = strings.Replace(localPath, doubleRegAsterisk+separator, doubleRegAsteriskWithSeparatorPrefix, -1)
	// Convert the last '/**' in the expression if exists : `/(.*)?` => `(/.*)?`
	if strings.HasSuffix(localPath, separator+doubleRegAsterisk) {
		localPath = strings.TrimSuffix(localPath, separator+doubleRegAsterisk) + doubleRegAsteriskWithSeparatorSuffix
	}
	return "^" + localPath + "$"
}

func getFileSeparator() string {
	if io.IsWindows() {
		return "\\\\"
	}
	return "/"
}

// Replaces matched regular expression from path to corresponding placeholder {i} at target.
// Example 1:
//      pattern = "repoA/1(.*)234" ; path = "repoA/1hello234" ; target = "{1}" ; ignoreRepo = false
//      returns "hello"
// Example 2:
//      pattern = "repoA/1(.*)234" ; path = "repoB/1hello234" ; target = "{1}" ; ignoreRepo = true
//      returns "hello"
// return (parsed target, placeholders replaced in target, error)
func BuildTargetPath(pattern, path, target string, ignoreRepo bool) (string, bool, error) {
	asteriskIndex := strings.Index(pattern, "*")
	slashIndex := strings.Index(pattern, "/")
	if shouldRemoveRepo(ignoreRepo, asteriskIndex, slashIndex) {
		// Removing the repository part of the path is required when working with virtual repositories, as the pattern
		// may contain the virtual-repository name, but the path contains the local-repository name.
		pattern = removeRepoFromPath(pattern)
		path = removeRepoFromPath(path)
	}
	pattern = addEscapingParentheses(pattern, target)
	pattern = stringutils.WildcardPatternToRegExp(pattern)
	if slashIndex < 0 {
		// If '/' doesn't exist, add an optional trailing-slash to support cases in which the provided pattern
		// is only the repository name.
		dollarIndex := strings.LastIndex(pattern, "$")
		pattern = pattern[:dollarIndex]
		pattern += "(/.*)?$"
	}

	r, err := regexp.Compile(pattern)
	err = errorutils.CheckError(err)
	if err != nil {
		return "", false, err
	}

	groups := r.FindStringSubmatch(path)
	if len(groups) > 0 {
		target, replaceOccurred := ReplacePlaceHolders(groups, target)
		return target, replaceOccurred, nil
	}
	return target, false, nil
}

// group - regular expression matched group to replace with placeholders
// toReplace - target pattern to replace
// Return - (parsed placeholders string, placeholders were  replaced)
func ReplacePlaceHolders(groups []string, toReplace string) (string, bool) {
	preReplaced := toReplace
	for i := 1; i < len(groups); i++ {
		group := strings.Replace(groups[i], "\\", "/", -1)
		toReplace = strings.Replace(toReplace, "{"+strconv.Itoa(i)+"}", group, -1)
	}
	replaceOccurred := preReplaced != toReplace
	return toReplace, replaceOccurred
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
	if io.IsWindows() {
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

func shouldRemoveRepo(ignoreRepo bool, asteriskIndex, slashIndex int) bool {
	if !ignoreRepo || slashIndex < 0 {
		return false
	}
	if asteriskIndex < 0 {
		return true
	}
	return IsSlashPrecedeAsterisk(asteriskIndex, slashIndex)
}

func IsSlashPrecedeAsterisk(asteriskIndex, slashIndex int) bool {
	return slashIndex < asteriskIndex && slashIndex >= 0
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

func AddProps(oldProps, additionalProps string) string {
	if len(oldProps) > 0 && !strings.HasSuffix(oldProps, ";") && len(additionalProps) > 0 {
		oldProps += ";"
	}
	return oldProps + additionalProps
}

type Artifact struct {
	LocalPath           string
	TargetPath          string
	SymlinkTargetPath   string
	TargetPathInArchive string
}

const (
	WildCardPattern PatternType = "wildcard"
	RegExp          PatternType = "regexp"
	AntPattern      PatternType = "ant"
)

type PatternType string

type PatternTypes struct {
	RegExp bool
	Ant    bool
}

func GetPatternType(patternTypes PatternTypes) PatternType {
	if patternTypes.RegExp {
		return RegExp
	}
	if patternTypes.Ant {
		return AntPattern
	}
	return WildCardPattern
}

type Sha256Summary struct {
	sha256    string
	succeeded bool
}

func NewSha256Summary() *Sha256Summary {
	return &Sha256Summary{}
}

func (bps *Sha256Summary) IsSucceeded() bool {
	return bps.succeeded
}

func (bps *Sha256Summary) SetSucceeded(succeeded bool) *Sha256Summary {
	bps.succeeded = succeeded
	return bps
}

func (bps *Sha256Summary) GetSha256() string {
	return bps.sha256
}

func (bps *Sha256Summary) SetSha256(sha256 string) *Sha256Summary {
	bps.sha256 = sha256
	return bps
}

// Represents a file transfer from SourcePath to TargetPath.
// Each of the paths can be on the local machine (full or relative) or in Artifactory (without Artifactory URL).
// The file's Sha256 is calculated by Artifactory during the upload. we read the sha256 from the HTTP's response body.
type FileTransferDetails struct {
	SourcePath string `json:"sourcePath,omitempty"`
	TargetPath string `json:"targetPath,omitempty"`
	RtUrl      string `json:"rtUrl,omitempty"`
	Sha256     string `json:"sha256,omitempty"`
}

// Represent deployed artifact's details returned from build-info project for maven and gradle.
type DeployableArtifactDetails struct {
	SourcePath       string `json:"sourcePath,omitempty"`
	ArtifactDest     string `json:"artifactDest,omitempty"`
	Sha256           string `json:"sha256,omitempty"`
	DeploySucceeded  bool   `json:"deploySucceeded,omitempty"`
	TargetRepository string `json:"targetRepository,omitempty"`
}

func (details *DeployableArtifactDetails) CreateFileTransferDetails(rtUrl, targetRepository string) (FileTransferDetails, error) {
	targetUrl, err := url.Parse(path.Join(targetRepository, details.ArtifactDest))
	if err != nil {
		return FileTransferDetails{}, err
	}
	return FileTransferDetails{SourcePath: details.SourcePath, TargetPath: targetUrl.String(), Sha256: details.Sha256, RtUrl: rtUrl}, nil
}

type UploadResponseBody struct {
	Checksums entities.Checksum `json:"checksums,omitempty"`
}

func SaveFileTransferDetailsInTempFile(filesDetails *[]FileTransferDetails) (filePath string, err error) {
	tempFile, err := fileutils.CreateTempFile()
	if err != nil {
		return "", err
	}
	defer func() {
		e := tempFile.Close()
		if err == nil {
			err = errorutils.CheckError(e)
		}
	}()
	filePath = tempFile.Name()
	return filePath, SaveFileTransferDetailsInFile(filePath, filesDetails)
}

func SaveFileTransferDetailsInFile(filePath string, details *[]FileTransferDetails) error {
	// Marshal and save files details to a file.
	// The details will be saved in a json format in an array with key "files" for printing later
	finalResult := struct {
		Files *[]FileTransferDetails `json:"files"`
	}{}
	finalResult.Files = details
	files, err := json.Marshal(finalResult)
	if err != nil {
		return errorutils.CheckError(err)
	}
	return errorutils.CheckError(ioutil.WriteFile(filePath, files, 0700))
}

// Extract sha256 of the uploaded file (calculated by artifactory) from the response's body.
// In case of uploading archive with "--explode" the response body will be empty and sha256 won't be shown at
// the detailed summary.
func ExtractSha256FromResponseBody(body []byte) (string, error) {
	if len(body) > 0 {
		responseBody := new(UploadResponseBody)
		err := json.Unmarshal(body, &responseBody)
		if errorutils.CheckError(err) != nil {
			return "", err
		}
		return responseBody.Checksums.Sha256, nil
	}
	return "", nil
}
