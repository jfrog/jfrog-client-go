package entitlements

import (
	"errors"
	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/bintray/services/versions"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"path"
	"strings"
)

func BuildEntitlementsUrl(bintrayDetails auth.BintrayDetails, details *versions.Path) string {
	return bintrayDetails.GetApiUrl() + createBintrayPath(details) + "/entitlements"
}

func BuildEntitlementUrl(bintrayDetails auth.BintrayDetails, details *versions.Path, entId string) string {
	return BuildEntitlementsUrl(bintrayDetails, details) + "/" + entId
}

func CreateVersionDetails(versionStr string) (*versions.Path, error) {
	parts := strings.Split(versionStr, "/")
	if len(parts) == 1 {
		err := errorutils.CheckError(errors.New("Argument format should be subject/repository or subject/repository/package or subject/repository/package/version. Got " + versionStr))
		if err != nil {
			return nil, err
		}
	}
	return versions.CreatePath(versionStr)
}

func createBintrayPath(details *versions.Path) string {
	if details.Version == "" {
		if details.Package == "" {
			return path.Join("repos", details.Subject, details.Repo)
		}
		return path.Join("packages", details.Subject, details.Repo, details.Package)
	}
	return path.Join("packages", details.Subject, details.Repo, details.Package, "versions", details.Version)

}
