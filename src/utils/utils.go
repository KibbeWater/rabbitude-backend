package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"main/structures"
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

func GetCacheDir() (string, error) {
	exeDir, err := GetExecutableDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(exeDir, "cache")
	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return cacheDir, nil

}

func CreateHash(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// Gets a cached .wav tts audio file for any given prompt
// Also a pointer is probably stupid, but it works and probably won't change performance in the long run
func GetCachedPromptTTS(prefix string, prompt string) (*structures.ProviderAudioResponse, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}
	hash := CreateHash([]byte(prompt))

	if len(hash) < 7 {
		return nil, fmt.Errorf("bad hash")
	}

	// The file name is the prefix + the first 7 characters of the hash dot wav
	fileName := filepath.Join(cacheDir, prefix+hash[:7]+".wav")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, nil
	}

	// Metadata filena
	metadataFileName := filepath.Join(cacheDir, prefix+hash[:7]+".txt")
	if _, err := os.Stat(metadataFileName); os.IsNotExist(err) {
		return nil, nil
	}

	// Read the file
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	metadata, err := os.ReadFile(metadataFileName)
	if err != nil {
		return nil, err
	}

	return &structures.ProviderAudioResponse{
		Audio:        data,
		TextMetadata: string(metadata),
	}, nil
}

// Caches a .wav tts audio file for any given prompt
// TODO: Put data into a gob and gzip to save on storage space
func CachePromptTTS(prefix string, prompt string, data []byte, textMetadata string) error {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}
	hash := CreateHash([]byte(prompt))

	if len(hash) < 7 {
		return fmt.Errorf("bad hash")
	}

	// The file name is the prefix + the first 7 characters of the hash dot wav
	fileName := filepath.Join(cacheDir, prefix+hash[:7]+".wav")
	err = os.WriteFile(fileName, data, os.ModePerm)
	if err != nil {
		return err
	}

	metadataFileName := filepath.Join(cacheDir, prefix+hash[:7]+".txt")
	err = os.WriteFile(metadataFileName, []byte(textMetadata), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
