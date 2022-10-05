package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"simplydash/internal"
	"simplydash/internal/config"
	"simplydash/internal/models"
	"strings"
)

func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

func GetConfig(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, cfg)
	}
}

func GetCategories(get func() []models.Category) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, get())
	}
}

func ItemHealthcheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Query("url")
		if strings.TrimSpace(url) == "" {
			c.Status(http.StatusNotFound)
			return
		}

		resp, err := http.Get(url)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Failed to close response body err = %+v", err)
			}
		}(resp.Body)

		c.Status(resp.StatusCode)
	}
}

func GetIcon(iconCache *internal.IconCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		urlStr := c.Query("url")
		icon, err := iconCache.GetIcon(urlStr)
		if err != nil {
			log.Printf("Failed to get icon for url = %s err = %s", urlStr, err)
			c.Status(http.StatusNotFound)
			return
		}

		c.File(icon)
	}
}
