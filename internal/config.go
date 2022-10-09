package internal

import (
	"gopkg.in/yaml.v3"
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

func NewConfig(path string, reader func(string) ([]byte, error), writer func(string, []byte) error) *Config {
	cfg := DefaultConfig()
	err := cfg.Load(path, reader)
	if err != nil {
		if os.IsNotExist(err) {
			Log.Info.Printf("No config file found at path %s - creating with defaults", path)
			err = createConfigFile(cfg, path, writer)
			if err != nil {
				Log.Error.Printf("Failed to create new config file with defaults at path %s err = %+v", path, err)
			}
			return cfg
		}
		Log.Error.Printf("Failed to load config from path %s err = %+v", path, err)
	}
	Log.Debug.Printf("Successfully loaded config = %+v", cfg)
	return cfg
}

func (c *Config) Load(path string, reader func(string) ([]byte, error)) error {
	bytes, err := reader(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return err
	}

	err = c.validate()
	if err != nil {
		return err
	}
	return nil
}

func createConfigFile(cfg *Config, path string, writer func(string, []byte) error) error {
	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err = writer(path, bytes); err != nil {
		return err
	}

	return nil
}

func (c *Config) validate() error {
	return nil
}
