package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Save(user *domain.User) error {

	var exist bool
	err := ur.db.Get(&exist, `SELECT EXIST(SELECT 1 FROM user WHERE email = $1)`, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	}

	if exist {
		return apperrors.ErrUserAlreadyExist
	}

	_, err = ur.db.Exec(`INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`, user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

func (ur *UserRepository) FindUserByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := ur.db.Get(&user, `SELECT * FROM users WHERE email = $1`, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrUserAlreadyExist
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}
