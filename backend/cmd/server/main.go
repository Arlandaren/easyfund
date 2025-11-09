package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/Arlandaren/easyfund/internal/config"
	"github.com/Arlandaren/easyfund/internal/handlers"
	"github.com/Arlandaren/easyfund/internal/logger"
	"github.com/Arlandaren/easyfund/internal/middleware"
	"github.com/Arlandaren/easyfund/internal/repos"
	"github.com/Arlandaren/easyfund/internal/services"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Логгер
	logger.Init(cfg.Logger.Level, cfg.Logger.File)
	defer logger.Log.Sync()

	// БД
	db, err := config.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Gin
	router := gin.Default()


	// CORSMiddleware: используем без аргументов (см. internal/middleware/cors.go ниже)
	router.Use(middleware.CORSMiddleware())

	// Репозитории
	userRepo := repos.NewUserRepository(db)
	loanRepo := repos.NewLoanRepository(db)
	accountRepo := repos.NewUserBankAccountRepository(db)
	paymentRepo := repos.NewLoanPaymentRepository(db)
	transactionRepo := repos.NewTransactionRepository(db)
	applicationRepo := repos.NewCreditApplicationRepository(db)

	// Сервисы
	userService := services.NewUserService(userRepo)                // содержит GetUserByEmail (исправлено)
	tokenService := services.NewTokenService(&cfg.JWT)
	loanService := services.NewLoanService(loanRepo, accountRepo, paymentRepo)
	accountService := services.NewUserBankAccountService(accountRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	applicationService := services.NewCreditApplicationService(applicationRepo, loanRepo)

	// Хэндлеры
	userHandler := handlers.NewUserHandler(userService)
	loanHandler := handlers.NewLoanHandler(loanService, accountService)
	accountHandler := handlers.NewUserBankAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	applicationHandler := handlers.NewCreditApplicationHandler(applicationService)
	authHandler := handlers.NewAuthHandler(userService, tokenService)

	// Роутинг
	handlers.RegisterRoutes(
		router,
		userHandler,
		loanHandler,
		accountHandler,
		transactionHandler,
		applicationHandler,
		authHandler,
		cfg.JWT.Secret,
	)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Log.Infof("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
