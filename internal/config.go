package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers Providers `json:"providers" yaml:"providers"`
	App       AppConfig `json:"app"       yaml:"app"`
}

type AppConfig struct {
	Name   string   `json:"name"   yaml:"name"`
	Groups []string `json:"groups" yaml:"groups"`
}

type Providers struct {
	File   map[string]FileProviderConfig   `json:"file"   yaml:"file"`
	Docker map[string]DockerProviderConfig `json:"docker" yaml:"docker"`
}

func DefaultConfig() Config {
	return Config{
		Providers: Providers{
			File:   map[string]FileProviderConfig{},
			Docker: map[string]DockerProviderConfig{},
		},
		App: AppConfig{
			Name:   "simplydash",
			Groups: []string{},
		},
	}
}

func GetConfig(args Args) (Config, error) {
	config := DefaultConfig()

	configFile, err := os.ReadFile(args.ConfigFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
