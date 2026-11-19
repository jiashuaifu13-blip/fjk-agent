package conversation

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu            sync.RWMutex
	conversations map[string]*Conversation
	messages      map[string][]*Message
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		conversations: make(map[string]*Conversation),
		messages:      make(map[string][]*Message),
	}
}

func (r *MemoryRepository) CreateConversation(
	ctx context.Context,
	item *Conversation,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *item
	r.conversations[item.ID] = &copy

	return nil
}

func (r *MemoryRepository) UpdateConversation(
	ctx context.Context,
	item *Conversation,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.conversations[item.ID]; !ok {
		return ErrConversationNotFound
	}

	copy := *item
	r.conversations[item.ID] = &copy

	return nil
}

func (r *MemoryRepository) GetConversation(
	ctx context.Context,
	id string,
) (*Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.conversations[id]

	if !ok {
		return nil, ErrConversationNotFound
	}

	copy := *item

	return &copy, nil
}

func (r *MemoryRepository) ListConversations(
	ctx context.Context,
) ([]*Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(
		[]*Conversation,
		0,
		len(r.conversations),
	)

	for _, item := range r.conversations {
		copy := *item

		result = append(
			result,
			&copy,
		)
	}

	return result, nil
}

func (r *MemoryRepository) CreateMessage(
	ctx context.Context,
	message *Message,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *message

	r.messages[message.ConversationID] =
		append(
			r.messages[message.ConversationID],
			&copy,
		)

	return nil
}

func (r *MemoryRepository) ListMessages(
	ctx context.Context,
	conversationID string,
) ([]*Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	source := r.messages[conversationID]

	result := make(
		[]*Message,
		0,
		len(source),
	)

	for _, message := range source {
		copy := *message

		result = append(
			result,
			&copy,
		)
	}

	return result, nil
}
