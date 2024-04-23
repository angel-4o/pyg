package persistence

import (
	"database/sql"
	"fmt"

	"pyg.com/api/domain"
)

type accountRepo struct {
	db *sql.DB
}

func MakeAccountRepo(db *sql.DB) *accountRepo {
	return &accountRepo{
		db: db,
	}
}

func (repo *accountRepo) Create(account *domain.Account) (domain.AccountId, error) {
	// `insert into accounts (email, username, first_name, last_name, password)
	//  select $1, $2, $3, $4, $5
	//  where not exists (select null from accounts where email = $1)`,

	err := repo.db.QueryRow(
		`insert into accounts (email, username, first_name, last_name, password)
		 values ($1, $2, $3, $4, $5)
		 on conflict do nothing
		 returning id`,
		account.Email,
		account.Username,
		account.FirstName,
		account.LastName,
		account.Password,
	).Scan(&account.Id)

	if err != nil {
		return domain.AccountId(0), fmt.Errorf("%w: %w", domain.ErrAccountExists, err)
	}

	return account.Id, nil
}

func (repo *accountRepo) FindById(accountId domain.AccountId) (domain.Account, error) {
	account := domain.Account{}
	err := repo.db.
		QueryRow(`
			select id, email, username, first_name, last_name, password
			from accounts
			where id = $1`,
			accountId).
		Scan(
			&account.Id,
			&account.Email,
			&account.Username,
			&account.FirstName,
			&account.LastName,
			&account.Password,
		)

	if err != nil {
		return domain.Account{}, fmt.Errorf("%w: %w", domain.ErrAccountNotFound, err)
	}

	return account, nil
}

func (repo *accountRepo) FindByUsername(username string) (domain.Account, error) {
	account := domain.Account{}
	err := repo.db.
		QueryRow(`
			select id, email, username, first_name, last_name, password
			from accounts
			where username = $1`,
			username).
		Scan(
			&account.Id,
			&account.Email,
			&account.Username,
			&account.FirstName,
			&account.LastName,
			&account.Password,
		)

	if err != nil {
		return domain.Account{}, fmt.Errorf("%w: %w", domain.ErrAccountNotFound, err)
	}

	return account, nil
}
