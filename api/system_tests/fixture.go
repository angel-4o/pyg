package systemtest

import (
	"database/sql"
	"log"
	"net/http"

	"pyg.com/api/infra"
)

func newDbAndHandler() (*sql.DB, http.Handler) {
	db := connectAndCleanDatabase()
	return db, infra.NewRootHandler(db)
}

func connectAndCleanDatabase() *sql.DB {
	db := infra.ConnectToDatabase()
	cleanDb(db)
	return db
}

func cleanDb(db *sql.DB) {
	var err error

	_, err = db.Exec("delete from developers")
	if err != nil {
		log.Fatalln("unexpected failure deleting developers: %w", err)
	}

	_, err = db.Exec("delete from sessions")
	if err != nil {
		log.Fatalln("unexpected failure deleting sessions: %w", err)
	}

	_, err = db.Exec("delete from accounts")
	if err != nil {
		log.Fatalln("unexpected failure deleting accounts: %w", err)
	}
}
