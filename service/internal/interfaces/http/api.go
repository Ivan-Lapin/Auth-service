package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func JSONError(w http.ResponseWriter, status int, publicMsg string, logger *zap.Logger) {
	logger.Info(publicMsg)
}
