package event

import (
	"context"
	"testing"
)

func TestBusPublish(t *testing.T) {
	bus := NewBus()

	called := false

	bus.Subscribe(
		TypeTaskCreated,
		func(ctx context.Context, evt Event) {
			called = true

			if evt.EventType != TypeTaskCreated {
				t.Fatalf(
					"unexpected event type: %s",
					evt.EventType,
				)
			}
		},
	)

	bus.Publish(
		context.Background(),
		New(
			TypeTaskCreated,
			Metadata{
				TaskID: "task_test",
			},
			nil,
		),
	)

	if !called {
		t.Fatal("event handler was not called")
	}
}
