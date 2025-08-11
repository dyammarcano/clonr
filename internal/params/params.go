package params

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	once       sync.Once
	AppdataDir string
)

func init() {
	once.Do(func() {
		var err error
		AppdataDir, err = getAppDataDir("clonr")
		if err != nil {
			panic(err)
		}

		if err := os.MkdirAll(AppdataDir, os.ModePerm); err != nil {
			panic(err)
		}
	})
}

func getAppDataDir(appName string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, appName), nil
}
