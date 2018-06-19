package services

import (
	"github.com/jfrog/jfrog-client-go/bintray/services/utils/tests"
	"github.com/jfrog/jfrog-client-go/bintray/services/versions"
	"testing"
)

func TestDownloadVersion(t *testing.T) {
	var err error
	params := &DownloadVersionParams{Params: &versions.Params{}}
	params.IncludeUnpublished = false
	params.Path, err = CreateVersionDetailsForDownloadVersion("test-subject/test-repo/test-package/ver-1.2")
	if err != nil {
		t.Error(err.Error())
	}

	url := buildDownloadVersionUrl(tests.CreateBintrayDetails().GetApiUrl(), params)
	expected := "https://api.bintray.com/packages/test-subject/test-repo/test-package/versions/ver-1.2/files"
	if expected != url {
		t.Error("Got unexpected url from BuildDownloadVersionUrl. Expected: " + expected + " Got " + url)
	}
}
