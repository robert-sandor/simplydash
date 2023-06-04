package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type App struct {
	Name        string         `yaml:"name"        json:"name"`
	Description string         `yaml:"description" json:"description"`
	Group       string         `yaml:"group"       json:"group"`
	Link        string         `yaml:"link"        json:"link"`
	Icon        string         `yaml:"icon"        json:"icon"`
	Healthcheck AppHealthcheck `yaml:"healthcheck" json:"healthcheck"`
}

type AppHealthcheck struct {
	Disable               bool   `json:"disable"                  yaml:"disable"`
	Link                  string `json:"link"                     yaml:"link"`
	PollIntervalInSeconds int    `json:"poll_interval_in_seconds" yaml:"poll_interval_in_seconds"`
}

var spaceRegexp = regexp.MustCompile(`\s+`)

func defaultIcon(name string) string {
	return fmt.Sprintf(
		"https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/%s.png",
		normalizeName(name),
	)
}

func normalizeName(name string) string {
	return spaceRegexp.ReplaceAllString(strings.TrimSpace(strings.ToLower(name)), "-")
}

func (app *App) validate() error {
	if "" == app.Name {
		return errors.New("name is required")
	}

	if "" == app.Link {
		return errors.New("link is required")
	}

	if "" == app.Description {
		app.Description = app.Link
	}

	if "" == app.Icon {
		app.Icon = defaultIcon(app.Name)
	}

	return app.Healthcheck.validate(app.Link)
}

func (heathcheck *AppHealthcheck) validate(appLink string) error {
	if "" == heathcheck.Link {
		heathcheck.Link = appLink
	}

	if 0 == heathcheck.PollIntervalInSeconds {
		heathcheck.PollIntervalInSeconds = 60
	}

	return nil
}

func (app *App) UnmarshalYAML(value *yaml.Node) error {
	type a App
	err := value.Decode((*a)(app))
	if err != nil {
		return err
	}
	err = app.validate()
	if err != nil {
		return err
	}
	return nil
}
