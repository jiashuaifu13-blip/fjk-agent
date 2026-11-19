package agent

import (
	"context"
	"testing"

	"github.com/fjk-dev/fjk-agent/core/event"
	"github.com/fjk-dev/fjk-agent/core/planner"
	"github.com/fjk-dev/fjk-agent/core/task"
)

func TestAgentRunEvents(
	t *testing.T,
) {
	ctx :=
		context.Background()

	bus :=
		event.NewBus()

	taskManager :=
		task.NewManager(
			task.NewMemoryRepository(),
			bus,
		)

	var received []event.Type

	handler :=
		func(
			ctx context.Context,
			evt event.Event,
		) {
			received = append(
				received,
				evt.EventType,
			)
		}

	bus.Subscribe(
		event.TypeAgentRunStarted,
		handler,
	)

	bus.Subscribe(
		event.TypeAgentPlanCreated,
		handler,
	)

	bus.Subscribe(
		event.TypeAgentStepStarted,
		handler,
	)

	bus.Subscribe(
		event.TypeAgentStepCompleted,
		handler,
	)

	bus.Subscribe(
		event.TypeAgentRunCompleted,
		handler,
	)

	runner :=
		NewRunnerWithEvents(
			taskManager,
			planner.NewStaticPlanner(),
			nil,
			bus,
		)

	run, err :=
		runner.Run(
			ctx,
			"session_event_test",
			"hello events",
		)

	if err != nil {
		t.Fatal(err)
	}

	if run.Status !=
		RunCompleted {
		t.Fatalf(
			"unexpected status: %s",
			run.Status,
		)
	}

	expected :=
		[]event.Type{
			event.TypeAgentRunStarted,
			event.TypeAgentPlanCreated,
			event.TypeAgentStepStarted,
			event.TypeAgentStepCompleted,
			event.TypeAgentRunCompleted,
		}

	if len(received) != len(expected) {
		t.Fatalf(
			"expected %d agent events, got %d: %v",
			len(expected),
			len(received),
			received,
		)
	}

	for index := range expected {

		if received[index] !=
			expected[index] {

			t.Fatalf(
				"event %d: expected %s, got %s",
				index,
				expected[index],
				received[index],
			)
		}
	}
}
