package app

import (
	"database/sql"

	"pyg.com/api/domain"
	"pyg.com/api/persistence"
)

// type RegisterParameters struct {
// 	Email     string `validate:"email,required"`
// 	Username  string `validate:"required"`
// 	FirstName string `validate:"required"`
// 	LastName  string `validate:"required"`
// 	Password  string `validate:"required,max=72"`
// }

func Register(email, username, firstName, lastName, password string, db *sql.DB) error {
	account, err := domain.MakeAccount(email, username, firstName, lastName, password)
	if err != nil {
		return err
	}

	accountRepo := persistence.MakeAccountRepo(db)
	_, err = accountRepo.Create(account)
	return err
}

func Register2(email, username, firstName, lastName, password string, profileType string, db *sql.DB) error {
	account, err := domain.MakeAccount2(email, username, firstName, lastName, password, profileType)
	if err != nil {
		return err
	}

	accountRepo := persistence.MakeAccountRepo(db)
	_, err = accountRepo.Create(account)
	// if err != nil {
	// 	// create dev or pub profile
	// 	if profileType == "developer" {
	// 		// create developer

	// 		CreateDeveloper("", username, db)
	// 	}
	// 	if profileType == "publisher" {
	// 		// create publisher
	// 		CreatePublisher("", username, db)
	// 	}

	// }
	return err
}
