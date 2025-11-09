package application

import (
	"errors"
	"fmt"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/infrastructure"
	"github.com/Ivan-Lapin/Auth-service/service/pkg/jwt"
	"go.uber.org/zap"
)

type AuthActions interface {
	Login(email, password string) (string, error)
}

type AuthService struct {
	userRepo infrastructure.UserRepository
	jwt      *jwt.JWT
	logger   *zap.Logger
}

func NewAuthService(userRepo infrastructure.UserRepository, jwt *jwt.JWT, logger *zap.Logger) *AuthService {
	return &AuthService{userRepo: userRepo, jwt: jwt, logger: logger}
}

func (as *AuthService) Login(email, password string) (string, error) {
	as.logger.Info("Loging the user...", zap.String("Email", email), zap.String("Password", "****"))

	user, err := as.userRepo.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundUser) {
			as.logger.Info("User not founded", zap.String("Email", email))
			return "", fmt.Errorf("user not founded: %w", err)
		}

		as.logger.Error("Login failed", zap.Error(err))
		return "", fmt.Errorf("login failed: %w", err)
	}

	if !user.Verify(password) {
		as.logger.Info("Invalid password", zap.String("Email", email))
		return "", apperrors.ErrInvalidPassword
	}

	token, err := as.jwt.Generate(user.ID)
	if err != nil {
		as.logger.Error("Failed to generate JWT token", zap.Error(err))
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	as.logger.Info("Successful user loging", zap.String("Email", email))

	return token, nil

}
