package event

import (
	"context"
	"testing"
)

func TestBusPublishSubscribe(t *testing.T) {
	bus := NewBus()

	received := false

	unsubscribe := bus.Subscribe(
		TypeTaskCreated,
		func(ctx context.Context, evt Event) {
			received = true

			if evt.EventType != TypeTaskCreated {
				t.Fatalf(
					"unexpected event type: %s",
					evt.EventType,
				)
			}
		},
	)

	defer unsubscribe()

	evt := New(
		TypeTaskCreated,
		Metadata{
			TaskID: "task_test_001",
		},
		nil,
	)

	bus.Publish(context.Background(), evt)

	if !received {
		t.Fatal("event was not received")
	}
}

func TestBusUnsubscribe(t *testing.T) {
	bus := NewBus()

	count := 0

	unsubscribe := bus.Subscribe(
		TypeTaskCreated,
		func(ctx context.Context, evt Event) {
			count++
		},
	)

	unsubscribe()

	bus.Publish(
		context.Background(),
		New(
			TypeTaskCreated,
			Metadata{},
			nil,
		),
	)

	if count != 0 {
		t.Fatalf(
			"expected 0 events, received %d",
			count,
		)
	}
}
