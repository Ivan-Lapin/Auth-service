package web

import (
	"net/http"
	"time"

	"github.com/Ivan-Lapin/Auth-service/service/internal/application"
	"github.com/Ivan-Lapin/Auth-service/service/internal/config"
	"go.uber.org/zap"
)

func NewServer(cfg *config.Config, logger *zap.Logger, usertSrv application.UserActions, authSrv *application.AuthService) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/test", MainPage(logger))
	mux.HandleFunc("/api/register", RegisterHandler(logger, usertSrv))
	mux.HandleFunc("/api/login", LoginHandler(logger, authSrv))

	return &http.Server{
		Addr:         cfg.Server.HTTPPort,
		Handler:      mux,
		ReadTimeout:  time.Duration(time.Second * 30),
		WriteTimeout: time.Duration(time.Second * 30),
	}
}
