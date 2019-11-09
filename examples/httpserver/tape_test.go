package main

import (
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver/helpers"
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver/objects"
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := helpers.CreateTempFile(t, "12345")
	defer clean()

	tape := objects.NewTape(file)

	tape.Write([]byte("abc"))

	file.Seek(0,0)
	newFileContents, _ := ioutil.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}