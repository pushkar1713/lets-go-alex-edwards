package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	mux := http.NewServeMux()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	app := &application{
		logger: logger,
	}

	addr := flag.String("addr", ":4000", "HTTP Network address")
	flag.Parse()

	fileServer := http.FileServer(http.Dir("./ui/static"))

	// logger.Info("this is a test log", "method", "put")

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.createSnippet)
	mux.HandleFunc("POST /snippet/save", app.saveSnippet)
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// log.Printf("starting on port %s", *addr)
	logger.Info("addr", "addr", *addr)

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
