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
