package repositories

import (
	"os"
	"path/filepath"
)

func getRepoConfigPath(configPath string) string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, configPath)
}
