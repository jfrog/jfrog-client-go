module github.com/jfrog/jfrog-client-go

require (
	github.com/asafgabai/build-info-go v0.0.0-20210930074455-ab8203bac58d
	github.com/buger/jsonparser v1.1.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/gookit/color v1.4.2
	github.com/jfrog/gofrog v1.0.7
	github.com/mholt/archiver/v3 v3.5.1-0.20210618180617-81fac4ba96e4
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/xanzy/ssh-agent v0.3.0
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
)

exclude (
	golang.org/x/text v0.3.3
	golang.org/x/text v0.3.4
)

//replace github.com/asafgabai/build-info-go => ../build-info-go

go 1.13
