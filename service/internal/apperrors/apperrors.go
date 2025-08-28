package apperrors

import "errors"

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrNotFoundUser     = errors.New("user doesn't exist")
	ErrInvalidPassword  = errors.New("ivalid password")
)
