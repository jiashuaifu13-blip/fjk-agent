package approval

import (
	"context"
	"sync"
	"time"
)

type MemoryRepository struct {
	mu sync.RWMutex

	items map[string]*Approval
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: make(
			map[string]*Approval,
		),
	}
}

func (r *MemoryRepository) Create(
	ctx context.Context,
	item *Approval,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := cloneApproval(item)

	r.items[item.ID] = copy

	return nil
}

func (r *MemoryRepository) Get(
	ctx context.Context,
	id string,
) (*Approval, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]

	if !ok {
		return nil, ErrNotFound
	}

	return cloneApproval(item), nil
}

func (r *MemoryRepository) Decide(
	ctx context.Context,
	id string,
	status Status,
) error {
	if status != StatusApproved &&
		status != StatusRejected {
		return ErrInvalidStatus
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.items[id]

	if !ok {
		return ErrNotFound
	}

	if item.Status != StatusPending {
		return ErrNotPending
	}

	now := time.Now().UTC()

	item.Status = status
	item.DecidedAt = &now

	return nil
}

func (r *MemoryRepository) Consume(
	ctx context.Context,
	id string,
	scopeHash string,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.items[id]

	if !ok {
		return ErrNotFound
	}

	if item.ScopeHash != scopeHash {
		return ErrScopeMismatch
	}

	if item.Status == StatusConsumed {
		return ErrAlreadyConsumed
	}

	if item.Status != StatusApproved {
		return ErrNotApproved
	}

	now := time.Now().UTC()

	item.Status = StatusConsumed
	item.ConsumedAt = &now

	return nil
}

func cloneApproval(
	item *Approval,
) *Approval {
	copy := *item

	if item.DecidedAt != nil {
		value := *item.DecidedAt
		copy.DecidedAt = &value
	}

	if item.ConsumedAt != nil {
		value := *item.ConsumedAt
		copy.ConsumedAt = &value
	}

	return &copy
}
