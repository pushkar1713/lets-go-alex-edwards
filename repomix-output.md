This file is a merged representation of the entire codebase, combined into a single document by Repomix.

# File Summary

## Purpose
This file contains a packed representation of the entire repository's contents.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Repository files (if enabled)
5. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Files are sorted by Git change count (files with more changes are at the bottom)

# Directory Structure
```
html/
  pages/
    home.html
  partials/
    nav.html
  base.html
  view.html
models/
  errors.go
  snippets.go
static/
  css/
    main.css
  js/
    main.js
web/
  handlers.go
  helpers.go
  main.go
  routes.go
  templates.go
```

# Files

## File: web/handlers.go
```go
package main

import (
	"errors"
	"fmt"
	"html/template"
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

	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%v\n", snippet)
	// }

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	data := templateData{
		Snippets: snippets,
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serveError(w, r, err)
		return
	}

	// err = ts.Execute(w, nil)
	err = ts.ExecuteTemplate(w, "base", data)
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

	snippets, err := app.snippets.Get(idInt)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serveError(w, r, err)
		}
		return
	}
	// fmt.Fprintf(w, "%+v", snippets)
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
		"./ui/html/view.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serveError(w, r, err)
		return
	}

	data := templateData{
		Snippet: snippets,
	}

	// err = ts.Execute(w, nil)
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serveError(w, r, err)
		return
	}
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
```

## File: web/helpers.go
```go
package main

import (
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
```

## File: web/main.go
```go
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"snippetbox.pushkar1713.dev/internal/models"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	addr := flag.String("addr", ":4000", "HTTP Network address")
	connString := flag.String("connString", os.Getenv("DATABASE_URL"), "your psql connection string")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	db, db_err := openDB(*connString)
	if db_err != nil {
		logger.Error(db_err.Error())
		os.Exit(1)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	fmt.Println(*connString)

	defer db.Close()

	// logger.Info("this is a test log", "method", "put")

	// log.Printf("starting on port %s", *addr)
	logger.Info("addr", "addr", *addr)

	err := http.ListenAndServe(*addr, app.routes())
	log.Fatal(err)
}

func openDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
```

## File: web/routes.go
```go
package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileserver := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("POST /snippet/create", app.createSnippet)
	mux.HandleFunc("POST /snippet/save", app.saveSnippet)

	return mux
}
```

## File: web/templates.go
```go
package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.pushkar1713.dev/internal/models"
)

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.html")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"./ui/html/base.html",
			"./ui/html/view.html",
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil

}
```

## File: models/errors.go
```go
package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models : no matching records found")
```

## File: models/snippets.go
```go
package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES ($1, $2, NOW(), NOW() + ($3 * INTERVAL '1 day'))`

	_, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// not supported by psql
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }

	// return int(id), nil
	return 1, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, expires from snippets
	WHERE expires > NOW() and id = $1`

	row := m.DB.QueryRow(stmt, id)

	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, expires from snippets WHERE expires > NOW() ORDER BY id desc limit 10`

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil

}
```

## File: html/pages/home.html
```html
{{define "title"}}Home{{end}} {{define "main"}}
<h2>Latest Snippets</h2>
{{if .Snippets}}
<table>
  <tr>
    <th>Title</th>
    <th>Created</th>
    <th>ID</th>
  </tr>
  {{range .Snippets}}
  <tr>
    <td><a href="/snippet/view/{{.ID}}">{{.Title}}</a></td>
    <td>{{.Created}}</td>
    <td>#{{.ID}}</td>
  </tr>
  {{end}}
</table>
{{else}}
<p>There's nothing to see here yet!</p>
{{end}} {{end}}
```

## File: html/partials/nav.html
```html
{{define "nav"}}
<nav>
  <a href="/">Home</a>
