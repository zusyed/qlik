package main

import (
	"log"
	"net/http"

	"github.com/zusyed/qlik/dal"
	"github.com/zusyed/qlik/handlers"
)

func main() {
	mDB, err := dal.NewMessageDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %+v", err)
	}

	h := handlers.NewHandler(mDB)
	router := handlers.NewRouter(h)

	log.Printf("Initialized HTTP server")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("Could not initiate a listen on http requests: %s", err)
	}
}
