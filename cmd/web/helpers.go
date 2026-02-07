package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serveError(w http.ResponseWriter, r *http.Request, err error) {
	var method = r.Method
	var uri = r.URL.RequestURI()
	var trace = string(debug.Stack())

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("this template doest not exist", page)
		app.serveError(w, r, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serveError(w, r, err)
	}
}
