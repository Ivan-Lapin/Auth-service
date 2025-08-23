package web

import (
	"net/http"
	"time"

	"github.com/Ivan-Lapin/Auth-service/service/internal/config"
	"go.uber.org/zap"
)

func NewServer(cfg *config.Config, logger *zap.Logger) *http.Server {
	mux := http.NewServeMux()
	return &http.Server{
		Addr:         cfg.Server.HTTPport,
		Handler:      mux,
		ReadTimeout:  time.Duration(time.Second * 30),
		WriteTimeout: time.Duration(time.Second * 30),
	}
}
