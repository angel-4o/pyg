package apperrors

import "errors"

var (
	ErrAccountNotFound = errors.New("login: account not found")
	ErrInvalidPassword = errors.New("login: invalid password")
)
