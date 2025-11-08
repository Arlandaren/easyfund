package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Arlandaren/easyfund/internal/logger"
	"github.com/Arlandaren/easyfund/internal/middleware"
	"github.com/Arlandaren/easyfund/internal/models"
	"github.com/Arlandaren/easyfund/internal/services"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GET /api/v1/users/random (публичный)
func (h *UserHandler) GetRandomUser(c *gin.Context) {
	user, err := h.service.GetRandomUser(c.Request.Context())
	if err != nil {
		logger.Log.Errorf("Failed to get random user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get random user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GET /api/v1/users/:id (защищенный)
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Проверяем, что пользователь может только смотреть свой профиль или он admin
	requestingUserID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Если не его профиль, возвращаем ошибку (в реальном приложении могут быть roles)
	if requestingUserID != userID {
		logger.Log.Warnf("User %s tried to access user %s profile", requestingUserID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Errorf("Failed to get user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// POST /api/v1/users (защищенный)
func (h *UserHandler) CreateUser(c *gin.Context) {
	// Получаем user_id из контекста
	requestingUserID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Используем ID из контекста (аутентифицированный пользователь)
	user.UserID = requestingUserID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = h.service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		logger.Log.Errorf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// PUT /api/v1/users/:id (защищенный)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Проверяем права доступа
	requestingUserID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if requestingUserID != userID {
		logger.Log.Warnf("User %s tried to update user %s", requestingUserID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user.UserID = userID
	user.UpdatedAt = time.Now()

	err = h.service.UpdateUser(c.Request.Context(), &user)
	if err != nil {
		logger.Log.Errorf("Failed to update user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DELETE /api/v1/users/:id (защищенный)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Проверяем права доступа
	requestingUserID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if requestingUserID != userID {
		logger.Log.Warnf("User %s tried to delete user %s", requestingUserID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	err = h.service.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Errorf("Failed to delete user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
