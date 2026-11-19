package task

import (
	"context"
	"fmt"
	"time"
)

func (r *SQLiteRepository) NextTaskID(
	ctx context.Context,
	now time.Time,
) (string, error) {
	day := now.UTC().Format(
		"20060102",
	)

	var sequence int

	err := r.db.QueryRowContext(
		ctx,
		`
INSERT INTO task_sequences (
day,
value
)
VALUES (?, 1)
ON CONFLICT(day)
DO UPDATE SET
value = value + 1
RETURNING value
`,
		day,
	).Scan(
		&sequence,
	)

	if err != nil {
		return "", fmt.Errorf(
			"allocate task id: %w",
			err,
		)
	}

	return fmt.Sprintf(
		"task_%s_%03d",
		day,
		sequence,
	), nil
}
