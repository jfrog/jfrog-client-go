package tests

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func CreateFileWithContent(fileName, relativePath string) (string, string, error) {
	var err error
	tempDirPath, err := ioutil.TempDir("", "tests")
	if err != nil {
		return tempDirPath, "", err
	}

	fullPath := ""
	if relativePath != "" {
		fullPath = filepath.Join(tempDirPath, relativePath)
		err = os.MkdirAll(fullPath, 0777)
		if err != nil {
			return tempDirPath, "", err
		}
	}
	fullPath = filepath.Join(fullPath, fileName)
	file, err := os.Create(fullPath)
	if err != nil {
		return tempDirPath, "", err
	}
	defer file.Close()
	_, err = file.Write([]byte(strconv.FormatInt(int64(time.Now().Unix()), 10)))
	return tempDirPath, fullPath, err
}
