package web

import (
	"errors"
	"net/http"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/application"
	"go.uber.org/zap"
)

func LoginHandler(logger *zap.Logger, authSvc application.AuthActions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := ReadBody(w, r, &req, logger); err != nil {
			return
		}

		if req.Email == "" || req.Password == "" {
			JSONError(w, http.StatusBadRequest, "invalid data request", logger)
			return
		}

		token, err := authSvc.Login(req.Email, req.Password)
		if err != nil {
			if errors.Is(err, apperrors.ErrNotFoundUser) {
				JSONError(w, http.StatusUnauthorized, "User not found", logger)
				return
			}
			if errors.Is(err, apperrors.ErrInvalidPassword) {
				JSONError(w, http.StatusUnauthorized, "Invalid password", logger)
				return
			}
			JSONError(w, http.StatusInternalServerError, "Login failed", logger)
			return
		}

		JSON(w, http.StatusOK, map[string]string{"token": token})

	}
}
