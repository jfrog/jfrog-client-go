package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

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

func (ss *SkillsService) SearchSkills(repoKey, query string, limit int) ([]SkillSearchResult, error) {
	log.Debug(fmt.Sprintf("Searching skills in repo '%s' with query '%s'...", repoKey, query))
	body, err := ss.sendGet(repoKey, fmt.Sprintf("search?q=%s&limit=%d", query, limit))
	if err != nil {
		return nil, err
	}

	var wrapper skillSearchResponse
	if err = json.Unmarshal(body, &wrapper); err != nil {
		return nil, errorutils.CheckErrorf("failed to parse skill search response: %s", err.Error())
	}
	return wrapper.Skills, nil
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
	Slug        string `json:"slug"`
	DisplayName string `json:"displayName,omitempty"`
	Summary     string `json:"summary,omitempty"`
}

type skillVersionsResponse struct {
	Items []SkillVersion `json:"items"`
}

type skillSearchResponse struct {
	Skills []SkillSearchResult `json:"skills"`
}
