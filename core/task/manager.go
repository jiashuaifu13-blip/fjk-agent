package task

import (
	"context"
	"fmt"

	"github.com/fjk-dev/fjk-agent/core/event"
)

type Manager struct {
	repository Repository
	events     *event.Bus
}

func NewManager(
	repository Repository,
	events *event.Bus,
) *Manager {
	return &Manager{
		repository: repository,
		events:     events,
	}
}

func (m *Manager) Create(
	ctx context.Context,
	sessionID string,
	title string,
	description string,
) (*Task, error) {
	item := New(
		NewID(),
		sessionID,
		title,
		description,
	)

	if err := m.repository.Create(ctx, item); err != nil {
		return nil, fmt.Errorf(
			"create task: %w",
			err,
		)
	}

	m.events.Publish(
		ctx,
		event.New(
			event.TypeTaskCreated,
			event.Metadata{
				TaskID:    item.ID,
				SessionID: item.SessionID,
			},
			map[string]any{
				"title": item.Title,
			},
		),
	)

	return item, nil
}

func (m *Manager) Transition(
	ctx context.Context,
	id string,
	next Status,
) (*Task, error) {
	item, err := m.repository.Get(ctx, id)

	if err != nil {
		return nil, fmt.Errorf(
			"get task: %w",
			err,
		)
	}

	if err := item.Transition(next); err != nil {
		return nil, err
	}

	if err := m.repository.Update(ctx, item); err != nil {
		return nil, fmt.Errorf(
			"update task: %w",
			err,
		)
	}

	eventType, ok := eventTypeForStatus(next)

	if ok {
		m.events.Publish(
			ctx,
			event.New(
				eventType,
				event.Metadata{
					TaskID:    item.ID,
					SessionID: item.SessionID,
				},
				map[string]any{
					"status": item.Status,
				},
			),
		)
	}

	return item, nil
}

func (m *Manager) Get(
	ctx context.Context,
	id string,
) (*Task, error) {
	return m.repository.Get(ctx, id)
}

func (m *Manager) List(
	ctx context.Context,
) ([]*Task, error) {
	return m.repository.List(ctx)
}

func eventTypeForStatus(status Status) (event.Type, bool) {
	switch status {
	case StatusQueued:
		return event.TypeTaskQueued, true

	case StatusRunning:
		return event.TypeTaskStarted, true

	case StatusWaitingApproval:
		return event.TypeTaskWaitingApproval, true

	case StatusPaused:
		return event.TypeTaskPaused, true

	case StatusCompleted:
		return event.TypeTaskCompleted, true

	case StatusFailed:
		return event.TypeTaskFailed, true

	case StatusCancelled:
		return event.TypeTaskCancelled, true

	default:
		return "", false
	}
}
