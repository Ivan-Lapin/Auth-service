package handler

import (
	"errors"
	"net/http"

	"github.com/Ivan-Lapin/Auth-service/service/internal/db"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.Logger
	storage db.Storage
}

type UserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewHandler(logger *zap.Logger, storage db.Storage) *Handler {
	return &Handler{
		logger:  logger,
		storage: storage,
	}
}

func (h *Handler) Register(c echo.Context) error {
	var ureq UserRequest

	if err := c.Bind(&ureq); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := c.Validate(ureq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	err := h.storage.CreateUser(h.logger, db.User{ureq.Username, ureq.Password})
	if err != nil {
		if errors.Is(err, db.ErrUsernameExist) {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Username already exists"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, nil)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {

}
