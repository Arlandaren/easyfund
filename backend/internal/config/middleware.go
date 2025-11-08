package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(cfg CORSConfig) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	
	if len(cfg.AllowedOrigins) > 0 {
		corsConfig.AllowOrigins = cfg.AllowedOrigins
	} else {
		corsConfig.AllowAllOrigins = true
	}
	
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 86400 // 24 часа
	
	return cors.New(corsConfig)
}
