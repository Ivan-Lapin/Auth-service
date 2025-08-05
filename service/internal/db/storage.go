package db

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	db *sql.DB
}

type User struct {
	Username string
	password string
}

func NewStorage(storagePath string, logger *zap.Logger) (*Storage, error) {
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		logger.Error("failed to opens a database: %w", zap.Error(err))
		return nil, fmt.Errorf("failed to opens a database: %w", err)
	}

	if err := db.Ping(); err != nil {
		logger.Error("failed to verifies a connection to the database: %w", zap.Error(err))
		return nil, fmt.Errorf("failed to verifies a connection to the database: %w", err)
	}

	var tableExist bool

	err = db.QueryRow(`SELECT EXIST (
	SELECT FROM information_shema.tables
	WHERE table_shema = 'public' AND table_name = 'users');`).Scan(&tableExist)

	if tableExist {
		logger.Warn("Users table already exists")
	} else {
		query, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
		);`)

		if err != nil {
			logger.Error("failed to create a prepared statement for later queries or executions: %w", zap.Error(err))
			return nil, fmt.Errorf("failed to create a prepared statement for later queries or executions: %w", err)
		}

		defer query.Close()

		_, err = query.Exec()
		if err != nil {
			logger.Error("failed to execute a prepared statement: %w", zap.Error(err))
			return nil, fmt.Errorf("failed to execute a prepared statement: %w", err)
		}

	}

	return &Storage{db: db}, err

}

func (st *Storage) CreateUser(logger *zap.Logger, user User) {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to generate hash from password: %w", zap.Error(err))
	}

	query := `INSERT INTO users (username, password) VALUES $1 $2;`
	_, err = st.db.Exec(query, user.Username, string(hashedPass))
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			logger.Warn("Username already exists", zap.String("username", user.Username))
		}

		logger.Error("Failed to create user", zap.Error(err), zap.String("username", user.Username))
	}

	logger.Info("User created successfully", zap.String("username", user.Username))
}
