package utils

import (
	"crypto/rand"
	"encoding/base64"
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

func GenerateUniqueID() string {
	b := make([]byte, 9) // Generates a 12-character string when base64 encoded
	_, err := rand.Read(b)
	if err != nil {
		return "" // Handle error or generate a fallback ID
	}
	return base64.URLEncoding.EncodeToString(b)
}
