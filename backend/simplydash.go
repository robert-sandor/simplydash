package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
)

type CliArguments struct {
	Host       string       `help:"Set the host to listen on" default:"0.0.0.0"`
	Port       int          `help:"Set the port to listen on" default:"8080"`
	ConfigPath string       `help:"Set path to config"        default:"/app/config/" type:"existingdir"`
	CachePath  string       `help:"Set path to image cache"   default:"/app/images/" type:"existingdir"`
	Log        LogArguments `                                                                           embed:"" prefix:"log-"`
}

type LogArguments struct {
	Level  string `help:"Set log level (debug / info / warn / error)" enum:"debug,info,warn,error" default:"info"`
	Format string `help:"Set log format (console / json)"             enum:"console,json"          default:"console"`
}

func main() {
	cliArgs := CliArguments{}
	kong.Parse(&cliArgs)
	setupLogging(cliArgs.Log)

	configService := NewConfigService(cliArgs.ConfigPath)
	err := configService.Init()
	if err != nil {
		logrus.WithField("err", err).Fatal("config service failed to start")
	}

	routes := NewRoutes(configService)

	httpServer := NewHttpServer(cliArgs, routes)
	httpServer.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("shutting down")

	httpServer.Stop()
	configService.Stop()
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
