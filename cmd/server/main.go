package main

import (
	"log"
	"net/http"

	"github.com/nshmdayo/nft-platform-sample/internal/config"
	"github.com/nshmdayo/nft-platform-sample/internal/database"
	"github.com/nshmdayo/nft-platform-sample/internal/handlers"
	"github.com/nshmdayo/nft-platform-sample/internal/repository"
	"github.com/nshmdayo/nft-platform-sample/internal/router"
	"github.com/nshmdayo/nft-platform-sample/internal/service"
	"github.com/nshmdayo/nft-platform-sample/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()
	logger.Info("Starting NFT Platform API Server")

	// Load configuration
	cfg := config.LoadConfig()
	logger.Info("Configuration loaded", "port", cfg.Server.Port)

	// Initialize database
	err := database.Connect(cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	logger.Info("Database connected successfully")

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	logger.Info("Database migrations completed")

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)
	paperRepo := repository.NewPaperRepository(database.DB)
	reviewRepo := repository.NewReviewRepository(database.DB)
	logger.Info("Repositories initialized")

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	paperService := service.NewPaperService(paperRepo)
	reviewService := service.NewReviewService(reviewRepo, paperRepo)
	logger.Info("Services initialized")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	paperHandler := handlers.NewPaperHandler(paperService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	logger.Info("Handlers initialized")

	// Initialize router
	r := router.NewRouter(cfg, authHandler, paperHandler, reviewHandler)
	handler := r.SetupRoutes()
	logger.Info("Router setup completed")

	// Start server
	port := ":" + cfg.Server.Port
	logger.Info("Server starting", "port", port)

	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
