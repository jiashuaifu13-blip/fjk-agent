package conversation

import (
	"context"
	"sync"
	"time"
)

type Manager struct {
	mu            sync.RWMutex
	conversations map[string]*Conversation
	messages      map[string][]*Message
}

func NewManager() *Manager {
	return &Manager{
		conversations: make(
			map[string]*Conversation,
		),
		messages: make(
			map[string][]*Message,
		),
	}
}

func (m *Manager) Create(
	ctx context.Context,
	title string,
) *Conversation {
	now := time.Now().UTC()

	item := &Conversation{
		ID:        NewConversationID(),
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.conversations[item.ID] = item

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
		ID:             NewMessageID(),
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		CreatedAt:      now,
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages[conversationID] = append(
		m.messages[conversationID],
		message,
	)

	if conversation, ok :=
		m.conversations[conversationID]; ok {
		conversation.UpdatedAt = now
	}

	return message
}

func (m *Manager) Messages(
	ctx context.Context,
	conversationID string,
) []*Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	source := m.messages[conversationID]

	result := make(
		[]*Message,
		len(source),
	)

	copy(result, source)

	return result
}
