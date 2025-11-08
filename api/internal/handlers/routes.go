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

    // Auth
    protected.GET("/auth/me", authHandler.GetMe)

    // Пользователи
    protected.GET("/users/:id", userHandler.GetUser)
    protected.POST("/users", userHandler.CreateUser)
    protected.PUT("/users/:id", userHandler.UpdateUser)
    protected.DELETE("/users/:id", userHandler.DeleteUser)

    // Счета и баланс (унифицируем :id)
    protected.GET("/users/:id/balance", accountHandler.GetTotalBalance)
    protected.GET("/users/:id/accounts", accountHandler.GetUserAccounts)

    // Кредиты (унифицируем :id для пользователя)
    protected.POST("/loans", loanHandler.CreateLoan)
    protected.GET("/loans/:id", loanHandler.GetLoanDetail)
    protected.GET("/users/:id/loans", loanHandler.ListUserLoans)
    protected.GET("/users/:id/debt", loanHandler.GetTotalDebt)
    protected.POST("/loans/:id/payment", loanHandler.MakePayment)

    // Транзакции (унифицируем :id для пользователя; bank_id остаётся)
    protected.GET("/users/:id/transactions", transactionHandler.GetUserTransactionHistory)
    protected.GET("/users/:id/banks/:bank_id/transactions", transactionHandler.GetBankTransactionHistory)

    // Заявки на кредиты (унифицируем :id для пользователя)
    protected.POST("/applications", applicationHandler.SubmitApplication)
    protected.GET("/users/:id/applications", applicationHandler.GetUserApplications)
    protected.POST("/applications/:id/approve", applicationHandler.ApproveApplication)
    protected.POST("/applications/:id/reject", applicationHandler.RejectApplication)
}
