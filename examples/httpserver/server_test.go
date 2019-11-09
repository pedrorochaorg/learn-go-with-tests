package main

import (
	"encoding/json"
	"fmt"
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver/helpers"
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver/objects"
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver/stores"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)


func TestGETPlayers(t *testing.T) {
	store := stores.StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		helpers.AssertResponseStatusCode(t, response.Code, http.StatusOK)
		helpers.AssertResponseBody(t, response.Body.String(), "20")

	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		helpers.AssertResponseStatusCode(t, response.Code, http.StatusOK)
		helpers.AssertResponseBody(t, response.Body.String(), "10")

	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		helpers.AssertResponseStatusCode(t, response.Code, http.StatusNotFound)
	})

}

func TestStoreWins(t *testing.T) {
	store := stores.StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("it returns accepted on POST", func(t *testing.T) {
		request := newPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		helpers.AssertResponseStatusCode(t, response.Code, http.StatusAccepted)

		if len(store.WinCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.WinCalls), 1)
		}
	})

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		helpers.AssertResponseStatusCode(t, response.Code, http.StatusAccepted)

		if store.WinCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.WinCalls[0], player)
		}
	})
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := helpers.CreateTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := stores.NewFileSystemPlayerStore(database)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v", err)
	}

	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		helpers.AssertResponseStatusCode(t, response.Code, http.StatusOK)
		helpers.AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		helpers.AssertResponseStatusCode(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := objects.League{
			{"Pepper", 3},
		}
		helpers.AssertLeague(t, got, want)
	})
}

func TestRecordingWinsConcurrentlyAndRetrieveThem(t *testing.T) {
	players := []struct {
		Name   string
		Hits   int
		Result string
	}{
		{"Pepper", 3, "3"},
		{"Pepper", 5, "8"},
		{"John", 5, "5"},
	}

	database, cleanDatabase := helpers.CreateTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := stores.NewFileSystemPlayerStore(database)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v", err)
	}

	server := NewPlayerServer(store)

	for _, test := range players {
		t.Run(fmt.Sprintf("testing concurrent win calls to %s", test.Name), func(t *testing.T) {

			var wg sync.WaitGroup
			wg.Add(test.Hits)

			for i := 0; i < test.Hits; i++ {
				go func(name string, w *sync.WaitGroup) {
					server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(name))
					wg.Done()
				}(test.Name, &wg)
			}

			wg.Wait()
			response := httptest.NewRecorder()
			server.ServeHTTP(response, newGetScoreRequest(test.Name))
			helpers.AssertResponseStatusCode(t, response.Code, http.StatusOK)

			helpers.AssertResponseBody(t, response.Body.String(), test.Result)

		})
	}

}

func TestLeague(t *testing.T) {

	t.Run("it returns 200 0on /league", func(t *testing.T) {
		store := stores.StubPlayerStore{}
		server := NewPlayerServer(&store)

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got objects.League

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
		}

		helpers.AssertResponseStatusCode(t, response.Code, http.StatusOK)
	})

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := objects.League{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := stores.StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		helpers.AssertResponseStatusCode(t, response.Code, http.StatusOK)
		helpers.AssertLeague(t, got, wantedLeague)

		helpers.AssertContentType(t, response, jsonContentType)

	})
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func getLeagueFromResponse(t *testing.T, body io.Reader) (league objects.League) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}
