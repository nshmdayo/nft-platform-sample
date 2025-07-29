package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nshmdayo/nft-platform-sample/internal/config"
	"github.com/nshmdayo/nft-platform-sample/internal/handlers"
	"github.com/nshmdayo/nft-platform-sample/internal/models"
	"github.com/nshmdayo/nft-platform-sample/internal/repository"
	"github.com/nshmdayo/nft-platform-sample/internal/router"
	"github.com/nshmdayo/nft-platform-sample/internal/service"
	"github.com/nshmdayo/nft-platform-sample/pkg/logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() http.Handler {
	// Initialize logger for tests
	logger.Init()
	// Use in-memory SQLite for testing
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// Run migrations
	db.AutoMigrate(&models.User{}, &models.Paper{}, &models.Review{}, &models.NFTMetadata{})

	// Initialize test configuration
	cfg := &config.Config{
		App: config.AppConfig{
			Environment: "test",
		},
		JWT: config.JWTConfig{
			Secret:    "test-secret",
			ExpiresIn: "24h",
		},
	}

	// Initialize repositories, services, and handlers
	userRepo := repository.NewUserRepository(db)
	paperRepo := repository.NewPaperRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	paperService := service.NewPaperService(paperRepo)
	reviewService := service.NewReviewService(reviewRepo, paperRepo)

	authHandler := handlers.NewAuthHandler(authService)
	paperHandler := handlers.NewPaperHandler(paperService)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	r := router.NewRouter(cfg, authHandler, paperHandler, reviewHandler)
	return r.SetupRoutes()
}

func TestHealthEndpoint(t *testing.T) {
	handler := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "ok", data["status"])
}

func TestRegisterEndpoint(t *testing.T) {
	handler := setupTestRouter()

	user := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}

	jsonData, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])

	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "token")
	assert.Contains(t, data, "user")
}
