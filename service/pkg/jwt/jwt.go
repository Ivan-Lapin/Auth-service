package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type JWT struct {
	secter []byte
	logger *zap.Logger
}

func NewJWT(secret []byte, logger *zap.Logger) *JWT {
	return &JWT{
		secter: secret, logger: logger,
	}
}

func (j *JWT) Generate(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secter)
	if err != nil {
		j.logger.Error("Failed to create and returns a complete, signed JWT", zap.Error(err))
		return "", fmt.Errorf("failed to create and returns a complete, signed JWT: %w", err)
	}

	return signed, nil
}
