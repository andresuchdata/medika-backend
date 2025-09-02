package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"medika-backend/internal/infrastructure/config"
	"medika-backend/internal/infrastructure/database"
	"medika-backend/internal/infrastructure/observability"
	"medika-backend/internal/infrastructure/redis"
	"medika-backend/internal/infrastructure/server"
	"medika-backend/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.New()
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(ctx, "Failed to load configuration", "error", err)
	}

	// Initialize observability (tracing, metrics)
	cleanup := observability.Initialize(cfg.Observability)
	defer cleanup()

	log.Info(ctx, "🚀 Starting Medika API Server")

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal(ctx, "Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal(ctx, "Failed to run migrations", "error", err)
	}

	// Initialize Redis
	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		log.Fatal(ctx, "Failed to connect to Redis", "error", err)
	}
	defer rdb.Close()

	// Initialize and start server
	srv := server.New(cfg.Server, db, rdb, log)
	
	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		log.Info(ctx, "Shutting down server...")
		
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error(ctx, "Server shutdown error", "error", err)
		}
		
		cancel()
	}()

	// Start server
	if err := srv.Start(ctx); err != nil {
		log.Fatal(ctx, "Server failed to start", "error", err)
	}

	log.Info(ctx, "Server stopped gracefully")
}
