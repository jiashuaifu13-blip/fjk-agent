package approval

import (
	"context"
	"fmt"
	"time"

	"github.com/fjk-dev/fjk-agent/core/event"
)

type Manager struct {
	repository Repository
	events     *event.Bus
}

func NewManager(
	repository Repository,
) *Manager {
	return &Manager{
		repository: repository,
	}
}

func NewManagerWithEvents(
	repository Repository,
	events *event.Bus,
) *Manager {
	return &Manager{
		repository: repository,
		events:     events,
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

	m.publish(
		ctx,
		event.TypeApprovalRequested,
		item,
	)

	return item, nil
}

func (m *Manager) Approve(
	ctx context.Context,
	id string,
) error {
	if err := m.repository.Decide(
		ctx,
		id,
		StatusApproved,
	); err != nil {
		return err
	}

	item, err :=
		m.repository.Get(
			ctx,
			id,
		)

	if err != nil {
		return err
	}

	m.publish(
		ctx,
		event.TypeApprovalApproved,
		item,
	)

	return nil
}

func (m *Manager) Reject(
	ctx context.Context,
	id string,
) error {
	if err := m.repository.Decide(
		ctx,
		id,
		StatusRejected,
	); err != nil {
		return err
	}

	item, err :=
		m.repository.Get(
			ctx,
			id,
		)

	if err != nil {
		return err
	}

	m.publish(
		ctx,
		event.TypeApprovalRejected,
		item,
	)

	return nil
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

	if err := m.repository.Consume(
		ctx,
		id,
		scopeHash,
	); err != nil {
		return err
	}

	item, err :=
		m.repository.Get(
			ctx,
			id,
		)

	if err != nil {
		return err
	}

	m.publish(
		ctx,
		event.TypeApprovalConsumed,
		item,
	)

	return nil
}

func (m *Manager) publish(
	ctx context.Context,
	eventType event.Type,
	item *Approval,
) {
	if m.events == nil {
		return
	}

	m.events.Publish(
		ctx,
		event.New(
			eventType,
			event.Metadata{
				TaskID:    item.TaskID,
				SessionID: item.SessionID,
			},
			map[string]any{
				"approval_id": item.ID,
				"tool_name":   item.ToolName,
				"status":      item.Status,
				"scope_hash":  item.ScopeHash,
			},
		),
	)
}
