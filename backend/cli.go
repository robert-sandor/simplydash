package main

import (
	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
)

type CliArguments struct {
	Host       string       `help:"Set the host to listen on" default:"0.0.0.0"      env:"HOST"`
	Port       int          `help:"Set the port to listen on" default:"8080"         env:"PORT"`
	ConfigPath string       `help:"Set path to config"        default:"/app/config/" env:"CONFIG_PATH" type:"existingdir"`
	ImagePath  string       `help:"Set path to image cache"   default:"/app/images/" env:"IMAGE_PATH"  type:"existingdir"`
	Log        LogArguments `                                                                                             embed:"" prefix:"log-" envprefix:"LOG_"`
}

type LogArguments struct {
	Level  string `help:"Set log level (debug / info / warn / error)" enum:"debug,info,warn,error" default:"info"    env:"LEVEL"`
	Format string `help:"Set log format (console / json)"             enum:"console,json"          default:"console" env:"FORMAT"`
}

func parseCliArgs() CliArguments {
	cliArgs := CliArguments{}
	kong.Parse(&cliArgs)
	setupLogging(cliArgs.Log)
	return cliArgs
}

func setupLogging(logArgs LogArguments) {
	logrusLevel, _ := logrus.ParseLevel(logArgs.Level)
	logrus.SetLevel(logrusLevel)

	if "json" == logArgs.Format {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			QuoteEmptyFields: true,
		})
	}
}
