package conversation

import (
	"context"
	"fmt"
	"time"
)

type Manager struct {
	repository Repository
}

func NewManager() *Manager {
	return NewManagerWithRepository(
		NewMemoryRepository(),
	)
}

func NewManagerWithRepository(
	repository Repository,
) *Manager {
	return &Manager{
		repository: repository,
	}
}

func (m *Manager) Create(
	ctx context.Context,
	title string,
) *Conversation {
	now := time.Now().UTC()

	item := &Conversation{
		ID: NewConversationID(),

		Title: title,

		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := m.repository.CreateConversation(
		ctx,
		item,
	); err != nil {
		return nil
	}

	return item
}

func (m *Manager) AddMessage(
	ctx context.Context,
	conversationID string,
	role Role,
	content string,
) *Message {
	now := time.Now().UTC()

	message := &Message{
		ID: NewMessageID(),

		ConversationID: conversationID,

		Role: role,

		Content: content,

		CreatedAt: now,
	}

	if err := m.repository.CreateMessage(
		ctx,
		message,
	); err != nil {
		return nil
	}

	conversation, err :=
		m.repository.GetConversation(
			ctx,
			conversationID,
		)

	if err == nil {
		conversation.UpdatedAt = now

		_ = m.repository.UpdateConversation(
			ctx,
			conversation,
		)
	}

	return message
}

func (m *Manager) Messages(
	ctx context.Context,
	conversationID string,
) []*Message {
	messages, err :=
		m.repository.ListMessages(
			ctx,
			conversationID,
		)

	if err != nil {
		return nil
	}

	return messages
}

func (m *Manager) Get(
	ctx context.Context,
	id string,
) (*Conversation, error) {
	item, err :=
		m.repository.GetConversation(
			ctx,
			id,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"get conversation: %w",
			err,
		)
	}

	return item, nil
}

func (m *Manager) List(
	ctx context.Context,
) ([]*Conversation, error) {
	return m.repository.ListConversations(
		ctx,
	)
}
