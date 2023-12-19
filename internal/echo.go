package internal

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouting(e *echo.Echo, websocketServer *WebsocketServer, imageService ImageService, config Config) {
	e.Static("/", "./web/build")

	e.GET("/ws", handleWebsocket(websocketServer))
	e.GET("/image", getImage(imageService))
	e.GET("/settings", getSettings(config))
}

func getSettings(config Config) func(c echo.Context) error {
	return func(c echo.Context) error {
		err := c.JSON(http.StatusOK, config.App)
		if err != nil {
			_ = c.NoContent(http.StatusInternalServerError)
		}
		return nil
	}
}

func getImage(imageService ImageService) func(c echo.Context) error {
	return func(c echo.Context) error {
		filePath, err := imageService.Get(c.QueryParam("url"))
		if err != nil {
			slog.Error("image not found", slog.Any("error", err))
			return c.NoContent(http.StatusNotFound)
		}

		return c.File(filePath)
	}
}

func handleWebsocket(websocketServer *WebsocketServer) func(c echo.Context) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			slog.Error("failed upgrade to websocket", slog.Any("error", err))
			return err
		}

		websocketServer.Connect(time.Now().String(), ws)
		return nil
	}
}

func CreateEcho(args Args) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	if args.AccessLogs {
		e.Use(loggingFunc())
	}
	return e
}

func loggingFunc() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
	})
}
