package task

import (
	"context"
	"testing"

	"github.com/fjk-dev/fjk-agent/core/event"
)

func TestManagerCreateAndTransition(t *testing.T) {
	ctx := context.Background()

	bus := event.NewBus()
	repository := NewMemoryRepository()
	manager := NewManager(
		repository,
		bus,
	)

	createdReceived := false
	queuedReceived := false

	bus.Subscribe(
		event.TypeTaskCreated,
		func(ctx context.Context, evt event.Event) {
			createdReceived = true
		},
	)

	bus.Subscribe(
		event.TypeTaskQueued,
		func(ctx context.Context, evt event.Event) {
			queuedReceived = true
		},
	)

	item, err := manager.Create(
		ctx,
		"session_test",
		"Test Task",
		"Task Manager test",
	)

	if err != nil {
		t.Fatal(err)
	}

	if !createdReceived {
		t.Fatal(
			"task.created event was not received",
		)
	}

	item, err = manager.Transition(
		ctx,
		item.ID,
		StatusQueued,
	)

	if err != nil {
		t.Fatal(err)
	}

	if item.Status != StatusQueued {
		t.Fatalf(
			"expected QUEUED, got %s",
			item.Status,
		)
	}

	if !queuedReceived {
		t.Fatal(
			"task.queued event was not received",
		)
	}
}

func TestManagerRejectsInvalidTransition(t *testing.T) {
	ctx := context.Background()

	manager := NewManager(
		NewMemoryRepository(),
		event.NewBus(),
	)

	item, err := manager.Create(
		ctx,
		"",
		"Invalid Transition",
		"",
	)

	if err != nil {
		t.Fatal(err)
	}

	_, err = manager.Transition(
		ctx,
		item.ID,
		StatusCompleted,
	)

	if err == nil {
		t.Fatal(
			"expected invalid transition error",
		)
	}
}
