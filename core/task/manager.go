package task

import (
	"context"
	"fmt"

	"github.com/fjk-dev/fjk-agent/core/event"
)

type Manager struct {
	repository Repository
	events     *event.Bus
	deviceID   string
}

func NewManager(
	repository Repository,
	events *event.Bus,
	deviceID string,
) *Manager {
	return &Manager{
		repository: repository,
		events:     events,
		deviceID:   deviceID,
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

	if err := m.repository.Create(item); err != nil {
		return nil, fmt.Errorf(
			"create task: %w",
			err,
		)
	}

	m.publish(
		ctx,
		item,
		event.TypeTaskCreated,
	)

	return item, nil
}

func (m *Manager) Transition(
	ctx context.Context,
	taskID string,
	next Status,
) (*Task, error) {
	item, err := m.repository.Get(taskID)
	if err != nil {
		return nil, fmt.Errorf(
			"get task: %w",
			err,
		)
	}

	if err := item.Transition(next); err != nil {
		return nil, err
	}

	if err := m.repository.Update(item); err != nil {
		return nil, fmt.Errorf(
			"update task: %w",
			err,
		)
	}

	m.publish(
		ctx,
		item,
		eventTypeForStatus(next),
	)

	return item, nil
}

func (m *Manager) Get(taskID string) (*Task, error) {
	return m.repository.Get(taskID)
}

func (m *Manager) List() ([]*Task, error) {
	return m.repository.List()
}

func (m *Manager) publish(
	ctx context.Context,
	item *Task,
	eventType event.Type,
) {
	if m.events == nil {
		return
	}

	m.events.Publish(
		ctx,
		event.New(
			eventType,
			event.Metadata{
				TaskID:    item.ID,
				SessionID: item.SessionID,
				DeviceID:  m.deviceID,
			},
			map[string]any{
				"status": item.Status,
				"title":  item.Title,
			},
		),
	)
}

func eventTypeForStatus(status Status) event.Type {
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
		return event.TypeTaskProgress
	}
}
