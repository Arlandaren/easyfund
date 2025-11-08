package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Arlandaren/easyfund/internal/logger"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Log.Errorf("Handler error: %v", err)
			}

			// Отправляем последнюю ошибку
			lastErr := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": lastErr.Error(),
			})
		}
	}
}