</nav>
{{end}}
```

## File: html/base.html
```html
{{ define "base" }}
  <!doctype html>
  <html lang="en">
    <head>
      <meta charset="utf-8" />
      <title>{{ template "title" . }} - Snippetbox</title>
      <!-- Link to the CSS stylesheet and favicon -->
      <link rel="stylesheet" href="/static/css/main.css" />
      <link
        rel="shortcut icon"
        href="/static/img/favicon.ico"
        type="image/x-icon"
      />
      <!-- Also link to some fonts hosted by Google -->
      <link
        rel="stylesheet"
        href="https://fonts.googleapis.com/css?family=Ubuntu+Mono:400,700"
      />
    </head>
    <body>
      <header>
        <h1><a href="/">Snippetbox</a></h1>
      </header>
      {{ template "nav" . }}
      <main>
        {{ template "main" . }}
      </main>
      <footer>Powered by <a href="https://golang.org/">Go</a></footer>
      <!-- And include the JavaScript file -->
      <script src="/static/js/main.js" type="text/javascript"></script>
    </body>
  </html>
{{ end }}
```

## File: html/view.html
```html
{{define "title"}}Snippet #{{.Snippet.ID}}{{end}} {{define "main"}} {{with
.Snippet}}
<div class="snippet">
  <div class="metadata">
    <strong>{{.Title}}</strong>
    <span>#{{.ID}}</span>
  </div>
  <pre><code>{{.Content}}</code></pre>
  <div class="metadata">
    <time>Created: {{.Created}}</time>
    <time>Expires: {{.Expires}}</time>
  </div>
</div>
{{end}} {{end}}
```

## File: static/css/main.css
```css
* {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
    font-size: 18px;
    font-family: "Ubuntu Mono", monospace;
}

html, body {
    height: 100%;
}

body {
    line-height: 1.5;
    background-color: #F1F3F6;
    color: #34495E;
    overflow-y: scroll;
}

header, nav, main, footer {
    padding: 2px calc((100% - 800px) / 2) 0;
}

main {
    margin-top: 54px;
    margin-bottom: 54px;
    min-height: calc(100vh - 345px);
    overflow: auto;
}

h1 a {
    font-size: 36px;
    font-weight: bold;
    background-image: url("/static/img/logo.png");
    background-repeat: no-repeat;
    background-position: 0px 0px;
    height: 36px;
    padding-left: 50px;
    position: relative;
}

h1 a:hover {
    text-decoration: none;
    color: #34495E;
}

h2 {
    font-size: 22px;
    margin-bottom: 36px;
    position: relative;
    top: -9px;
}

a {
    color: #62CB31;
    text-decoration: none;
}

a:hover {
    color: #4EB722;
    text-decoration: underline;
}

textarea, input:not([type="submit"]) {
    font-size: 18px;
    font-family: "Ubuntu Mono", monospace;
}

