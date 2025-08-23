package service

import (
	"fmt"
	"time"

	"github.com/Ivan-Lapin/Auth-service/service/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthService interface {
	GenertaeToken(username string) (string, error)
	ValidateToken(tokenString string) (Claims, error)
}

type authService struct {
	secret string
	logger *zap.Logger
}

type App struct {
	Config *config.Config
	Logger *zap.Logger
	AS     AuthService
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewApp(config *config.Config, logger *zap.Logger) *App {
	return &App{Config: config, Logger: logger}
}

func (a *App) NewAuthService() AuthService {
	return &authService{secret: a.Config.JWTSecretKey, logger: a.Logger}
}

func (as *authService) GenertaeToken(username string) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(as.secret)
	if err != nil {
		as.logger.Error("Failed to sign JWT token", zap.Error(err), zap.String("username", username))
		return "", fmt.Errorf("Failed to sign JWT token: %w", err)
	}

	as.logger.Info("JWT token generated", zap.String("username", username))

	return tokenString, nil
}

func (as *authService) ValidateToken(tokenString string) (Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			as.logger.Warn("Unexpected singning method", zap.String("alg", token.Header["alg"].(string)))
			return claims, fmt.Errorf("Unexpected singning method: %v", token.Header["alg"].(string))
		}
		return []byte(as.secret), nil
	})

	if err != nil {
		return Claims{}, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return Claims{}, fmt.Errorf("token is not valid")
	}

	return *claims, nil
}
