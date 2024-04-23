package persistence

import (
	"database/sql"
	"errors"
	"time"

	"pyg.com/api/app/data"
	"pyg.com/api/domain"
)

var (
	ErrSessionNotFound = errors.New("session_data_source: session not found")
)

type SessionDataSource struct {
	db *sql.DB
}

func MakeSessionDataSource(db *sql.DB) *SessionDataSource {
	return &SessionDataSource{
		db: db,
	}
}

func (dataSource *SessionDataSource) CreateSession(session data.Session) error {
	result, err := dataSource.db.Exec(
		`insert into sessions (account_id, token, expiry_time)
	     values ($1, $2, $3)
		 on conflict do nothing`,
		session.AccountId,
		session.Token,
		session.ExpiryTime,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} 
	
	if rowsAffected != 1 {
		return domain.ErrAccountExists
	}

	return nil
}

func (dataSource *SessionDataSource) FindSessionByToken(sessionToken data.SessionToken) (data.Session, error) {
	session := data.Session{}

	err := dataSource.db.QueryRow(
		`select account_id, token, expiry_time
		 from sessions
		 where token = $1`,
		sessionToken,
	).Scan(&session.AccountId, &session.Token, &session.ExpiryTime)

	if err != nil || time.Now().After(session.ExpiryTime) {
		return data.Session{}, ErrSessionNotFound
	}

	return session, nil
}