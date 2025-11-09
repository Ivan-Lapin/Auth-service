package application

import (
	"errors"
	"testing"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/domain"
	"go.uber.org/zap"
)

type mockUserRepo struct {
	users   map[string]*domain.User
	saveErr error
	findErr error
}

func NewMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[string]*domain.User),
	}
}

func (mur *mockUserRepo) Save(user *domain.User) error {
	if mur.saveErr != nil {
		return mur.saveErr
	}

	if _, exist := mur.users[user.Email]; exist {
		return apperrors.ErrUserAlreadyExist
	}

	mur.users[user.Email] = user

	return nil
}

func (mur *mockUserRepo) FindUserByEmail(email string) (*domain.User, error) {
	if mur.findErr != nil {
		return nil, mur.findErr
	}

	user, exist := mur.users[email]
	if !exist {
		return nil, apperrors.ErrNotFoundUser
	}

	return user, nil
}

func TestUserServiceRegisterSuccess(t *testing.T) {
	name, email, password := "testName", "test@email.ru", "testPass"

	testRepo := NewMockUserRepo()
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Errorf("Failed to create logger: %v", err)
	}
	userSrv := NewUserService(testRepo, logger)

	user, err := userSrv.Register(name, email, password)
	if err != nil {
		t.Fatalf("Register error: %v", err)
	}

	if user.Email != email {
		t.Fatalf("expected email: %v got emal: %v", email, user.Email)
	}
}

func TestUserServiceRegisterWrong(t *testing.T) {
	name, email, password := "testName", "test@email.ru", "testPass"

	testRepo := NewMockUserRepo()
	logger := zap.NewNop()

	userSrv := NewUserService(testRepo, logger)

	_, _ = userSrv.Register(name, email, password)

	newName := "otherName"

	_, err := userSrv.Register(newName, email, password)
	if err == nil {
		t.Fatal("expected error for existing email")
	}

	if !errors.Is(err, apperrors.ErrUserAlreadyExist) {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestUserServiceRegisterSaveError(t *testing.T) {
	testRepo := NewMockUserRepo()
	testRepo.saveErr = errors.New("db failure")

	logger := zap.NewNop() // отключаем вывод логов

	userSrv := NewUserService(testRepo, logger)

	_, err := userSrv.Register("name", "email", "password")

	if err == nil {
		t.Fatal("expected error from repo.Save")
	} else {
		t.Logf("Received expected error: %v", err)
	}
}
