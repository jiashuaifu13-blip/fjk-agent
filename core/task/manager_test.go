package task

import (
	"context"
	"testing"

	"github.com/fjk-dev/fjk-agent/core/event"
)

func TestManagerCreateAndTransition(t *testing.T) {
	bus := event.NewBus()
	repository := NewMemoryRepository()

	manager := NewManager(
		repository,
		bus,
		"device_test",
	)

	createdEvents := 0
	startedEvents := 0

	bus.Subscribe(
		event.TypeTaskCreated,
		func(ctx context.Context, evt event.Event) {
			createdEvents++
		},
	)

	bus.Subscribe(
		event.TypeTaskStarted,
		func(ctx context.Context, evt event.Event) {
			startedEvents++
		},
	)

	item, err := manager.Create(
		context.Background(),
		"session_test",
		"Test task",
		"",
	)

	if err != nil {
		t.Fatalf(
			"create task failed: %v",
			err,
		)
	}

	if createdEvents != 1 {
		t.Fatalf(
			"expected 1 created event, got %d",
			createdEvents,
		)
	}

	if _, err := manager.Transition(
		context.Background(),
		item.ID,
		StatusQueued,
	); err != nil {
		t.Fatalf(
			"queue task failed: %v",
			err,
		)
	}

	if _, err := manager.Transition(
		context.Background(),
		item.ID,
		StatusRunning,
	); err != nil {
		t.Fatalf(
			"start task failed: %v",
			err,
		)
	}

	if startedEvents != 1 {
		t.Fatalf(
			"expected 1 started event, got %d",
			startedEvents,
		)
	}
}
