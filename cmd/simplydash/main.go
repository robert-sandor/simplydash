package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/robert-sandor/simplydash/internal"
)

func main() {
	args := internal.GetArgs()
	internal.SetupLogging(args)
	internal.SetupSlog(args)

	slog.LogAttrs(context.Background(), slog.LevelDebug, "loading config",
		slog.String("configFile", args.ConfigFile))

	config, err := internal.GetConfig(args)
	if err != nil {
		slog.Error("invalid config", slog.Any("error", err))
	}

	healthCheckService := internal.NewHealthcheckService()
	healthCheckService.Init()

	slog.Debug("initializing app service")
	appService := internal.NewAppService(config, healthCheckService)
	appService.Init()

	slog.Debug("initializing websocket server")
	websocketServer := internal.NewWebsocketServer(appService)
	websocketServer.Init()

	imageService := internal.NewImageService(args.ImageCacheDir)

	slog.Debug("initializing http server")
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	if args.AccessLogs {
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogValuesFunc: func(_ echo.Context, values middleware.RequestLoggerValues) error {
				slog.Info("http request", slog.Group("request",
					slog.String("uri", values.URI),
					slog.Any("latency", values.Latency),
					slog.String("protocol", values.Protocol),
					slog.String("remoteIp", values.RemoteIP),
					slog.String("method", values.Method),
					slog.Int("status", values.Status),
					slog.String("contentLength", values.ContentLength),
					slog.Int64("responseSize", values.ResponseSize),
				))
				return nil
			},
			LogLatency:       true,
			LogProtocol:      true,
			LogRemoteIP:      true,
			LogMethod:        true,
			LogURI:           true,
			LogStatus:        true,
			LogContentLength: true,
			LogResponseSize:  true,
		}))
	}

	e.Static("/", "./web/build")

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	e.GET("/ws", func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			slog.Error("failed upgrade to websocket", slog.Any("error", err))
			return err
		}

		websocketServer.Connect(time.Now().String(), ws)
		return nil
	})

	e.GET("/image", func(c echo.Context) error {
		filePath, err := imageService.Get(c.QueryParam("url"))
		if err != nil {
			slog.Error("image not found", slog.Any("error", err))
			return c.NoContent(http.StatusNotFound)
		}

		return c.File(filePath)
	})

	e.GET("/settings", func(c echo.Context) error {
		err := c.JSON(http.StatusOK, config.App)
		if err != nil {
			_ = c.NoContent(http.StatusInternalServerError)
		}
		return nil
	})

	e.GET("/timeout", func(c echo.Context) error {
		time.Sleep(5 * time.Minute)
		_ = c.NoContent(http.StatusOK)
		return nil
	})

	err = e.Start(fmt.Sprintf("%s:%s", args.Host, args.Port))
	if err != nil {
		slog.Error("starting http server", slog.Any("error", err))
	}
}
