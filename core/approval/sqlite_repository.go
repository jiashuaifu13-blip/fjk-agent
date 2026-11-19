package approval

import (
	"context"
	"database/sql"
	"errors"
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
	item *Approval,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
INSERT INTO approvals (
approval_id,
task_id,
session_id,
tool_name,
scope_hash,
status,
created_at,
decided_at,
consumed_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, NULL, NULL)
`,
		item.ID,
		item.TaskID,
		item.SessionID,
		item.ToolName,
		item.ScopeHash,
		string(item.Status),
		formatTime(item.CreatedAt),
	)

	if err != nil {
		return fmt.Errorf(
			"insert approval: %w",
			err,
		)
	}

	return nil
}

func (r *SQLiteRepository) Get(
	ctx context.Context,
	id string,
) (*Approval, error) {
	var (
		item Approval

		status string

		createdAt string

		decidedAt sql.NullString

		consumedAt sql.NullString
	)

	err := r.db.QueryRowContext(
		ctx,
		`
SELECT
approval_id,
task_id,
session_id,
tool_name,
scope_hash,
status,
created_at,
decided_at,
consumed_at
FROM approvals
WHERE approval_id = ?
`,
		id,
	).Scan(
		&item.ID,
		&item.TaskID,
		&item.SessionID,
		&item.ToolName,
		&item.ScopeHash,
		&status,
		&createdAt,
		&decidedAt,
		&consumedAt,
	)

	if errors.Is(
		err,
		sql.ErrNoRows,
	) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	item.Status = Status(status)

	item.CreatedAt, err =
		time.Parse(
			time.RFC3339Nano,
			createdAt,
		)

	if err != nil {
		return nil, err
	}

	if decidedAt.Valid {
		value, parseErr :=
			time.Parse(
				time.RFC3339Nano,
				decidedAt.String,
			)

		if parseErr != nil {
			return nil, parseErr
		}

		item.DecidedAt = &value
	}

	if consumedAt.Valid {
		value, parseErr :=
			time.Parse(
				time.RFC3339Nano,
				consumedAt.String,
			)

		if parseErr != nil {
			return nil, parseErr
		}

		item.ConsumedAt = &value
	}

	return &item, nil
}

func (r *SQLiteRepository) Decide(
	ctx context.Context,
	id string,
	status Status,
) error {
	if status != StatusApproved &&
		status != StatusRejected {
		return ErrInvalidStatus
	}

	now := time.Now().UTC()

	result, err := r.db.ExecContext(
		ctx,
		`
UPDATE approvals
SET
status = ?,
decided_at = ?
WHERE
approval_id = ?
AND status = ?
`,
		string(status),
		formatTime(now),
		id,
		string(StatusPending),
	)

	if err != nil {
		return err
	}

	rows, err :=
		result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotPending
	}

	return nil
}

func (r *SQLiteRepository) Consume(
	ctx context.Context,
	id string,
	scopeHash string,
) error {
	now := time.Now().UTC()

	result, err := r.db.ExecContext(
		ctx,
		`
UPDATE approvals
SET
status = ?,
consumed_at = ?
WHERE
approval_id = ?
AND scope_hash = ?
AND status = ?
`,
		string(StatusConsumed),
		formatTime(now),
		id,
		scopeHash,
		string(StatusApproved),
	)

	if err != nil {
		return err
	}

	rows, err :=
		result.RowsAffected()

	if err != nil {
		return err
	}

	if rows > 0 {
		return nil
	}

	item, err :=
		r.Get(
			ctx,
			id,
		)

	if err != nil {
		return err
	}

	if item.ScopeHash != scopeHash {
		return ErrScopeMismatch
	}

	if item.Status == StatusConsumed {
		return ErrAlreadyConsumed
	}

	return ErrNotApproved
}

func formatTime(
	value time.Time,
) string {
	return value.UTC().Format(
		time.RFC3339Nano,
	)
}
