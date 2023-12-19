package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/robert-sandor/simplydash/internal"
)

func main() {
	args := internal.GetArgs()
	internal.SetupSlog(args)

	slog.LogAttrs(context.Background(), slog.LevelDebug, "loading config", slog.String("configFile", args.ConfigFile))
	config, err := internal.GetConfig(args)
	logErrorAndExit(err, "invalid config")

	healthCheckService := internal.NewHealthcheckService()
	healthCheckService.Init()

	slog.Debug("initializing app service")
	appService := internal.NewAppService(config, healthCheckService)
	appService.Init()

	slog.Debug("initializing websocket server")
	websocketServer := internal.NewWebsocketServer(appService)
	websocketServer.Init()

	slog.Debug("initializing image service")
	imageService := internal.NewImageService(args.ImageCacheDir)

	slog.Debug("initializing echo")
	echo := internal.CreateEcho(args)
	internal.SetupRouting(echo, websocketServer, imageService, config)

	err = echo.Start(fmt.Sprintf("%s:%s", args.Host, args.Port))
	logErrorAndExit(err, "server shut down")
}

func logErrorAndExit(err error, msg string) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}
