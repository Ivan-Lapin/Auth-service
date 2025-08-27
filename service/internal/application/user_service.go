package application

import (
	"errors"
	"fmt"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/domain"
	"github.com/Ivan-Lapin/Auth-service/service/internal/infrastructure"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo infrastructure.UserRepository
	logger   *zap.Logger
}

func NewUserService(ur infrastructure.UserRepository, loggerr *zap.Logger) *UserService {
	return &UserService{
		userRepo: ur, logger: loggerr,
	}
}

func (us *UserService) Register(name, email, password string) (*domain.User, error) {
	us.logger.Info("Registering User")

	id := uuid.New()
	user := domain.User{
		ID:       id,
		Username: name,
		Email:    email,
	}

	if err := user.SetPassword(password); err != nil {
		us.logger.Error("Failed to set password for user", zap.String("Username", user.Username), zap.String("Email", user.Email), zap.Error(err))
		return nil, fmt.Errorf("failed to set password for user: %w", err)
	}

	err := us.userRepo.Save(&user)
	if errors.Is(err, apperrors.ErrUserAlreadyExist) {
		us.logger.Info("User already exist", zap.String("Email", user.Email))
		return nil, fmt.Errorf("user already exist: %w", err)
	}

	if err != nil {
		us.logger.Error("Failed to save user into DB", zap.String("Username", user.Username), zap.String("Email", user.Email), zap.Error(err))
		return nil, fmt.Errorf("failed to save user into DB: %w", err)
	}

	us.logger.Info("Successful user addition", zap.String("Username", user.Username), zap.String("Email", user.Email))

	return &user, nil

}
