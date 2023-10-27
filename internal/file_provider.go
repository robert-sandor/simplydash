package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type appConfig struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Group       string            `yaml:"group"`
	Link        string            `yaml:"link"`
	Icon        string            `yaml:"icon"`
	Healthcheck healthcheckConfig `yaml:"healthcheck"`
}

type healthcheckConfig struct {
	Link     string        `yaml:"link"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
	Enable   bool          `yaml:"enable"`
}

type FileProviderConfig struct {
	Path string `yaml:"path" json:"path"`
}

type FileProvider struct {
	watcher          *fsnotify.Watcher
	logger           *slog.Logger
	notificationChan chan<- string
	id               string
	path             string
	apps             []App
}

func NewFileProvider(name string, config FileProviderConfig, notificationChan chan<- string) Provider {
	id := fmt.Sprintf("file-%s", name)
	return &FileProvider{
		id:               id,
		path:             config.Path,
		apps:             make([]App, 0),
		notificationChan: notificationChan,
		logger:           slog.With("id", id),
	}
}

func (fp *FileProvider) ID() string {
	return fp.id
}

func (fp *FileProvider) Apps() []App {
	return fp.apps
}

func (fp *FileProvider) Init() error {
	absPath, err := filepath.Abs(fp.path)
	if err != nil {
		return err
	}
	fp.path = absPath

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	fp.watcher = watcher
	err = fp.watcher.Add(absPath)
	if err != nil {
		return err
	}

	go fp.parseFile()
	go fp.watch()
	return nil
}

func (fp *FileProvider) watch() {
	defer func(w *fsnotify.Watcher) {
		_ = w.Close()
	}(fp.watcher)

	for {
		select {
		case err := <-fp.watcher.Errors:
			fp.logger.Error("fsnotify", "error", err)
			return
		case event := <-fp.watcher.Events:
			if event.Has(fsnotify.Write) {
				fp.parseFile()
			}
		}
	}
}

func (fp *FileProvider) parseFile() {
	bytes, err := os.ReadFile(fp.path)
	if err != nil {
		fp.logger.Error("reading file", "path", fp.path, "error", err)
		return
	}

	appConfigs := make([]appConfig, 0)
	err = yaml.Unmarshal(bytes, &appConfigs)
	if err != nil {
		fp.logger.Error("parsing yaml", "path", fp.path, "error", err)
		return
	}

	apps := make([]App, 0)
	for _, appConfig := range appConfigs {
		app := appConfig.toApp()
		errs := app.Validate()
		if len(errs) > 0 {
			fp.logger.Error("invalid app config", "appConfig", appConfig, "error", errors.Join(errs...))
		} else {
			apps = insertOrdered(apps, app)
		}
	}

	if !reflect.DeepEqual(fp.apps, apps) {
		fp.apps = apps
		fp.notificationChan <- fp.id
	}
}

func (cfg appConfig) toApp() App {
	return App{
		Name:        cfg.Name,
		Description: cfg.Description,
		Link:        cfg.Link,
		Icon:        cfg.Icon,
		Group:       cfg.Group,
		Healthcheck: AppHealthcheck{
			Health:   Unknown,
			Enabled:  cfg.Healthcheck.Enable,
			Interval: cfg.Healthcheck.Interval,
			Timeout:  cfg.Healthcheck.Timeout,
		},
	}
}
