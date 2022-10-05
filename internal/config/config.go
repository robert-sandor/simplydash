package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Settings      Settings             `json:"settings" yaml:"settings"`
	FileProviders []FileProviderConfig `json:"files" yaml:"files"`
}

type FileProviderConfig struct {
	Path  string `json:"path" yaml:"path"`
	Watch bool   `json:"watch" yaml:"watch"`
}

func DefaultConfig() *Config {
	return &Config{
		Settings:      DefaultSettings(),
		FileProviders: []FileProviderConfig{},
	}
}

func LoadConfig(path string, reader func(string) ([]byte, error)) (*Config, error) {
	bytes, err := reader(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Config file not found at %s - creating with defaults", path)
			return createConfigFileWithDefaults(path)
		}
		return nil, err
	}

	config := DefaultConfig()
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	err = validateConfig(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func createConfigFileWithDefaults(path string) (*Config, error) {
	config := DefaultConfig()
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return nil, err
	}

	if _, err = file.Write(bytes); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *Config) error {
	return nil
}
