package internal

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Router struct {
	ginRouter *gin.Engine
}

func NewRouter(cfg *Config, svc *AggregatorService, iconCache *IconCache, notifier *WebsocketNotifier) *Router {
	r := router(cfg, svc, iconCache, notifier)
	return &Router{ginRouter: r}
}

func router(cfg *Config, svc *AggregatorService, iconCache *IconCache, notifier *WebsocketNotifier) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"content-type"},
	}))

	r.GET("/health", Health())
	r.GET("/api/config", GetConfig(cfg))
	r.GET("/api/categories", GetCategories(svc.Get))
	r.GET("/api/url/health", ItemHealthcheck())
	r.GET("/api/icon", GetIcon(iconCache))
	r.GET("/api/ws", Websocket(upgrader(), notifier))

	r.StaticFile("/", "web/index.html")
	r.Static("/assets", "web/assets")

	return r
}

func upgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		HandshakeTimeout: 60 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func (r *Router) Run(port string) {
	err := r.ginRouter.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Printf("Server stopped err = %+v", err)
	}
}
