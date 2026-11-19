package approval

import (
	"context"
	"fmt"
	"time"
)

type Manager struct {
	repository Repository
}

func NewManager(
	repository Repository,
) *Manager {
	return &Manager{
		repository: repository,
	}
}

func (m *Manager) Request(
	ctx context.Context,
	taskID string,
	sessionID string,
	toolName string,
	arguments map[string]any,
) (*Approval, error) {
	scopeHash, err :=
		ScopeHash(
			taskID,
			sessionID,
			toolName,
			arguments,
		)

	if err != nil {
		return nil, err
	}

	item := &Approval{
		ID: NewID(),

		TaskID: taskID,

		SessionID: sessionID,

		ToolName: toolName,

		ScopeHash: scopeHash,

		Status: StatusPending,

		CreatedAt: time.Now().UTC(),
	}

	if err := m.repository.Create(
		ctx,
		item,
	); err != nil {
		return nil, fmt.Errorf(
			"create approval: %w",
			err,
		)
	}

	return item, nil
}

func (m *Manager) Approve(
	ctx context.Context,
	id string,
) error {
	return m.repository.Decide(
		ctx,
		id,
		StatusApproved,
	)
}

func (m *Manager) Reject(
	ctx context.Context,
	id string,
) error {
	return m.repository.Decide(
		ctx,
		id,
		StatusRejected,
	)
}

func (m *Manager) Get(
	ctx context.Context,
	id string,
) (*Approval, error) {
	return m.repository.Get(
		ctx,
		id,
	)
}

func (m *Manager) ValidateAndConsume(
	ctx context.Context,
	id string,
	taskID string,
	sessionID string,
	toolName string,
	arguments map[string]any,
) error {
	scopeHash, err :=
		ScopeHash(
			taskID,
			sessionID,
			toolName,
			arguments,
		)

	if err != nil {
		return err
	}

	return m.repository.Consume(
		ctx,
		id,
		scopeHash,
	)
}
