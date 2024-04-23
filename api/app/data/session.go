package data

import (
	"time"

	"pyg.com/api/domain"
)

type SessionToken string

func (token *SessionToken) IsValid() bool {
	return *token != ""
}

type Session struct {
	AccountId domain.AccountId
	Token SessionToken
	ExpiryTime time.Time
}