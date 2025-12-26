package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"snippetbox.pushkar1713.dev/internal/models"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
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

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
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
