package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/Arlandaren/easyfund/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	userHandler *UserHandler,
	loanHandler *LoanHandler,
	accountHandler *UserBankAccountHandler,
	transactionHandler *TransactionHandler,
	applicationHandler *CreditApplicationHandler,
	authHandler *AuthHandler,
	jwtSecret string,
) {
	v1 := router.Group("/api/v1")

	// Публичные endpoints
	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)
	v1.POST("/auth/refresh", authHandler.RefreshToken)
	v1.GET("/users/random", userHandler.GetRandomUser)

	// Защищённые endpoints
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))

	// === 1. Сначала все маршруты с 3+ сегментами (чтобы не перехватывались :id) ===
	protected.GET("/auth/me", authHandler.GetMe)

	// Счета и баланс
	protected.GET("/users/:id/balance", accountHandler.GetTotalBalance)
	protected.GET("/users/:id/accounts", accountHandler.GetUserAccounts)

	// Кредиты пользователя
	protected.GET("/users/:id/loans", loanHandler.ListUserLoans)
	protected.GET("/users/:id/debt", loanHandler.GetTotalDebt)

	// Транзакции
	protected.GET("/users/:id/transactions", transactionHandler.GetUserTransactionHistory)
	protected.GET("/users/:id/banks/:bank_id/transactions", transactionHandler.GetBankTransactionHistory)

	// Заявки на кредит
	protected.GET("/users/:id/applications", applicationHandler.GetUserApplications)

	// === 2. Потом — маршруты с двумя сегментами (включая /users/:id) ===
	protected.GET("/users/:id", userHandler.GetUser)
	protected.PUT("/users/:id", userHandler.UpdateUser)
	protected.DELETE("/users/:id", userHandler.DeleteUser)

	// === 3. Остальные маршруты (без :id или с другим паттерном) ===
	protected.POST("/users", userHandler.CreateUser)

	// Глобальные кредиты и заявки
	protected.POST("/loans", loanHandler.CreateLoan)
	protected.GET("/loans/:id", loanHandler.GetLoanDetail)
	protected.POST("/loans/:id/payment", loanHandler.MakePayment)

	protected.POST("/applications", applicationHandler.SubmitApplication)
	protected.POST("/applications/:id/approve", applicationHandler.ApproveApplication)
	protected.POST("/applications/:id/reject", applicationHandler.RejectApplication)
}