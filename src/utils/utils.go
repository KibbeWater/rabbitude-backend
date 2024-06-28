package utils

import (
	"os"
	"path/filepath"
)

func GetExecutableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}
