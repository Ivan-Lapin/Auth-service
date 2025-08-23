package apperrors

import "errors"

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrNotFoundUser     = errors.New("user doesn't exist")
	ErrIvalidPassword   = errors.New("ivalid password")
)
