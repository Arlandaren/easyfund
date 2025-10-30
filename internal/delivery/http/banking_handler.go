package http

import (
    "net/http"
    "strconv"

    "github.com/Arlandaren/easyfund/internal/usecase"
    "github.com/Arlandaren/easyfund/pkg/banking"
    "github.com/gin-gonic/gin"
)

type BankingHandler struct {
    bankingUsecase *usecase.BankingUsecase
    vbankClient    banking.VBankAPI
}

func NewBankingHandler(bankingUsecase *usecase.BankingUsecase, vbankClient banking.VBankAPI) *BankingHandler {
    return &BankingHandler{
        bankingUsecase: bankingUsecase,
        vbankClient:    vbankClient,
    }
}

func (h *BankingHandler) extractClaims(c *gin.Context) (*banking.JWTClaims, bool) {
    claims, exists := c.Get("claims")
    if !exists || claims == nil {
        c.JSON(http.StatusUnauthorized, banking.ErrorResponse{Error: "No authentication claims found", Code: http.StatusUnauthorized})
        return nil, false
    }
    userClaims, ok := claims.(*banking.JWTClaims)
    if !ok || userClaims == nil {
        c.JSON(http.StatusUnauthorized, banking.ErrorResponse{Error: "Invalid authentication claims", Code: http.StatusUnauthorized})
        return nil, false
    }
    return userClaims, true
}

func (h *BankingHandler) resolveClientID(userClaims *banking.JWTClaims) (string, bool) {
    // Используем client_id из токена; если нет — fallback на person_id
    if userClaims.ClientID != "" {
        return userClaims.ClientID, true
    }
    if userClaims.PersonID != "" {
        return userClaims.PersonID, true
    }
    return "", false
}

func (h *BankingHandler) getConsentID(c *gin.Context, clientID string) (string, bool) {
    if h.bankingUsecase != nil {
        if consentID, err := h.bankingUsecase.GetConsentID(c.Request.Context(), clientID); err == nil && consentID != "" {
            return consentID, true
        }
    }
    if hdr := c.GetHeader("x-consent-id"); hdr != "" {
        return hdr, true
    }
    c.JSON(http.StatusBadRequest, banking.ErrorResponse{Error: "Consent not found for client_id; provide x-consent-id or save consent", Code: http.StatusBadRequest})
    return "", false
}

func (h *BankingHandler) GetAccounts(c *gin.Context) {
    userClaims, ok := h.extractClaims(c); if !ok { return }
    clientID, ok := h.resolveClientID(userClaims); if !ok {
        c.JSON(http.StatusBadRequest, banking.ErrorResponse{Error: "client_id not found in JWT claims", Code: http.StatusBadRequest}); return
    }
    consentID, ok := h.getConsentID(c, clientID); if !ok { return }

    data, err := h.vbankClient.GetAccountsWithConsent(clientID, consentID)
    if err != nil {
        c.JSON(http.StatusBadGateway, banking.ErrorResponse{Error: "Failed to get accounts", Message: err.Error(), Code: http.StatusBadGateway})
        return
    }
    c.JSON(http.StatusOK, data)
}

func (h *BankingHandler) GetTransactions(c *gin.Context) {
    userClaims, ok := h.extractClaims(c); if !ok { return }
    clientID, ok := h.resolveClientID(userClaims); if !ok {
        c.JSON(http.StatusBadRequest, banking.ErrorResponse{Error: "client_id not found in JWT claims", Code: http.StatusBadRequest}); return
    }
    accountID := c.Param("account_id")
    if accountID == "" {
        c.JSON(http.StatusBadRequest, banking.ErrorResponse{Error: "account_id is required", Code: http.StatusBadRequest}); return
    }
    consentID, ok := h.getConsentID(c, clientID); if !ok { return }
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

    data, err := h.vbankClient.GetTransactionsWithConsent(accountID, clientID, consentID, page, limit)
    if err != nil {
        c.JSON(http.StatusBadGateway, banking.ErrorResponse{Error: "Failed to get transactions", Message: err.Error(), Code: http.StatusBadGateway})
        return
    }
    c.JSON(http.StatusOK, data)
}

func (h *BankingHandler) GetBalances(c *gin.Context) {
    userClaims, ok := h.extractClaims(c); if !ok { return }
    clientID, ok := h.resolveClientID(userClaims); if !ok {
        c.JSON(http.StatusBadRequest, banking.ErrorResponse{Error: "client_id not found in JWT claims", Code: http.StatusBadRequest}); return
    }
    accountID := c.Param("account_id")
    if accountID == "" {
        c.JSON(http.StatusBadRequest, banking.ErrorResponse{Error: "account_id is required", Code: http.StatusBadRequest}); return
    }
    consentID, ok := h.getConsentID(c, clientID); if !ok { return }

    data, err := h.vbankClient.GetBalancesWithConsent(accountID, clientID, consentID)
    if err != nil {
        c.JSON(http.StatusBadGateway, banking.ErrorResponse{Error: "Failed to get balances", Message: err.Error(), Code: http.StatusBadGateway})
        return
    }
    c.JSON(http.StatusOK, data)
}

func (h *BankingHandler) GetFinancialInsights(c *gin.Context) {
    userClaims, ok := h.extractClaims(c); if !ok { return }
    insights, err := h.bankingUsecase.GetFinancialInsights(userClaims)
    if err != nil {
        c.JSON(http.StatusInternalServerError, banking.ErrorResponse{Error: "Failed to get financial insights", Message: err.Error(), Code: http.StatusInternalServerError})
        return
    }
    c.JSON(http.StatusOK, gin.H{"insights": insights, "message": "Financial insights generated successfully"})
}

type CreateConsentRequest struct {
    ClientID           string `json:"client_id" binding:"required"`
    RequestingBank     string `json:"requesting_bank" binding:"required"`
    RequestingBankName string `json:"requesting_bank_name" binding:"required"`
}

func (h *BankingHandler) CreateConsent(c *gin.Context) {
    if _, ok := h.extractClaims(c); !ok { return }
    var req CreateConsentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    resp, err := h.vbankClient.CreateConsent(req.ClientID, req.RequestingBank, req.RequestingBankName)
    if err != nil {
        c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
        return
    }
    // при наличии репозитория: _ = h.bankingUsecase.SaveConsentID(c.Request.Context(), req.ClientID, resp.ID)
    c.JSON(http.StatusOK, resp)
}
