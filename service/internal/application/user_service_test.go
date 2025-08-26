package application

import "github.com/Ivan-Lapin/Auth-service/service/internal/domain"

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
