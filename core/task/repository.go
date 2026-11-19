package task

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("task not found")

type Repository interface {
	Create(task *Task) error
	Update(task *Task) error
	Get(id string) (*Task, error)
	List() ([]*Task, error)
}

type MemoryRepository struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks: make(map[string]*Task),
	}
}

func (r *MemoryRepository) Create(item *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[item.ID] = cloneTask(item)

	return nil
}

func (r *MemoryRepository) Update(item *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[item.ID]; !exists {
		return ErrNotFound
	}

	r.tasks[item.ID] = cloneTask(item)

	return nil
}

func (r *MemoryRepository) Get(id string) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.tasks[id]
	if !exists {
		return nil, ErrNotFound
	}

	return cloneTask(item), nil
}

func (r *MemoryRepository) List() ([]*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*Task, 0, len(r.tasks))

	for _, item := range r.tasks {
		items = append(items, cloneTask(item))
	}

	return items, nil
}

func cloneTask(item *Task) *Task {
	copy := *item
	return &copy
}
