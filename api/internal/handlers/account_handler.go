package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/Arlandaren/easyfund/internal/logger"
    "github.com/Arlandaren/easyfund/internal/middleware"
    "github.com/Arlandaren/easyfund/internal/services"
)

type UserBankAccountHandler struct {
    service services.UserBankAccountService
}

func NewUserBankAccountHandler(service services.UserBankAccountService) *UserBankAccountHandler {
    return &UserBankAccountHandler{service: service}
}

// GET /api/v1/users/:id/balance (защищенный)
func (h *UserBankAccountHandler) GetTotalBalance(c *gin.Context) {
    requestingUserID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    userIDStr := c.Param("id")
    userID, err := strconv.ParseInt(userIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Проверяем права доступа
    if requestingUserID != userID {
        logger.Log.Warnf("User %d tried to get balance of user %d", requestingUserID, userID)
        c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
        return
    }

    balance, err := h.service.GetTotalBalance(c.Request.Context(), userID)
    if err != nil {
        logger.Log.Errorf("Failed to get total balance: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total balance"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"total_balance": balance})
}

// GET /api/v1/users/:id/accounts (защищенный)
func (h *UserBankAccountHandler) GetUserAccounts(c *gin.Context) {
    requestingUserID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    userIDStr := c.Param("id")
    userID, err := strconv.ParseInt(userIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Проверяем права доступа
    if requestingUserID != userID {
        logger.Log.Warnf("User %d tried to get accounts of user %d", requestingUserID, userID)
        c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
        return
    }

    accounts, err := h.service.GetUserAccounts(c.Request.Context(), userID)
    if err != nil {
        logger.Log.Errorf("Failed to get user accounts: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user accounts"})
        return
    }

    c.JSON(http.StatusOK, accounts)
}