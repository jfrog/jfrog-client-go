package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// Xray gate status values returned by the Skills xrayStatus endpoint.
const (
	SkillXrayStatusNotInEntitlement = "NOT_IN_ENTITLEMENT"
	SkillXrayStatusDisabledForRepo  = "XRAY_DISABLED_FOR_REPO"
	SkillXrayStatusScanInProgress   = "SCAN_IN_PROGRESS"
	SkillXrayStatusBlocked          = "BLOCKED"
	SkillXrayStatusApproved         = "APPROVED"
)

// SkillXrayStatusResponse is the response from the Skills Xray gate status endpoint.
type SkillXrayStatusResponse struct {
	Status  string `json:"status"`
	RepoKey string `json:"repoKey"`
	Path    string `json:"path"`
}

type SkillsService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewSkillsService(client *jfroghttpclient.JfrogHttpClient) *SkillsService {
	return &SkillsService{client: client}
}

func (ss *SkillsService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return ss.client
}

func (ss *SkillsService) ListVersions(repoKey, slug string) ([]SkillVersion, error) {
	log.Debug(fmt.Sprintf("Listing versions for skill '%s' in repo '%s'...", slug, repoKey))
	body, err := ss.sendGet(repoKey, fmt.Sprintf("skills/%s/versions", slug))
	if err != nil {
		return nil, err
	}

	var wrapper skillVersionsResponse
	if err = json.Unmarshal(body, &wrapper); err != nil {
		return nil, errorutils.CheckErrorf("failed to parse skill versions response: %s", err.Error())
	}
	return wrapper.Items, nil
}

func (ss *SkillsService) ListSkills(repoKey string, limit int, cursor, sortBy string) ([]SkillListItem, string, error) {
	log.Debug(fmt.Sprintf("Listing skills in repo '%s'...", repoKey))
	sort := "updated"
	if sortBy == "downloads" {
		sort = "downloads"
	}
	endpoint := fmt.Sprintf("skills?limit=%d&cursor=%s&sort=%s", limit, url.QueryEscape(cursor), sort)
	body, err := ss.sendGet(repoKey, endpoint)
	if err != nil {
		return nil, "", err
	}

	var wrapper skillListResponse
	if err = json.Unmarshal(body, &wrapper); err != nil {
		return nil, "", errorutils.CheckErrorf("failed to parse skill list response: %s", err.Error())
	}
	return wrapper.Items, wrapper.NextCursor, nil
}

func (ss *SkillsService) GetSkillDetail(repoKey, slug string) (*SkillDetail, error) {
	log.Debug(fmt.Sprintf("Getting skill detail for '%s' in repo '%s'...", slug, repoKey))
	body, err := ss.sendGet(repoKey, fmt.Sprintf("skills/%s", slug))
	if err != nil {
		return nil, err
	}
	var detail SkillDetail
	if err = json.Unmarshal(body, &detail); err != nil {
		return nil, errorutils.CheckErrorf("failed to parse skill detail response: %s", err.Error())
	}
	return &detail, nil
}

func (ss *SkillsService) SearchSkills(repoKey, query string, limit int) ([]SkillSearchResult, error) {
	log.Debug(fmt.Sprintf("Searching skills in repo '%s' with query '%s'...", repoKey, query))
	body, err := ss.sendGet(repoKey, fmt.Sprintf("search?q=%s&limit=%d", url.QueryEscape(query), limit))
	if err != nil {
		return nil, err
	}

	var wrapper skillSearchResponse
	if err = json.Unmarshal(body, &wrapper); err != nil {
		return nil, errorutils.CheckErrorf("failed to parse skill search response: %s", err.Error())
	}
	return wrapper.Results, nil
}

func (ss *SkillsService) VersionExists(repoKey, slug, version string) (bool, error) {
	versions, err := ss.ListVersions(repoKey, slug)
	if err != nil {
		return false, err
	}
	for _, v := range versions {
		if v.Version == version {
			return true, nil
		}
	}
	return false, nil
}

