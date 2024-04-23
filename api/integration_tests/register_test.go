package integrationtest

import (
	"errors"
	"testing"

	"pyg.com/api/app"
	"pyg.com/api/domain"
	"pyg.com/api/infra"
)

func TestRegister(t *testing.T) {
	db := infra.ConnectToDatabase()
	defer db.Close()
	cleanDb(db)

	err := app.Register("user@email.com", "username", "FirstName", "LastName", "Pass", db) 
	if err != nil { t.Fatal(err) }

	err = app.Register("user@email.com", "username", "FirstName", "LastName", "Pass", db) 
	if !errors.Is(err, domain.ErrAccountExists) {
		t.Fatal("account registered multiple times")
	}
}
