package utils

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"golang.org/x/crypto/openpgp"
	"os"
	"path/filepath"
	"strings"
)

type RbGpgValidator struct {
	client            *jfroghttpclient.JfrogHttpClient
	artDetails        *auth.ServiceDetails
	rbName            string
	rbVersion         string
	publicKeyFilePath string
	// Map containing the artifacts which are associated with a specific release bundle.
	// This map is used for validating that the downloaded files have the same checksum as in the release bundle manifest. This is done for security reasons.
	// The key is the path of the artifact in Artifactory and the value is it's sha256.
	artifactsMap map[string]string
}

func (r *RbGpgValidator) ArtifactsMap() map[string]string {
	return r.artifactsMap
}

func (r *RbGpgValidator) SetRbName(rbName string) *RbGpgValidator {
	r.rbName = rbName
	return r
}

func (r *RbGpgValidator) SetRbVersion(rbVersion string) *RbGpgValidator {
	r.rbVersion = rbVersion
	return r
}

func (r *RbGpgValidator) SetPublicKey(path string) *RbGpgValidator {
	r.publicKeyFilePath = path
	return r
}

func (r *RbGpgValidator) SetClient(client *jfroghttpclient.JfrogHttpClient) *RbGpgValidator {
	r.client = client
	return r
}

func (r *RbGpgValidator) SetAtrifactoryDetails(artDetails *auth.ServiceDetails) *RbGpgValidator {
	r.artDetails = artDetails
	return r
}

func NewRbGpgValidator() *RbGpgValidator {
	return &RbGpgValidator{}
}

func GetTestResourcesPath() string {
	dir, _ := os.Getwd()
	return filepath.ToSlash(dir + "/testdata/")
}

// Validate gets a signed release bundle from Artifactory, validates the signature with the public gpg key and saves the release bundle's artifacts in a map
func (r *RbGpgValidator) Validate() error {
	httpClientsDetails := (*r.artDetails).CreateHttpClientDetails()
	// Release bundle's details return in a JWS format, so we can validate the signature of the signed release bundle with the provided public key.
	request := (*r.artDetails).GetUrl() + "api/release/bundles/" + r.rbName + "/" + r.rbVersion + "?format=jws"
	_, body, _, err := r.client.SendGet(request, true, &httpClientsDetails)
	if err != nil {
		return err
	}
	verifiedRB, err := r.verifyJwtToken(string(body))
	if err != nil {
		return err
	}
	// Save all release bundle's artifacts in a map.
	elementMap := make(map[string]string)
	for _, artifact := range verifiedRB.Artifacts {
		elementMap[artifact.RepoPath] = artifact.Checksum
	}
	r.artifactsMap = elementMap
	return nil
}

func (r *RbGpgValidator) VerifyArtifact(artifactPath, sha256 string) error {
	if r.artifactsMap[artifactPath] == sha256 {
		return nil
	}
	return errorutils.CheckErrorf("GPG validation failed: artifact in not signed with the provided key - %s ", artifactPath)
}

func (r *RbGpgValidator) verifyJwtToken(bundleTokenStr string) (*ReleaseBundleModel, error) {
	model := &ReleaseBundleModel{}
	token, err := jwt.ParseWithClaims(bundleTokenStr, model, func(token *jwt.Token) (interface{}, error) {
		key, err := os.ReadFile(filepath.Join(r.publicKeyFilePath))
		if err != nil {
			return nil, errorutils.CheckError(err)
		}
		reader := strings.NewReader(string(key))
		entities, err := openpgp.ReadArmoredKeyRing(reader)
		if err != nil {
			return nil, errorutils.CheckError(err)
		}
		return entities[0].PrimaryKey.PublicKey, nil
	})
	if errorutils.CheckError(err) != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errorutils.CheckErrorf("release bundle token is invalid")
	}
	if claims, ok := token.Claims.(*ReleaseBundleModel); ok {
		return claims, nil
	}
	return nil, errorutils.CheckErrorf("failed casting the token payload to the release bundle model")
}

type ReleaseBundleModel struct {
	Id           string              `json:"id"`
	Name         string              `json:"name"`
	Version      string              `json:"version"`
	Created      string              `json:"created"`
	Description  string              `json:"description"`
	Artifacts    []*ReleaseArtifact  `json:"artifacts"`
	ReleaseNotes *BundleReleaseNotes `json:"release_notes"`
	Signature    string              `json:"signature"`
	Type         string              `json:"type"`
}

func (rbm *ReleaseBundleModel) Valid() error {
	return nil
}

func (rbm *ReleaseBundleModel) GetOrCalculateId() string {
	if rbm.Id != "" {
		return rbm.Id
	}
	// in practice, Id doesn't seem to be in use, so we calculate one
	return fmt.Sprintf("%s:%s", rbm.Name, rbm.Version)
}

type BundleReleaseNotes struct {
	Content string `json:"content"`
	Syntax  string `json:"syntax"`
}

type ReleaseArtifact struct {
	RepoPath      string `json:"repo_path"`
	Checksum      string `json:"checksum"`
	DownloadToken string `json:"download_token"`
}
