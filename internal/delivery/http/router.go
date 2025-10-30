package http

import (
	"github.com/Arlandaren/easyfund/internal/config"
	"github.com/Arlandaren/easyfund/internal/usecase"
	"github.com/Arlandaren/easyfund/pkg/banking"
	"github.com/Arlandaren/easyfund/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(config *config.Config) *gin.Engine {
	// Инициализируем клиенты и сервисы
	vbankClient := banking.NewVBankClient(
		config.VBankBaseURL,
		config.VBankClientID,
		config.VBankClientSecret,
	)

	// Инициализируем use cases
	authUsecase := usecase.NewAuthUsecase(config, vbankClient)
	bankingUsecase := usecase.NewBankingUsecase(vbankClient)

	// Инициализируем handlers
	authHandler := NewAuthHandler(authUsecase)
	bankingHandler := NewBankingHandler(bankingUsecase, *vbankClient)
	healthHandler := NewHealthHandler()

	// Настраиваем роутер
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware(config))
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Health check
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API v1
	v1 := r.Group("/api/v1")

	// Auth routes (публичные)
	auth := v1.Group("/auth")
	{
		auth.GET("/random-demo-client", authHandler.GetRandomDemoClient)
		auth.POST("/login-demo-client", authHandler.LoginDemoClient)
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	// Banking routes (защищенные)
	banking := v1.Group("/banking")
	banking.Use(middleware.AuthMiddleware(authUsecase))
	{
		banking.GET("/accounts", bankingHandler.GetAccounts)
		banking.GET("/accounts/:account_id/transactions", bankingHandler.GetTransactions)
		banking.GET("/accounts/:account_id/balances", bankingHandler.GetBalances)
		banking.GET("/insights", bankingHandler.GetFinancialInsights)
		banking.POST("/create-consent", bankingHandler.CreateConsent)

	}

	return r
}
