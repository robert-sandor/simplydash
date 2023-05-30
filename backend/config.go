package main

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Docker map[string]DockerConfig `json:"docker" yaml:"docker"`
	Web    WebConfig               `json:"web"    yaml:"web"`
}

type WebConfig struct {
	Name            string `json:"name"             yaml:"name"`
	Theme           string `json:"theme"            yaml:"theme"`
	Background      string `json:"background"       yaml:"background"`
	BackgroundImage string `json:"background_image" yaml:"background_image"`
}

func DefaultConfig() Config {
	return Config{
		Docker: make(map[string]DockerConfig),
		Web: WebConfig{
			Name:            "simplydash",
			Theme:           "dark",
			Background:      "#0a0a0a",
			BackgroundImage: "",
		},
	}
}

type DockerConfig struct {
	Socket string `json:"socket" yaml:"socket"`
	Host   string `json:"host"   yaml:"host"`
	Port   int    `json:"port"   yaml:"port"`
}

const (
	configFileName string = "config.yml"
)

type ConfigService interface {
	Init() error
}

type ConfigServiceImpl struct {
	filePath string
}

func NewConfigService(configPath string) ConfigService {
	return &ConfigServiceImpl{
		filePath: path.Join(configPath, configFileName),
	}
}

func (configService *ConfigServiceImpl) Init() error {
	logrus.WithField("configFile", configService.filePath).Debug("checking if config file exists")
	if _, err := os.Stat(configService.filePath); err == nil {
		return configService.loadConfigFile()
	}

	return configService.createConfigFile()
}

func (configService *ConfigServiceImpl) createConfigFile() error {
	config := DefaultConfig()

	logrus.WithField("configFile", configService.filePath).Debug("creating new config file")
	file, err := os.Create(configService.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	logrus.WithField("config", config).Debug("marshalling config")
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	logrus.WithField("yaml", string(bytes)).Debug("writing default config")
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (configService *ConfigServiceImpl) loadConfigFile() error {
	return nil
}
