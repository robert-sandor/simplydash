package main

import (
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
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

func DefaultConfig() *Config {
	return &Config{
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
	Stop()
	Get() Config
}

type ConfigServiceImpl struct {
	filePath  string
	config    *Config
	stopWatch chan struct{}
}

func NewConfigService(configPath string) ConfigService {
	return &ConfigServiceImpl{
		filePath:  path.Join(configPath, configFileName),
		config:    DefaultConfig(),
		stopWatch: make(chan struct{}, 1),
	}
}

func (configService *ConfigServiceImpl) Init() error {
	logrus.WithField("configFile", configService.filePath).Debug("checking if config file exists")
	if _, err := os.Stat(configService.filePath); err == nil {
		configService.loadConfigFile()
		if err != nil {
			return err
		}
	} else {
		err := configService.createConfigFile()
		if err != nil {
			return err
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Add(configService.filePath)
	if err != nil {
		return err
	}

	go configService.watchFile(watcher)
	return nil
}

func (configService *ConfigServiceImpl) Stop() {
	configService.stopWatch <- struct{}{}
}

func (configService *ConfigServiceImpl) Get() Config {
	return *configService.config
}

func (configService *ConfigServiceImpl) createConfigFile() error {
	logrus.WithField("configFile", configService.filePath).Debug("creating new config file")
	file, err := os.Create(configService.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	logrus.WithField("config", configService.config).Debug("marshalling config")
	bytes, err := yaml.Marshal(configService.config)
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
	bytes, err := os.ReadFile(configService.filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, configService.config)
	if err != nil {
		return err
	}

	return nil
}

func (configService *ConfigServiceImpl) watchFile(watcher *fsnotify.Watcher) {
	defer func() {
		err := watcher.Close()
		if err != nil {
			logrus.WithField("err", err).Error("failed to close fsnotify")
		}
	}()

	for {
		select {
		case <-configService.stopWatch:
			logrus.Info("received stop signal")
			return
		case err := <-watcher.Errors:
			logrus.WithField("err", err).Error("fsnotify")
			return
		case event := <-watcher.Events:
			logrus.WithField("event", event).Debug("fsnotify")
			if event.Op.Has(fsnotify.Write) {
				configService.loadConfigFile()
			}
		}
	}
}
