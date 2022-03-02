module github.com/jfrog/jfrog-client-go

require (
	github.com/buger/jsonparser v1.1.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/gookit/color v1.4.2
	github.com/jfrog/gofrog v1.0.6
	github.com/mholt/archiver/v3 v3.5.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/xanzy/ssh-agent v0.3.0
	golang.org/x/crypto v0.0.0-20211202192323-5770296d904e
	golang.org/x/text v0.3.7 // indirect
)

exclude (
	golang.org/x/text v0.3.3
	golang.org/x/text v0.3.4
	golang.org/x/text v0.3.6
)

go 1.13
