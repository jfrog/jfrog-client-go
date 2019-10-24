# jfrog-client-go

|Branch|Status|
|:---:|:---:|
|master|[![Build status](https://ci.appveyor.com/api/projects/status/2wkemson2sj4skyh/branch/master?svg=true)](https://ci.appveyor.com/project/jfrog-ecosystem/jfrog-client-go/branch/master)
|dev|[![Build status](https://ci.appveyor.com/api/projects/status/2wkemson2sj4skyh/branch/dev?svg=true)](https://ci.appveyor.com/project/jfrog-ecosystem/jfrog-client-go/branch/dev)

## General

*jfrog-client-go* is a library which provides Go APIs to performs actions on JFrog Artifactory or Bintray from your Go application.
The project is still relatively new, and its APIs may therefore change frequently between releases.
The library can be used as a go-module, which should be added to your project's go.mod file. As a reference you may look at [JFrog CLI](https://github.com/jfrog/jfrog-cli-go)'s [go.mod](https://github.com/jfrog/jfrog-cli-go/blob/master/go.mod) file, which uses this library as a dependency.

## Pull Requests
We welcome pull requests from the community.

### Guidelines
* Before creating your first pull request, please join our contributors community by signing [JFrog's CLA](https://secure.echosign.com/public/hostedForm?formid=5IYKLZ2RXB543N).
* If the existing tests do not already cover your changes, please add tests.
* Pull requests should be created on the **dev** branch.
* Please use gofmt for formatting the code before submitting the pull request.

## General APIs
### Set logger
```
var file *os.File
...
log.SetLogger(log.NewLogger(log.INFO, file))
```

### Setting the temp dir
The default temp dir used is  'os.TempDir()'. Use the following API to set a new temp dir:
```
    fileutils.SetTempDirBase(filepath.Join("my", "temp", "path"))
```

## Artifactory APIs
### Creating a Service Manager
#### Creating Artifactory Details
```
    rtDetails := auth.NewArtifactoryDetails()
    rtDetails.SetUrl("http://localhost:8081/artifactory")
    rtDetails.SetSshKeysPath("path/to/.ssh/")
    rtDetails.SetApiKey("apikey")
    rtDetails.SetUser("user")
    rtDetails.SetPassword("password")
    rtDetails.SetAccessToken("accesstoken")
```
#### Creating Service Config
```
    serviceConfig, err := artifactory.NewConfigBuilder().
        SetArtDetails(rtDetails).
        SetCertificatesPath(certPath).
        SetThreads(threads).
        SetDryRun(false).
        Build()
```
#### Creating New Service Manager
```
    rtManager, err := artifactory.New(&rtDetails, serviceConfig)
```

### Using Services
#### Uploading Files to Artifactory
```
    params := services.NewUploadParams()
    params.Pattern = "repo/*/*.zip"
    params.Target = "repo/path/"
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

    rtManager.UploadFiles(params)
```

#### Downloading Files from Artifactory
```
    params := services.NewDownloadParams()
    params.Pattern = "repo/*/*.zip"
    params.Target = "target/path/"
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

    rtManager.DownloadFiles(params)
```

#### Copying Files in Artifactory
```
    params := services.NewMoveCopyParams()
    params.Pattern = "repo/*/*.zip"
    params.Target = "target/path/"
    params.Recursive = true
    params.Flat = false

    rtManager.Copy(params)
```

#### Moving Files in Artifactory
```
    params := services.NewMoveCopyParams()
    params.Pattern = "repo/*/*.zip"
    params.Target = "target/path/"
    params.Recursive = true
    params.Flat = false

    rtManager.Move(params)
```

#### Deleting Files from Artifactory
```
    params := services.NewDeleteParams()
    params.Pattern = "repo/*/*.zip"
    params.Recursive = true

    pathsToDelete := rtManager.GetPathsToDelete(params)
    rtManager.DeleteFiles(pathsToDelete)
```

#### Searching Files in Artifactory
```
    params := services.NewSearchParams()
    params.Pattern = "repo/*/*.zip"
    params.Recursive = true

    rtManager.SearchFiles(params)
```

#### Setting Properties on Files in Artifactory
```
    searchParams = services.NewSearchParams()
    searchParams.Recursive = true
    searchParams.IncludeDirs = false

    resultItems = rtManager.SearchFiles(searchParams)

    propsParams = services.NewPropsParams()
    propsParams.Pattern = "repo/*/*.zip"
    propsParams.Items = resultItems
    propsParams.Props = "key=value"

    rtManager.SetProps(propsParams)
```

#### Deleting Properties from Files in Artifactory
```
    searchParams = services.NewSearchParams()
    searchParams.Recursive = true
    searchParams.IncludeDirs = false

    resultItems = rtManager.SearchFiles(searchParams)

    propsParams = services.NewPropsParams()
    propsParams.Pattern = "repo/*/*.zip"
    propsParams.Items = resultItems
    propsParams.Props = "key=value"

    rtManager.DeleteProps(propsParams)
```

#### Publishing Build Info to Artifactory
```
    buildInfo := &buildinfo.BuildInfo{}
    ...
    rtManager.PublishBuildInfo(buildInfo)
```

#### Fetching Build Info from Artifactory
```
    buildInfoParams := services.NewBuildInfoParams{}
    buildInfoParams.BuildName = "buildName"
    buildInfoParams.BuildNumber = "LATEST"

    rtManager.GetBuildInfo(buildInfoParams)
```

#### Promoting Published Builds in Artifactory
```
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
```
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
```
    params := services.NewXrayScanParams()
    params.BuildName = buildName
    params.BuildNumber = buildNumber

    rtManager.XrayScanBuild(params)
```

#### Discarding Old Builds
```
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
```
    params := services.NewGitLfsCleanParams()
    params.Refs = "refs/remotes/*"
    params.Repo = "my-project-lfs"
    params.GitPath = "path/to/git"

    filesToDelete := rtManager.GetUnreferencedGitLfsFiles(params)
    rtManager.DeleteFiles(filesToDelete)
```

#### Executing AQLs
```
    rtManager.Aql(aql string)
```

#### Reading Files in Artifactory
```
    rtManager.ReadRemoteFile(FilePath string)
```

#### Creating an access token
```
    params := services.NewCreateTokenParams()
    params.Scope = "api:* member-of-groups:readers"
    params.Username = "user"
    params.ExpiresIn = 3600
    params.GrantType = "client_credentials"
    params.Refreshable = true
    params.Audience = "jfrt@<serviceID1> jfrt@<serviceID2>"
    results, err := rtManager.CreateToken(params)
```

#### Fetching access tokens
```
    results, err := rtManager.GetTokens()
```

#### Refreshing an access token
```
    params := services.NewRefreshTokenParams()
    params.AccessToken = "<access token>"
    params.RefreshToken = "<refresh token>"
    params.Token.Scope = "api:*"
    params.Token.ExpiresIn = 3600
    results, err := rtManager.RefreshToken(params)
```

#### Revoking an access token
```
    params := services.NewRevokeTokenParams()

    // Provide either TokenId or Token
    params.TokenId = "<token id>"
    // params.Token = "access token"

    err := rtManager.RevokeToken(params)
```

## Bintray APIs
### Creating Bintray Details
 ```
    btDetails := auth.NewBintrayDetails()
    btDetails.SetUser("user")
    btDetails.SetKey("key")
    btDetails.SetDefPackageLicense("Apache 2.0")
 ```

### Creating a Service Manager
```
    serviceConfig := bintray.NewConfigBuilder().
        SetBintrayDetails(btDetails).
        SetDryRun(false).
        SetThreads(threads).
        Build()

    btManager, err := bintray.New(serviceConfig)
```
### Using Services
#### Uploading a Single File to Bintray
```
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
```
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
```
    params := services.NewDownloadVersionParams()
    params.Path, err = versions.CreatePath("subject/repo/pkg/version")

    params.IncludeUnpublished = false
    params.TargetPath = "target/path/"

    btManager.DownloadVersion(params)
```

#### Showing / Deleting a Bintray Package
```
    pkgPath, err := packages.CreatePath("subject/repo/pkg")

    btManager.ShowPackage(pkgPath)
    btManager.DeletePackage(pkgPath)
```

#### Creating / Updating a Bintray Package
```
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
```
    versionPath, err := versions.CreatePath("subject/repo/pkg/version")

    btManager.ShowVersion(versionPath)
    btManager.DeleteVersion(versionPath)
```

#### Creating / Updating a Bintray Version
```
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
```
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
```
    versionPath, err := versions.CreatePath("subject/repo/pkg/version")

    btManager.ShowAllEntitlements(versionPath)
    btManager.ShowEntitlement("entitelmentID", versionPath)
    btManager.DeleteEntitlement("entitelmentID", versionPath)
```

#### Creating / Updating Access Keys
```
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
```
    btManager.ShowAllAccessKeys("org")
    btManager.ShowAccessKey("org", "KeyID")
    btManager.DeleteAccessKey("org", "KeyID")
```

#### Signing a URL
```
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
```
    path, err := utils.CreatePathDetails("subject/repository/file-path")

    btManager.GpgSignFile(path, "passphrase") 
```
	
#### GPG Signing Version Files
```
    path, err := versions.CreatePath("subject/repo/pkg/version")

    btManager.GpgSignVersion(path, "passphrase")
```

#### Listing Logs
```
    path, err := versions.CreatePath("subject/repo/pkg/version")

    btManager.LogsList(versionPath)
```

#### Downloading Logs
```
    path, err := versions.CreatePath("subject/repo/pkg/version")

    btManager.DownloadLog(path, "logName")
```

## Tests
To run tests on the source code, you'll need a running JFrog Artifactory Pro instance.
Use the following command with the below options to run the tests.
````
go test -v github.com/jfrog/jfrog-client-go/tests
````
Optional flags:

| Flag | Description |
| --- | --- |
| `-rt.url` | [Default: http://localhost:8081/artifactory] Artifactory URL. |
| `-rt.user` | [Default: admin] Artifactory username. |
| `-rt.password` | [Default: password] Artifactory password. |
| `-rt.apikey` | [Optional] Artifactory API key. |
| `-rt.sshKeyPath` | [Optional] Ssh key file path. Should be used only if the Artifactory URL format is ssh://[domain]:port |
| `-rt.sshPassphrase` | [Optional] Ssh key passphrase. |
| `-rt.accessToken` | [Optional] Artifactory access token. |
| `-log-level` | [Default: INFO] Sets the log level. |


* The tests create an Artifactory repository named *jfrog-client-tests-repo1*.<br/>
  Once the tests are completed, the content of this repository is deleted.
