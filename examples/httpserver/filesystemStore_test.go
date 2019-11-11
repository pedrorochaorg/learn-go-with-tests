package poker_test

import (
	poker "github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver"
	"log"
	"testing"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("/league from a reader", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)

		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}

		got := store.GetLeague()

		want := poker.League{
			{"Chris", 33},
			{"Cleo", 10},
		}

		poker.AssertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		poker.AssertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)

		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}


		got := store.GetPlayerScore("Chris")

		want := 33

		poker.AssertScoreEquals(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)

		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}

		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		poker.AssertScoreEquals(t, got, want)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)

		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}


		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		poker.AssertScoreEquals(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := poker.NewFileSystemPlayerStore(database)

		poker.AssertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)

		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}

		got := store.GetLeague()

		want := []poker.Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		poker.AssertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		poker.AssertLeague(t, got, want)
	})
}
