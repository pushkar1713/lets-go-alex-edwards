package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go from Mac")
	// w.Write([]byte("hello from puhskar"))

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serveError(w, r, err)
		return
	}
	// err = ts.Execute(w, nil)
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serveError(w, r, err)
		return
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || id < "1" {
		http.NotFound(w, r)
		return
	}
	msg := fmt.Sprintf("this is the snippet for the id %d...", idInt)
	w.Write([]byte(msg))
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("create a new snippet here"))
	io.WriteString(w, "create a new snippet here")
}

func (app *application) saveSnippet(w http.ResponseWriter, r *http.Request) {
	data := r.Body
	message := fmt.Sprintf("this is the response body we are saving this %s", data)
	// w.WriteHeader(201)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(message))
}
