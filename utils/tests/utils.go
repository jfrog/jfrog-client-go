package tests

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	biutils "github.com/jfrog/build-info-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
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
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, errors.New("couldn't assert listener address to tcpAddr")
	}
	return tcpAddr.Port, nil
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
	testsPackages = append([]string{"test", "-v", "-p", "1"}, testsPackages...)
	cmd := exec.Command("go", testsPackages...)

	if hideUnitTestsLog {
		tempDirPath := filepath.Join(os.TempDir(), "jfrog_tests_logs")
		exitOnErr(fileutils.CreateDirIfNotExist(tempDirPath))

		f, err := os.Create(filepath.Join(tempDirPath, "unit_tests.log"))
		exitOnErr(err)

		cmd.Stdout, cmd.Stderr = f, f
		if err = cmd.Run(); err != nil {
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

func exitOnErr(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func InitVcsSubmoduleTestDir(t *testing.T, srcPath, tmpDir string) (submodulePath string) {
	var err error
	assert.NoError(t, biutils.CopyDir(srcPath, tmpDir, true, nil))
	if found, err := fileutils.IsDirExists(filepath.Join(tmpDir, "gitdata"), false); found {
		assert.NoError(t, err)
		assert.NoError(t, fileutils.RenamePath(filepath.Join(tmpDir, "gitdata"), filepath.Join(tmpDir, ".git")))
	}
	submoduleDst := filepath.Join(tmpDir, "subdir", "submodule")
	assert.NoError(t, biutils.CopyFile(submoduleDst, filepath.Join(tmpDir, "gitSubmoduleData")))
	assert.NoError(t, fileutils.MoveFile(filepath.Join(submoduleDst, "gitSubmoduleData"), filepath.Join(submoduleDst, ".git")))
	submodulePath, err = filepath.Abs(submoduleDst)
	assert.NoError(t, err)
	return submodulePath
}

func InitVcsWorktreeTestDir(t *testing.T, srcPath, tmpDir string) (worktreePath string) {
	var err error
	assert.NoError(t, biutils.CopyDir(srcPath, tmpDir, true, nil))
	if found, err := fileutils.IsDirExists(filepath.Join(tmpDir, "gitdata"), false); found {
		assert.NoError(t, err)
		assert.NoError(t, fileutils.RenamePath(filepath.Join(tmpDir, "gitdata"), filepath.Join(tmpDir, "bare.git")))
	}
	worktreeDst := filepath.Join(tmpDir, "worktree_repo")
	worktreePath, err = filepath.Abs(worktreeDst)
	assert.NoError(t, fileutils.MoveFile(filepath.Join(worktreeDst, "gitWorktreeData"), filepath.Join(worktreeDst, ".git")))
	assert.NoError(t, err)
	return worktreePath
}

func ChangeDirAndAssert(t *testing.T, dirPath string) {
	assert.NoError(t, os.Chdir(dirPath), "Couldn't change dir to "+dirPath)
}

// ChangeDirWithCallback changes working directory to the given path and return function that change working directory back to the original path.
func ChangeDirWithCallback(t *testing.T, originWd, destinationWd string) func() {
	ChangeDirAndAssert(t, destinationWd)
	return func() {
		ChangeDirAndAssert(t, originWd)
	}
}

func RemoveAndAssert(t *testing.T, path string) {
	assert.NoError(t, os.Remove(path), "Couldn't remove: "+path)
}

func RemoveAllAndAssert(t *testing.T, path string) {
	assert.NoError(t, fileutils.RemoveTempDir(path), "Couldn't removeAll: "+path)
}

func RemoveAllQuietly(t *testing.T, path string) {
	err := fileutils.RemoveTempDir(path)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to remove %s: %+v", path, err))
	}
}

func CloseQuietly(t *testing.T, closer io.Closer) {
	if closer != nil {
		err := closer.Close()
		if err != nil {
			log.Warn(fmt.Sprintf("Failed to close %+v", err))
		}
	}
}

func SetEnvAndAssert(t *testing.T, key, value string) {
	assert.NoError(t, os.Setenv(key, value), "Failed to set env: "+key)
}

func SetEnvWithCallbackAndAssert(t *testing.T, key, value string) func() {
	oldValue, exist := os.LookupEnv(key)
	SetEnvAndAssert(t, key, value)

	if exist {
		return func() {
			SetEnvAndAssert(t, key, oldValue)
		}
	}

	return func() {
		UnSetEnvAndAssert(t, key)
	}
}

func UnSetEnvAndAssert(t *testing.T, key string) {
	assert.NoError(t, os.Unsetenv(key), "Failed to unset env: "+key)
}

func GetLocalArtifactoryTokenIfNeeded(url string) (adminToken string) {
	if strings.Contains(url, "localhost:808") {
		adminToken = os.Getenv("JFROG_TESTS_LOCAL_ACCESS_TOKEN")
	}
	return
}

// Set new logger with output redirection to a null logger. This is useful for negative tests.
// Caller is responsible to set the old log back.
func RedirectLogOutputToNil() (previousLog log.Log) {
	previousLog = log.Logger
	newLog := log.NewLogger(log.INFO, nil)
	newLog.SetOutputWriter(io.Discard)
	newLog.SetLogsWriter(io.Discard, 0)
	log.SetLogger(newLog)
	return previousLog
}
