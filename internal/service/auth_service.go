package service

import (
	"errors"
	"time"

	"github.com/nshmdayo/nft-platform-sample/internal/config"
	"github.com/nshmdayo/nft-platform-sample/internal/models"
	"github.com/nshmdayo/nft-platform-sample/internal/repository"
	"github.com/nshmdayo/nft-platform-sample/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Name        string `json:"name" binding:"required"`
	Institution string `json:"institution"`
	WalletAddr  string `json:"wallet_address"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func NewAuthService(userRepo *repository.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:       req.Email,
		Password:    hashedPassword,
		Name:        req.Name,
		Institution: req.Institution,
		WalletAddr:  req.WalletAddr,
		Role:        "researcher",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT token
	expiresIn, _ := time.ParseDuration(s.config.JWT.ExpiresIn)
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, s.config.JWT.Secret, expiresIn)
	if err != nil {
		return nil, err
	}

	// Clear password from response
	user.Password = ""

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check password
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	expiresIn, _ := time.ParseDuration(s.config.JWT.ExpiresIn)
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, s.config.JWT.Secret, expiresIn)
	if err != nil {
		return nil, err
	}

	// Clear password from response
	user.Password = ""

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}
