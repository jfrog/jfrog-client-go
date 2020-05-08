package mavensync

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/bintray/services/utils/tests"
	"github.com/jfrog/jfrog-client-go/bintray/services/versions"
)

func TestGenerateMavenCentralSyncPath(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			"valid full path should generate valid url and no error",
			"my-subject/my-repo/my-pkg/ver-1.9.1",
			"https://api.bintray.com/maven_central_sync/my-subject/my-repo/my-pkg/versions/ver-1.9.1",
			false,
		},
		{
			"only package path should result in error",
			"my-subject/my-repo/my-pkg",
			"",
			true,
		},
		{
			"only repo path should result in error",
			"my-subject/my-repo",
			"",
			true,
		},
		{
			"only subject path should result in error",
			"my-subject",
			"",
			true,
		},
		{
			"empty path should result in error",
			"",
			"",
			true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			path, _ := versions.CreatePath(tt.path)

			got, err := buildSyncURL(tests.CreateBintrayDetails(), path)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildSyncURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("buildSyncURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
