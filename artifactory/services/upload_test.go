package services

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebianProperties(t *testing.T) {
	var debianPaths = []struct {
		in       string
		expected string
	}{
		{"dist/comp/arch", ";deb.distribution=dist;deb.component=comp;deb.architecture=arch"},
		{"dist1,dist2/comp/arch", ";deb.distribution=dist1,dist2;deb.component=comp;deb.architecture=arch"},
		{"dist/comp1,comp2/arch", ";deb.distribution=dist;deb.component=comp1,comp2;deb.architecture=arch"},
		{"dist/comp/arch1,arch2", ";deb.distribution=dist;deb.component=comp;deb.architecture=arch1,arch2"},
		{"dist1,dist2/comp1,comp2/arch1,arch2", ";deb.distribution=dist1,dist2;deb.component=comp1,comp2;deb.architecture=arch1,arch2"},
	}

	for _, v := range debianPaths {
		result := getDebianProps(v.in)
		if result != v.expected {
			t.Errorf("getDebianProps(\"%s\") => '%s', want '%s'", v.in, result, v.expected)
		}
	}
}

func TestBuildUploadUrls(t *testing.T) {
	var testsParams = []struct {
		targetPath                  string
		targetProps                 string
		buildProps                  string
		expectedTargetPathWithProps string
	}{
		{"repo1/file1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file1;k1=v1;k2=v2"},
		{"repo1/file@1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file@1;k1=v1;k2=v2"},
		{"repo1/file;1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file%3B1;k1=v1;k2=v2"},
		{"repo1/file,1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file,1;k1=v1;k2=v2"},
		{"repo1/file^1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file%5E1;k1=v1;k2=v2"},
		{"repo1/file:1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file:1;k1=v1;k2=v2"},
		{"repo1/file1", "", "k2=v2", "http://localhost:8881/artifactory/repo1/file1;k2=v2"},
		{"repo1/file1", "k1=v1", "", "http://localhost:8881/artifactory/repo1/file1;k1=v1"},
		{"repo1/file1", "", "", "http://localhost:8881/artifactory/repo1/file1"},
	}

	for _, v := range testsParams {
		targetProps, e := utils.ParseProperties(v.targetProps)
		assert.NoError(t, e)
		actualTargetPathWithProps, e := buildUploadUrls("http://localhost:8881/artifactory/", v.targetPath, v.buildProps, "", targetProps)
		assert.NoError(t, e)
		assert.Equal(t, v.expectedTargetPathWithProps, actualTargetPathWithProps)
	}
}

func TestAddEscapingParenthesesWithTargetInArchive(t *testing.T) {
	type args struct {
		pattern         string
		target          string
		targetInArchive string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty parentheses", args{"()", "", "{2}"}, "\\(\\)"},
		{"empty parentheses", args{"()", "", "{}"}, "\\(\\)"},
		{"empty parentheses", args{"()", "", "{1}"}, "()"},
		{"empty parentheses", args{")(", "", "{1}"}, "\\)\\("},
		{"first parentheses", args{"(a)/(b)/(c)", "", "{2}/{3}"}, "\\(a\\)/(b)/(c)"},
		{"second parentheses", args{"(a)/(b)/(c)", "", "{1}/{3}"}, "(a)/\\(b\\)/(c)"},
		{"third parentheses", args{"(a)/(b)/(c)", "", "{1}/{2}"}, "(a)/(b)/\\(c\\)"},
		{"empty placeholders", args{"(a)/(b)/(c)", "", ""}, "\\(a\\)/\\(b\\)/\\(c\\)"},
		{"un-symmetric parentheses", args{")a)/(b)/(c(", "", ""}, "\\)a\\)/\\(b\\)/\\(c\\("},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, addEscapingParenthesesForUpload(tt.args.pattern, tt.args.target, tt.args.targetInArchive), "AddEscapingParentheses(%v, %v)", tt.args.pattern, tt.args.target)
		})
	}
}

func TestAddEscapingParenthesesWithTargetAndTargetInArchive(t *testing.T) {
	type args struct {
		pattern         string
		target          string
		targetInArchive string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty parentheses", args{"()", "{2}", "{3}"}, "\\(\\)"},
		{"empty parentheses", args{"()", "{}", "{}"}, "\\(\\)"},
		{"empty parentheses", args{"()()", "{2}", "{1}"}, "()()"},
		{"empty parentheses", args{"))(((", "{2}", "{1}"}, "\\)\\)\\(\\(\\("},
		{"first parentheses", args{"(a)/(b)/(c)/(d)", "{4}", "{2}/{3}"}, "\\(a\\)/(b)/(c)/(d)"},
		{"second parentheses", args{"(a)/(b)/(c)/(d)", "{1}/{4}", "{1}/{3}"}, "(a)/\\(b\\)/(c)/(d)"},
		{"last parentheses", args{"(a)/(b)/(c)/(d)", "{1}/{3}", "{2}/{3}"}, "(a)/(b)/(c)/\\(d\\)"},
		{"mixed parentheses", args{"(a)/(b)/(c)/(d)/(e)/(f)", "{5}", "{1}/{2}"}, "(a)/(b)/\\(c\\)/\\(d\\)/(e)/\\(f\\)"},
		{"out of range placeholders", args{"(a)/(b)/(c)", "{5}", "{4}"}, "\\(a\\)/\\(b\\)/\\(c\\)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, addEscapingParenthesesForUpload(tt.args.pattern, tt.args.target, tt.args.targetInArchive), "AddEscapingParentheses(%v, %v)", tt.args.pattern, tt.args.target)
		})
	}
}

