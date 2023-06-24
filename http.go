package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	cliArgs CliArguments
	routes  *Routes
	server  *http.Server
}

func NewHttpServer(cliArgs CliArguments, routes *Routes) *HttpServer {
	if cliArgs.Log.Level != logrus.DebugLevel.String() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(LoggingMiddleware())

	routes.setup(engine)

	return &HttpServer{
		cliArgs: cliArgs,
		routes:  routes,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cliArgs.Host, cliArgs.Port),
			Handler: engine,
		},
	}
}

func (httpServer *HttpServer) Start() {
	go func() {
		err := httpServer.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.WithField("err", err).Fatal("failed to start server")
		}
	}()
}

func (httpServer *HttpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.server.Shutdown(ctx); err != nil {
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
