package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Arlandaren/easyfund/internal/logger"
	"github.com/Arlandaren/easyfund/internal/middleware"
	"github.com/Arlandaren/easyfund/internal/models"
	"github.com/Arlandaren/easyfund/internal/services"
)

type LoanHandler struct {
	loanService    services.LoanService
	accountService services.UserBankAccountService
}

func NewLoanHandler(loanService services.LoanService, accountService services.UserBankAccountService) *LoanHandler {
	return &LoanHandler{
		loanService:    loanService,
		accountService: accountService,
	}
}

type CreateLoanRequest struct {
	OriginalAmount string             `json:"original_amount" binding:"required"`
	InterestRate   string             `json:"interest_rate" binding:"required"`
	Purpose        string             `json:"purpose"`
	Splits         []map[int16]string `json:"splits" binding:"required"`
}

// POST /api/v1/loans (защищенный)
func (h *LoanHandler) CreateLoan(c *gin.Context) {
	// Получаем user_id из контекста
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateLoanRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	loan := &models.Loan{
		UserID:        userID,
		OriginalAmount: req.OriginalAmount,
		TakenAt:       time.Now(),
		InterestRate:  req.InterestRate,
		Status:        "ACTIVE",
		Purpose:       req.Purpose,
		CreatedAt:     time.Now(),
	}

	detail, err := h.loanService.CreateLoan(c.Request.Context(), loan, req.Splits)
	if err != nil {
		logger.Log.Errorf("Failed to create loan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create loan"})
		return
	}

	logger.Log.Infof("User %s created loan with ID %d", userID, detail.Loan.LoanID)
	c.JSON(http.StatusCreated, detail)
}

// GET /api/v1/loans/:id (защищенный)
func (h *LoanHandler) GetLoanDetail(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	loanIDStr := c.Param("id")
	loanID, err := strconv.ParseInt(loanIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	detail, err := h.loanService.GetLoanDetail(c.Request.Context(), loanID)
	if err != nil {
		logger.Log.Errorf("Failed to get loan detail: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get loan detail"})
		return
	}

	// Проверяем, что это кредит пользователя
	if detail.Loan.UserID != userID {
		logger.Log.Warnf("User %s tried to access loan %d of user %s", userID, loanID, detail.Loan.UserID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, detail)
}

// GET /api/v1/users/:user_id/loans (защищенный)
func (h *LoanHandler) ListUserLoans(c *gin.Context) {
	requestingUserID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Проверяем права доступа
	if requestingUserID != userID {
		logger.Log.Warnf("User %s tried to list loans of user %s", requestingUserID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	loans, err := h.loanService.ListUserLoans(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Errorf("Failed to list user loans: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list user loans"})
		return
	}

	c.JSON(http.StatusOK, loans)
}

// GET /api/v1/users/:user_id/debt (защищенный)
func (h *LoanHandler) GetTotalDebt(c *gin.Context) {
	requestingUserID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Проверяем права доступа
	if requestingUserID != userID {
		logger.Log.Warnf("User %s tried to get debt of user %s", requestingUserID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	totalDebt, err := h.loanService.GetTotalDebt(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Errorf("Failed to get total debt: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total debt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_debt": totalDebt})
}

type MakePaymentRequest struct {
	LoanID      int64                     `json:"loan_id" binding:"required"`
	TotalAmount string                   `json:"total_amount" binding:"required"`
	Comment     string                   `json:"comment"`
	Allocations []map[string]interface{} `json:"allocations" binding:"required"`
}

// POST /api/v1/loans/:id/payment (защищенный)
func (h *LoanHandler) MakePayment(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	loanIDStr := c.Param("id")
	loanID, err := strconv.ParseInt(loanIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	var req MakePaymentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Проверяем, что это кредит пользователя
	loan, err := h.loanService.GetLoanDetail(c.Request.Context(), loanID)
	if err != nil {
		logger.Log.Errorf("Failed to get loan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get loan"})
		return
	}

	if loan.Loan.UserID != userID {
		logger.Log.Warnf("User %s tried to pay loan %d of user %s", userID, loanID, loan.Loan.UserID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	payment := &models.LoanPayment{
		LoanID:      loanID,
		UserID:      userID,
		PaidAt:      time.Now(),
		TotalAmount: req.TotalAmount,
		Comment:     req.Comment,
	}

	// Парсим allocations
	var allocations []models.PaymentAllocation
	for _, alloc := range req.Allocations {
		splitID := int64(alloc["split_id"].(float64))
		principalPaid := alloc["principal_paid"].(string)
		interestPaid := alloc["interest_paid"].(string)

		allocations = append(allocations, models.PaymentAllocation{
			SplitID:       splitID,
			PrincipalPaid: principalPaid,
			InterestPaid:  interestPaid,
		})
	}

	err = h.loanService.MakePayment(c.Request.Context(), payment, allocations)
	if err != nil {
		logger.Log.Errorf("Failed to make payment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make payment"})
		return
	}

	logger.Log.Infof("User %s made payment for loan %d", userID, loanID)
	c.JSON(http.StatusCreated, gin.H{"message": "Payment recorded successfully"})
}
