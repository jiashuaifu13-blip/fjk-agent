package tool

import (
	"fmt"
	"sync"
)

type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Register(
	item Tool,
) error {
	if item == nil {
		return ErrInvalidToolDefinition
	}

	definition := item.Definition()

	if definition.Name == "" {
		return ErrInvalidToolDefinition
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[definition.Name]; exists {
		return fmt.Errorf(
			"%w: %s",
			ErrToolAlreadyRegistered,
			definition.Name,
		)
	}

	r.tools[definition.Name] = item

	return nil
}

func (r *Registry) Get(
	name string,
) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.tools[name]

	if !exists {
		return nil, fmt.Errorf(
			"%w: %s",
			ErrToolNotFound,
			name,
		)
	}

	return item, nil
}

func (r *Registry) Definitions() []Definition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(
		[]Definition,
		0,
		len(r.tools),
	)

	for _, item := range r.tools {
		result = append(
			result,
			item.Definition(),
		)
	}

	return result
}
