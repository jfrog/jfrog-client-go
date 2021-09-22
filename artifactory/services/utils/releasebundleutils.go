package utils

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/openpgp"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type RbGPGValidation struct {
	client            *jfroghttpclient.JfrogHttpClient
	artDetails        *auth.ServiceDetails
	rbName            string
	rbVersion         string
	publicKeyFilePath string
	artifactsMap      map[string]string
}

func (r *RbGPGValidation) ArtifactsMap() map[string]string {
	return r.artifactsMap
}

func (r *RbGPGValidation) SetRbName(rbName string) *RbGPGValidation {
	r.rbName = rbName
	return r
}

func (r *RbGPGValidation) SetRbVersion(rbVersion string) *RbGPGValidation {
	r.rbVersion = rbVersion
	return r
}

func (r *RbGPGValidation) SetPublicKey(path string) *RbGPGValidation {
	r.publicKeyFilePath = path
	return r
}

func (r *RbGPGValidation) SetClient(client *jfroghttpclient.JfrogHttpClient) *RbGPGValidation {
	r.client = client
	return r
}

func (r *RbGPGValidation) SetAtrifactoryDetails(artDetails *auth.ServiceDetails) *RbGPGValidation {
	r.artDetails = artDetails
	return r
}

func NewRbGPGValidation() (*RbGPGValidation, error) {
	return &RbGPGValidation{}, nil
}

func GetTestResourcesPath() string {
	dir, _ := os.Getwd()
	return filepath.ToSlash(dir + "/testdata/")
}

// Validate gets a signed release bundle from artifactory, validate the signature with the public gpg key and save the release bundle's artifacts in a map
func (r *RbGPGValidation) Validate() error {
	httpClientsDetails := (*r.artDetails).CreateHttpClientDetails()
	// Release bundle's details return in a JWS format, so we can validate the signature of the signed release bundle with the provided public key.
	request := (*r.artDetails).GetUrl() + "api/release/bundles/" + r.rbName + "/" + r.rbVersion + "?format=jws"
	_, body, _, err := r.client.SendGet(request, true, &httpClientsDetails)
	if err != nil {
		return err
	}
	verifiedRB, err := VerifyJwtToken(string(body), r.publicKeyFilePath)
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

func (r *RbGPGValidation) VerifySpecificArtifact(artifactPath, sha256 string) bool {
	if r.artifactsMap[artifactPath] == sha256 {
		return true
	}
	return false
}

func VerifyJwtToken(bundleTokenStr, publicKeyPath string) (*ReleaseBundleModel, error) {
	model := &ReleaseBundleModel{}
	token, err := jwt.ParseWithClaims(bundleTokenStr, model, func(token *jwt.Token) (interface{}, error) {
		key, err := ioutil.ReadFile(filepath.Join(publicKeyPath))
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
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errorutils.CheckError(errors.New("release bundle token is invalid"))
	}
	if claims, ok := token.Claims.(*ReleaseBundleModel); ok {
		return claims, nil
	}
	return nil, errorutils.CheckError(errors.New("token payload could not be cast to release bundle model"))
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
	// in practice Id doesn't seem to be in use, so we calculate one
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
