package util

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	gap "github.com/muesli/go-app-paths"
)

func initDirectory(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}
		return err
	}
	return nil
}

var ErrDataDirectory = errors.New("couldnt create data directory for ctrl database")

func DataPath() (string, error) {
	scope := gap.NewScope(gap.User, "ctrl")

	dirs, err := scope.DataDirs()
	if err != nil {
		return "", fmt.Errorf("%w :: %s", ErrDataDirectory, err.Error())
	}

	var dataDir string
	if len(dirs) > 0 {
		dataDir = dirs[0]
	} else {
		dataDir, err = os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("%w :: %s", ErrDataDirectory, err.Error())
		}
	}

	if err := initDirectory(dataDir); err != nil {
		return "", fmt.Errorf("%w :: %s", ErrDataDirectory, err.Error())
	}

	return dataDir, nil
}

var ErrConfigDirectory = errors.New("couldnt create config directory for ctrl database")

func ConfigPath() (string, error) {
	scope := gap.NewScope(gap.User, "ctrl")

	dirs, err := scope.ConfigDirs()
	if err != nil {
		return "", fmt.Errorf("%w :: %s", ErrConfigDirectory, err.Error())
	}

	var configDir string
	if len(dirs) > 0 {
		configDir = dirs[0]
	} else {
		configDir, err = os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("%w :: %s", ErrConfigDirectory, err.Error())
		}
	}

	if err := initDirectory(configDir); err != nil {
		return "", fmt.Errorf("%w :: %s", ErrConfigDirectory, err.Error())
	}

	return configDir, nil
}

func GetOS() string {
	if runtime.GOOS == "windows" {
		return "windows"
	}
	return "linux"
}