func TestSkipDirUpload(t *testing.T) {
	data := []struct {
		targetFiles []string
		sourceDirs  []string
		targetDir   string
		sourceDir   string
		includeDirs bool
		result      bool
	}{
		{[]string{}, []string{}, "cli-rt1-1671381032/b", "testdata/a/b", true, false},
		{[]string{"dirdir/b/"}, []string{}, "dirdir", "testdata/a/b", true, true},
		{[]string{"cli-rt1-1671381032/b/"}, []string{}, "cli-rt1-1671381032/b", "testdata/a/b", true, true},
		{[]string{"cli-rt1-1671383851/c", "cli-rt1-1671383851/b3.in"}, []string{filepath.Join("testdata", "a", "b", "c")}, "cli-rt1-1671383851/b", filepath.Join("testdata", "a", "b"), true, true},
	}
	for _, d := range data {
		got := skipDirUpload(d.targetFiles, d.sourceDirs, d.targetDir, d.sourceDir, d.includeDirs)
		assert.Equal(t, d.result, got)
	}
}

// The checksum-deploy PUT that finalizes a multipart upload must carry the token returned by Artifactory,
// also when /api/v1/uploads/complete was retried after a RETRYABLE_ERROR status.
func TestMultipartUploadChecksumDeployTokenAfterRetriedComplete(t *testing.T) {
	const (
		checksumDeployToken = "checksum-deploy-token"
		targetPath          = "repo/file.bin"
	)
	var completeCalls, statusCalls, abortCalls, deployCalls atomic.Int32

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/uploads/config":
			_, err := w.Write([]byte(`{"supported":true}`))
			assert.NoError(t, err)
		case "/api/v1/uploads/create":
			_, err := w.Write([]byte(`{"token":"multipart-token"}`))
			assert.NoError(t, err)
		case "/api/v1/uploads/urlPart":
			_, err := w.Write([]byte(`{"url":"` + server.URL + `/storage/part"}`))
			assert.NoError(t, err)
		case "/storage/part":
			w.WriteHeader(http.StatusOK)
		case "/api/v1/uploads/complete":
			completeCalls.Add(1)
			w.WriteHeader(http.StatusAccepted)
		case "/api/v1/uploads/status":
			// The first status is a retryable error without a token, the status after the retried complete is finished with a token.
			response := `{"status":"RETRYABLE_ERROR","error":"Injected S3 SlowDown","progress":0}`
			if statusCalls.Add(1) > 1 {
				response = `{"status":"FINISHED","progress":100,"checksumToken":"` + checksumDeployToken + `"}`
			}
			_, err := w.Write([]byte(response))
			assert.NoError(t, err)
		case "/api/v1/uploads/abort":
			abortCalls.Add(1)
			w.WriteHeader(http.StatusOK)
		case "/" + targetPath:
			deployCalls.Add(1)
			assert.Equal(t, http.MethodPut, r.Method)
			assert.Equal(t, "true", r.Header.Get("X-Checksum-Deploy"))
			assert.Equal(t, checksumDeployToken, r.Header.Get("X-Checksum-Deploy-Token"), "the checksum deploy must carry the token returned after the retried complete")
			w.WriteHeader(http.StatusCreated)
		default:
			assert.Fail(t, "unexpected request", r.Method+" "+r.URL.Path)
		}
	}))
	defer server.Close()

	localPath := filepath.Join(t.TempDir(), "file.bin")
	require.NoError(t, os.WriteFile(localPath, []byte("multipart upload content"), 0o600))

	// The retries value sets the attempts budget of the multipart upload, so it must be above zero.
	client, err := jfroghttpclient.JfrogClientBuilder().SetRetries(3).Build()
	require.NoError(t, err)
	uploadService := NewUploadService(client)
	uploadService.SetServiceDetails(&dummyServiceDetails{version: "7.133.0"})
	uploadService.MultipartUpload = utils.NewMultipartUpload(client, &httputils.HttpClientDetails{}, server.URL)

	uploadParams := NewUploadParams()
	uploadParams.MinChecksumDeploy = 1 << 30
	uploadParams.MinSplitSize = 1
	uploadParams.SplitCount = 1

	resp, _, _, checksumDeployed, err := uploadService.doUpload(clientutils.Artifact{LocalPath: localPath, TargetPath: targetPath},
		server.URL+"/"+targetPath, "", httputils.HttpClientDetails{Headers: map[string]string{}}, uploadParams)
	require.NoError(t, err)
	assert.False(t, checksumDeployed)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.EqualValues(t, 2, completeCalls.Load(), "complete should be retried once")
	assert.EqualValues(t, 2, statusCalls.Load())
	assert.EqualValues(t, 1, deployCalls.Load())
	assert.Zero(t, abortCalls.Load())
}

// A service details stub with a fixed Artifactory version, so the version is not requested from the server.
type dummyServiceDetails struct {
	auth.CommonConfigFields
	version string
}

func (d *dummyServiceDetails) GetVersion() (string, error) {
	return d.version, nil
}
