package infrastructure

import (
	"testing"

	"github.com/Ivan-Lapin/Auth-service/service/internal/apperrors"
	"github.com/Ivan-Lapin/Auth-service/service/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	dsn := "host=127.0.0.1 port=8080 user=pinchik password=qwerty1 dbname=postgres sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	return db
}

func TestSaveAndFind(t *testing.T) {
	db := setupTestDB(t)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	userSrv := UserRepo{db: db, logger: logger}
	defer db.Close()

	newID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("NewUUID failed: %v", err)
	}
	user := domain.User{
		ID:       newID,
		Username: "testUsername",
		Email:    "test@email.ru",
	}

	password := "testPassword"

	if err := user.SetPassword(password); err != nil {
		t.Fatalf("SetPassword failed: %v", err)
	}

	if err := userSrv.Save(&user); err != nil {
		if err == apperrors.ErrUserAlreadyExist {
			t.Errorf("ErrUserAlreadyExist: %v", err)
		} else {
			t.Fatalf("Save error: %v", err)
		}
	}

	testEmail := "test@email.ru"

	outputUser, err := userSrv.FindUserByEmail(testEmail)
	if err != nil {
		if err == apperrors.ErrNotFoundUser {
			t.Errorf("ErrNotFoundUser: %v", err)
		} else {
			t.Fatalf("FindUserByEmail error: %v", err)
		}
	}

	if user.Email != outputUser.Email {
		t.Fatalf("expected email %v, got %v", user.Email, outputUser.Email)
	}

}
