package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gin-gonic/gin"
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
	logrus.WithField("args", cliArgs).Info("starting...")

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(LoggingMiddleware())
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cliArgs.Host, cliArgs.Port),
		Handler: engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.WithField("err", err).Fatal("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatal("forcing server shutdown")
	}

	logrus.Info("server shut down")
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		startTime := time.Now()
		context.Next()
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		method := context.Request.Method
		uri := context.Request.RequestURI
		statusCode := context.Writer.Status()
		clientIP := context.ClientIP()
		logrus.WithFields(logrus.Fields{
			"method":  method,
			"uri":     uri,
			"status":  statusCode,
			"latency": latency,
			"client":  clientIP,
		}).Info("HTTP request")
		context.Next()
	}
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
