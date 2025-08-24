package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func JSONError(w http.ResponseWriter, status int, publicMsg string, logger *zap.Logger) {
	logger.Info("returning error to client", zap.String("Error Msg", publicMsg), zap.Int("Status Code", status))
	JSON(w, status, map[string]string{"error": publicMsg})
}

func ReadBody(w http.ResponseWriter, r *http.Request, encoded_data interface{}, logger *zap.Logger) error {
	if w.Header().Get("Content-Type") != "application/json" {
		JSONError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json", logger)
		return errors.New("ivalid Content-Type")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "failed to read request Body", logger)
		return fmt.Errorf("failed to read request Body: %w", err)
	}

	err = json.Unmarshal(body, encoded_data)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "failed to parse the JSON-encoded data and stores the result", logger)
		return fmt.Errorf("failed to parse the JSON-encoded data and stores the result: %w", err)
	}

	return nil

}
