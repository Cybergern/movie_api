package main

import (
	"log"
	"movies"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	finalHandler := http.HandlerFunc(movies.GetMovies)
	mux.Handle("/movies", movies.HeaderMethodCheck(finalHandler))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
