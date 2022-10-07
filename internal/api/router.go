package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"simplydash/internal"
	"simplydash/internal/config"
	"time"
)

type Router struct {
	cfg        *config.Config
	svc        *internal.Service
	router     *gin.Engine
	wsUpgrader *websocket.Upgrader
	iconCache  *internal.IconCache
}

func NewRouter(
	cfg *config.Config,
	svc *internal.Service,
	iconCache *internal.IconCache,
) *Router {
	u := upgrader()
	r := router(cfg, svc, iconCache, u)
	return &Router{
		cfg:        cfg,
		svc:        svc,
		router:     r,
		wsUpgrader: u,
		iconCache:  iconCache,
	}
}

func router(cfg *config.Config, svc *internal.Service, iconCache *internal.IconCache, upgrader *websocket.Upgrader) *gin.Engine {
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
	r.GET("/api/ws", Websocket(upgrader, svc))
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
	err := r.router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Printf("Server stopped err = %+v", err)
	}
}
