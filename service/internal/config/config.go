package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	ConnToDB string
	HTTPport string
	GRPCport string
}

func LoadConfig(logger *zap.Logger) *Config {
	configpath := os.Getenv("CONFIG_PATH")
	logger.Info("Using congif path from env", zap.String("config path", configpath))

	viper.SetConfigFile(configpath)

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("failed to discover and to load the configuration file from disk and key/value stores: %w", zap.Error(err))
	}

	config := &Config{
		ConnToDB: viper.GetString("connectToDB"),
		HTTPport: viper.GetString("httpPort"),
		GRPCport: viper.GetString("grpcPort"),
	}

	log.Println("Load file configuration", zap.Any("Config", config))

	return config
}
