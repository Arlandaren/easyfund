package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Arlandaren/easyfund/internal/logger"
	"github.com/Arlandaren/easyfund/internal/middleware"
	"github.com/Arlandaren/easyfund/internal/models"
	"github.com/Arlandaren/easyfund/internal/services"
)

type CreditApplicationHandler struct {
	service services.CreditApplicationService
}

func NewCreditApplicationHandler(service services.CreditApplicationService) *CreditApplicationHandler {
	return &CreditApplicationHandler{service: service}
}

type SubmitApplicationRequest struct {
	BankID         int16  `json:"bank_id" binding:"required"`
	TypeCode       string `json:"type_code" binding:"required"`
	RequestedAmount string `json:"requested_amount" binding:"required"`
}

// POST /api/v1/applications (защищенный)
func (h *CreditApplicationHandler) SubmitApplication(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req SubmitApplicationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	app := &models.CreditApplication{
		UserID:         userID,
		BankID:         req.BankID,
		TypeCode:       req.TypeCode,
		StatusCode:     "PENDING",
		RequestedAmount: req.RequestedAmount,
		SubmittedAt:    time.Now(),
		UpdatedAt:      time.Now(),
	}

	appID, err := h.service.SubmitApplication(c.Request.Context(), app)
	if err != nil {
		logger.Log.Errorf("Failed to submit application: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit application"})
		return
	}

	logger.Log.Infof("User %d submitted credit application %d", userID, appID)
	c.JSON(http.StatusCreated, gin.H{"application_id": appID, "status": "PENDING"})
}

// GET /api/v1/users/:id/applications (защищенный)
func (h *CreditApplicationHandler) GetUserApplications(c *gin.Context) {
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
		logger.Log.Warnf("User %d tried to get applications of user %d", requestingUserID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	apps, err := h.service.GetApplications(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Errorf("Failed to get applications: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get applications"})
		return
	}

	c.JSON(http.StatusOK, apps)
}

type ApproveApplicationRequest struct {
	Splits []map[int16]string `json:"splits" binding:"required"`
}

// POST /api/v1/applications/:id/approve (администратор/защищенный)
func (h *CreditApplicationHandler) ApproveApplication(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	appIDStr := c.Param("id")
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	var req ApproveApplicationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.service.ApproveApplication(c.Request.Context(), appID, req.Splits)
	if err != nil {
		logger.Log.Errorf("Failed to approve application: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve application"})
		return
	}

	logger.Log.Infof("User %d (admin) approved application %d", userID, appID)
	c.JSON(http.StatusOK, gin.H{"message": "Application approved"})
}

// POST /api/v1/applications/:id/reject (администратор/защищенный)
func (h *CreditApplicationHandler) RejectApplication(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	appIDStr := c.Param("id")
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	err = h.service.RejectApplication(c.Request.Context(), appID)
	if err != nil {
		logger.Log.Errorf("Failed to reject application: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject application"})
		return
	}

	logger.Log.Infof("User %d (admin) rejected application %d", userID, appID)
	c.JSON(http.StatusOK, gin.H{"message": "Application rejected"})
}
