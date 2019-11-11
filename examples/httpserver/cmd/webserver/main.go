package main

import (
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver"
	"log"
	"net/http"
)

const dbFileName = "game.db.json"

func main() {

	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)


	server, err := poker.NewPlayerServer(store, game)

	if err != nil {
		log.Fatal(err)
	}


	log.Println("Webserver listening on port 500")
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}

}
