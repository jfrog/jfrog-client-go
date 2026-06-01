package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePropSearchURI(t *testing.T) {
	tests := []struct {
		name   string
		uri    string
		want   SkillPropertySearchResult
		wantOK bool
	}{
		{
			name:   "valid URI",
			uri:    "https://example.jfrog.io/artifactory/api/storage/rafi-skills/4chan-reader/1.0.0/4chan-reader-1.0.0.zip",
			want:   SkillPropertySearchResult{Repo: "rafi-skills", Name: "4chan-reader", Version: "1.0.0", URI: "https://example.jfrog.io/artifactory/api/storage/rafi-skills/4chan-reader/1.0.0/4chan-reader-1.0.0.zip"},
			wantOK: true,
		},
		{
			name:   "valid URI with different version",
			uri:    "https://host.com/artifactory/api/storage/my-repo/my-skill/2.3.1/my-skill-2.3.1.zip",
			want:   SkillPropertySearchResult{Repo: "my-repo", Name: "my-skill", Version: "2.3.1", URI: "https://host.com/artifactory/api/storage/my-repo/my-skill/2.3.1/my-skill-2.3.1.zip"},
			wantOK: true,
		},
		{
			name:   "no api/storage segment",
			uri:    "https://host.com/artifactory/rafi-skills/4chan-reader/1.0.0/file.zip",
			want:   SkillPropertySearchResult{},
			wantOK: false,
		},
		{
			name:   "too few path segments",
			uri:    "https://host.com/artifactory/api/storage/repo/slug",
			want:   SkillPropertySearchResult{},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parsePropSearchURI(tt.uri)
			assert.Equal(t, tt.wantOK, ok)
			if ok {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
