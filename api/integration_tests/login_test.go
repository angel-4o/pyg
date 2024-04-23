package integrationtest

import (
	"log"
	"testing"

	"pyg.com/api/app"
	"pyg.com/api/infra"
)

func TestLogin(t *testing.T) {
	db := infra.ConnectToDatabase()
	defer db.Close()
	cleanDb(db)

	err := app.Register("user@email.com", "username", "FirstName", "LastName", "Pass", db) 
	if err != nil { log.Fatal(err) }

	// Success with correct password
	_, err = app.Login("username", "Pass", db)
	if err != nil { t.Fatal("failed to login") }

	// Failure with incorrect password
	_, err = app.Login("username", "Wrong", db)
	if err == nil { t.Fatal("login erroneously succeeded") }
}
