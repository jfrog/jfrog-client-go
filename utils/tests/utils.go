package tests

import (
	"bufio"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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
	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}} {{.TestGoFiles}}", searchPattern)
	packages, _ := cmd.Output()

	scanner := bufio.NewScanner(strings.NewReader(string(packages)))
	var unitTests []string
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		// Skip if package does not contain test files.
		if len(fields) > 1 && len(fields[1]) > 2 {
			unitTests = append(unitTests, fields[0])
		}
	}
	return unitTests
}

func ExcludeTestsPackage(packages []string, packageToExclude string) []string {
	var res []string
	for _, packageName := range packages {
		if packageName != packageToExclude {
			res = append(res, packageName)
		}
	}
	return res
}

func RunTests(testsPackages []string, hideUnitTestsLog bool) error {
	if len(testsPackages) == 0 {
		return nil
	}
	testsPackages = append([]string{"test", "-v"}, testsPackages...)
	cmd := exec.Command("go", testsPackages...)

	if hideUnitTestsLog {
		tempDirPath, err := getTestsLogsDir()
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

	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error("Unit tests failed")
		exitOnErr(err)
	}

	return nil
}

func getTestsLogsDir() (string, error) {
	tempDirPath := filepath.Join(os.TempDir(), "jfrog_tests_logs")
	return tempDirPath, fileutils.CreateDirIfNotExist(tempDirPath)
}

func exitOnErr(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func InitVcsSubmoduleTestDir(t *testing.T, srcPath string) (submodulePath, tmpDir string) {
	var err error
	tmpDir, err = fileutils.CreateTempDir()
	assert.NoError(t, err)

	err = fileutils.CopyDir(srcPath, tmpDir, true, nil)
	assert.NoError(t, err)
	if found, err := fileutils.IsDirExists(filepath.Join(tmpDir, "gitdata"), false); found {
		assert.NoError(t, err)
		err := fileutils.RenamePath(filepath.Join(tmpDir, "gitdata"), filepath.Join(tmpDir, ".git"))
		assert.NoError(t, err)
	}
	submoduleDst := filepath.Join(tmpDir, "subdir", "submodule")
	err = fileutils.CopyFile(submoduleDst, filepath.Join(tmpDir, "gitSubmoduleData"))
	assert.NoError(t, err)
	err = fileutils.MoveFile(filepath.Join(submoduleDst, "gitSubmoduleData"), filepath.Join(submoduleDst, ".git"))
	assert.NoError(t, err)
	submodulePath, err = filepath.Abs(submoduleDst)
	assert.NoError(t, err)
	return submodulePath, tmpDir
}
