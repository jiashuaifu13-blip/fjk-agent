package conversation

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

func (r *SQLiteRepository) CreateConversation(
	ctx context.Context,
	item *Conversation,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
INSERT INTO conversations (
conversation_id,
title,
created_at,
updated_at
)
VALUES (?, ?, ?, ?)
`,
		item.ID,
		item.Title,
		formatTime(item.CreatedAt),
		formatTime(item.UpdatedAt),
	)

	if err != nil {
		return fmt.Errorf(
			"insert conversation: %w",
			err,
		)
	}

	return nil
}

func (r *SQLiteRepository) UpdateConversation(
	ctx context.Context,
	item *Conversation,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`
UPDATE conversations
SET
title = ?,
updated_at = ?
WHERE conversation_id = ?
`,
		item.Title,
		formatTime(item.UpdatedAt),
		item.ID,
	)

	if err != nil {
		return fmt.Errorf(
			"update conversation: %w",
			err,
		)
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrConversationNotFound
	}

	return nil
}

func (r *SQLiteRepository) GetConversation(
	ctx context.Context,
	id string,
) (*Conversation, error) {
	var (
		item      Conversation
		createdAt string
		updatedAt string
	)

	err := r.db.QueryRowContext(
		ctx,
		`
SELECT
conversation_id,
title,
created_at,
updated_at
FROM conversations
WHERE conversation_id = ?
`,
		id,
	).Scan(
		&item.ID,
		&item.Title,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrConversationNotFound
	}

	if err != nil {
		return nil, err
	}

	var parseErr error

	item.CreatedAt, parseErr =
		time.Parse(
			time.RFC3339Nano,
			createdAt,
		)

	if parseErr != nil {
		return nil, parseErr
	}

	item.UpdatedAt, parseErr =
		time.Parse(
			time.RFC3339Nano,
			updatedAt,
		)

	if parseErr != nil {
		return nil, parseErr
	}

	return &item, nil
}

func (r *SQLiteRepository) ListConversations(
	ctx context.Context,
) ([]*Conversation, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
SELECT
conversation_id,
title,
created_at,
updated_at
FROM conversations
ORDER BY updated_at DESC
`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(
		[]*Conversation,
		0,
	)

	for rows.Next() {
		var (
			item      Conversation
			createdAt string
			updatedAt string
		)

		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		item.CreatedAt, err =
			time.Parse(
				time.RFC3339Nano,
				createdAt,
			)

		if err != nil {
			return nil, err
		}

		item.UpdatedAt, err =
			time.Parse(
				time.RFC3339Nano,
				updatedAt,
			)

		if err != nil {
			return nil, err
		}

		result = append(
			result,
			&item,
		)
	}

	return result, rows.Err()
}

func (r *SQLiteRepository) CreateMessage(
	ctx context.Context,
	message *Message,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
INSERT INTO messages (
message_id,
conversation_id,
role,
content,
created_at
)
VALUES (?, ?, ?, ?, ?)
`,
		message.ID,
		message.ConversationID,
		string(message.Role),
		message.Content,
		formatTime(message.CreatedAt),
	)

	if err != nil {
		return fmt.Errorf(
			"insert message: %w",
			err,
		)
	}

	return nil
}

func (r *SQLiteRepository) ListMessages(
	ctx context.Context,
	conversationID string,
) ([]*Message, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
SELECT
message_id,
conversation_id,
role,
content,
created_at
FROM messages
WHERE conversation_id = ?
ORDER BY created_at ASC
`,
		conversationID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(
		[]*Message,
		0,
	)

	for rows.Next() {
		var (
			item      Message
			role      string
			createdAt string
		)

		if err := rows.Scan(
			&item.ID,
			&item.ConversationID,
			&role,
			&item.Content,
			&createdAt,
		); err != nil {
			return nil, err
		}

		item.Role = Role(role)

		item.CreatedAt, err =
			time.Parse(
				time.RFC3339Nano,
				createdAt,
			)

		if err != nil {
			return nil, err
		}

		result = append(
			result,
			&item,
		)
	}

	return result, rows.Err()
}

func formatTime(
	value time.Time,
) string {
	return value.UTC().Format(
		time.RFC3339Nano,
	)
}
