package apperrors

import "errors"

var (
	ErrCannotMarshalMembers = errors.New("cannot marshal members")
	ErrCannotUnmarshalMembers = errors.New("cannot unmarshal members")
)