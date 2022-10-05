package internal

import (
	"flag"
	"os"
)

type Args struct {
	ConfigPath    arg
	IconCachePath arg
}

type arg struct {
	fromCmd      *string
	envVar       string
	defaultValue string
}

func NewArgs() *Args {
	args := Args{
		ConfigPath: arg{
			fromCmd:      flag.String("config", "", "path to config file"),
			envVar:       "CONFIG_FILE_PATH",
			defaultValue: "config.yml",
		},
		IconCachePath: arg{
			fromCmd:      flag.String("icon-cache", "", "path to icon cache"),
			envVar:       "ICON_CACHE_PATH",
			defaultValue: "icons",
		},
	}
	flag.Parse()
	return &args
}

func (a *arg) Get() string {
	if *a.fromCmd != "" {
		return *a.fromCmd
	}

	if env := os.Getenv(a.envVar); env != "" {
		return env
	}

	return a.defaultValue
}
