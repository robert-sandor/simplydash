package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
)

type FileProvider struct {
	Path     string
	cache    []Category
	reader   func(string) ([]byte, error)
	notifier *WebsocketNotifier
}

func NewFileProvider(config FileProviderConfig, notifier *WebsocketNotifier, reader func(string) ([]byte, error)) (*FileProvider, error) {
	err := validateConfig(config)
	if err != nil {
		return nil, err
	}
	return &FileProvider{Path: config.Path, notifier: notifier, reader: reader, cache: []Category{}}, nil
}

func (f *FileProvider) Init() error {
	f.updateCache()
	return nil
}

func (f *FileProvider) Get() []Category {
	return f.cache
}

func (f *FileProvider) Update() {
	f.updateCache()
	f.notifier.Update(UpdateCategories)
}

func (f *FileProvider) updateCache() {
	bytes, err := f.reader(f.Path)
	if err != nil {
		log.Printf("Failed to read from path = %s err = %+v", f.Path, err)
		return
	}

	var categories []Category
	err = yaml.Unmarshal(bytes, &categories)
	if err != nil {
		log.Printf("Failed to read yaml from path = %s err = %+v", f.Path, err)
		return
	}

	f.cache = categories
}

func validateConfig(config FileProviderConfig) error {
	if !FileExists(config.Path) {
		return fmt.Errorf("file with path = %s does not exist, or is a directory", config.Path)
	}

	return nil
}
