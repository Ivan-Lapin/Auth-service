package web

import (
	"time"

	"github.com/Ivan-Lapin//Auth-service/service/internal/interfaces/http"
	"github.com/Ivan-Lapin/Auth-service/service/internal/application"
	"github.com/Ivan-Lapin/Auth-service/service/internal/config"
	"go.uber.org/zap"
)

func NewServer(cfg *config.Config, logger *zap.Logger, usertSrv *application.UserService, authSrv *application.AuthService) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("POST api/register", http.RegisterHandler(logger, usertSrv))

	return &http.Server{
		Addr:         cfg.Server.HTTPport,
		Handler:      mux,
		ReadTimeout:  time.Duration(time.Second * 30),
		WriteTimeout: time.Duration(time.Second * 30),
	}
}
