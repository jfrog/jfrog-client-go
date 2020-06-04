# jfrog-client-go

| Branch |                                                                                        Status                                                                                         |
| :----: | :-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |
| master | [![Build status](https://ci.appveyor.com/api/projects/status/2wkemson2sj4skyh/branch/master?svg=true)](https://ci.appveyor.com/project/jfrog-ecosystem/jfrog-client-go/branch/master) |
|  dev   |    [![Build status](https://ci.appveyor.com/api/projects/status/2wkemson2sj4skyh/branch/dev?svg=true)](https://ci.appveyor.com/project/jfrog-ecosystem/jfrog-client-go/branch/dev)    |

## General

_jfrog-client-go_ is a library which provides Go APIs to performs actions on JFrog Artifactory or Bintray from your Go application.
The project is still relatively new, and its APIs may therefore change frequently between releases.
The library can be used as a go-module, which should be added to your project's go.mod file. As a reference you may look at [JFrog CLI](https://github.com/jfrog/jfrog-cli-go)'s [go.mod](https://github.com/jfrog/jfrog-cli-go/blob/master/go.mod) file, which uses this library as a dependency.

## Pull Requests

We welcome pull requests from the community.

### Guidelines

- Before creating your first pull request, please join our contributors community by signing [JFrog's CLA](https://secure.echosign.com/public/hostedForm?formid=5IYKLZ2RXB543N).
- If the existing tests do not already cover your changes, please add tests.
- Pull requests should be created on the **dev** branch.
- Please use gofmt for formatting the code before submitting the pull request.

## General APIs

### Set logger

```go
var file *os.File
...
log.SetLogger(log.NewLogger(log.INFO, file))
```

### Setting the temp dir

The default temp dir used is 'os.TempDir()'. Use the following API to set a new temp dir:

```go
fileutils.SetTempDirBase(filepath.Join("my", "temp", "path"))
```

## Artifactory APIs

### Creating a Service Manager

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

#### Creating Service Config

```go
serviceConfig, err := config.NewConfigBuilder().
    SetServiceDetails(rtDetails).
    SetCertificatesPath(certPath).
    SetThreads(threads).
    SetDryRun(false).
    Build()
```

#### Creating New Service Manager

```go
rtManager, err := artifactory.New(&rtDetails, serviceConfig)
```

### Using Services

#### Uploading Files to Artifactory

```go
params := services.NewUploadParams()
params.Pattern = "repo/*/*.zip"
params.Target = "repo/path/"
// Attach properties to the uploaded files.
params.Props = "key1=val1;key2=val2"
params.Recursive = true
params.Regexp = false
params.IncludeDirs = false
params.Flat = true
params.Explode = false
params.Deb = ""
params.Symlink = false
params.AddVcsProps = false
// Retries default value: 3
params.Retries = 5
// MinChecksumDeploy default value: 10400
params.MinChecksumDeploy = 15360

rtManager.UploadFiles(params)
```

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
Use `reader.NextRecord()` and `FileInfo` type from `utils` package to iterate over the download results.
Use `reader.Close()` method (preferably using `defer`), to remove the results reader after it is used:
``` 
defer reader.Close()
var file utils.FileInfo
for e := resultReader.NextRecord(&file); e == nil; e = resultReader.NextRecord(&file) {
    fmt.Printf("Download source: %s\n", file.ArtifactoryPath)
    fmt.Printf("Download target: %s\n", file.LocalPath)
    fmt.Printf("SHA1: %s\n", file.Sha1)
    fmt.Printf("MD5: %s\n", file.Md5) 
}
```

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

pathsToDelete := rtManager.GetPathsToDelete(params)
rtManager.DeleteFiles(pathsToDelete)
```

#### Searching Files in Artifactory

```go
params := services.NewSearchParams()
params.Pattern = "repo/*/*.zip"
// Filter the files by properties. 
params.Props = "key1=val1;key2=val2"
params.Recursive = true

rtManager.SearchFiles(params)
```

#### Setting Properties on Files in Artifactory

```go
searchParams = services.NewSearchParams()
searchParams.Recursive = true
searchParams.IncludeDirs = false

resultItems = rtManager.SearchFiles(searchParams)

propsParams = services.NewPropsParams()
propsParams.Pattern = "repo/*/*.zip"
propsParams.Items = resultItems
// Filter the files by properties. 
propsParams.Props = "key=value"

rtManager.SetProps(propsParams)
```

#### Deleting Properties from Files in Artifactory

```go
searchParams = services.NewSearchParams()
searchParams.Recursive = true
searchParams.IncludeDirs = false

resultItems = rtManager.SearchFiles(searchParams)

propsParams = services.NewPropsParams()
propsParams.Pattern = "repo/*/*.zip"
propsParams.Items = resultItems
// Filter the files by properties. 
propsParams.Props = "key=value"

rtManager.DeleteProps(propsParams)
```

#### Publishing Build Info to Artifactory

```go
buildInfo := &buildinfo.BuildInfo{}
...
rtManager.PublishBuildInfo(buildInfo)
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

filesToDelete := rtManager.GetUnreferencedGitLfsFiles(params)
rtManager.DeleteFiles(filesToDelete)
```

#### Executing AQLs

```go
rtManager.Aql(aql string)
```

#### Reading Files in Artifactory

```go
rtManager.ReadRemoteFile(FilePath string)
```

#### Creating an access token

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

#### Fetching access tokens

```go
results, err := rtManager.GetTokens()
```

#### Refreshing an access token

```go
params := services.NewRefreshTokenParams()
params.AccessToken = "<access token>"
params.RefreshToken = "<refresh token>"
params.Token.Scope = "api:*"
params.Token.ExpiresIn = 3600
results, err := rtManager.RefreshToken(params)
```

#### Revoking an access token

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
pparams.Key = "generic-repo"
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

#### Creating and Updating Repository Replication

Example of creating repository replication:

```go
params := services.NewCreateReplicationParams()
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

Updating local repository replication:

```go
params := services.NewUpdateReplicationParams()
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

#### Fetch Artifactory's version

```go
version, err := servicesManager.GetVersion()
```

#### Fetch Artifactory's service id

```go
serviceId, err := servicesManager.GetServiceId()
```

## Distribution APIs

### Creating a Service Manager

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

#### Creating Service Config

```go
serviceConfig, err := config.NewConfigBuilder().
    SetServiceDetails(rtDetails).
    SetCertificatesPath(certPath).
    SetThreads(threads).
    SetDryRun(false).
    Build()
```

#### Creating New Service Manager

```go
distManager, err := distribution.New(&distDetails, serviceConfig)
```

### Using Services

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

#### Distributing a Release Bundle

```go
params := services.NewDistributeReleaseBundleParams("bundle-name", "1")
distributionRules := utils.DistributionCommonParams{SiteName: "Swamp-1", "CityName": "Tel-Aviv", "CountryCodes": []string{"123"}}}
params.DistributionRules = []*utils.DistributionCommonParams{distributionRules}

err := distManager.DistributeReleaseBundle(params)
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

### Creating a Service Manager

```go
serviceConfig := bintray.NewConfigBuilder().
    SetBintrayDetails(btDetails).
    SetDryRun(false).
    SetThreads(threads).
    Build()

btManager, err := bintray.New(serviceConfig)
```

### Using Services

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

#### Showing / Deleting a Bintray Package

```go
pkgPath, err := packages.CreatePath("subject/repo/pkg")

btManager.ShowPackage(pkgPath)
btManager.DeletePackage(pkgPath)
```

#### Creating / Updating a Bintray Package

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

#### Showing / Deleting a Bintray Version

```go
versionPath, err := versions.CreatePath("subject/repo/pkg/version")

btManager.ShowVersion(versionPath)
btManager.DeleteVersion(versionPath)
```

#### Creating / Updating a Bintray Version

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

#### Creating / Updating Entitlements

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

#### Showing / Deleting Entitlements

```go
versionPath, err := versions.CreatePath("subject/repo/pkg/version")

btManager.ShowAllEntitlements(versionPath)
btManager.ShowEntitlement("entitelmentID", versionPath)
btManager.DeleteEntitlement("entitelmentID", versionPath)
```

#### Creating / Updating Access Keys

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

#### Showing / Deleting Access Keys

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
| `-rt.apikey`        | [Optional] Artifactory API key.                                                                        |
| `-rt.sshKeyPath`    | [Optional] Ssh key file path. Should be used only if the Artifactory URL format is ssh://[domain]:port |
| `-rt.sshPassphrase` | [Optional] Ssh key passphrase.                                                                         |
| `-rt.accessToken`   | [Optional] Artifactory access token.                                                                   |
| `-log-level`        | [Default: INFO] Sets the log level.                                                                    |

- The tests create an Artifactory repository named _jfrog-client-tests-repo1_.<br/>
  Once the tests are completed, the content of this repository is deleted.
