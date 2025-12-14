package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", createSnippet)
	mux.HandleFunc("POST /snippet/save", saveSnippet)
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	log.Print("starting on port 4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
