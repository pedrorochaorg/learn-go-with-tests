package main

import (
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver/stores"
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {

	// Open the file in the same directory of the exec
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := stores.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v", err)
	}

	server := NewPlayerServer(store)

	// start's the web server on port 5000 and add's the previous created handler to it.
	// if the operation fails, the error message returned by the command ListenAndServe of the http package
	// must be printed out to the console
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
