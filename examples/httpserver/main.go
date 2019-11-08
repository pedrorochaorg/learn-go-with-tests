package main

import (
	"log"
	"net/http"
)


func main() {

	server := &PlayerServer{NewInMemoryPlayerStore()}

	// start's the web server on port 5000 and add's the previous created handler to it.
	// if the operation fails, the error message returned by the command ListenAndServe of the http package
	// must be printed out to the console
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
