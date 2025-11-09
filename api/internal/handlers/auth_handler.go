package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/Arlandaren/easyfund/internal/logger"
	"github.com/Arlandaren/easyfund/internal/middleware"
	"github.com/Arlandaren/easyfund/internal/models"
	"github.com/Arlandaren/easyfund/internal/services"
)

type AuthHandler struct {
	userService  services.UserService
	tokenService services.TokenService
}

func NewAuthHandler(userService services.UserService, tokenService services.TokenService) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string       `json:"token"`
	User      *models.User `json:"user"`
	ExpiresIn int64        `json:"expires_in"`
}

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required,min=6"`
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		logger.Log.Warnf("Invalid login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.userService.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil || user == nil || user.UserID == 0 {
		logger.Log.Warnf("Invalid credentials for email: %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.Log.Warnf("Invalid password for user: %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := h.tokenService.GenerateToken(user)
	if err != nil {
		logger.Log.Errorf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		User:      user,
		ExpiresIn: int64(24 * 60 * 60),
	})
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		logger.Log.Warnf("Invalid register request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Errorf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user := &models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hashBytes),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.userService.CreateUser(c.Request.Context(), user); err != nil {
		logger.Log.Errorf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := h.tokenService.GenerateToken(user)
	if err != nil {
		logger.Log.Errorf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, LoginResponse{
		Token:     token,
		User:      user,
		ExpiresIn: int64(24 * 60 * 60),
	})
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}
	oldToken := authHeader[7:]

	newToken, err := h.tokenService.RefreshToken(oldToken)
	if err != nil {
		logger.Log.Warnf("Failed to refresh token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken, "expires_in": int64(24 * 60 * 60)})
}

// GET /api/v1/auth/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		logger.Log.Errorf("Failed to get user_id from context: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Errorf("Failed to get user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}