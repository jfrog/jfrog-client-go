package tests

import (
	"github.com/jfrog/jfrog-client-go/bintray/auth"
)

func CreateBintrayDetails() auth.BintrayDetails {
	btDetails := auth.NewBintrayDetails()
	btDetails.SetApiUrl("https://api.bintray.com/")
	btDetails.SetDownloadServerUrl("https://dl.bintray.com/")
	btDetails.SetUser("user")
	btDetails.SetKey("api-key")
	btDetails.SetDefPackageLicense("Apache-2.0")
	return btDetails
}
