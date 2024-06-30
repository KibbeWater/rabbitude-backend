package config

import (
	"encoding/json"
	"main/structures"
	"os"
	"path/filepath"
)

type Config = structures.Config

var ConfigData *Config

func GetConfig() *Config {
	// Load the config if it hasn't been loaded yet
	if ConfigData == nil {
		loadConfig()
	}

	if ConfigData == nil {
		ConfigData = defaultConfig()
		SaveConfig()
	}

	return ConfigData
}

func SaveConfig() {
	configData := GetConfig()

	// Marshal the config data into JSON
	data, err := json.MarshalIndent(configData, "", "    ")
	if err != nil {
		return
	}

	exeDir, err := getExecutableDir()
	if err != nil {
		return
	}

	configPath := exeDir + "/config.json"

	// Write the JSON data to the config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		println("Unable to save config file")
		return
	}
}

func loadConfig() {
	exeDir, err := getExecutableDir()
	if err != nil {
		return
	}

	configPath := exeDir + "/config.json"

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return
	}

	// Unmarshal the JSON data into the Config struct
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return
	}

	// Set the global config variable
	ConfigData = &cfg
}

func GetProviderConfig(provider string) *structures.ProviderConfig {
	config := GetConfig()

	for _, providerConfig := range config.Providers {
		if providerConfig.ProviderName == provider {
			return &providerConfig
		}
	}

	return nil
}

func SaveProviderConfig(providerConfig *structures.ProviderConfig) {
	config := GetConfig()

	for i, pc := range config.Providers {
		if pc.ProviderName == providerConfig.ProviderName {
			config.Providers[i] = *providerConfig
			SaveConfig()
			return
		}
	}

	config.Providers = append(config.Providers, *providerConfig)
	SaveConfig()
}

func defaultConfig() *Config {
	return &Config{
		Version: 1,
		General: structures.GeneralConfig{
			Port:           8080,
			BaseProvider:   "",
			LLMProvider:    "",
			SpeechProvider: "",
			TTSProvider:    "",
			SearchProvider: "",
		},
	}
}

func getExecutableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}
