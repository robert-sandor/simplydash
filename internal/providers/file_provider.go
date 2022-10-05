package providers

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"simplydash/internal/config"
	"simplydash/internal/models"
	"simplydash/internal/utils"
)

type FileProvider struct {
	Path  string
	Watch bool
	cache []models.Category
}

func NewFileProvider(config config.FileProviderConfig) (*FileProvider, error) {
	err := validateConfig(config)
	if err != nil {
		return nil, err
	}
	return &FileProvider{Path: config.Path, Watch: config.Watch, cache: []models.Category{}}, nil
}

func (f *FileProvider) Load() error {
	categories, err := f.getItems(utils.FileReader)
	if err != nil {
		return err
	}
	f.cache = categories
	return nil
}

func (f *FileProvider) Get() []models.Category {
	return f.cache
}

func (f *FileProvider) getItems(reader func(string) ([]byte, error)) ([]models.Category, error) {
	bytes, err := reader(f.Path)
	if err != nil {
		return nil, err
	}

	var categories []models.Category
	err = yaml.Unmarshal(bytes, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func validateConfig(config config.FileProviderConfig) error {
	if !utils.FileExists(config.Path) {
		return fmt.Errorf("file with Path = %s does not exist, or is a directory", config.Path)
	}

	return nil
}
