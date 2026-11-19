package task

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(
	db *sql.DB,
) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Create(
	ctx context.Context,
	item *Task,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
INSERT INTO tasks (
task_id,
session_id,
title,
description,
status,
created_at,
updated_at,
started_at,
completed_at,
failure_message
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`,
		item.ID,
		item.SessionID,
		item.Title,
		item.Description,
		string(item.Status),
		formatTime(item.CreatedAt),
		formatTime(item.UpdatedAt),
		formatTimePtr(item.StartedAt),
		formatTimePtr(item.CompletedAt),
		item.FailureMessage,
	)

	if err != nil {
		return fmt.Errorf(
			"insert task: %w",
			err,
		)
	}

	return nil
}

func (r *SQLiteRepository) Update(
	ctx context.Context,
	item *Task,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`
UPDATE tasks
SET
session_id = ?,
title = ?,
description = ?,
status = ?,
updated_at = ?,
started_at = ?,
completed_at = ?,
failure_message = ?
WHERE task_id = ?
`,
		item.SessionID,
		item.Title,
		item.Description,
		string(item.Status),
		formatTime(item.UpdatedAt),
		formatTimePtr(item.StartedAt),
		formatTimePtr(item.CompletedAt),
		item.FailureMessage,
		item.ID,
	)

	if err != nil {
		return fmt.Errorf(
			"update task: %w",
			err,
		)
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *SQLiteRepository) Get(
	ctx context.Context,
	id string,
) (*Task, error) {
	row := r.db.QueryRowContext(
		ctx,
		`
SELECT
task_id,
session_id,
title,
description,
status,
created_at,
updated_at,
started_at,
completed_at,
failure_message
FROM tasks
WHERE task_id = ?
`,
		id,
	)

	item, err := scanTask(row)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *SQLiteRepository) List(
	ctx context.Context,
) ([]*Task, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
SELECT
task_id,
session_id,
title,
description,
status,
created_at,
updated_at,
started_at,
completed_at,
failure_message
FROM tasks
ORDER BY created_at DESC
`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := make([]*Task, 0)

	for rows.Next() {
		item, err := scanTask(rows)

		if err != nil {
			return nil, err
		}

		items = append(
			items,
			item,
		)
	}

	return items, rows.Err()
}

type scanner interface {
	Scan(...any) error
}

func scanTask(
	row scanner,
) (*Task, error) {
	var (
		item        Task
		status      string
		createdAt   string
		updatedAt   string
		startedAt   sql.NullString
		completedAt sql.NullString
	)

	err := row.Scan(
		&item.ID,
		&item.SessionID,
		&item.Title,
		&item.Description,
		&status,
		&createdAt,
		&updatedAt,
		&startedAt,
		&completedAt,
		&item.FailureMessage,
	)

	if err != nil {
		return nil, err
	}

	item.Status = Status(status)

	var parseErr error

	item.CreatedAt, parseErr = time.Parse(
		time.RFC3339Nano,
		createdAt,
	)

	if parseErr != nil {
		return nil, parseErr
	}

	item.UpdatedAt, parseErr = time.Parse(
		time.RFC3339Nano,
		updatedAt,
	)

	if parseErr != nil {
		return nil, parseErr
	}

	if startedAt.Valid {
		value, err := time.Parse(
			time.RFC3339Nano,
			startedAt.String,
		)

		if err != nil {
			return nil, err
		}

		item.StartedAt = &value
	}

	if completedAt.Valid {
		value, err := time.Parse(
			time.RFC3339Nano,
			completedAt.String,
		)

		if err != nil {
			return nil, err
		}

		item.CompletedAt = &value
	}

	return &item, nil
}

func formatTime(
	value time.Time,
) string {
	return value.UTC().Format(
		time.RFC3339Nano,
	)
}

func formatTimePtr(
	value *time.Time,
) any {
	if value == nil {
		return nil
	}

	return formatTime(*value)
}
