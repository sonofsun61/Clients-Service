package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/repository"
	"github.com/AI-Hackathon-2026/Clients-Service/pkg/hashutil"
	"github.com/AI-Hackathon-2026/Clients-Service/pkg/jwtutil"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) error
	Login(ctx context.Context, req model.LoginRequest) (model.TokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (model.TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	repo      repository.AuthRepository
	logger    *slog.Logger
	jwtSecret string
}

func NewAuthService(repo repository.AuthRepository, logger *slog.Logger, jwtSecret string) AuthService {
	return &authService{
		repo:      repo,
		logger:    logger,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, req model.RegisterRequest) error {
	passwordHash, err := hashutil.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return errors.New("internal registration error")
	}
	userData := model.UserData{
		UserID:       uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		StreakStart:  time.Now(),
		LastActivity: time.Now(),
	}
	if err := s.repo.Register(ctx, userData); err != nil {
		s.logger.Error("failed to register user in repo", "error", err, "email", req.Email)
		return fmt.Errorf("repository error: %w", err)
	}
	return nil
}

func (s *authService) Login(ctx context.Context, req model.LoginRequest) (model.TokenResponse, error) {
	userData, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("login failed: user not found", "email", req.Email)
		return model.TokenResponse{}, errors.New("invalid email or password")
	}
	if !hashutil.CheckPasswordHash(req.Password, userData.PasswordHash) {
		s.logger.Warn("login failed: wrong password", "email", req.Email)
		return model.TokenResponse{}, errors.New("invalid email or password")
	}
	accessToken, err := jwtutil.GenerateToken(userData.UserID, s.jwtSecret, 15*time.Minute)
	if err != nil {
		s.logger.Error("failed to generate access token", "error", err)
		return model.TokenResponse{}, errors.New("internal server error")
	}
	refreshToken, err := jwtutil.GenerateToken(userData.UserID, s.jwtSecret, 30*24*time.Hour)
	return model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (model.TokenResponse, error) {
	claims, err := jwtutil.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return model.TokenResponse{}, errors.New("invalid refresh token")
	}
	newAccess, _ := jwtutil.GenerateToken(claims.UserID, s.jwtSecret, 15*time.Minute)
	newRefresh, _ := jwtutil.GenerateToken(claims.UserID, s.jwtSecret, 30*24*time.Hour)
	return model.TokenResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return nil
}
