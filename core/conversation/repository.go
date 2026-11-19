package conversation

import (
	"context"
	"errors"
)

var (
	ErrConversationNotFound = errors.New(
		"conversation not found",
	)
	ErrMessageNotFound = errors.New(
		"message not found",
	)
)

type Repository interface {
	CreateConversation(
		context.Context,
		*Conversation,
	) error

	UpdateConversation(
		context.Context,
		*Conversation,
	) error

	GetConversation(
		context.Context,
		string,
	) (*Conversation, error)

	ListConversations(
		context.Context,
	) ([]*Conversation, error)

	CreateMessage(
		context.Context,
		*Message,
	) error

	ListMessages(
		context.Context,
		string,
	) ([]*Message, error)
}
