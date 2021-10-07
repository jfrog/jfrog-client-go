module github.com/jfrog/jfrog-client-go

require (
	github.com/buger/jsonparser v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-git/go-git/v5 v5.4.2
	github.com/gookit/color v1.4.2
	github.com/jfrog/gofrog v1.0.6
	github.com/mholt/archiver/v3 v3.5.1-0.20210618180617-81fac4ba96e4
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/xanzy/ssh-agent v0.3.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

exclude (
	golang.org/x/text v0.3.3
	golang.org/x/text v0.3.4
)

go 1.13
