package app

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"pyg.com/api/app/data"
	apperrors "pyg.com/api/app/errors"
	"pyg.com/api/domain"
	"pyg.com/api/persistence"
)

func makeRandomToken() string {
	buffer := make([]byte, 8)
	rand.Read(buffer)
	return fmt.Sprintf("%x", buffer)
}

func Login(username, password string, db *sql.DB) (data.SessionToken, error) {
	accountRepo := persistence.MakeAccountRepo(db)
	account, err := accountRepo.FindByUsername(username)
	if err != nil {
		return "", fmt.Errorf("%w: %s", apperrors.ErrAccountNotFound, username)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("%w: %s", apperrors.ErrInvalidPassword, username)
	}

	sessionToken := data.SessionToken(makeRandomToken())
	session := data.Session{
		AccountId:  account.Id,
		Token:      sessionToken,
		ExpiryTime: time.Now().Add(time.Duration(1 * time.Hour)),
	}

	sessionDataSource := persistence.MakeSessionDataSource(db)
	err = sessionDataSource.CreateSession(session)
	if err != nil {
		return data.SessionToken(""), err
	}

	return sessionToken, nil
}

func Login2(username, password string, db *sql.DB) (domain.Account, data.SessionToken, error) {
	accountRepo := persistence.MakeAccountRepo(db)
	account, err := accountRepo.FindByUsername(username)
	if err != nil {
		return account, "", fmt.Errorf("%w: %s", apperrors.ErrAccountNotFound, username)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return account, "", fmt.Errorf("%w: %s", apperrors.ErrInvalidPassword, username)
	}

	sessionToken := data.SessionToken(makeRandomToken())
	session := data.Session{
		AccountId:  account.Id,
		Token:      sessionToken,
		ExpiryTime: time.Now().Add(time.Duration(1 * time.Hour)),
	}

	sessionDataSource := persistence.MakeSessionDataSource(db)
	err = sessionDataSource.CreateSession(session)
	if err != nil {
		return account, data.SessionToken(""), err
	}

	return account, sessionToken, nil
}
