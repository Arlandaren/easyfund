package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(cors.Default())

	// Metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	r.GET("/ping", HealthCheck)

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "register endpoint"})
		})
		api.POST("/auth/login", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "login endpoint"})
		})
	}

	return r
}
