package infrastructure

import (
	"fmt"

	"github.com/Ivan-Lapin/Auth-service/service/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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

func RunMigrations(db *sqlx.DB, migrationsPath string, logger *zap.Logger) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.Error("Failed to create postgres driver", zap.Error(err))
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
	if err != nil {
		logger.Error("Failed to return a new Migrate instance", zap.Error(err))
		return fmt.Errorf("failed to return a new Migrate instance: %w", err)
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Error("Failed to looks all active migration version or will migrate it", zap.Error(err))
		return fmt.Errorf("failed to looks all active migration version or will migrate it: %w", err)
	}

	logger.Info("Migrations applied successfully")
	return nil
}
