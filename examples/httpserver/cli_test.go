package poker_test

import (
	"bytes"
	poker "github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver"
	"strings"
	"testing"
	"time"
)

var dummySpyAlerter = &poker.SpyBlindAlerter{}
var dummyBlindAlerter = &poker.SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

func TestCLI(t *testing.T) {
	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &poker.GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("3", "Chris wins")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		poker.AssertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		poker.AssertGameStartedWith(t, game, 3)
		poker.AssertFinishCalledWith(t, game, "Chris")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &poker.GameSpy{}

		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		poker.AssertGameStartedWith(t, game, 8)
		poker.AssertFinishCalledWith(t, game, "Cleo")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		game := &poker.GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("pies")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		poker.AssertGameNotStarted(t, game)
		poker.AssertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})

	t.Run("it prints an error when we don't send wins after the user name", func(t *testing.T) {
		game := &poker.GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("8", "Loyd is a killer")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		poker.AssertGameStartedWith(t, game, 8)
		poker.AssertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerWinInputErrMsg)
	})
}

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(5, dummyStdOut)

		cases := []poker.ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 10 * time.Minute, Amount: 200},
			{At: 20 * time.Minute, Amount: 300},
			{At: 30 * time.Minute, Amount: 400},
			{At: 40 * time.Minute, Amount: 500},
			{At: 50 * time.Minute, Amount: 600},
			{At: 60 * time.Minute, Amount: 800},
			{At: 70 * time.Minute, Amount: 1000},
			{At: 80 * time.Minute, Amount: 2000},
			{At: 90 * time.Minute, Amount: 4000},
			{At: 100 * time.Minute, Amount: 8000},
		}

		poker.CheckSchedulingCases(t, blindAlerter.Alerts, cases)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(7, dummyStdOut)

		cases := []poker.ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 12 * time.Minute, Amount: 200},
			{At: 24 * time.Minute, Amount: 300},
			{At: 36 * time.Minute, Amount: 400},
			{At: 48 * time.Minute, Amount: 500},
			{At: 60 * time.Minute, Amount: 600},
			{At: 72 * time.Minute, Amount: 800},
			{At: 84 * time.Minute, Amount: 1000},
			{At: 96 * time.Minute, Amount: 2000},
			{At: 108 * time.Minute, Amount: 4000},
			{At: 120 * time.Minute, Amount: 8000},
		}

		poker.CheckSchedulingCases(t, blindAlerter.Alerts, cases)
	})

}

func TestGame_Finish(t *testing.T) {
	store := &poker.StubPlayerStore{}
	game := poker.NewTexasHoldem(dummyBlindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	poker.AssertPlayerWin(t, store, winner)
}

func userSends(messages ...string) *strings.Reader {
	message := strings.Join(messages, "\n")
	return strings.NewReader(message)

}
