package main

import (
	"fmt"
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)
	cli := poker.NewCLI(os.Stdin, os.Stdout, game)
	cli.PlayPoker()

}