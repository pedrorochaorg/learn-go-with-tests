package poker_test

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestGETPlayers(t *testing.T) {
	store := poker.StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := mustMakePlayerServer(t, &store, &poker.GameSpy{})

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatusCode(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "20")

	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatusCode(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "10")

	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatusCode(t, response.Code, http.StatusNotFound)
	})

}

func TestStoreWins(t *testing.T) {
	store := poker.StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := mustMakePlayerServer(t, &store, &poker.GameSpy{})

	t.Run("it returns accepted on POST", func(t *testing.T) {
		request := newPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatusCode(t, response.Code, http.StatusAccepted)

		if len(store.WinCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.WinCalls), 1)
		}
	})

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatusCode(t, response.Code, http.StatusAccepted)

		if store.WinCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.WinCalls[0], player)
		}
	})
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := poker.CreateTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := poker.NewFileSystemPlayerStore(database)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v", err)
	}

	server := mustMakePlayerServer(t, store, &poker.GameSpy{})
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		poker.AssertResponseStatusCode(t, response.Code, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		poker.AssertResponseStatusCode(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := poker.League{
			{"Pepper", 3},
		}
		poker.AssertLeague(t, got, want)
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

	database, cleanDatabase := poker.CreateTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := poker.NewFileSystemPlayerStore(database)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v", err)
	}

	server := mustMakePlayerServer(t, store, &poker.GameSpy{})

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
			poker.AssertResponseStatusCode(t, response.Code, http.StatusOK)

			poker.AssertResponseBody(t, response.Body.String(), test.Result)

		})
	}

}

func TestLeague(t *testing.T) {

	t.Run("it returns 200 0on /league", func(t *testing.T) {
		store := poker.StubPlayerStore{}
		server := mustMakePlayerServer(t, &store, &poker.GameSpy{})

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got poker.League

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
		}

		poker.AssertResponseStatusCode(t, response.Code, http.StatusOK)
	})

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := poker.League{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := poker.StubPlayerStore{nil, nil, wantedLeague}
		server := mustMakePlayerServer(t, &store, &poker.GameSpy{})

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		poker.AssertResponseStatusCode(t, response.Code, http.StatusOK)
		poker.AssertLeague(t, got, wantedLeague)

		poker.AssertContentType(t, response, poker.JsonContentType)

	})
}

func TestGame(t *testing.T) {
	tenMs := 10 * time.Millisecond
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &poker.StubPlayerStore{}, &poker.GameSpy{})

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseStatusCode(t, response.Code, http.StatusOK)
	})

	t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
		store := &poker.StubPlayerStore{}
		game := &poker.GameSpy{}
		winner := "Ruth"
		server := httptest.NewServer(mustMakePlayerServer(t, store, game))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws := mustDialWS(t, wsURL)
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		time.Sleep(tenMs)

		poker.AssertFinishCalledWith(t, game, winner)
	})

	t.Run("start a game with 3 players and declare Ruth the winner", func(t *testing.T) {
		game := &poker.GameSpy{}
		winner := "Ruth"
		server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
		defer server.Close()

		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		time.Sleep(tenMs)
		poker.AssertGameStartedWith(t, game, 3)
		poker.AssertFinishCalledWith(t, game, winner)
	})

	t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		winner := "Ruth"

		game := &poker.GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		time.Sleep(tenMs)

		poker.AssertGameStartedWith(t, game, 3)
		poker.AssertFinishCalledWith(t, game, winner)
		within(t, tenMs, func() { poker.AssertWebsocketGotMsg(t, ws, wantedBlindAlert) })
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

func newGameRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return req
}

func getLeagueFromResponse(t *testing.T, body io.Reader) (league poker.League) {
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

func mustMakePlayerServer(t *testing.T, store poker.PlayerStore, game poker.Game) *poker.PlayerServer {
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}

func writeWSMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func within(t *testing.T, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}