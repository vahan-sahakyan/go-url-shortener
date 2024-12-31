package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS url(
        id INTEGER PRIMARY KEY,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
  `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES(?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		// TODO: refactor this
		if sqliteErr, ok := err.(sqlite3.Error); ok && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias=?`)
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement %w", op, err)
	}

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare(`DELETE FROM url WHERE alias=?`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURLs() ([]map[string]string, error) {
	const op = "storage.sqlite.GetURLs"

	stmt, err := s.db.Prepare(`SELECT alias, url FROM url`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement %w", op, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			_ = fmt.Errorf("%s: close rows: %w", op, err)
		}
	}(rows)

	var urls []map[string]string
	for rows.Next() {
		var alias string
		var url string
		if err := rows.Scan(&alias, &url); err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		fmt.Println(alias, url)
		urls = append(urls, map[string]string{
			"alias":    alias,
			"url":      url,
			"details":  "http://localhost:8082/url/" + alias,
			"redirect": "http://localhost:8082/" + alias,
		})
	}
	return urls, nil
}
