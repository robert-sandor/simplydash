package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type App struct {
	Name        string `yaml:"name"        json:"name"`
	Description string `yaml:"description" json:"description"`
	Link        string `yaml:"link"        json:"link"`
	Icon        string `yaml:"icon"        json:"icon"`
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
