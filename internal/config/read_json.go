package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	var cfg Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(homeDirectory + "/" + configFileName)
	return path, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}

const configFileName = ".gatorconfig.json"
