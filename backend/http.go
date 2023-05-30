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

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func healthcheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(200)
	}
}

func setupRoutes(engine *gin.Engine) {
	engine.GET("/health", healthcheck())
}

func startServer(cliArgs CliArguments) {
	logrus.WithField("args", cliArgs).Info("starting...")

	if cliArgs.Log.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(LoggingMiddleware())

	setupRoutes(engine)

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
