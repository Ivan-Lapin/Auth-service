package main

import (
	"log"
	"time"

	"github.com/Ivan-Lapin/Auth-service/service/internal/config"
	"github.com/Ivan-Lapin/Auth-service/service/internal/db"
	"github.com/Ivan-Lapin/Auth-service/service/internal/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := zap.NewDevelopmentConfig()

	// Устанавливаем TimeKey, чтобы включить временную метку
	cfg.EncoderConfig.TimeKey = "time" // Ключ для временной метки
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("02.01.2006. 15:04:05"))
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatal("failed to build constructs a logger from the config and options: %w", err)
	}

	defer logger.Sync()

	config := config.LoadConfig(logger)

	_, err = db.NewStoragePostgre(config.ConnToDB, logger)

	e := echo.New()

	e.Validator = middleware.NewCustomValidator()
}
