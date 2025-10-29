package middleware

import (
	"github.com/Arlandaren/easyfund/internal/config"
	"github.com/Arlandaren/easyfund/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authUsecase *usecase.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
				"code":  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Извлекаем токен из "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
				"code":  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := authUsecase.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
				"code":  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Устанавливаем claims в контекст
		c.Set("claims", claims)
		c.Next()
	}
}

func CORSMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", strings.Join(config.CORSAllowedOrigins, ","))
		c.Header("Access-Control-Allow-Methods", strings.Join(config.CORSAllowedMethods, ","))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.CORSAllowedHeaders, ","))
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
