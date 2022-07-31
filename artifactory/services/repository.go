package services

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type RepositoryService struct {
	repoType   string
	isUpdate   bool
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewRepositoryService(repoType string, client *jfroghttpclient.JfrogHttpClient, isUpdate bool) *RepositoryService {
	return &RepositoryService{repoType: repoType, client: client, isUpdate: isUpdate}
}

func (rs *RepositoryService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return rs.client
}

func (rs *RepositoryService) performRequest(params interface{}, repoKey string) error {
	content, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return err
	}
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.repositories."+strings.ToTitle(rs.repoType)+"RepositoryConfiguration+json", &httpClientsDetails.Headers)
	var url = rs.ArtDetails.GetUrl() + "api/repositories/" + url.PathEscape(repoKey)
	var operationString string
	var resp *http.Response
	var body []byte
	if rs.isUpdate {
		log.Info("Updating " + strings.ToLower(rs.repoType) + " repository...")
		operationString = "updating"
		resp, body, err = rs.client.SendPost(url, content, &httpClientsDetails)
	} else {
		log.Info("Creating " + strings.ToLower(rs.repoType) + " repository...")
		operationString = "creating"
		resp, body, err = rs.client.SendPut(url, content, &httpClientsDetails)
	}
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done", operationString, "repository", "'"+repoKey+"'.")
	return nil
}

type RepositoryBaseParams struct {
	Rclass          string   `json:"rclass"`
	Key             string   `json:"key,omitempty"`
	ProjectKey      string   `json:"projectKey,omitempty"`
	Environments    []string `json:"environments,omitempty"`
	PackageType     string   `json:"packageType,omitempty"`
	Description     string   `json:"description,omitempty"`
	Notes           string   `json:"notes,omitempty"`
	IncludesPattern string   `json:"includesPattern,omitempty"`
	ExcludesPattern string   `json:"excludesPattern,omitempty"`
	RepoLayoutRef   string   `json:"repoLayoutRef,omitempty"`
}

type AdditionalRepositoryBaseParams struct {
	BlackedOut         *bool    `json:"blackedOut,omitempty"`
	XrayIndex          *bool    `json:"xrayIndex,omitempty"`
	PropertySets       []string `json:"propertySets,omitempty"`
	DownloadRedirect   *bool    `json:"downloadRedirect,omitempty"`
	PriorityResolution *bool    `json:"priorityResolution,omitempty"`
}

type CargoRepositoryParams struct {
	CargoAnonymousAccess *bool `json:"cargoAnonymousAccess,omitempty"`
}

type DebianRepositoryParams struct {
	DebianTrivialLayout             *bool    `json:"debianTrivialLayout,omitempty"`
	OptionalIndexCompressionFormats []string `json:"optionalIndexCompressionFormats,omitempty"`
}

type DockerRepositoryParams struct {
	MaxUniqueTags       int    `json:"maxUniqueTags,omitempty"`
	DockerApiVersion    string `json:"dockerApiVersion,omitempty"`
	BlockPushingSchema1 *bool  `json:"blockPushingSchema1,omitempty"`
	DockerTagRetention  int    `json:"dockerTagRetention,omitempty"`
}

type JavaPackageManagersRepositoryParams struct {
	MaxUniqueSnapshots           int    `json:"maxUniqueSnapshots,omitempty"`
	HandleReleases               *bool  `json:"handleReleases,omitempty"`
	HandleSnapshots              *bool  `json:"handleSnapshots,omitempty"`
	SuppressPomConsistencyChecks *bool  `json:"suppressPomConsistencyChecks,omitempty"`
	SnapshotVersionBehavior      string `json:"snapshotVersionBehavior,omitempty"`
	ChecksumPolicyType           string `json:"checksumPolicyType,omitempty"`
}

type KeyPairRefsRepositoryParams struct {
	PrimaryKeyPairRef   string `json:"primaryKeyPairRef,omitempty"`
	SecondaryKeyPairRef string `json:"secondaryKeyPairRef,omitempty"`
}

type NugetRepositoryParams struct {
	MaxUniqueSnapshots       int   `json:"maxUniqueSnapshots,omitempty"`
	ForceNugetAuthentication *bool `json:"forceNugetAuthentication,omitempty"`
}

type RpmRepositoryParams struct {
	YumRootDepth            int    `json:"yumRootDepth,omitempty"`
	CalculateYumMetadata    *bool  `json:"calculateYumMetadata,omitempty"`
	EnableFileListsIndexing *bool  `json:"enableFileListsIndexing,omitempty"`
	YumGroupFileNames       string `json:"yumGroupFileNames,omitempty"`
}
