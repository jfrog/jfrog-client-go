package services

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/bintray/services/utils"
	"github.com/jfrog/jfrog-client-go/bintray/services/versions"
	"github.com/jfrog/jfrog-client-go/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func NewUploadService(client *httpclient.HttpClient) *UploadService {
	us := &UploadService{client: client}
	return us
}

func NewUploadParams() *UploadParams {
	return &UploadParams{Params: &versions.Params{}}
}

type UploadService struct {
	client         *httpclient.HttpClient
	BintrayDetails auth.BintrayDetails
	DryRun         bool
	Threads        int
}

type UploadParams struct {
	// Files pattern to be uploaded
	Pattern string   `yaml:"file_pattern,omitempty"`
	Files   []string `yaml:"files,omitempty"`

	// Target version details
	*versions.Params `yaml:"params"`

	// Target local path
	TargetPath string `yaml:"target"`

	UseRegExp bool   `yaml:"regex"`
	Flat      bool   `yaml:"flat,omitempty"`
	Recursive bool   `yaml:"recursive,omitempty"`
	Explode   bool   `yaml:"explode,omitempty"`
	Override  bool   `yaml:"override,omitempty"`
	Publish   bool   `yaml:"publish,omitempty"`
	Deb       string `yaml:"pkgdeb"`
}

func (us *UploadService) Upload(uploadDetails *UploadParams) (totalUploaded, totalFailed int, err error) {
	if us.BintrayDetails.GetUser() == "" {
		us.BintrayDetails.SetUser(uploadDetails.Subject)
	}

	// Get the list of artifacts to be uploaded to:
	var artifacts []clientutils.Artifact
	artifacts, err = us.getFilesToUpload(uploadDetails)
	if err != nil {
		return
	}

	baseUrl := us.BintrayDetails.GetApiUrl() + path.Join("content", uploadDetails.Subject, uploadDetails.Repo, uploadDetails.Package, uploadDetails.Version)
	totalUploaded, totalFailed, err = us.uploadFiles(uploadDetails, artifacts, baseUrl)
	return
}

