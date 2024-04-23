package infra

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	defaultConnectionString = "postgres://pyg:pyg@localhost:5432/pyg?sslmode=disable"
)

func ConnectToDatabase() *sql.DB {
	connectionString := os.Getenv("PYG_CONNECTION_STRING")
	if connectionString == "" {
		connectionString = defaultConnectionString
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalln(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return db
}
