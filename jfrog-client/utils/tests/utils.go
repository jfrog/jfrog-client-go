package tests

import (
	"bufio"
	"github.com/jfrog/jfrog-cli-go/jfrog-cli/utils/tests"
	"github.com/jfrog/jfrog-cli-go/jfrog-client/utils/log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type HttpServerHandlers map[string]func(w http.ResponseWriter, r *http.Request)

func StartHttpServer(handlers HttpServerHandlers) (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	go func() {
		httpMux := http.NewServeMux()
		for k, v := range handlers {
			httpMux.HandleFunc(k, v)
		}
		err = http.Serve(listener, httpMux)
		if err != nil {
			panic(err)
		}
	}()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func GetTestPackages(searchPattern string) []string {
	// Get all packages with test files.
	rootDir := FindRoot()
	cmd := exec.Command("go", "list", "-f", "{{.Dir}} {{.TestGoFiles}}", searchPattern)
	cmd.Dir = rootDir
	packages, _ := cmd.Output()

	scanner := bufio.NewScanner(strings.NewReader(string(packages)))
	var unitTests []string
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		// Skip if package does not contain test files.
		if len(fields) > 1 && len(fields[1]) > 2 {
			unitTests = append(unitTests, "." + strings.TrimPrefix(fields[0], rootDir))
		}
	}
	return unitTests
}

func FindRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Error("Cannot retrieve current dir path")
		exitOnErr(err)
	}
	origDir := dir
	for len(dir) > 2 {
		if _, err := os.Stat(dir + "/go.mod"); err == nil {
			log.Info("Found root Dir at:", dir)
			return dir
		}
		dir = filepath.Dir(dir)
	}
	log.Error("Did not find root dir with go.mod file under", origDir)
	os.Exit(1)
	return origDir
}

func ExcludeTestsPackage(packages []string, packageToExclude string) []string {
	var res []string
	for _, packageName := range packages {
		if packageName != packageToExclude {
			res = append(res, packageName)
		}
	}
	log.Info("Executing unit tests in packages:", res)
	return res
}

func RunTests(testsPackages []string) error {
	if len(testsPackages) == 0 {
		return nil
	}
	testsPackages = append([]string{"test", "-v"}, testsPackages...)
	cmd := exec.Command("vgo", testsPackages...)
	cmd.Dir = FindRoot()

	tempDirPath, err := tests.GetTestsLogsDir()
	exitOnErr(err)

	f, err := os.Create(filepath.Join(tempDirPath, "unit_tests.log"))
	exitOnErr(err)

	cmd.Stdout, cmd.Stderr = f, f
	if err := cmd.Run(); err != nil {
		log.Error("Unit tests failed, full report available at the following path:", f.Name())
		exitOnErr(err)
	}
	log.Info("Full unit testing report available at the following path:", f.Name())
	return nil
}

func exitOnErr(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