func (us *UploadService) uploadFiles(uploadDetails *UploadParams, artifacts []clientutils.Artifact, baseUrl string) (totalUploaded, totalFailed int, err error) {
	size := len(artifacts)
	var wg sync.WaitGroup

	// Create an map where the key is a threadId so each entry is tied to a specific thread
	// this avoids us needing to use a Mutex or sync.Map within the go routine.
	uploadedArtifacts := make(map[int][]clientutils.Artifact, us.Threads)
	matrixParams := getMatrixParams(uploadDetails)
	for i := 0; i < us.Threads; i++ {
		wg.Add(1)
		// The threadId key MUST exist in the map first so that the map itself is not being mutated inside a go routine
		uploadedArtifacts[i] = make([]clientutils.Artifact, 0)
		go func(threadId int) {
			logMsgPrefix := clientutils.GetLogMsgPrefix(threadId, us.DryRun)
			for j := threadId; j < size; j += us.Threads {
				if !us.DryRun {
					url := baseUrl + "/" + artifacts[j].TargetPath + matrixParams
					uploaded, e := uploadFile(artifacts[j], url, logMsgPrefix, us.BintrayDetails)
					if e != nil {
						log.Error(logMsgPrefix, "Failed uploading artifact:", artifacts[j].LocalPath, ":", e)
					}
					if uploaded {
						uploadedArtifacts[threadId] = append(uploadedArtifacts[threadId], artifacts[j])
					}
				} else {
					log.Info("[Dry Run] Uploading artifact:", artifacts[j].LocalPath)
					uploadedArtifacts[threadId] = append(uploadedArtifacts[threadId], artifacts[j])
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	if uploadDetails.ShowInDownloadList {
		// Even though we are not running the list for download in go routines we need this outer loop
		// since we are using a thread specific key in the uploadedArtifacts map to get around needing to use
		// a Mutex or sync.Map when adding entries to the map.
		for i := 0; i < us.Threads; i++ {
			for _, artifact := range uploadedArtifacts[i] {
				if !us.DryRun {
					listUrl := us.BintrayDetails.GetApiUrl() + path.Join(
						"file_metadata",
						uploadDetails.Subject,
						uploadDetails.Repo, artifact.TargetPath)

					var listed bool
					// Retry loop, will retry to list uploaded artifacts.
					for j := 0; j < 30; j++ {
						if listed, err = SownInDownloadList(listUrl, us.BintrayDetails); listed || err != nil {
							if err != nil {
								log.Error(err)
							}
							break
						}
						time.Sleep(1 * time.Second)
					}
					if listed {
						log.Info("Listed for download", artifact.TargetPath)
					} else {
						log.Error("Failed listing for download", artifact.TargetPath)
					}
				} else {
					log.Info("[Dry Run] Listed for download", artifact.TargetPath)
				}
			}
		}
	}

	totalUploaded = 0
	for _, i := range uploadedArtifacts {
		totalUploaded += len(i)
	}
	log.Debug("Uploaded", strconv.Itoa(totalUploaded), "artifacts.")
	totalFailed = size - totalUploaded
	if totalFailed > 0 {
		log.Error("Failed uploading", strconv.Itoa(totalFailed), "artifacts.")
	}
	return
}

func getMatrixParams(uploadDetails *UploadParams) string {
	params := ""
	if uploadDetails.Publish {
		params += ";publish=1"
	}
	if uploadDetails.Override {
		params += ";override=1"
	}
	if uploadDetails.Explode {
		params += ";explode=1"
	}
	if uploadDetails.Deb != "" {
		params += getDebianMatrixParams(uploadDetails.Deb)
	}
	return params
}

func getDebianMatrixParams(debianPropsStr string) string {
	debProps := strings.Split(debianPropsStr, "/")
	return ";deb_distribution=" + debProps[0] +
		";deb_component=" + debProps[1] +
		";deb_architecture=" + debProps[2]
}

func getDebianDefaultPath(debianPropsStr, packageName string) string {
	debProps := strings.Split(debianPropsStr, "/")
	component := strings.Split(debProps[1], ",")[0]
	return path.Join("pool", component, packageName[0:1], packageName) + "/"
}

func uploadFile(artifact clientutils.Artifact, url, logMsgPrefix string, bintrayDetails auth.BintrayDetails) (bool, error) {
	log.Info(logMsgPrefix+"Uploading artifact:", artifact.LocalPath)

	httpClientsDetails := bintrayDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return false, err
	}
	resp, body, err := client.UploadFile(artifact.LocalPath, url, logMsgPrefix, httpClientsDetails, utils.BintrayUploadRetries, nil)
	if err != nil {
		return false, err
	}
	log.Debug(logMsgPrefix+"Bintray response:", resp.Status)
	fmt.Println(resp.Status, artifact)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Error(logMsgPrefix + "Bintray response: " + resp.Status + "\n" + clientutils.IndentJson(body))
	}

	return resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK, nil
}

func SownInDownloadList(url string, bintrayDetails auth.BintrayDetails) (bool, error) {
	httpClientsDetails := bintrayDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return false, err
	}
	resp, body, err := client.SendPut(url, []byte(`{"list_in_downloads":true}`), httpClientsDetails)
	if err != nil {
		return false, err
	}
	log.Debug("Bintray response: " + resp.Status + "\n" + clientutils.IndentJson(body))

	return resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK, nil
}

func getSingleFileToUpload(rootPath, targetPath string, flat bool) clientutils.Artifact {
	var uploadPath string
	rootPathOrig := rootPath
	if targetPath != "" && !strings.HasSuffix(targetPath, "/") {
		rootPath = targetPath
		targetPath = ""
	}

	if flat {
		uploadPath, _ = fileutils.GetFileAndDirFromPath(rootPath)
		uploadPath = targetPath + uploadPath
	} else {
		uploadPath = targetPath + rootPath
		uploadPath = clientutils.TrimPath(uploadPath)
	}
	return clientutils.Artifact{LocalPath: rootPathOrig, TargetPath: uploadPath}
}

func (us *UploadService) getFilesToUpload(uploadDetails *UploadParams) ([]clientutils.Artifact, error) {
	var debianDefaultPath string
	if uploadDetails.TargetPath == "" && uploadDetails.Deb != "" {
		debianDefaultPath = getDebianDefaultPath(uploadDetails.Deb, uploadDetails.Package)
	}

	rootPath := clientutils.GetRootPath(uploadDetails.Pattern, uploadDetails.UseRegExp)
	if !fileutils.IsPathExists(rootPath, false) {
		err := errorutils.CheckError(errors.New("Path does not exist: " + rootPath))
		if err != nil {
			return nil, err
		}
	}
	localPath := clientutils.ReplaceTildeWithUserHome(uploadDetails.Pattern)
	localPath = clientutils.PrepareLocalPathForUpload(localPath, uploadDetails.UseRegExp)

	artifacts := []clientutils.Artifact{}
	// If the path is a single file then return it
	dir, err := fileutils.IsDirExists(rootPath, false)
	if err != nil {
		return nil, err
	}

	if !dir {
		artifact := getSingleFileToUpload(rootPath, uploadDetails.TargetPath, uploadDetails.Flat)
		return append(artifacts, artifact), nil
	}

	r, err := regexp.Compile(localPath)
	err = errorutils.CheckError(err)
	if err != nil {
		return nil, err
	}

	log.Info("Collecting files for upload...")
	paths, err := us.listFiles(uploadDetails.Recursive, rootPath)
	if err != nil {
		return nil, err
	}

	for _, filePath := range paths {
		dir, err := fileutils.IsDirExists(filePath, false)
		if err != nil {
			return nil, err
		}
		if dir {
			continue
		}

		groups := r.FindStringSubmatch(filePath)
		size := len(groups)
		target := uploadDetails.TargetPath

		if size > 0 {
			for i := 1; i < size; i++ {
				group := strings.Replace(groups[i], "\\", "/", -1)
				target = strings.Replace(target, "{"+strconv.Itoa(i)+"}", group, -1)
			}

			if target == "" || strings.HasSuffix(target, "/") {
				if target == "" {
					target = debianDefaultPath
				}
				if uploadDetails.Flat {
					fileName, _ := fileutils.GetFileAndDirFromPath(filePath)
					target += fileName
				} else {
					uploadPath := clientutils.TrimPath(filePath)
					target += uploadPath
				}
			}

			artifacts = append(artifacts, clientutils.Artifact{LocalPath: filePath, TargetPath: target, Symlink: ""})
		}
	}
	return artifacts, nil
}

func (us *UploadService) listFiles(recursive bool, rootPath string) ([]string, error) {
	if recursive {
		return fileutils.ListFilesRecursiveWalkIntoDirSymlink(rootPath, false)
	}
	return fileutils.ListFiles(rootPath, false)
}
