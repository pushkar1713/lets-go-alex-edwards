package models

import (
	"database/sql"
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
	return Snippet{}, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
