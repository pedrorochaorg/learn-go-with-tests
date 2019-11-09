package stores

import "github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver/objects"

type StubPlayerStore struct {
	Scores map[string]int
	WinCalls []string
	League objects.League
}

func (s *StubPlayerStore) GetPlayerScore(name string) (score int) {
	score = s.Scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.WinCalls = append(s.WinCalls, name)
}
func (s *StubPlayerStore) GetLeague() objects.League {
	return s.League
}
