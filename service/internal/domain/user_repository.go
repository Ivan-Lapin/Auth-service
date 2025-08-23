package domain

import "github.com/google/uuid"

type UserRepository interface {
	Save(user *User) (uuid.UUID, error)
	FindUserByEmail(email string) (*User, error)
}
