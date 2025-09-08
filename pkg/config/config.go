package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config defines the structure of the YAML file
type Config struct {
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	Database struct {
		Source        string `yaml:"source"`
		MigrationPath string `yaml:"migration_path"`
	} `yaml:"database"`
}

// LoadConfig reads and parses the YAML config file
func LoadConfig(filePath string) (*Config, error) {
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	// Parse YAML into the Config struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}

	return &config, nil
}
