package internal

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type App struct {
	Name        string         `json:"name"`
	Link        string         `json:"link"`
	Group       string         `json:"group"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"`
	Healthcheck AppHealthcheck `json:"healthcheck"`
}

type AppHealthcheck struct {
	Enabled  bool          `json:"enabled"`
	Health   AppHealth     `json:"health"`
	Interval time.Duration `json:"poll_interval"`
	Timeout  time.Duration `json:"timeout"`
}

type AppHealth uint32

func (a AppHealth) MarshalText() ([]byte, error) {
	switch a {
	case Healthy:
		return []byte("healthy"), nil
	case Timeout:
		return []byte("timeout"), nil
	case Warning:
		return []byte("warning"), nil
	case Error:
		return []byte("error"), nil
	case Unknown:
		return []byte("unknown"), nil
	}

	return nil, fmt.Errorf("invalid value %d", a)
}

func (a AppHealth) String() string {
	if bytes, err := a.MarshalText(); err == nil {
		return string(bytes)
	}
	return "unknown"
}

const (
	Healthy AppHealth = iota
	Timeout
	Warning
	Error
	Unknown
)

type AppGroup struct {
	Name string `json:"name"`
	Apps []App  `json:"apps"`
}

func NewAppGroup(name string) AppGroup {
	return AppGroup{
		Name: name,
		Apps: make([]App, 0),
	}
}

func (app *App) Validate() (errs []error) {
	errs = make([]error, 0)
	if strings.TrimSpace(app.Name) == "" {
		errs = append(errs, errors.New("name is required"))
	}

	if strings.TrimSpace(app.Link) == "" {
		errs = append(errs, errors.New("link is required"))
	}

	if strings.TrimSpace(app.Group) == "" {
		errs = append(errs, errors.New("group is required"))
	}

	if strings.TrimSpace(app.Description) == "" {
		app.Description = app.Link
	}

	if strings.TrimSpace(app.Icon) == "" {
		app.Icon = app.Name
	}
	app.resolveIconUrl()

	errs = append(errs, app.Healthcheck.Validate()...)
	return
}

func (app *App) resolveIconUrl() {
	if url, err := url.ParseRequestURI(app.Icon); err == nil {
		log.Debug(url)
		return
	}

	regexMatcher := regexp.MustCompile(`[\s_-]+`)
	normalized := strings.ToLower(regexMatcher.ReplaceAllString(strings.TrimSpace(app.Icon), "-"))
	app.Icon = fmt.Sprintf(
		"https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/svg/%s.svg",
		normalized,
	)
}

func (appHealth *AppHealthcheck) Validate() (errs []error) {
	errs = make([]error, 0)

	if appHealth.Interval <= 0 {
		appHealth.Interval = time.Minute
	}

	if appHealth.Timeout <= 0 {
		appHealth.Timeout = 30 * time.Second
	}

	return
}
