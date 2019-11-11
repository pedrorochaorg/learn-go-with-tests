package poker_test

import (
	poker "github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver"
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := poker.CreateTempFile(t, "12345")
	defer clean()

	tape := poker.NewTape(file)

	tape.Write([]byte("abc"))

	file.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