header {
    background-image: -webkit-linear-gradient(left, #34495e, #34495e 25%, #9b59b6 25%, #9b59b6 35%, #3498db 35%, #3498db 45%, #62cb31 45%, #62cb31 55%, #ffb606 55%, #ffb606 65%, #e67e22 65%, #e67e22 75%, #e74c3c 85%, #e74c3c 85%, #c0392b 85%, #c0392b 100%);
    background-image: -moz-linear-gradient(left, #34495e, #34495e 25%, #9b59b6 25%, #9b59b6 35%, #3498db 35%, #3498db 45%, #62cb31 45%, #62cb31 55%, #ffb606 55%, #ffb606 65%, #e67e22 65%, #e67e22 75%, #e74c3c 85%, #e74c3c 85%, #c0392b 85%, #c0392b 100%);
    background-image: -ms-linear-gradient(left, #34495e, #34495e 25%, #9b59b6 25%, #9b59b6 35%, #3498db 35%, #3498db 45%, #62cb31 45%, #62cb31 55%, #ffb606 55%, #ffb606 65%, #e67e22 65%, #e67e22 75%, #e74c3c 85%, #e74c3c 85%, #c0392b 85%, #c0392b 100%);
    background-image: linear-gradient(to right, #34495e, #34495e 25%, #9b59b6 25%, #9b59b6 35%, #3498db 35%, #3498db 45%, #62cb31 45%, #62cb31 55%, #ffb606 55%, #ffb606 65%, #e67e22 65%, #e67e22 75%, #e74c3c 85%, #e74c3c 85%, #c0392b 85%, #c0392b 100%);
    background-size: 100% 6px;
    background-repeat: no-repeat;
    border-bottom: 1px solid #E4E5E7;
    overflow: auto;
    padding-top: 33px;
    padding-bottom: 27px;
    text-align: center;
}

header a {
    color: #34495E;
    text-decoration: none;
}

nav {
    border-bottom: 1px solid #E4E5E7;
    padding-top: 17px;
    padding-bottom: 15px;
    background: #F7F9FA;
    height: 60px;
    color: #6A6C6F;
}

nav a {
    margin-right: 1.5em;
    display: inline-block;
}

nav form {
    display: inline-block;
    margin-left: 1.5em;
}

nav div {
    width: 50%;
    float: left;
}

nav div:last-child {
    text-align: right;
}

nav div:last-child a {
    margin-left: 1.5em;
    margin-right: 0;
}

nav a.live {
    color: #34495E;
    cursor: default;
}

nav a.live:hover {
    text-decoration: none;
}

nav a.live:after {
    content: '';
    display: block;
    position: relative;
    left: calc(50% - 7px);
    top: 9px;
    width: 14px;
    height: 14px;
    background: #F7F9FA;
    border-left: 1px solid #E4E5E7;
    border-bottom: 1px solid #E4E5E7;
    -moz-transform: rotate(45deg);
    -webkit-transform: rotate(-45deg);
}

a.button, input[type="submit"] {
    background-color: #62CB31;
    border-radius: 3px;
    color: #FFFFFF;
    padding: 18px 27px;
    border: none;
    display: inline-block;
    margin-top: 18px;
    font-weight: 700;
}

a.button:hover, input[type="submit"]:hover {
    background-color: #4EB722;
    color: #FFFFFF;
    cursor: pointer;
    text-decoration: none;
}

form div {
    margin-bottom: 18px;
}

form div:last-child {
    border-top: 1px dashed #E4E5E7;
}

form input[type="radio"] {
    margin-left: 18px;
}

form input[type="text"], form input[type="password"], form input[type="email"] {
    padding: 0.75em 18px;
    width: 100%;
}

form input[type=text], form input[type="password"], form input[type="email"], textarea {
    color: #6A6C6F;
    background: #FFFFFF;
    border: 1px solid #E4E5E7;
    border-radius: 3px;
}

form label {
    display: inline-block;
    margin-bottom: 9px;
}

.error {
    color: #C0392B;
    font-weight: bold;
    display: block;
}

.error + textarea, .error + input {
    border-color: #C0392B !important;
    border-width: 2px !important;
}

textarea {
    padding: 18px;
    width: 100%;
    height: 266px;
}

button {
    background: none;
    padding: 0;
    border: none;
    color: #62CB31;
    text-decoration: none;
}

button:hover {
    color: #4EB722;
    text-decoration: underline;
    cursor: pointer;
}

.snippet {
    background-color: #FFFFFF;
    border: 1px solid #E4E5E7;
    border-radius: 3px;
}

.snippet pre {
    padding: 18px;
    border-top: 1px solid #E4E5E7;
    border-bottom: 1px solid #E4E5E7;
}

.snippet .metadata {
    background-color: #F7F9FA;
    color: #6A6C6F;
    padding: 0.75em 18px;
    overflow: auto;
}

.snippet .metadata span {
    float: right;
}

.snippet .metadata strong {
    color: #34495E;
}

.snippet .metadata time {
    display: inline-block;
}

.snippet .metadata time:first-child {
    float: left;
}

.snippet .metadata time:last-child {
    float: right;
}

div.flash {
    color: #FFFFFF;
    font-weight: bold;
    background-color: #34495E;
    padding: 18px;
    margin-bottom: 36px;
    text-align: center;
}

div.error {
    color: #FFFFFF;
    background-color: #C0392B;
    padding: 18px;
    margin-bottom: 36px;
    font-weight: bold;
    text-align: center;
}

table {
    background: white;
    border: 1px solid #E4E5E7;
    border-collapse: collapse;
    width: 100%;
}

td, th {
    text-align: left;
    padding: 9px 18px;
}

th:last-child, td:last-child {
    text-align: right;
    color: #6A6C6F;
}

tr {
    border-bottom: 1px solid #E4E5E7;
}

tr:nth-child(2n) {
    background-color: #F7F9FA;
}

footer {
    border-top: 1px solid #E4E5E7;
    padding-top: 17px;
    padding-bottom: 15px;
    background: #F7F9FA;
    height: 60px;
    color: #6A6C6F;
    text-align: center;
}
```

## File: static/js/main.js
```javascript
var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}
```
