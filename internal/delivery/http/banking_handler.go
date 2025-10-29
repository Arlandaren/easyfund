package http

import (
    "github.com/Arlandaren/easyfund/internal/usecase"
    "github.com/Arlandaren/easyfund/pkg/banking"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type BankingHandler struct {
    bankingUsecase *usecase.BankingUsecase
}

func NewBankingHandler(bankingUsecase *usecase.BankingUsecase) *BankingHandler {
    return &BankingHandler{
        bankingUsecase: bankingUsecase,
    }
}

// GetAccounts получает список счетов клиента
func (h *BankingHandler) GetAccounts(c *gin.Context) {
    claims, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, banking.ErrorResponse{
            Error: "No authentication claims found",
            Code:  http.StatusUnauthorized,
        })
        return
    }
    userClaims := claims.(*banking.JWTClaims)

    accounts, err := h.bankingUsecase.GetAccounts(userClaims)
    if err != nil {
        c.JSON(http.StatusInternalServerError, banking.ErrorResponse{
            Error:   "Failed to get accounts",
            Message: err.Error(),
            Code:    http.StatusInternalServerError,
        })
        return
    }

    c.JSON(http.StatusOK, accounts)
}

// GetTransactions получает транзакции по счету
func (h *BankingHandler) GetTransactions(c *gin.Context) {
    claims, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, banking.ErrorResponse{
            Error: "No authentication claims found",
            Code:  http.StatusUnauthorized,
        })
        return
    }
    userClaims := claims.(*banking.JWTClaims)
    accountID := c.Param("account_id")

    pageStr := c.DefaultQuery("page", "1")
    limitStr := c.DefaultQuery("limit", "50")

    page, _ := strconv.Atoi(pageStr)
    limit, _ := strconv.Atoi(limitStr)

    transactions, err := h.bankingUsecase.GetTransactions(userClaims, accountID, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, banking.ErrorResponse{
            Error:   "Failed to get transactions",
            Message: err.Error(),
            Code:    http.StatusInternalServerError,
        })
        return
    }

    c.JSON(http.StatusOK, transactions)
}

// GetBalances получает баланс счета
func (h *BankingHandler) GetBalances(c *gin.Context) {
    claims, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, banking.ErrorResponse{
            Error: "No authentication claims found",
            Code:  http.StatusUnauthorized,
        })
        return
    }
    userClaims := claims.(*banking.JWTClaims)
    accountID := c.Param("account_id")

    balances, err := h.bankingUsecase.GetBalances(userClaims, accountID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, banking.ErrorResponse{
            Error:   "Failed to get balances",
            Message: err.Error(),
            Code:    http.StatusInternalServerError,
        })
        return
    }

    c.JSON(http.StatusOK, balances)
}

// GetFinancialInsights получает финансовые инсайты клиента
func (h *BankingHandler) GetFinancialInsights(c *gin.Context) {
    claims, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, banking.ErrorResponse{
            Error: "No authentication claims found",
            Code:  http.StatusUnauthorized,
        })
        return
    }
    userClaims := claims.(*banking.JWTClaims)

    insights, err := h.bankingUsecase.GetFinancialInsights(userClaims)
    if err != nil {
        c.JSON(http.StatusInternalServerError, banking.ErrorResponse{
            Error:   "Failed to get financial insights",
            Message: err.Error(),
            Code:    http.StatusInternalServerError,
        })
        return
    }

    c.JSON(http.StatusOK, map[string]interface{}{
        "insights": insights,
        "message":  "Financial insights generated successfully",
    })
}
