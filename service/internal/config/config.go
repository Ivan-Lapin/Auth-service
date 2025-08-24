package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Server struct {
		HTTPport string
	}
	DB struct {
		ConnToDB       string
		MigrationsPath string
	}
	JWT struct {
		JWTSecretKey string
	}
	Logging struct {
		Level       string
		Development bool
	}
}

func LoadConfig() (*Config, error) {
	configpath := os.Getenv("CONFIG_PATH_AUTH_SERVICE")

	viper.SetConfigFile(configpath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("failed to discover and to load the configuration file from disk and key/value stores: %w", zap.Error(err))
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the config into the struct: %w", err)
	}

	log.Println("Load file configuration", zap.Any("Config", cfg))

	return &cfg, nil
}

func NewLogger(cfg *Config) (*zap.Logger, error) {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:      true,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	}

	return config.Build()
}
