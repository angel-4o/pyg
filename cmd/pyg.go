package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"pyg.com/api/infra"
)

const (
	defaultListenAddress = ":8080"
)

func main() {
	db := infra.ConnectToDatabase()
	defer db.Close()

	rootHandler := infra.NewRootHandler(db)

	// server
	listenAddress := os.Getenv("PYG_LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = defaultListenAddress
	}

	fmt.Printf("starting server on %s\n", listenAddress)
	err := http.ListenAndServe(listenAddress, rootHandler)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
