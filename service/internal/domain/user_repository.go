package domain

type UserRepository interface {
	Save(user *User) error
	FindUserByEmail(email string) (*User, error)
}
