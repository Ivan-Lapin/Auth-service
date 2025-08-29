package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"go.uber.org/zap"
)

type mockAuthService struct {
	returnToken string
	returnErr   error
}

func (mas *mockAuthService) Login(email, password string) (string, error) {
	return mas.returnToken, mas.returnErr
}

func TestLoginHandler_Success(t *testing.T) {
	authSrv := &mockAuthService{returnToken: "exampleToken"}
	logger := zap.NewNop()

	body := `{"email":"test@email.ru", "password":"testPassword"}`

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler := LoginHandler(logger, authSrv)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %v", w.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if res["token"] != "exampleToken" {
		t.Fatalf("expected token <exampleToken>, got %v", res["toker"])
	}
}

func TestLoginHandler_NotFound(t *testing.T) {
	authSrv := &mockAuthService{returnErr: apperrors.ErrNotFoundUser}
	logger := zap.NewNop()

	body := `{"email":"test@email.ru", "password":"testPassword"}`

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler := LoginHandler(logger, authSrv)
	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %v", w.Code)
	}

}

func TestLoginHandler_InvalidPassword(t *testing.T) {
	authSrv := &mockAuthService{returnErr: apperrors.ErrInvalidPassword}
	logger := zap.NewNop()

	body := `{"email":"test@email.ru", "password":"testPassword"}`

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler := LoginHandler(logger, authSrv)
	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %v", w.Code)
	}
}
