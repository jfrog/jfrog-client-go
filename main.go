package main // import "github.com/jfrog/jfrog-cli-go"

import (
	"github.com/jfrog/jfrog-cli-go/jfrog-cli/jfrog"
	"os"
)

func main() {
	args := os.Args
	jfrog.Run(args)
}
