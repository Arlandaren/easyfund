package main

import (
	"github.com/Arlandaren/easyfund/internal/config"
	"github.com/Arlandaren/easyfund/internal/delivery/http"
	"log"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Настраиваем роутер
	router := http.SetupRouter(cfg)

	// Запускаем сервер
	log.Printf("Starting EasyFund API server on port %s", cfg.Port)
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("VBank Base URL: %s", cfg.VBankBaseURL)
	
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