// SearchByProperty uses the Artifactory property search API to find skills
// by their skill.name property across all repositories.
func (ss *SkillsService) SearchByProperty(query string) ([]SkillPropertySearchResult, error) {
	log.Debug(fmt.Sprintf("Searching skills by property skill.name='%s'...", query))
	baseURL := utils.AddTrailingSlashIfNeeded(ss.ArtDetails.GetUrl())
	searchURL := fmt.Sprintf("%sapi/search/prop?skill.name=%s", baseURL, query)
	log.Debug("Property search request:", searchURL)

	httpDetails := ss.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet(searchURL, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}

	var wrapper propSearchResponse
	if err = json.Unmarshal(body, &wrapper); err != nil {
		return nil, errorutils.CheckErrorf("failed to parse property search response: %s", err.Error())
	}

	var results []SkillPropertySearchResult
	for _, item := range wrapper.Results {
		r, ok := parsePropSearchURI(item.URI)
		if !ok {
			log.Warn(fmt.Sprintf("Skipping property search result with unparseable URI: %s", item.URI))
			continue
		}
		results = append(results, r)
	}
	return results, nil
}

// parsePropSearchURI extracts repo, slug, and version from a URI like:
// https://host/artifactory/api/storage/{repo}/{slug}/{version}/{slug}-{version}.zip
func parsePropSearchURI(uri string) (SkillPropertySearchResult, bool) {
	idx := strings.Index(uri, "/api/storage/")
	if idx == -1 {
		return SkillPropertySearchResult{}, false
	}
	// path after /api/storage/ => {repo}/{slug}/{version}/{file}
	path := uri[idx+len("/api/storage/"):]
	parts := strings.SplitN(path, "/", 4)
	if len(parts) < 3 {
		return SkillPropertySearchResult{}, false
	}
	return SkillPropertySearchResult{
		Repo:    parts[0],
		Name:    parts[1],
		Version: parts[2],
		URI:     uri,
	}, true
}

// GetXrayStatus calls the Skills Xray gate endpoint to check the scan status of a skill artifact.
// This endpoint uses api/skills/{repoKey}/xrayStatus (no /api/v1/ segment).
func (ss *SkillsService) GetXrayStatus(repoKey, artifactPath string) (*SkillXrayStatusResponse, error) {
	baseURL := utils.AddTrailingSlashIfNeeded(ss.ArtDetails.GetUrl())
	requestURL := fmt.Sprintf("%sapi/skills/%s/xrayStatus?path=%s", baseURL, repoKey, url.QueryEscape(artifactPath))
	log.Debug("Skills Xray status request:", requestURL)

	httpDetails := ss.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet(requestURL, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}

	var result SkillXrayStatusResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, errorutils.CheckErrorf("failed to parse xray status response: %s", err.Error())
	}
	return &result, nil
}

func (ss *SkillsService) sendGet(repoKey, endpoint string) ([]byte, error) {
	baseURL := utils.AddTrailingSlashIfNeeded(ss.ArtDetails.GetUrl())
	url := fmt.Sprintf("%sapi/skills/%s/api/v1/%s", baseURL, repoKey, endpoint)
	log.Debug("Skills API request:", url)

	httpDetails := ss.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	log.Debug("Artifactory response:", resp.Status)
	return body, nil
}

type SkillVersion struct {
	Version   string `json:"version"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	Changelog string `json:"changelog,omitempty"`
}

type SkillSearchResult struct {
	Name        string `json:"name"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

type SkillPropertySearchResult struct {
	Repo    string `json:"repo"`
	Name    string `json:"name"`
	Version string `json:"version"`
	URI     string `json:"uri"`
}

type SkillDetail struct {
	Slug          string `json:"slug"`
	DisplayName   string `json:"displayName,omitempty"`
	Summary       string `json:"summary,omitempty"`
	LatestVersion string `json:"latestVersion,omitempty"`
	Versions      int    `json:"versions,omitempty"`
	Updated       string `json:"updated,omitempty"`
}

type SkillListItem struct {
	Slug          string        `json:"slug"`
	DisplayName   string        `json:"displayName,omitempty"`
	Summary       string        `json:"summary,omitempty"`
	LatestVersion *SkillVersion `json:"latestVersion,omitempty"`
	Versions      int           `json:"versions,omitempty"`
	Updated       string        `json:"updated,omitempty"`
}

type skillListResponse struct {
	Items      []SkillListItem `json:"items"`
	NextCursor string          `json:"nextCursor,omitempty"`
	Total      int             `json:"total,omitempty"`
}

type skillVersionsResponse struct {
	Items []SkillVersion `json:"items"`
}

type skillSearchResponse struct {
	Results []SkillSearchResult `json:"results"`
	Total   int                 `json:"total,omitempty"`
	Offset  int                 `json:"offset,omitempty"`
	Limit   int                 `json:"limit,omitempty"`
}

type propSearchResponse struct {
	Results []propSearchResultItem `json:"results"`
}

type propSearchResultItem struct {
	URI string `json:"uri"`
}
