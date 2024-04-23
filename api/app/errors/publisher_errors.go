package apperrors

import "errors"

var (
	ErrCreatePublisherFailed     = errors.New("publisher_repo: create failed")
	ErrUpdatePublisherFailed     = errors.New("publisher_repo: update failed")
	ErrPublisherNotFound         = errors.New("publisher_repo: publisher not found")
	ErrPublisherValidationFailed = errors.New("publisher validation failed")
)
