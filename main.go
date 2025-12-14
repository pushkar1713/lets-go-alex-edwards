package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Server", "Go from Mac")
	w.Write([]byte("hello from puhskar"))
}

func snippetView(w http.ResponseWriter, r *http.Request){
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if(err != nil || id < "1"){
		http.NotFound(w, r)
		return
	}
	msg := fmt.Sprintf("this is the snippet for the id %d...", idInt)
	w.Write([]byte(msg))
}

func createSnippet(w http.ResponseWriter, r *http.Request){
	// w.Write([]byte("create a new snippet here"))
	io.WriteString(w, "create a new snippet here")
}

func saveSnippet(w http.ResponseWriter, r *http.Request){
	data := r.Body
	message := fmt.Sprintf("this is the response body we are saving this %s", data)
	// w.WriteHeader(201)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(message))
}

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", createSnippet)
	mux.HandleFunc("POST /snippet/save", saveSnippet)

	log.Print("starting on port 4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}