package http

import (
	"errors"
	"net/http"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/application"
	"go.uber.org/zap"
)

func LoginHandler(logger *zap.Logger, authSvc *application.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			email    string `json:"email"`
			password string `json:"password"`
		}

		if err := ReadBody(w, r, &req, logger); err != nil {
			return
		}

		token, err := authSvc.Login(req.email, req.password)
		if err != nil {
			if errors.Is(err, apperrors.ErrNotFoundUser) {
				JSONError(w, http.StatusUnauthorized, "User not found", logger)
			}
			if errors.Is(err, apperrors.ErrIvalidPassword) {
				JSONError(w, http.StatusUnauthorized, "Invalid password", logger)
			}
			JSONError(w, http.StatusInternalServerError, "Login failed", logger)
		}

		JSON(w, http.StatusOK, map[string]string{"token": token})

	}
}
