package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/domain"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type UserRepository interface {
	Save(user *domain.User) error
	FindUserByEmail(email string) (*domain.User, error)
}

type UserRepo struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewUserRepository(db *sqlx.DB, logger *zap.Logger) *UserRepo {
	return &UserRepo{db: db, logger: logger}
}

func (ur *UserRepo) Save(user *domain.User) error {

	var exist bool
	err := ur.db.Get(&exist, `SELECT EXISTS(SELECT 1 FROM public.users WHERE email = $1)`, user.Email)
	if err != nil {
		ur.logger.Error("failed to check email uniqueness", zap.Error(err), zap.String("Email", user.Email))
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	}

	if exist {
		ur.logger.Info("user already exist", zap.String("Email", user.Email))
		return apperrors.ErrUserAlreadyExist
	}

	_, err = ur.db.Exec(`INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)`, user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		ur.logger.Error("failed to save user", zap.String("Username", user.Username), zap.String("Email", user.Email), zap.String("Password", user.Password), zap.Error(err))
		return fmt.Errorf("failed to save user: %w", err)
	}

	ur.logger.Info("Successfuly save User into DB", zap.String("Username", user.Username), zap.String("Email", user.Email))

	return nil
}

func (ur *UserRepo) FindUserByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := ur.db.Get(&user, `SELECT * FROM users WHERE email = $1`, email)
	if errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error("There are not user with that email", zap.String("Email", email))
		return nil, apperrors.ErrNotFoundUser
	}

	if err != nil {
		ur.logger.Error("Failed to find user by email", zap.String("Email", email))
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	ur.logger.Info("Successfuly find User by Email", zap.String("Username", user.Username), zap.String("Email", user.Email))

	return &user, nil
}
