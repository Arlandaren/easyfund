package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Arlandaren/easyfund/internal/logger"
	"github.com/Arlandaren/easyfund/internal/middleware"
	"github.com/Arlandaren/easyfund/internal/services"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// GET /api/v1/users/:user_id/transactions (защищенный)
func (h *TransactionHandler) GetUserTransactionHistory(c *gin.Context) {
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

	if requestingUserID != userID {
		logger.Log.Warnf("User %d tried to get transactions of user %d", requestingUserID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	history, err := h.service.GetUserTransactionHistory(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Errorf("Failed to get transaction history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transaction history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GET /api/v1/users/:user_id/banks/:bank_id/transactions (защищенный)
func (h *TransactionHandler) GetBankTransactionHistory(c *gin.Context) {
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

	if requestingUserID != userID {
		logger.Log.Warnf("User %d tried to get bank transactions of user %d", requestingUserID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	bankIDStr := c.Param("bank_id")
	bankID, err := strconv.ParseInt(bankIDStr, 10, 16)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bank ID"})
		return
	}

	transactions, err := h.service.GetBankTransactionHistory(c.Request.Context(), userID, int16(bankID))
	if err != nil {
		logger.Log.Errorf("Failed to get bank transaction history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bank transaction history"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}