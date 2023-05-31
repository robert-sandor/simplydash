package main

import "github.com/gin-gonic/gin"

type Routes struct{}

func (routes *Routes) setup(engine *gin.Engine) {
	engine.GET("/health", healthcheck())
}

func healthcheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(200)
	}
}
