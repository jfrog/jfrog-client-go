package bundle

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/version"
	"net/http"
	"path"
	"sort"
)

type BundleService struct {
	client     *httpclient.HttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewBundleService(client *httpclient.HttpClient) *BundleService {
	return &BundleService{client: client}
}

func (ds *BundleService) GetArtifactoryDetails() auth.ArtifactoryDetails {
	return ds.ArtDetails
}

func (ds *BundleService) SetArtifactoryDetails(rt auth.ArtifactoryDetails) {
	ds.ArtDetails = rt
}

func (ds *BundleService) GetJfrogHttpClient() *httpclient.HttpClient {
	return ds.client
}

func (ds *BundleService) GetBundleVersions(bundleName string) ([]Version, error) {
	bundlesUrl := ds.ArtDetails.GetUrl() + path.Join("api/release/bundles", bundleName)

	httpClientsDetails := ds.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := ds.client.SendGet(bundlesUrl, true, httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	var result versions
	err = json.Unmarshal(body, &result)
	if errorutils.CheckError(err) != nil {
		return nil, err
	}
	sort.Sort(SortableVersion(result.Versions))
	return result.Versions, nil

}

type SortableVersion []Version

type versions struct {
	Versions SortableVersion `json:"versions"`
}

type Version struct {
	Version string `json:"version"`
	Created string `json:"created,omitempty"`
	Status  string `json:"status"`
}

func (a SortableVersion) Len() int {
	return len(a)
}

func (a SortableVersion) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a SortableVersion) Less(i, j int) bool {
	return version.Compare(a[i].Version, a[j].Version) < 0
}
