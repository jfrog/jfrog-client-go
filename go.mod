module github.com/jfrog/jfrog-client-go

go 1.17

require (
	github.com/buger/jsonparser v1.1.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/gookit/color v1.4.2
	github.com/jfrog/build-info-go v0.1.0
	github.com/jfrog/gofrog v1.1.0
	github.com/mholt/archiver/v3 v3.5.1-0.20210618180617-81fac4ba96e4
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/xanzy/ssh-agent v0.3.0
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
)

require (
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/andybalholm/brotli v1.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/golang/snappy v0.0.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/klauspost/compress v1.11.4 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/ulikunitz/xz v0.5.9 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/net v0.0.0-20210326060303-6b1517762897 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

exclude (
	golang.org/x/text v0.3.3
	golang.org/x/text v0.3.4
)

//replace github.com/jfrog/build-info-go => github.com/jfrog/build-info-go v0.0.0-20211108124451-f5ad09ddfe92

//replace github.com/jfrog/gofrog => github.com/jfrog/gofrog v1.0.7-0.20211109140605-15e312b86c9f
