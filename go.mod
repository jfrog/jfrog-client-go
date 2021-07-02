module github.com/jfrog/jfrog-client-go

require (
	github.com/buger/jsonparser v0.0.0-20180910192245-6acdf747ae99
	github.com/go-git/go-git/v5 v5.4.2
	github.com/gookit/color v1.4.2
	github.com/jfrog/gofrog v1.0.6
	github.com/mholt/archiver/v3 v3.5.1-0.20210618180617-81fac4ba96e4
	github.com/stretchr/testify v1.7.0
	github.com/xanzy/ssh-agent v0.3.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	golang.org/x/text v0.3.6 // indirect
)

exclude (
	golang.org/x/text v0.3.3
	golang.org/x/text v0.3.4
)

go 1.13
