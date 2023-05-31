package main

import "github.com/gin-gonic/gin"

type Routes struct {
	configService ConfigService
}

func NewRoutes(configService ConfigService) *Routes {
	return &Routes{
		configService: configService,
	}
}

func (routes *Routes) setup(engine *gin.Engine) {
	engine.GET("/health", healthcheck())
	engine.GET("/config", getConfig(routes.configService))
}

func getConfig(configService ConfigService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, configService.Get())
	}
}

func healthcheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(200)
	}
}
