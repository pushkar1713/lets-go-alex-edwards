package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"snippetbox.pushkar1713.dev/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go from Mac")
	// w.Write([]byte("hello from puhskar"))
	//

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serveError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	// files := []string{
	// 	"./ui/html/base.html",
	// 	"./ui/html/partials/nav.html",
	// 	"./ui/html/pages/home.html",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serveError(w, r, err)
	// 	return
	// }
	// // err = ts.Execute(w, nil)
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serveError(w, r, err)
	// 	return
	// }
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || id < "1" {
		http.NotFound(w, r)
		return
	}

	snippets, err := app.snippets.Get(idInt)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serveError(w, r, err)
		}
		return
	}
	fmt.Fprintf(w, "%+v", snippets)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	_, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serveError(w, r, err)
	}

	io.WriteString(w, "create a new snippet here")
}

func (app *application) saveSnippet(w http.ResponseWriter, r *http.Request) {
	data := r.Body
	message := fmt.Sprintf("this is the response body we are saving this %s", data)
	// w.WriteHeader(201)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(message))
}
