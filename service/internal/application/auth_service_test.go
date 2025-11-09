package application

import (
	"testing"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/domain"
	"github.com/Ivan-Lapin/Auth-service/service/pkg/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type mockAuthRepo struct {
	users map[string]*domain.User
}

func NewMockAuthRepo() *mockAuthRepo {
	return &mockAuthRepo{
		users: map[string]*domain.User{},
	}
}

func (mar *mockAuthRepo) Save(user *domain.User) error {

	if _, exist := mar.users[user.Email]; exist {
		return apperrors.ErrUserAlreadyExist
	}

	mar.users[user.Email] = user

	return nil
}

func (mar *mockAuthRepo) FindUserByEmail(email string) (*domain.User, error) {

	user, exist := mar.users[email]
	if !exist {
		return nil, apperrors.ErrNotFoundUser
	}

	return user, nil
}

func TestAuthService_Login_Success(t *testing.T) {
	testRepo := NewMockAuthRepo()
	logger := zap.NewNop()
	testJWT := jwt.NewJWT([]byte("secret-key"), logger)
	authSrv := NewAuthService(testRepo, testJWT, logger)

	id, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to New UUID: %v", err)
	}
	user := domain.User{
		ID:       id,
		Username: "testName",
		Email:    "test@email.ru",
	}

	err = user.SetPassword("testPassword")
	if err != nil {
		t.Fatalf("Failed SetPassword: %v", err)
	}

	err = testRepo.Save(&user)
	if err != nil {
		t.Fatalf("Failed Save: %v", err)
	}

	token, err := authSrv.Login("test@email.ru", "testPassword")
	if err != nil {
		t.Fatalf("Failed Login: %v", err)
	}

	if token == "" {
		t.Fatalf("Expected token, got empty string")
	}

}
