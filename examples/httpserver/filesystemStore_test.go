package main

import (
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver/helpers"
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver/objects"
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver/stores"
	"log"
	"testing"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("/league from a reader", func(t *testing.T) {
		database, cleanDatabase := helpers.CreateTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := stores.NewFileSystemPlayerStore(database)

		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}

		got := store.GetLeague()

		want := objects.League{
			{"Cleo", 10},
			{"Chris", 33},
		}

		helpers.AssertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		helpers.AssertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := helpers.CreateTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := stores.NewFileSystemPlayerStore(database)

		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}


		got := store.GetPlayerScore("Chris")

		want := 33

		helpers.AssertScoreEquals(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := helpers.CreateTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := stores.NewFileSystemPlayerStore(database)

		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}

		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		helpers.AssertScoreEquals(t, got, want)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		database, cleanDatabase := helpers.CreateTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := stores.NewFileSystemPlayerStore(database)

		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}


		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		helpers.AssertScoreEquals(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := helpers.CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := stores.NewFileSystemPlayerStore(database)

		helpers.AssertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := helpers.CreateTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := stores.NewFileSystemPlayerStore(database)

		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}

		got := store.GetLeague()

		want := []objects.Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		helpers.AssertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		helpers.AssertLeague(t, got, want)
	})
}
