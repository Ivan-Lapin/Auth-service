package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/application"
	"go.uber.org/zap"
)

func MainPage(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Testing Hello...")
	}
}

func RegisterHandler(logger *zap.Logger, userSrvc application.UserActions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := ReadBody(w, r, &req, logger); err != nil {
			return
		}

		if req.Name == "" || req.Email == "" || req.Password == "" {
			JSONError(w, http.StatusBadRequest, "Ivalid data request", logger)
			return
		}

		user, err := userSrvc.Register(req.Name, req.Email, req.Password)
		if err != nil {
			if errors.Is(err, apperrors.ErrUserAlreadyExist) {
				JSONError(w, http.StatusConflict, "Email already exist", logger)
				return
			}
			JSONError(w, http.StatusInternalServerError, "Regestration failed", logger)
			return
		}

		JSON(w, http.StatusOK, user)

	}
}
