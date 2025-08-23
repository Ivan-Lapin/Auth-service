package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameExist = errors.New("Username already exist")
)

type DataBase interface {
	CreateUser(logger *zap.Logger, user User) error
	GetUserByUsername(logger *zap.Logger, username string) (User, error)
}

type Storage struct {
	db *sql.DB
}

type User struct {
	Username string
	Password string
}

func NewStoragePostgre(storagePath string, logger *zap.Logger) (*Storage, error) {
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		logger.Error("Failed to opens a database: %w", zap.Error(err))
		return nil, fmt.Errorf("Failed to opens a database: %w", err)
	}

	if err := db.Ping(); err != nil {
		logger.Error("Failed to verifies a connection to the database: %w", zap.Error(err))
		return nil, fmt.Errorf("Failed to verifies a connection to the database: %w", err)
	}

	var tableExist bool

	err = db.QueryRow(`SELECT EXISTS (
    SELECT FROM information_schema.tables
    WHERE table_schema = 'public' AND table_name = 'users'
	);`).Scan(&tableExist)

	if tableExist {
		logger.Warn("Users table already exists")
	} else {
		query, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		Password VARCHAR(255) NOT NULL
		);`)

		if err != nil {
			logger.Error("Failed to create a prepared statement for later queries or executions: %w", zap.Error(err))
			return nil, fmt.Errorf("Failed to create a prepared statement for later queries or executions: %w", err)
		}

		defer query.Close()

		_, err = query.Exec()
		if err != nil {
			logger.Error("Failed to execute a prepared statement: %w", zap.Error(err))
			return nil, fmt.Errorf("Failed to execute a prepared statement: %w", err)
		}

	}

	return &Storage{db: db}, err

}

func (st *Storage) CreateUser(logger *zap.Logger, user User) error {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to generate hash from Password: %w", zap.Error(err))
		return fmt.Errorf("Failed to generate hash from Password: %w", err)
	}

	query := `INSERT INTO users (username, password) VALUES $1 $2;`
	_, err = st.db.Exec(query, user.Username, string(hashedPass))
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			logger.Warn("Username already exists", zap.String("username", user.Username))
			return ErrUsernameExist
		}

		logger.Error("Failed to create user", zap.Error(err), zap.String("username", user.Username))
		return fmt.Errorf("Failed to create user %s: %w", user.Username, err)
	}

	logger.Info("User created successfully", zap.String("username", user.Username))
	return err
}

func (st *Storage) GetUserByUsername(logger *zap.Logger, username string) (User, error) {
	var user User

	query := `SELECT username, password FROM users WHERE username = $1`
	err := st.db.QueryRow(query, username).Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("User %s not found: %w", zap.String("username", username), zap.Error(err))
			return User{}, fmt.Errorf("User %s not found: %w", username, err)
		}

		logger.Error("failed to find user", zap.String("username", username), zap.Error(err))
		return User{}, fmt.Errorf("Failed to find user - %s: %w", username, err)
	}

	return user, nil
}
