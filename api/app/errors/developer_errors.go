package apperrors

import "errors"

var (
	ErrCreateDeveloperFailed     = errors.New("developer_repo: create developer failed")
	ErrUpdateDeveloperFailed     = errors.New("developer_repo: update developer failed")
	ErrDeveloperNotFound         = errors.New("developer_repo: developer not found")
	ErrUpdateUnauthorized        = errors.New("developer update unauthorized")
	ErrDeveloperValidationFailed = errors.New("developer validation failed")
)
