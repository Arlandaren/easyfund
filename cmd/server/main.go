package main

import (
    "context"
    "log"

    _ "github.com/lib/pq"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/sql"

    "github.com/Arlandaren/easyfund/ent"
    "github.com/Arlandaren/easyfund/internal/config"
    httpDelivery "github.com/Arlandaren/easyfund/internal/delivery/http"
    "github.com/Arlandaren/easyfund/pkg/logger"
    "go.uber.org/zap"
)

func main() {
    // Load config
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Init logger
    lg, err := logger.New()
    if err != nil {
        log.Fatalf("Failed to init logger: %v", err)
    }
    defer func() {
        if err := lg.Sync(); err != nil {
            log.Printf("Error syncing logger: %v", err)
        }
    }()

    // Connect to DB
    drv, err := sql.Open(dialect.Postgres, cfg.PostgresURL)
    if err != nil {
        lg.Fatal("Failed to connect to database", zap.Error(err))
    }

    client := ent.NewClient(ent.Driver(drv))
    defer func() {
        if err := client.Close(); err != nil {
            lg.Error("Error closing database client", zap.Error(err))
        }
    }()

    // Run migrations
    if err := client.Schema.Create(context.Background()); err != nil {
        lg.Fatal("Failed to run migrations", zap.Error(err))
    }

    lg.Info("Database migrations completed")

    // Start HTTP server
    router := httpDelivery.NewRouter()
    lg.Info("Starting HTTP server on " + cfg.HTTPAddress)

    if err := router.Run(cfg.HTTPAddress); err != nil {
        lg.Fatal("Failed to start HTTP server", zap.Error(err))
    }
}
