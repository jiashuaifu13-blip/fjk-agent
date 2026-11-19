package task

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks: make(map[string]*Task),
	}
}

func (r *MemoryRepository) Create(
	ctx context.Context,
	item *Task,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[item.ID] = cloneTask(item)

	return nil
}

func (r *MemoryRepository) Update(
	ctx context.Context,
	item *Task,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[item.ID]; !ok {
		return ErrNotFound
	}

	r.tasks[item.ID] = cloneTask(item)

	return nil
}

func (r *MemoryRepository) Get(
	ctx context.Context,
	id string,
) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.tasks[id]

	if !ok {
		return nil, ErrNotFound
	}

	return cloneTask(item), nil
}

func (r *MemoryRepository) List(
	ctx context.Context,
) ([]*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*Task, 0, len(r.tasks))

	for _, item := range r.tasks {
		items = append(items, cloneTask(item))
	}

	return items, nil
}

func cloneTask(item *Task) *Task {
	if item == nil {
		return nil
	}

	copy := *item

	return &copy
}
