package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(cfg CORSConfig) gin.HandlerFunc {
	c := cors.DefaultConfig()
	if len(cfg.AllowedOrigins) > 0 {
		c.AllowOrigins = cfg.AllowedOrigins
	} else {
		c.AllowAllOrigins = true
	}
	c.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	c.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	c.AllowCredentials = true
	return cors.New(c)
}
