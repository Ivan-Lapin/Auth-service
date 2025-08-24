package web

import (
	"errors"
	"net/http"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/application"
	"go.uber.org/zap"
)

func RegisterHandler(logger *zap.Logger, userSrvc *application.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			name     string `json:"name"`
			email    string `json:"email"`
			password string `json:"password"`
		}

		if err := ReadBody(w, r, &req, logger); err != nil {
			return
		}

		user, err := userSrvc.Register(req.name, req.email, req.password)
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
