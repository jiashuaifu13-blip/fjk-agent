package sqlite

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func (s *Store) Migrate() error {
	entries, err := fs.ReadDir(
		migrationFiles,
		"migrations",
	)

	if err != nil {
		return fmt.Errorf(
			"read migrations: %w",
			err,
		)
	}

	sort.Slice(
		entries,
		func(i, j int) bool {
			return entries[i].Name() <
				entries[j].Name()
		},
	)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		version := strings.TrimSuffix(
			name,
			".sql",
		)

		applied, err := s.migrationApplied(
			version,
		)

		if err != nil {
			return err
		}

		if applied {
			continue
		}

		content, err := migrationFiles.ReadFile(
			"migrations/" + name,
		)

		if err != nil {
			return fmt.Errorf(
				"read migration %s: %w",
				name,
				err,
			)
		}

		tx, err := s.db.Begin()

		if err != nil {
			return fmt.Errorf(
				"begin migration %s: %w",
				name,
				err,
			)
		}

		if _, err := tx.Exec(
			string(content),
		); err != nil {
			_ = tx.Rollback()

			return fmt.Errorf(
				"execute migration %s: %w",
				name,
				err,
			)
		}

		_, err = tx.Exec(
			`
INSERT INTO schema_migrations (
version,
applied_at
)
VALUES (?, ?)
`,
			version,
			time.Now().UTC().Format(
				time.RFC3339Nano,
			),
		)

		if err != nil {
			_ = tx.Rollback()

			return fmt.Errorf(
				"record migration %s: %w",
				name,
				err,
			)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf(
				"commit migration %s: %w",
				name,
				err,
			)
		}
	}

	return nil
}

func (s *Store) migrationApplied(
	version string,
) (bool, error) {
	var count int

	err := s.db.QueryRow(
		`
SELECT COUNT(*)
FROM schema_migrations
WHERE version = ?
`,
		version,
	).Scan(&count)

	if err != nil {
		return false, nil
	}

	return count > 0, nil
}
