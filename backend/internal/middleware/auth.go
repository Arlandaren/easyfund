package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/Arlandaren/easyfund/internal/logger"
)

// CustomClaims структура для хранения custom claims в JWT
type CustomClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
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

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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

		if !token.Valid {
			logger.Log.Warnf("Token is not valid")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is not valid",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			logger.Log.Errorf("Could not extract claims from token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Could not extract claims from token",
			})
			c.Abort()
			return
		}

		if claims.UserID == 0 {
			logger.Log.Warnf("User ID is empty in token claims")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID is empty in token claims",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		logger.Log.Infof("User %d authenticated successfully", claims.UserID)

		c.Next()
	}
}

// GetUserIDFromContext извлекает user_id из контекста
func GetUserIDFromContext(c *gin.Context) (int64, error) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("userID not found in context")
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		return 0, fmt.Errorf("userID has invalid type")
	}

	return userID, nil
}