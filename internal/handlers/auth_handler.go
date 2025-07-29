package handlers

import (
	"net/http"

	"github.com/nshmdayo/nft-platform-sample/internal/dto"
	"github.com/nshmdayo/nft-platform-sample/internal/errors"
	"github.com/nshmdayo/nft-platform-sample/internal/middleware"
	"github.com/nshmdayo/nft-platform-sample/internal/service"
	"github.com/nshmdayo/nft-platform-sample/internal/validation"
	"github.com/nshmdayo/nft-platform-sample/pkg/logger"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	BaseHandler
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := h.ValidateMethod(r, http.MethodPost); err != nil {
		h.SendError(w, err)
		return
	}

	var req dto.RegisterRequest
	if err := h.DecodeJSON(r, &req); err != nil {
		h.SendError(w, err)
		return
	}

	// Validate request
	validator := validation.NewValidator()
	validator.Required("email", req.Email).Email("email", req.Email)
	validator.Required("password", req.Password).MinLength("password", req.Password, 8)
	validator.Required("name", req.Name).MinLength("name", req.Name, 2).MaxLength("name", req.Name, 100)

	if err := validator.Validate(); err != nil {
		h.SendError(w, errors.Validation("Validation failed", err))
		return
	}

	// Convert DTO to service request
	serviceReq := &service.RegisterRequest{
		Email:       req.Email,
		Password:    req.Password,
		Name:        req.Name,
		Institution: req.Institution,
		WalletAddr:  "", // Will be set later when blockchain integration is added
	}

	response, err := h.authService.Register(serviceReq)
	if err != nil {
		logger.Error("Failed to register user", "error", err, "email", req.Email)
		h.SendError(w, errors.Conflict("User registration failed"))
		return
	}

	// Convert to DTO response
	authResponse := &dto.AuthResponse{
		Token: response.Token,
		User: &dto.UserInfo{
			ID:          response.User.ID,
			Email:       response.User.Email,
			Name:        response.User.Name,
			Role:        response.User.Role,
			Institution: response.User.Institution,
			CreatedAt:   response.User.CreatedAt,
		},
	}

	logger.Info("User registered successfully", "user_id", response.User.ID, "email", req.Email)
	h.SendResponse(w, http.StatusCreated, authResponse)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := h.ValidateMethod(r, http.MethodPost); err != nil {
		h.SendError(w, err)
		return
	}

	var req dto.LoginRequest
	if err := h.DecodeJSON(r, &req); err != nil {
		h.SendError(w, err)
		return
	}

	// Validate request
	validator := validation.NewValidator()
	validator.Required("email", req.Email).Email("email", req.Email)
	validator.Required("password", req.Password)

	if err := validator.Validate(); err != nil {
		h.SendError(w, errors.Validation("Validation failed", err))
		return
	}

	// Convert DTO to service request
	serviceReq := &service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := h.authService.Login(serviceReq)
	if err != nil {
		logger.Error("Failed to login user", "error", err, "email", req.Email)
		h.SendError(w, errors.Unauthorized("Invalid credentials"))
		return
	}

	// Convert to DTO response
	authResponse := &dto.AuthResponse{
		Token: response.Token,
		User: &dto.UserInfo{
			ID:          response.User.ID,
			Email:       response.User.Email,
			Name:        response.User.Name,
			Role:        response.User.Role,
			Institution: response.User.Institution,
			CreatedAt:   response.User.CreatedAt,
		},
	}

	logger.Info("User logged in successfully", "user_id", response.User.ID, "email", req.Email)
	h.SendResponse(w, http.StatusOK, authResponse)
}

// GetProfile handles getting user profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if err := h.ValidateMethod(r, http.MethodGet); err != nil {
		h.SendError(w, err)
		return
	}

	userID, err := GetUserIDFromContext(r)
	if err != nil {
		h.SendError(w, errors.Unauthorized("User not authenticated"))
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		logger.Error("Failed to get user profile", "error", err, "user_id", userID)
		h.SendError(w, errors.NotFound("User"))
		return
	}

	userInfo := &dto.UserInfo{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Role:        user.Role,
		Institution: user.Institution,
		CreatedAt:   user.CreatedAt,
	}

	h.SendResponse(w, http.StatusOK, userInfo)
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(r *http.Request) (uint, error) {
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		return 0, errors.Unauthorized("User not authenticated")
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return 0, errors.Unauthorized("Invalid user ID format")
	}

	return userIDUint, nil
}
