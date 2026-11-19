package model

import (
	"fmt"
	"sync"
)

type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

func (r *Registry) Register(
	provider Provider,
) error {
	if provider == nil {
		return fmt.Errorf(
			"provider cannot be nil",
		)
	}

	name := provider.Name()

	if name == "" {
		return fmt.Errorf(
			"provider name cannot be empty",
		)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[name]; exists {
		return fmt.Errorf(
			"provider already registered: %s",
			name,
		)
	}

	r.providers[name] = provider

	return nil
}

func (r *Registry) Get(
	name string,
) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok := r.providers[name]

	if !ok {
		return nil, fmt.Errorf(
			"provider not found: %s",
			name,
		)
	}

	return provider, nil
}
