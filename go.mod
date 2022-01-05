module github.com/jfrog/jfrog-client-go

go 1.17

require (
	github.com/buger/jsonparser v1.1.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/gookit/color v1.4.2
	github.com/jfrog/build-info-go v0.1.6
	github.com/jfrog/gofrog v1.1.1
	github.com/mholt/archiver/v3 v3.5.1-0.20210618180617-81fac4ba96e4
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/xanzy/ssh-agent v0.3.0
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
)

require golang.org/x/text v0.3.7 // indirect

exclude (
	golang.org/x/text v0.3.3
	golang.org/x/text v0.3.4
)

//replace github.com/jfrog/build-info-go => github.com/jfrog/build-info-go v0.1.5-0.20211209071650-c5f4d2e581c3

//replace github.com/jfrog/gofrog => github.com/jfrog/gofrog v1.0.7-0.20211128152632-e218c460d703
