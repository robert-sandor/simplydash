package main

import "github.com/gin-gonic/gin"

type Routes struct {
	configService ConfigService
	appService    AppService
}

func NewRoutes(configService ConfigService, appService AppService) *Routes {
	return &Routes{
		configService: configService,
		appService:    appService,
	}
}

func (routes *Routes) setup(engine *gin.Engine) {
	engine.GET("/health", healthcheck())
	engine.GET("/config", getConfig(routes.configService))
	engine.GET("/apps", getApps(routes.appService))
}

func getApps(appService AppService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, appService.GetApps())
	}
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
