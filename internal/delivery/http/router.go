package http

import (
    "github.com/Arlandaren/easyfund/internal/config"
    "github.com/Arlandaren/easyfund/internal/usecase"
    "github.com/Arlandaren/easyfund/pkg/banking"
    "github.com/Arlandaren/easyfund/pkg/middleware"

    "github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
    // Реальный клиент
    var vbankClient banking.VBankAPI = banking.NewVBankClient(
        cfg.VBankBaseURL,
        cfg.VBankClientID,
        cfg.VBankClientSecret,
    )

    authUsecase := usecase.NewAuthUsecase(cfg, vbankClient.(*banking.VBankClient))
    bankingUsecase := usecase.NewBankingUsecase(vbankClient)

    authHandler := NewAuthHandler(authUsecase)
    bankingHandler := NewBankingHandler(bankingUsecase, vbankClient)
    healthHandler := NewHealthHandler()

    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    r := gin.Default()

    r.Use(middleware.CORSMiddleware(cfg))
    r.Use(gin.Recovery())
    r.Use(gin.Logger())

    r.GET("/health", healthHandler.HealthCheck)
    r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })

    v1 := r.Group("/api/v1")

    auth := v1.Group("/auth")
    {
        auth.GET("/random-demo-client", authHandler.GetRandomDemoClient)
        auth.POST("/login-demo-client", authHandler.LoginDemoClient)
        auth.POST("/login", authHandler.Login)
        auth.POST("/register", authHandler.Register)
    }

    bankingGroup := v1.Group("/banking")
    bankingGroup.Use(middleware.AuthMiddleware(authUsecase))
    {
        bankingGroup.GET("/accounts", bankingHandler.GetAccounts)
        bankingGroup.GET("/accounts/:account_id/transactions", bankingHandler.GetTransactions)
        bankingGroup.GET("/accounts/:account_id/balances", bankingHandler.GetBalances)
        bankingGroup.GET("/insights", bankingHandler.GetFinancialInsights)
        bankingGroup.POST("/create-consent", bankingHandler.CreateConsent)
    }

    return r
}
