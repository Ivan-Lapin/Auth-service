package infrastructure

import (
	"fmt"

	"github.com/Ivan-Lapin/Auth-service/service/internal/config"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func NewDB(config *config.Config, logger *zap.Logger) (*sqlx.DB, error) {

	db, err := sqlx.Connect("postgres", config.DB.ConnToDB)
	if err != nil {
		logger.Error("Failed to connect to DB:", zap.String("dsn", config.DB.ConnToDB), zap.Error(err))
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	logger.Info("DB connection established")
	return db, nil
}
