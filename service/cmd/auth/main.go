package main

import (
	"log"
	"net/http"

	"github.com/Ivan-Lapin/Auth-service/service/internal/application"
	"github.com/Ivan-Lapin/Auth-service/service/internal/config"
	"github.com/Ivan-Lapin/Auth-service/service/internal/infrastructure"
	"github.com/Ivan-Lapin/Auth-service/service/internal/interfaces/web"
	"github.com/Ivan-Lapin/Auth-service/service/pkg/jwt"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to Load Config", err)
	}

	logger, err := config.NewLogger(cfg)
	if err != nil {
		log.Fatal("Failed to Create Logger", err)
	}

	defer logger.Sync()

	db, err := infrastructure.NewDB(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	err = infrastructure.RunMigrations(db, cfg.DB.MigrationsPath, logger)
	if err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	userRepo := infrastructure.NewUserRepository(db, logger)
	userService := application.NewUserService(userRepo, logger)
	JSONWebToken := jwt.NewJWT([]byte(cfg.JWT.JWTSecretKey), logger)
	authService := application.NewAuthService(userRepo, JSONWebToken, logger)

	server := web.NewServer(cfg, logger, userService, authService)

	logger.Info("Starting server...", zap.String("port", cfg.Server.HTTPport))
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server failed", zap.Error(err))
	}

}
