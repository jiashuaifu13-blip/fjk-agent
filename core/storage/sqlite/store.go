package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(
		filepath.Dir(path),
		0o755,
	); err != nil {
		return nil, fmt.Errorf(
			"create database directory: %w",
			err,
		)
	}

	db, err := sql.Open(
		"sqlite",
		path,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"open sqlite database: %w",
			err,
		)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf(
			"ping sqlite database: %w",
			err,
		)
	}

	store := &Store{
		db: db,
	}

	if err := store.configure(); err != nil {
		_ = db.Close()

		return nil, err
	}

	return store, nil
}

func (s *Store) configure() error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA busy_timeout = 5000;",
	}

	for _, pragma := range pragmas {
		if _, err := s.db.Exec(pragma); err != nil {
			return fmt.Errorf(
				"apply sqlite pragma: %w",
				err,
			)
		}
	}

	return nil
}

func (s *Store) DB() *sql.DB {
	return s.db
}

func (s *Store) Close() error {
	return s.db.Close()
}
