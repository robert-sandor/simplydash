package main

import (
	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
)

var CliArguments struct {
	Host       string `help:"Set the host to listen on" default:"0.0.0.0"`
	Port       int    `help:"Set the port to listen on" default:"8080"`
	ConfigPath string `help:"Set path to config" type:"existingdir" default:"/app/config/"`
	CachePath  string `help:"Set path to image cache" type:"existingdir" default:"/app/images/"`
	Log        struct {
		Level  string `help:"Set log level (debug / info / warn / error)" enum:"debug,info,warn,error" default:"info"`
		Format string `help:"Set log format (console / json)" enum:"console,json" default:"console"`
	} `embed:"" prefix:"log-"`
}

func main() {
	kong.Parse(&CliArguments)
	setupLogging(CliArguments.Log.Level, CliArguments.Log.Format)
	logrus.WithField("args", CliArguments).Info("starting...")
}

func setupLogging(level, format string) {
	logrusLevel, _ := logrus.ParseLevel(level)
	logrus.SetLevel(logrusLevel)

	if "json" == format {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			QuoteEmptyFields: true,
		})
	}
}
