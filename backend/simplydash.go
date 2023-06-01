package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	cliArgs := parseCliArgs()

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
