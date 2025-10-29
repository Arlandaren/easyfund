package http

import (
    "github.com/Arlandaren/easyfund/internal/usecase"
    "github.com/Arlandaren/easyfund/pkg/banking"
    "net/http"

    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
    return &AuthHandler{
        authUsecase: authUsecase,
    }
}

// GetRandomDemoClient получает случайного тестового клиента
func (h *AuthHandler) GetRandomDemoClient(c *gin.Context) {
    client, err := h.authUsecase.GetRandomDemoClient()
    if err != nil {
        c.JSON(http.StatusInternalServerError, banking.ErrorResponse{
            Error:   "Failed to get demo client",
            Message: err.Error(),
            Code:    http.StatusInternalServerError,
        })
        return
    }

    c.JSON(http.StatusOK, client)
}

// DemoLoginRequest структура запроса входа демо клиента
type DemoLoginRequest struct {
    PersonID string `json:"person_id" binding:"required"`
    Password string `json:"password" binding:"required"`
    FullName string `json:"full_name" binding:"required"`
}

// LoginDemoClient логин демо клиента
func (h *AuthHandler) LoginDemoClient(c *gin.Context) {
    var req DemoLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Передаем всю структуру запроса, а не отдельные аргументы
    authResp, err := h.authUsecase.LoginDemoClient(&banking.DemoClientLoginRequest{
        PersonID: req.PersonID,
        Password: req.Password,
        FullName: req.FullName,
    })
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, authResp)
}

// Login обычный логин
func (h *AuthHandler) Login(c *gin.Context) {
    var req banking.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, banking.ErrorResponse{
            Error:   "Invalid request",
            Message: err.Error(),
            Code:    http.StatusBadRequest,
        })
        return
    }

    authResp, err := h.authUsecase.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, banking.ErrorResponse{
            Error:   "Authentication failed",
            Message: err.Error(),
            Code:    http.StatusUnauthorized,
        })
        return
    }

    c.JSON(http.StatusOK, authResp)
}

// Register регистрация пользователя
func (h *AuthHandler) Register(c *gin.Context) {
    var req banking.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, banking.ErrorResponse{
            Error:   "Invalid request",
            Message: err.Error(),
            Code:    http.StatusBadRequest,
        })
        return
    }

    authResp, err := h.authUsecase.Register(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, banking.ErrorResponse{
            Error:   "Registration failed",
            Message: err.Error(),
            Code:    http.StatusInternalServerError,
        })
        return
    }

    c.JSON(http.StatusCreated, authResp)
}
