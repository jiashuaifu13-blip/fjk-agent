package task

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
	now := time.Now().UTC()

	id, err :=
		allocateID(
			ctx,
			m.repository,
			now,
		)

	if err != nil {
		return nil, err
	}

	item := &Task{
		ID:          id,
		SessionID:   sessionID,
		Title:       title,
		Description: description,
		Status:      StatusCreated,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err :=
		m.repository.Create(
			ctx,
			item,
		); err != nil {

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
				"title":       item.Title,
				"description": item.Description,
				"status":      item.Status,
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
	item, err :=
		m.repository.Get(
			ctx,
			id,
		)

	if err != nil {
		return nil, err
	}

	if err := item.Transition(
		next,
	); err != nil {
		return nil, err
	}

	if err :=
		m.repository.Update(
			ctx,
			item,
		); err != nil {

		return nil, err
	}

	eventType :=
		eventForStatus(
			next,
		)

	if eventType != "" {
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
	return m.repository.Get(
		ctx,
		id,
	)
}

func (m *Manager) List(
	ctx context.Context,
) ([]*Task, error) {
	return m.repository.List(
		ctx,
	)
}

func eventForStatus(
	status Status,
) event.Type {
	switch status {
	case StatusQueued:
		return event.TypeTaskQueued

	case StatusRunning:
		return event.TypeTaskStarted

	case StatusWaitingApproval:
		return event.TypeTaskWaitingApproval

	case StatusPaused:
		return event.TypeTaskPaused

	case StatusCompleted:
		return event.TypeTaskCompleted

	case StatusFailed:
		return event.TypeTaskFailed

	case StatusCancelled:
		return event.TypeTaskCancelled

	default:
		return ""
	}
}
