# jfrog-client-go

| Branch |                                                                                        Status                                                                                         |
| :----: | :-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |
| master | [![Build status](https://ci.appveyor.com/api/projects/status/2wkemson2sj4skyh/branch/master?svg=true)](https://ci.appveyor.com/project/jfrog-ecosystem/jfrog-client-go/branch/master) |
|  dev   |    [![Build status](https://ci.appveyor.com/api/projects/status/2wkemson2sj4skyh/branch/dev?svg=true)](https://ci.appveyor.com/project/jfrog-ecosystem/jfrog-client-go/branch/dev)    |

## Table of Contents
- [jfrog-client-go](#jfrog-client-go)
  - [Table of Contents](#table-of-contents)
  - [General](#general)
  - [Pull Requests](#pull-requests)
    - [Guidelines](#guidelines)
  - [Tests](#tests)
  - [General APIs](#general-apis)
    - [Setting the Logger](#setting-the-logger)
    - [Setting the Temp Dir](#setting-the-temp-dir)
  - [Artifactory APIs](#artifactory-apis)
    - [Creating Artifactory Service Manager](#creating-artifactory-service-manager)
      - [Creating Artifactory Details](#creating-artifactory-details)
      - [Creating Artifactory Service Config](#creating-artifactory-service-config)
      - [Creating New Artifactory Service Manager](#creating-new-artifactory-service-manager)
    - [Using Artifactory Services](#using-artifactory-services)
      - [Uploading Files to Artifactory](#uploading-files-to-artifactory)
        - [Uploading Files:](#uploading-files)
        - [Uploading Files with Results Reader:](#uploading-files-with-results-reader)
      - [Downloading Files from Artifactory](#downloading-files-from-artifactory)
        - [Downloading Files:](#downloading-files)
        - [Downloading Files with Results Reader:](#downloading-files-with-results-reader)
      - [Copying Files in Artifactory](#copying-files-in-artifactory)
      - [Moving Files in Artifactory](#moving-files-in-artifactory)
      - [Deleting Files from Artifactory](#deleting-files-from-artifactory)
      - [Searching Files in Artifactory](#searching-files-in-artifactory)
      - [Setting Properties on Files in Artifactory](#setting-properties-on-files-in-artifactory)
      - [Deleting Properties from Files in Artifactory](#deleting-properties-from-files-in-artifactory)
      - [Publishing Build Info to Artifactory](#publishing-build-info-to-artifactory)
      - [Fetching Build Info from Artifactory](#fetching-build-info-from-artifactory)
      - [Promoting Published Builds in Artifactory](#promoting-published-builds-in-artifactory)
      - [Promoting a Docker Image in Artifactory](#promoting-a-docker-image-in-artifactory)
      - [Distributing Published Builds to JFrog Bintray](#distributing-published-builds-to-jfrog-bintray)
      - [Triggering Build Scanning with JFrog Xray](#triggering-build-scanning-with-jfrog-xray)
      - [Discarding Old Builds](#discarding-old-builds)
      - [Cleaning Unreferenced Git LFS Files from Artifactory](#cleaning-unreferenced-git-lfs-files-from-artifactory)
      - [Executing AQLs](#executing-aqls)
      - [Reading Files in Artifactory](#reading-files-in-artifactory)
      - [Creating an Access Token](#creating-an-access-token)
      - [Fetching Access Tokens](#fetching-access-tokens)
      - [Refreshing an Access Token](#refreshing-an-access-token)
      - [Revoking an Access Token](#revoking-an-access-token)
      - [Regenerate API Key](#regenerate-api-key)
      - [Creating and Updating Local Repository](#creating-and-updating-local-repository)
      - [Creating and Updating Remote Repository](#creating-and-updating-remote-repository)
      - [Creating and Updating Virtual Repository](#creating-and-updating-virtual-repository)
      - [Removing a Repository](#removing-a-repository)
      - [Getting Repository Details](#getting-repository-details)
      - [Getting All Repositories](#getting-all-repositories)
      - [Creating and Updating Repository Replications](#creating-and-updating-repository-replications)
      - [Getting a Repository Replication](#getting-a-repository-replication)
      - [Removing a Repository Replication](#removing-a-repository-replication)
      - [Creating and Updating Permission Targets](#creating-and-updating-permission-targets)
      - [Removing a Permission Target](#removing-a-permission-target)
      - [Fetching Artifactory's Version](#fetching-artifactorys-version)
      - [Fetching Artifactory's Service ID](#fetching-artifactorys-service-id)
  - [Distribution APIs](#distribution-apis)
    - [Creating Distribution Service Manager](#creating-distribution-service-manager)
      - [Creating Distribution Details](#creating-distribution-details)
      - [Creating Distribution Service Config](#creating-distribution-service-config)
      - [Creating New Distribution Service Manager](#creating-new-distribution-service-manager)
    - [Using Distribution Services](#using-distribution-services)
      - [Setting Distribution Signing Key](#setting-distribution-signing-key)
      - [Creating a Release Bundle](#creating-a-release-bundle)
      - [Updating a Release Bundle](#updating-a-release-bundle)
      - [Signing a Release Bundle](#signing-a-release-bundle)
      - [Async Distributing a Release Bundle](#async-distributing-a-release-bundle)
      - [Sync Distributing a Release Bundle](#sync-distributing-a-release-bundle)
      - [Getting Distribution Status](#getting-distribution-status)
      - [Deleting a Remote Release Bundle](#deleting-a-remote-release-bundle)
      - [Deleting a Local Release Bundle](#deleting-a-local-release-bundle)
  - [Bintray APIs](#bintray-apis)
    - [Creating Bintray Details](#creating-bintray-details)
    - [Creating Bintray Service Manager](#creating-bintray-service-manager)
    - [Using Bintray Services](#using-bintray-services)
      - [Uploading a Single File to Bintray](#uploading-a-single-file-to-bintray)
      - [Downloading a Single File from Bintray](#downloading-a-single-file-from-bintray)
      - [Downloading Version Files from Bintray](#downloading-version-files-from-bintray)
      - [Showing and Deleting a Bintray Package](#showing-and-deleting-a-bintray-package)
      - [Creating and Updating a Bintray Package](#creating-and-updating-a-bintray-package)
      - [Showing and Deleting a Bintray Version](#showing-and-deleting-a-bintray-version)
      - [Creating and Updating a Bintray Version](#creating-and-updating-a-bintray-version)
      - [Creating and Updating Entitlements](#creating-and-updating-entitlements)
      - [Showing and Deleting Entitlements](#showing-and-deleting-entitlements)
      - [Creating and Updating Access Keys](#creating-and-updating-access-keys)
      - [Showing and Deleting Access Keys](#showing-and-deleting-access-keys)
      - [Signing a URL](#signing-a-url)
      - [GPG Signing a File](#gpg-signing-a-file)
      - [GPG Signing Version Files](#gpg-signing-version-files)
      - [Listing Logs](#listing-logs)
      - [Downloading Logs](#downloading-logs)
      - [Syncing Content To Maven Central](#syncing-content-to-maven-central)
  - [Using ContentReader](#using-contentreader)
  - [Xray APIs](#xray-apis)
    - [Creating Xray Service Manager](#creating-xray-service-manager)
      - [Creating Xray Details](#creating-xray-details)
      - [Creating Xray Service Config](#creating-xray-service-config)
      - [Creating New Xray Service Manager](#creating-new-xray-service-manager)
    - [Using Xray Services](#using-xray-services)
      - [Creating an Xray Watch](#creating-an-xray-watch)
      - [Get an Xray Watch](#get-an-xray-watch)
      - [Update an Xray Watch](#update-an-xray-watch)
      - [Delete an Xray Watch](#delete-an-xray-watch)

## General
_jfrog-client-go_ is a library which provides Go APIs to performs actions on JFrog Artifactory or Bintray from your Go application.
The project is still relatively new, and its APIs may therefore change frequently between releases.
The library can be used as a go-module, which should be added to your project's go.mod file. As a reference you may look at [JFrog CLI](https://github.com/jfrog/jfrog-cli-go)'s [go.mod](https://github.com/jfrog/jfrog-cli-go/blob/master/go.mod) file, which uses this library as a dependency.

## Pull Requests
We welcome pull requests from the community.

### Guidelines
- If the existing tests do not already cover your changes, please add tests.
- Pull requests should be created on the **dev** branch.
- Please use gofmt for formatting the code before submitting the pull request.

## Tests
To run tests on the source code, you'll need a running JFrog Artifactory Pro instance.
Use the following command with the below options to run the tests.
```sh
go test -v github.com/jfrog/jfrog-client-go/tests
```
Optional flags:

| Flag                | Description                                                                                            |
| ------------------- | ------------------------------------------------------------------------------------------------------ |
| `-rt.url`           | [Default: http://localhost:8081/artifactory] Artifactory URL.                                          |
| `-rt.user`          | [Default: admin] Artifactory username.                                                                 |
| `-rt.password`      | [Default: password] Artifactory password.                                                              |
| `-rt.distUrl`       | [Optional] JFrog Distribution URL.                                                                     |
| `-rt.xrayUrl`       | [Optional] JFrog Xray URL.                                                                     |
| `-rt.apikey`        | [Optional] Artifactory API key.                                                                        |
| `-rt.sshKeyPath`    | [Optional] Ssh key file path. Should be used only if the Artifactory URL format is ssh://[domain]:port |
| `-rt.sshPassphrase` | [Optional] Ssh key passphrase.                                                                         |
| `-rt.accessToken`   | [Optional] Artifactory access token.                                                                   |
| `-log-level`        | [Default: INFO] Sets the log level.                                                                    |

- The tests create an Artifactory repository named _jfrog-client-tests-repo1_.<br/>
  Once the tests are completed, the content of this repository is deleted.

## General APIs
### Setting the Logger
```go
var file *os.File
...
log.SetLogger(log.NewLogger(log.INFO, file))
```

### Setting the Temp Dir
The default temp dir used is 'os.TempDir()'. Use the following API to set a new temp dir:
```go
fileutils.SetTempDirBase(filepath.Join("my", "temp", "path"))
```

## Artifactory APIs
### Creating Artifactory Service Manager
#### Creating Artifactory Details
```go
rtDetails := auth.NewArtifactoryDetails()
rtDetails.SetUrl("http://localhost:8081/artifactory")
rtDetails.SetSshKeysPath("path/to/.ssh/")
rtDetails.SetApiKey("apikey")
rtDetails.SetUser("user")
rtDetails.SetPassword("password")
rtDetails.SetAccessToken("accesstoken")
// if client certificates are required
rtDetails.SetClientCertPath("path/to/.cer")
rtDetails.SetClientCertKeyPath("path/to/.key")
```

#### Creating Artifactory Service Config
```go
serviceConfig, err := config.NewConfigBuilder().
    SetServiceDetails(rtDetails).
    SetCertificatesPath(certPath).
    SetThreads(threads).
    SetDryRun(false).
    // Add [Context](https://golang.org/pkg/context/)
    SetContext(ctx).
    Build()
```

#### Creating New Artifactory Service Manager
```go
rtManager, err := artifactory.New(&rtDetails, serviceConfig)
```

### Using Artifactory Services
#### Uploading Files to Artifactory
##### Uploading Files:
Using the `UploadFiles()` function, we can upload files and get the general statistics of the action (The actual number of successful and failed uploads), and the error value if it occurred.
```go
params := services.NewUploadParams()
params.Pattern = "repo/*/*.zip"
params.Target = "repo/path/"
// Attach properties to the uploaded files.
params.Props = "key1=val1;key2=val2"
params.AddVcsProps = false
params.BuildProps = "build.name=buildName;build.number=17;build.timestamp=1600856623553"
params.Recursive = true
params.Regexp = false
params.IncludeDirs = false
params.Flat = true
params.Explode = false
params.Deb = ""
params.Symlink = false
// Retries default value: 3
params.Retries = 5
// MinChecksumDeploy default value: 10400
params.MinChecksumDeploy = 15360

totalUploaded, totalFailed, err := rtManager.UploadFiles(params)
```
##### Uploading Files with Results Reader:
Similar to `UploadFlies()`, but returns a reader, which allows iterating over the details of the uploaded files. Only files which were successfully uploaded are available by the reader.
```go
params := services.NewUploadParams()
params.Pattern = "repo/*/*.zip"
params.Target = "repo/path/"
// Attach properties to the uploaded files.
params.Props = "key1=val1;key2=val2"
params.AddVcsProps = false
params.BuildProps = "build.name=buildName;build.number=17;build.timestamp=1600856623553"
params.Recursive = true
params.Regexp = false
params.IncludeDirs = false
params.Flat = true
params.Explode = false
params.Deb = ""
params.Symlink = false
// Retries default value: 3
params.Retries = 5
// MinChecksumDeploy default value: 10400
params.MinChecksumDeploy = 15360

reader, totalUploaded, totalFailed, err := rtManager.UploadFilesWithResultReader(params)
```
Read more about [ContentReader](#using-contentReader).

#### Downloading Files from Artifactory
##### Downloading Files:
Using the `DownloadFiles()` function, we can download files and get the general statistics of the action (The actual number of files downloaded, and the number of files we expected to download), and the error value if it occurred.
```go
params := services.NewDownloadParams()
params.Pattern = "repo/*/*.zip"
params.Target = "target/path/"
// Filter the downloaded files by properties.
params.Props = "key1=val1;key2=val2"
params.Recursive = true
params.IncludeDirs = false
params.Flat = false
params.Explode = false
params.Symlink = true
params.ValidateSymlink = false
// Retries default value: 3
params.Retries = 5
// SplitCount default value: 3
params.SplitCount = 2
// MinSplitSize default value: 5120
params.MinSplitSize = 7168

totalDownloaded, totalExpected, err := rtManager.DownloadFiles(params)
```

##### Downloading Files with Results Reader:
Similar to `DownloadFiles()`, but returns a reader, which allows iterating over the details of the downloaded files. Only files which were successfully downloaded are available by the reader.
```go
params := services.NewDownloadParams()
params.Pattern = "repo/*/*.zip"
params.Target = "target/path/"
// Filter the downloaded files by properties.
params.Props = "key1=val1;key2=val2"
params.Recursive = true
params.IncludeDirs = false
params.Flat = false
params.Explode = false
params.Symlink = true
params.ValidateSymlink = false
params.Retries = 5
params.SplitCount = 2
params.MinSplitSize = 7168

reader, totalDownloaded, totalExpected, err := rtManager.DownloadFilesWithResultReader(params)
```
Read more about [ContentReader](#using-contentReader).

#### Copying Files in Artifactory
```go
params := services.NewMoveCopyParams()
params.Pattern = "repo/*/*.zip"
params.Target = "target/path/"
// Filter the files by properties.
params.Props = "key1=val1;key2=val2"
params.Recursive = true
params.Flat = false

rtManager.Copy(params)
```

#### Moving Files in Artifactory
```go
params := services.NewMoveCopyParams()
params.Pattern = "repo/*/*.zip"
params.Target = "target/path/"
// Filter the files by properties.
params.Props = "key1=val1;key2=val2"
params.Recursive = true
params.Flat = false

rtManager.Move(params)
```

#### Deleting Files from Artifactory
```go
params := services.NewDeleteParams()
params.Pattern = "repo/*/*.zip"
// Filter the files by properties.
params.Props = "key1=val1;key2=val2"
params.Recursive = true

pathsToDelete, err := rtManager.GetPathsToDelete(params)
if err != nil {
    return err
}
defer pathsToDelete.Close()
rtManager.DeleteFiles(pathsToDelete)
```
Read more about [ContentReader](#using-contentReader).

#### Searching Files in Artifactory
```go
params := services.NewSearchParams()
params.Pattern = "repo/*/*.zip"
// Filter the files by properties.
params.Props = "key1=val1;key2=val2"
params.Recursive = true

reader, err := rtManager.SearchFiles(params)
if err != nil {
    return err
}
defer reader.Close()
```
Read more about [ContentReader](#using-contentReader).

#### Setting Properties on Files in Artifactory
```go
searchParams = services.NewSearchParams()
searchParams.Recursive = true
searchParams.IncludeDirs = false

reader, err = rtManager.SearchFiles(searchParams)
if err != nil {
    return err
}
defer reader.Close()
propsParams = services.NewPropsParams()
propsParams.Pattern = "repo/*/*.zip"
propsParams.Reader = reader
// Filter the files by properties.
propsParams.Props = "key=value"

rtManager.SetProps(propsParams)
```
Read more about [ContentReader](#using-contentReader).

#### Deleting Properties from Files in Artifactory
```go
searchParams = services.NewSearchParams()
searchParams.Recursive = true
searchParams.IncludeDirs = false

resultItems, err = rtManager.SearchFiles(searchParams)
if err != nil {
    return err
}
defer reader.Close()
propsParams = services.NewPropsParams()
propsParams.Pattern = "repo/*/*.zip"
propsParams.Reader = reader
// Filter the files by properties.
propsParams.Props = "key=value"

rtManager.DeleteProps(propsParams)
```
Read more about [ContentReader](#using-contentReader).

#### Publishing Build Info to Artifactory
```go
buildInfo := &buildinfo.BuildInfo{}
// Optional Artifactory project name
project := "my-project"
...
rtManager.PublishBuildInfo(buildInfo, project)
```

#### Fetching Build Info from Artifactory
```go
buildInfoParams := services.NewBuildInfoParams{}
buildInfoParams.BuildName = "buildName"
buildInfoParams.BuildNumber = "LATEST"

rtManager.GetBuildInfo(buildInfoParams)
```

#### Promoting Published Builds in Artifactory
```go
params := services.NewPromotionParams()
params.BuildName = "buildName"
params.BuildNumber = "10"
params.TargetRepo = "target-repo"
params.Status = "status"
params.Comment = "comment"
params.Copy = true
params.IncludeDependencies = false
params.SourceRepo = "source-repo"

rtManager.DownloadFiles(params)
```

#### Promoting a Docker Image in Artifactory
```go
sourceDockerImage := "hello-world"
sourceRepo := "docker-local-1"
targetRepo := "docker-local-2"
params := services.NewDockerPromoteParams(sourceDockerImage, sourceRepo, targetRepo)

// Optional parameters:
params.TargetDockerImage = "target-docker-image"
params.SourceTag = "42"
params.TargetTag = "43"
params.Copy = true

rtManager.PromoteDocker(params)
```

#### Distributing Published Builds to JFrog Bintray
```go
params := services.NewBuildDistributionParams()
params.SourceRepos = "source-repo"
params.TargetRepo = "target-repo"
params.GpgPassphrase = "GpgPassphrase"
params.Publish = false
params.OverrideExistingFiles = false
params.Async = true
params.BuildName = "buildName"
params.BuildNumber = "10"
params.Pattern = "repo/*/*.zip"

rtManager.DistributeBuild(params)
```

#### Triggering Build Scanning with JFrog Xray
```go
params := services.NewXrayScanParams()
params.BuildName = buildName
params.BuildNumber = buildNumber

rtManager.XrayScanBuild(params)
```

#### Discarding Old Builds
```go
params := services.NewDiscardBuildsParams()
params.BuildName = "buildName"
params.MaxDays = "max-days"
params.MaxBuilds = "max-builds"
params.ExcludeBuilds = "1,2"
params.DeleteArtifacts = false
params.Async = false

rtManager.DiscardBuilds(params)
```

#### Cleaning Unreferenced Git LFS Files from Artifactory
```go
params := services.NewGitLfsCleanParams()
params.Refs = "refs/remotes/*"
params.Repo = "my-project-lfs"
params.GitPath = "path/to/git"

reader,err := rtManager.GetUnreferencedGitLfsFiles(params)

defer reader.Close()
rtManager.DeleteFiles(reader)
```

#### Executing AQLs
```go
rtManager.Aql(aql string)
```

#### Reading Files in Artifactory
```go
rtManager.ReadRemoteFile(FilePath string)
```

#### Creating an Access Token
```go
params := services.NewCreateTokenParams()
params.Scope = "api:* member-of-groups:readers"
params.Username = "user"
params.ExpiresIn = 3600 // default -1 (use server default)
params.GrantType = "client_credentials"
params.Refreshable = true
params.Audience = "jfrt@<serviceID1> jfrt@<serviceID2>"

results, err := rtManager.CreateToken(params)
```

#### Fetching Access Tokens
```go
results, err := rtManager.GetTokens()
```

#### Refreshing an Access Token
```go
params := services.NewRefreshTokenParams()
params.AccessToken = "<access token>"
params.RefreshToken = "<refresh token>"
params.Token.Scope = "api:*"
params.Token.ExpiresIn = 3600
results, err := rtManager.RefreshToken(params)
```

#### Revoking an Access Token
```go
params := services.NewRevokeTokenParams()

// Provide either TokenId or Token
params.TokenId = "<token id>"
// params.Token = "access token"

err := rtManager.RevokeToken(params)
```

#### Regenerate API Key
```go
apiKey, err := rtManager.RegenerateAPIKey()
```

#### Creating and Updating Local Repository
You can create and update a local repository for the following package types:

Maven, Gradle, Ivy, Sbt, Helm, Cocoapods, Opkg, Rpm, Nuget, Cran, Gems, Npm, Bower, Debian, Composer, Pypi, Docker,
Vagrant, Gitlfs, Go, Yum, Conan, Chef, Puppet and Generic.

Each package type has it's own parameters struct, can be created using the method
`New<packageType>LocalRepositoryParams()`.

Example for creating local Generic repository:
```go
params := services.NewGenericLocalRepositoryParams()
params.Key = "generic-repo"
params.Description = "This is a public description for generic-repo"
params.Notes = "These are internal notes for generic-repo"
params.RepoLayoutRef = "simple-default"
params.ArchiveBrowsingEnabled = true
params.XrayIndex = true
params.IncludesPattern = "**/*"
params.ExcludesPattern = "excludedDir/*"
params.DownloadRedirect = true

err = servicesManager.CreateLocalRepository().Generic(params)
```

Updating local Generic repository:
```go
err = servicesManager.UpdateLocalRepository().Generic(params)
```

#### Creating and Updating Remote Repository
You can create and update a remote repository for the following package types:

Maven, Gradle, Ivy, Sbt, Helm, Cocoapods, Opkg, Rpm, Nuget, Cran, Gems, Npm, Bower, Debian, Composer, Pypi, Docker,
Gitlfs, Go, Yum, Conan, Chef, Puppet, Conda, P2, Vcs and Generic.

Each package type has it's own parameters struct, can be created using the method
`New<packageType>RemoteRepositoryParams()`.

Example for creating remote Maven repository:
```go
params := services.NewMavenRemoteRepositoryParams()
params.Key = "jcenter-remote"
params.Url = "http://jcenter.bintray.com"
params.RepoLayoutRef = "maven-2-default"
params.Description = "A caching proxy repository for a JFrog's jcenter"
params.HandleSnapshot = false
params.HandleReleases = true
params.FetchJarsEagerly = true
params.XrayIndex = true
params.AssumedOfflinePeriodSecs = 600
params.SuppressPomConsistencyChecks = true
params.RemoteRepoChecksumPolicyType = "pass-thru"

err = servicesManager.CreateRemoteRepository().Maven(params)
```

Updating remote Maven repository:
```go
err = servicesManager.UpdateRemoteRepository().Maven(params)
```

#### Creating and Updating Virtual Repository
You can create and update a virtual repository for the following package types:

Maven, Gradle, Ivy, Sbt, Helm, Rpm, Nuget, Cran, Gems, Npm, Bower, Debian, Pypi, Docker, Gitlfs, Go, Yum, Conan,
Chef, Puppet, Conda, P2 and Generic

Each package type has it's own parameters struct, can be created using the method
`New<packageType>VirtualRepositoryParams()`.

Example for creating virtual Go repository:
```go
params := services.NewGoVirtualRepositoryParams()
params.Description = "This is an aggregated repository for several go repositories"
params.RepoLayoutRef = "go-default"
params.Repositories = {"gocenter-remote", "go-local"}
params.DefaultDeploymentRepo = "go-local"
params.ExternalDependenciesEnabled = true
params.ExternalDependenciesPatterns = {"**/github.com/**", "**/golang.org/**", "**/gopkg.in/**"}
params.ArtifactoryRequestsCanRetrieveRemoteArtifacts = true

err = servicesManager.CreateVirtualRepository().Go(params)
```

Updating remote Maven repository:
```go
err = servicesManager.UpdateVirtualRepository().Go(params)
```

#### Removing a Repository
You can remove a repository from Artifactory using its key:
```go
servicesManager.DeleteRepository("generic-repo")
```

#### Getting Repository Details
You can get repository details from Artifactory using its key:
```go
servicesManager.GetRepository("generic-repo")
```

#### Getting All Repositories
You can get all repositories from Artifactory:
```go
servicesManager.GetAllRepositories()
```

#### Creating and Updating Repository Replications
Example of creating a repository replication:
```go
params := services.NewCreateReplicationParams()
// Source replication repository.
params.RepoKey = "my-repository"
params.CronExp = "0 0 12 * * ?"
params.Username = "admin"
params.Password = "password"
params.Url = "http://localhost:8081/artifactory/remote-repo"
params.Enabled = true
params.SocketTimeoutMillis = 15000
params.EnableEventReplication = true
params.SyncDeletes = true
params.SyncProperties = true
params.SyncStatistics = true
params.PathPrefix = "/path/to/repo"

err = servicesManager.CreateReplication(params)
```

Example of updating a local repository replication:
```go
params := services.NewUpdateReplicationParams()
// Source replication repository.
params.RepoKey = "my-repository"
params.CronExp = "0 0 12 * * ?"
params.Enabled = true
params.SocketTimeoutMillis = 15000
params.EnableEventReplication = true
params.SyncDeletes = true
params.SyncProperties = true
params.SyncStatistics = true
params.PathPrefix = "/path/to/repo"

err = servicesManager.UpdateReplication(params)
```

#### Getting a Repository Replication
You can get a repository replication configuration from Artifactory using its key:
```go
replicationConfiguration, err := servicesManager.GetReplication("my-repository")
```

#### Removing a Repository Replication
You can remove a repository replication configuration from Artifactory using its key:
```go
err := servicesManager.DeleteReplication("my-repository")
```

#### Creating and Updating Permission Targets
You can create or update a permission target in Artifactory.
Permissions are set according to the following conventions:
`read, write, annotate, delete, manage, managedXrayMeta, distribute`
For repositories You can specify the name `"ANY"` in order to apply to all repositories, `"ANY REMOTE"` for all remote repositories or `"ANY LOCAL"` for all local repositories.

Creating a new permission target :
```go
params := services.NewPermissionTargetParams()
params.Name = "java-developers"
params.Repo.Repositories = []string{"ANY REMOTE", "local-repo1", "local-repo2"}
params.Repo.ExcludePatterns = []string{"dir/*"}
params.Repo.Actions.Users = map[string][]string {
	"user1" : {"read", "write"},
    "user2" : {"write","annotate", "read"},
}
params.Repo.Actions.Groups = map[string][]string {
	"group1" : {"manage","read","annotate"},
}
// This is the default value that cannot be changed
params.Build.Repositories = []string{"artifactory-build-info"}
params.Build.Actions.Groups = map[string][]string {
	"group1" : {"manage","read","write","annotate","delete"},
	"group2" : {"read"},

}

err = servicesManager.CreatePermissionTarget(params)
```
Updating an existing permission target :
```go
err = servicesManager.UpdatePermissionTarget(params)
```

#### Removing a Permission Target
You can remove a permission target from Artifactory using its name:
```go
servicesManager.DeletePermissionTarget("java-developers")
```

#### Fetching Artifactory's Version
```go
version, err := servicesManager.GetVersion()
```

#### Fetching Artifactory's Service ID
```go
serviceId, err := servicesManager.GetServiceId()
```

## Distribution APIs
### Creating Distribution Service Manager
#### Creating Distribution Details
```go
distDetails := auth.NewDistributionDetails()
distDetails.SetUrl("http://localhost:8081/distribution")
distDetails.SetSshKeysPath("path/to/.ssh/")
distDetails.SetApiKey("apikey")
distDetails.SetUser("user")
distDetails.SetPassword("password")
distDetails.SetAccessToken("accesstoken")
// if client certificates are required
distDetails.SetClientCertPath("path/to/.cer")
distDetails.SetClientCertKeyPath("path/to/.key")
```

#### Creating Distribution Service Config
```go
serviceConfig, err := config.NewConfigBuilder().
    SetServiceDetails(rtDetails).
    SetCertificatesPath(certPath).
    SetThreads(threads).
    SetDryRun(false).
    // Add [Context](https://golang.org/pkg/context/)
    SetContext(ctx).
    Build()
```

#### Creating New Distribution Service Manager
```go
distManager, err := distribution.New(&distDetails, serviceConfig)
```

### Using Distribution Services
#### Setting Distribution Signing Key
```go
params := services.NewSetSigningKeyParams("private-gpg-key", "public-gpg-key")

err := distManager.SetSigningKey(params)
```

#### Creating a Release Bundle
```go
params := services.NewCreateReleaseBundleParams("bundle-name", "1")
params.SpecFiles = []*utils.ArtifactoryCommonParams{{Pattern: "repo/*/*.zip"}}
params.Description = "Description"
params.ReleaseNotes = "Release notes"
params.ReleaseNotesSyntax = "plain_text"

err := distManager.CreateReleaseBundle(params)
```

#### Updating a Release Bundle
```go
params := services.NewUpdateReleaseBundleParams("bundle-name", "1")
params.SpecFiles = []*utils.ArtifactoryCommonParams{{Pattern: "repo/*/*.zip"}}
params.Description = "New Description"
params.ReleaseNotes = "New Release notes"
params.ReleaseNotesSyntax = "plain_text"

err := distManager.CreateReleaseBundle(params)
```

#### Signing a Release Bundle
```go
params := services.NewSignBundleParams("bundle-name", "1")
params.GpgPassphrase = "123456"

err := distManager.SignReleaseBundle(params)
```

#### Async Distributing a Release Bundle
```go
params := services.NewDistributeReleaseBundleParams("bundle-name", "1")
distributionRules := utils.DistributionCommonParams{SiteName: "Swamp-1", "CityName": "Tel-Aviv", "CountryCodes": []string{"123"}}}
params.DistributionRules = []*utils.DistributionCommonParams{distributionRules}

err := distManager.DistributeReleaseBundle(params)
```

#### Sync Distributing a Release Bundle
```go
params := services.NewDistributeReleaseBundleParams("bundle-name", "1")
distributionRules := utils.DistributionCommonParams{SiteName: "Swamp-1", "CityName": "Tel-Aviv", "CountryCodes": []string{"123"}}}
params.DistributionRules = []*utils.DistributionCommonParams{distributionRules}
// Wait up to 120 minutes for the release bundle distribution
err := distManager.DistributeReleaseBundleSync(params, 120)
```

#### Getting Distribution Status
```go
params := services.NewDistributionStatusParams()
// Optional parameters:
// If missing, get status for all distributions
params.Name = "bundle-name"
// If missing, get status for all versions of "bundle-name"
params.Version = "1"
// If missing, get status for all "bundle-name" with version "1"
params.TrackerId = "123456789"

status, err := distributeBundleService.GetStatus(params)
```

#### Deleting a Remote Release Bundle
```go
params := services.NewDeleteReleaseBundleParams("bundle-name", "1")
params.DeleteFromDistribution = true
distributionRules := utils.DistributionCommonParams{SiteName: "Swamp-1", "CityName": "Tel-Aviv", "CountryCodes": []string{"123"}}}
params.DistributionRules = []*utils.DistributionCommonParams{distributionRules}

err := distManager.DeleteReleaseBundle(params)
```

#### Deleting a Local Release Bundle
```go
params := services.NewDeleteReleaseBundleParams("bundle-name", "1")

err := distManager.DeleteLocalReleaseBundle(params)
```

## Bintray APIs
### Creating Bintray Details
```go
btDetails := auth.NewBintrayDetails()
btDetails.SetUser("user")
btDetails.SetKey("key")
btDetails.SetDefPackageLicense("Apache 2.0")
```

### Creating Bintray Service Manager
```go
serviceConfig := bintray.NewConfigBuilder().
    SetBintrayDetails(btDetails).
    SetDryRun(false).
    SetThreads(threads).
    Build()

btManager, err := bintray.New(serviceConfig)
```

### Using Bintray Services
#### Uploading a Single File to Bintray
```go
params := services.NewUploadParams()
params.Pattern = "*/*.zip"
params.Path = versions.CreatePath("subject/repo/pkg/version")
params.TargetPath = "path/to/files"
params.Deb = "distribution/component/architecture"
params.Recursive = true
params.Flat = true
params.Publish = false
params.Override = false
params.Explode = false
params.UseRegExp = false
params.ShowInDownloadList = false

btManager.UploadFiles(params)
```

#### Downloading a Single File from Bintray
```go
params := services.NewDownloadFileParams()
params.Flat = false
params.IncludeUnpublished = false
params.PathDetails = "path/to/file"
params.TargetPath = "target/path/"
// SplitCount default value: 3
params.SplitCount = 2
// MinSplitSize default value: 5120
params.MinSplitSize = 7168

btManager.DownloadFile(params)
```

#### Downloading Version Files from Bintray
```go
params := services.NewDownloadVersionParams()
params.Path, err = versions.CreatePath("subject/repo/pkg/version")

params.IncludeUnpublished = false
params.TargetPath = "target/path/"

btManager.DownloadVersion(params)
```

#### Showing and Deleting a Bintray Package
```go
pkgPath, err := packages.CreatePath("subject/repo/pkg")

btManager.ShowPackage(pkgPath)
btManager.DeletePackage(pkgPath)
```

#### Creating and Updating a Bintray Package
```go
params := packages.NewPackageParams()
params.Path, err = packages.CreatePath("subject/repo/pkg")

params.Desc = "description"
params.Labels = "labels"
params.Licenses = "licences"
params.CustomLicenses = "custum-licenses"
params.VcsUrl = "https://github.com/jfrog/jfrog-cli-go"
params.WebsiteUrl = "https://jfrog.com"
params.IssueTrackerUrl = "https://github.com/bintray/bintray-client-java/issues"
params.GithubRepo = "bintray/bintray-client-java"
params.GithubReleaseNotesFile = "RELEASE_1.2.3.txt" "github-rel-notes"
params.PublicDownloadNumbers = "true"
params.PublicStats = "true"

btManager.CreatePackage(params)
btManager.UpdatePackage(params)
```

#### Showing and Deleting a Bintray Version
```go
versionPath, err := versions.CreatePath("subject/repo/pkg/version")

btManager.ShowVersion(versionPath)
btManager.DeleteVersion(versionPath)
```

#### Creating and Updating a Bintray Version
```go
params := versions.NewVersionParams()
params.Path, err = versions.CreatePath("subject/repo/pkg/version")

params.Desc = "description"
params.VcsTag = "1.1.5"
params.Released = "true"
params.GithubReleaseNotesFile = "RELEASE_1.2.3.txt"
params.GithubUseTagReleaseNotes = "false"

btManager.CreateVersion(params)
btManager.UpdateVersion(params)
```

#### Creating and Updating Entitlements
```go
params := entitlements.NewEntitlementsParams()
params.VersionPath, err = versions.CreatePath("subject/repo/pkg/version")

params.Path = "a/b/c"
params.Access = "rw"
params.Keys = "keys"

btManager.CreateEntitlement(params)

params.Id = "entitlementID"
btManager.UpdateEntitlement(params)
```

#### Showing and Deleting Entitlements
```go
versionPath, err := versions.CreatePath("subject/repo/pkg/version")

btManager.ShowAllEntitlements(versionPath)
btManager.ShowEntitlement("entitelmentID", versionPath)
btManager.DeleteEntitlement("entitelmentID", versionPath)
```

#### Creating and Updating Access Keys
```go
params := accesskeys.NewAccessKeysParams()
params.Password = "password"
params.Org = "org"
params.Expiry = time.Now() + time.Hour * 10
params.ExistenceCheckUrl = "http://callbacks.myci.org/username=:username,password=:password"
params.ExistenceCheckCache = 60
params.WhiteCidrs = "127.0.0.1/22,193.5.0.1/92"
params.BlackCidrs = "127.0.0.1/22,193.5.0.1/92"
params.ApiOnly = true

btManager.CreateAccessKey(params)

params.Id = "KeyID"
btManager.UpdateAccessKey(params)
```

#### Showing and Deleting Access Keys
```go
btManager.ShowAllAccessKeys("org")
btManager.ShowAccessKey("org", "KeyID")
btManager.DeleteAccessKey("org", "KeyID")
```

#### Signing a URL
```go
params := url.NewURLParams()
params.PathDetails, err = utils.CreatePathDetails("subject/repository/file-path")
// Check for errors
params.Expiry = time.Now() + time.Hour * 10
params.ValidFor = 60
params.CallbackId = "callback-id"
params.CallbackEmail = "callback-email"
params.CallbackUrl = "callback-url"
params.CallbackMethod = "callback-method"

btManager.SignUrl(params)
```

#### GPG Signing a File
```go
path, err := utils.CreatePathDetails("subject/repository/file-path")

btManager.GpgSignFile(path, "passphrase")
```

#### GPG Signing Version Files
```go
path, err := versions.CreatePath("subject/repo/pkg/version")

btManager.GpgSignVersion(path, "passphrase")
```

#### Listing Logs
```go
path, err := versions.CreatePath("subject/repo/pkg/version")

btManager.LogsList(versionPath)
```

#### Downloading Logs
```go
path, err := versions.CreatePath("subject/repo/pkg/version")

btManager.DownloadLog(path, "logName")
```

#### Syncing Content To Maven Central
```go
params := mavensync.NewParams("user","password", false)
path, err = versions.CreatePath("subject/repo/pkg/version")

btManager.MavenCentralContentSync(params, path)
```

## Using ContentReader
Some APIs return a ```content.ContentReader``` struct, which allows reading the API's output. ```content.ContentReader``` provides access to large amounts of data safely, without loading all of it into the memory.
Here's an example for how ```content.ContentReader``` should be used:

```go
reader, err := servicesManager.SearchFiles(searchParams)
if err != nil {
    return err
}

// Remove the data file used by the reader.
defer func() {
    if reader != nil {
        err = reader.Close()
    }
}()

// Iterate over the results.
for currentResult := new(ResultItem); reader.NextRecord(currentResult) == nil; currentResult = new(ResultItem)  {
    fmt.Printf("Found artifact: %s of type: %s\n", searchResult.Name, searchResult.Type)
}
if err := resultReader.GetError(); err != nil {
    return err
}

// Resets the reader pointer back to the beginning of the output. Make sure not to call this method after the reader had been closed using ```reader.Close()```
reader.Reset()
```

* `reader.NextRecord(currentResult)` reads the next record from the reader into `currentResult` of type `ResultItem`.

* `reader.Close()` removes the file used by the reader after it is used (preferably using `defer`).

* `reader.GetError()` returns any error that might have occurd during `NextRecord()`.

* `reader.Reset()` resets the reader back to the beginning of the output.

## Xray APIs
### Creating Xray Service Manager
#### Creating Xray Details
```go
xrayDetails := auth.NewXrayDetails()
xrayDetails.SetUrl("http://localhost:8081/xray")
xrayDetails.SetSshKeysPath("path/to/.ssh/")
xrayDetails.SetApiKey("apikey")
xrayDetails.SetUser("user")
xrayDetails.SetPassword("password")
xrayDetails.SetAccessToken("accesstoken")
// if client certificates are required
xrayDetails.SetClientCertPath("path/to/.cer")
xrayDetails.SetClientCertKeyPath("path/to/.key")
```

#### Creating Xray Service Config
```go
serviceConfig, err := config.NewConfigBuilder().
    SetServiceDetails(xrayDetails).
    SetCertificatesPath(certPath).
    Build()
```

#### Creating New Xray Service Manager
```go
xrayManager, err := xray.New(&xrayDetails, serviceConfig)
```

### Using Xray Services
#### Creating an Xray Watch

This uses API version 2.

You are able to configure repositories and builds on a watch.
However, bundles are not supported.

```go
params := utils.NewWatchParams()
params.Name = "example-watch-all"
params.Description = "All Repos"
params.Active = true

params.Repositories.Type = utils.WatchRepositoriesAll
params.Repositories.All.Filters.PackageTypes = []string{"Npm", "maven"}
params.Repositories.ExcludePatterns = []string{"excludePath1", "excludePath2"}
params.Repositories.IncludePatterns = []string{"includePath1", "includePath2"}

params.Builds.Type = utils.WatchBuildAll
params.Builds.All.Bin_Mgr_ID = "default"

params.Policies = []utils.AssignedPolicy{
  {
    Name: policy1Name,
    Type: "security",
  },
  {
    Name: policy2Name,
    Type: "security",
  },
}

resp, err := xrayManager.CreateWatch(*params)
```

#### Get an Xray Watch
```go
watch, resp, err := xrayManager.GetWatch("example-watch-all")
```

#### Update an Xray Watch
```go
watch, resp, err := xrayManager.GetWatch("example-watch-all")
watch.Description = "Updated description"

resp, err := xrayManager.UpdateWatch(*watch)
```

#### Delete an Xray Watch
```go
resp, err := xrayManager.DeleteWatch("example-watch-all")
```
