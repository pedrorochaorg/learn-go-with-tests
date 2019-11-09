package helpers

import (
	"github.com/pedrorochaorg/learn-go-with-tests/examples/httpserver/objects"
	"net/http/httptest"
	"reflect"
	"testing"
)

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func AssertResponseStatusCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response status code is wrong, got %d want %d", got, want)
	}
}

func AssertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("score isn't as it was expected, got %d want %d", got, want)
	}
}

func AssertLeague(t *testing.T, got, want objects.League) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}