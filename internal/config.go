package internal

import (
	"log/slog"
	"os"
	"path"

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

	if _, err := os.Stat(args.ConfigFile); err != nil {
		slog.Warn("no config file found", "configPath", args.ConfigFile)
		createConfigFile(args.ConfigFile, config)
		return config, nil
	}

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

func createConfigFile(configPath string, config Config) {
	err := os.MkdirAll(path.Dir(configPath), 0755)
	if err != nil {
		slog.Error("failed to create config dir", "configPath", configPath, "err", err)
		return
	}

	file, err := os.Create(configPath)
	if err != nil {
		slog.Error("failed to create config file", "configPath", configPath, "err", err)
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	yamlContent, err := yaml.Marshal(config)
	if err != nil {
		slog.Error("failed to serialize default config", "configPath", configPath, "err", err)
		return
	}

	_, err = file.WriteString(string(yamlContent))
	if err != nil {
		slog.Error("failed to write config file", "configPath", configPath, "err", err)
		return
	}
}
