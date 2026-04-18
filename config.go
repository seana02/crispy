package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DBPath string
	// CurrencyAPI  string
	// AutoConfirm  bool
}

func LoadConfig() (*Config, error) {
	dbPath, err := getConfigFilePath("crispy", "data.db")
	if err != nil {
		return nil, fmt.Errorf("Could not get database path: %s", err)
	}

	cfg := &Config{
		DBPath: dbPath,
	}

	return cfg, nil
}

func getConfigFilePath(appName, dbName string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, appName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, dbName), nil
}
