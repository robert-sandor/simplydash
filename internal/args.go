package internal

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	ConfigPath    arg
	IconCachePath arg
	Port          arg
	LogLevel      arg
}

func (a *Args) String() string {
	return fmt.Sprintf("Args: log_level = %s , port = %s , config_path = %s , icon_cache = %s",
		a.LogLevel.Get(), a.Port.Get(), a.ConfigPath.Get(), a.IconCachePath.Get())
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
		Port: arg{
			fromCmd:      flag.String("port", "", "port to use"),
			envVar:       "PORT",
			defaultValue: "8080",
		},
		LogLevel: arg{
			fromCmd:      flag.String("log-level", "", "logging level"),
			envVar:       "LOG_LEVEL",
			defaultValue: "WARN",
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
