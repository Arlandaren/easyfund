package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/Arlandaren/easyfund/internal/logger"
)

// CustomClaims структура для хранения custom claims в JWT
type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware проверяет JWT токен и извлекает user_id
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Log.Warnf("Missing authorization header for %s %s", c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing authorization header",
			})
			c.Abort()
			return
		}

		// Парсим "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Log.Warnf("Invalid authorization header format for %s %s", c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Парсим и валидируем токен
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Проверяем алгоритм подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			logger.Log.Warnf("Failed to parse JWT token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// Проверяем, валиден ли токен
		if !token.Valid {
			logger.Log.Warnf("Token is not valid")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is not valid",
			})
			c.Abort()
			return
		}

		// Извлекаем claims
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			logger.Log.Errorf("Could not extract claims from token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Could not extract claims from token",
			})
			c.Abort()
			return
		}

		// Проверяем, что user_id присутствует
		if claims.UserID == uuid.Nil {
			logger.Log.Warnf("User ID is empty in token claims")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID is empty in token claims",
			})
			c.Abort()
			return
		}

		// Устанавливаем user_id и email в контекст
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		logger.Log.Infof("User %s authenticated successfully", claims.UserID.String())

		c.Next()
	}
}

// GetUserIDFromContext извлекает user_id из контекста
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user_id not found in context")
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id has invalid type")
	}

	return userID, nil
}

// GetEmailFromContext извлекает email из контекста
func GetEmailFromContext(c *gin.Context) (string, error) {
	emailInterface, exists := c.Get("email")
	if !exists {
		return "", fmt.Errorf("email not found in context")
	}

	email, ok := emailInterface.(string)
	if !ok {
		return "", fmt.Errorf("email has invalid type")
	}

	return email, nil
}
