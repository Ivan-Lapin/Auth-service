package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/domain"
	"go.uber.org/zap"
)

type mockUserService struct {
	registredCall bool
	returnUser    *domain.User
	returnError   error
}

func (mus *mockUserService) Register(name, email, password string) (*domain.User, error) {
	mus.registredCall = true

	return mus.returnUser, mus.returnError
}

func TestRegisterHandler_Success(t *testing.T) {
	userSrv := &mockUserService{
		returnUser: &domain.User{Username: "testName", Email: "test@email.ru"},
	}
	logger := zap.NewNop()

	body := `{"name":"testName", "email":"test@email.ru", "password":"testPassword"}`

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler := RegisterHandler(logger, userSrv)

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if response["email"] != "test@email.ru" {
		t.Fatalf("expected email test@email.ru, got %s", response["email"])
	}
}

func TestRegisterHandler_EmailExist(t *testing.T) {
	userSvc := &mockUserService{
		returnError: apperrors.ErrUserAlreadyExist,
	}
	logger := zap.NewNop()

	body := `{"name":"testName", "email":"test@email.ru", "password":"testPass"}`

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler := RegisterHandler(logger, userSvc)
	handler(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %v", w.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("ivalid json response: %v", err)
	}

	if _, ok := res["error"]; !ok {
		t.Fatalf("exepcted error in response")
	}

}

func TestRegisterHandler_BadRequest(t *testing.T) {
	userSrv := &mockUserService{}
	logger := zap.NewNop()

	body := `{"name":"name", "email":"email", "password":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler := RegisterHandler(logger, userSrv)
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %v", w.Code)
	}

}
