package domain

import (
	"errors"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type AccountId int64

type Account struct {
	Id        AccountId
	Email     string `validate:"email,required"`
	Username  string `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Password  string `validate:"required,max=72"`
}

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrAccountExists    = errors.New("account already exists")
	ErrAccountNotFound  = errors.New("account not found")
)

func MakeAccount(email, username, firstName, lastName, password string) (*Account, error) {
	account := Account{
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
	}

	validator := validator.New()
	err := validator.Struct(account)
	if err != nil {
		return nil, ErrValidationFailed
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	account.Password = string(passwordHash)
	return &account, nil
}
