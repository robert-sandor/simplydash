package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/robert-sandor/simplydash/internal"
)

func main() {
	args := internal.GetArgs()
	internal.SetupLogging(args)

	log.WithField("configFile", args.ConfigFile).Debug("loading config")
	config, err := internal.GetConfig(args)
	if err != nil {
		log.WithError(err).Fatal("invalid config")
	}

	healthCheckService := internal.NewHealthcheckService()
	healthCheckService.Init()

	log.Debug("init app service")
	appService := internal.NewAppService(config, healthCheckService)
	appService.Init()

	log.Debug("init websocket server")
	websocketServer := internal.NewWebsocketServer(appService)
	websocketServer.Init()

	imageService := internal.NewImageService(args.ImageCacheDir)

	log.Debug("init echo")
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	if args.AccessLogs {
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogValuesFunc: func(_ echo.Context, values middleware.RequestLoggerValues) error {
				log.WithFields(log.Fields{
					"uri":           values.URI,
					"latency":       values.Latency,
					"protocol":      values.Protocol,
					"remoteIp":      values.RemoteIP,
					"method":        values.Method,
					"status":        values.Status,
					"contentLength": values.ContentLength,
					"responseSize":  values.ResponseSize,
				}).Info("request")

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
			log.WithError(err).Error("upgrading connection")
			return err
		}

		websocketServer.Connect(time.Now().String(), ws)
		return nil
	})

	e.GET("/image", func(c echo.Context) error {
		filePath, err := imageService.Get(c.QueryParam("url"))
		if err != nil {
			log.WithError(err).Warn("image not found")
			return c.NoContent(http.StatusNotFound)
		}

		return c.File(filePath)
	})

	e.GET("/settings", func(c echo.Context) error {
		c.JSON(http.StatusOK, config.App)
		return nil
	})

	e.GET("/timeout", func(c echo.Context) error {
		time.Sleep(5 * time.Minute)
		c.NoContent(http.StatusOK)
		return nil
	})

	err = e.Start(fmt.Sprintf("%s:%s", args.Host, args.Port))
	if err != nil {
		log.WithError(err).Fatal("http server")
	}
}
