package agent

import (
	"context"
	"testing"

	"github.com/fjk-dev/fjk-agent/core/event"
	"github.com/fjk-dev/fjk-agent/core/planner"
	"github.com/fjk-dev/fjk-agent/core/task"
)

func TestStaticAgentRun(
	t *testing.T,
) {
	ctx := context.Background()

	taskManager :=
		task.NewManager(
			task.NewMemoryRepository(),
			event.NewBus(),
		)

	runner :=
		NewRunner(
			taskManager,
			planner.NewStaticPlanner(),
			nil,
		)

	run, err :=
		runner.Run(
			ctx,
			"session_test",
			"hello fjk agent",
		)

	if err != nil {
		t.Fatal(err)
	}

	if run.Status != RunCompleted {
		t.Fatalf(
			"expected COMPLETED, got %s",
			run.Status,
		)
	}

	if run.Output !=
		"hello fjk agent" {
		t.Fatalf(
			"unexpected output: %s",
			run.Output,
		)
	}

	stored, err :=
		taskManager.Get(
			ctx,
			run.TaskID,
		)

	if err != nil {
		t.Fatal(err)
	}

	if stored.Status !=
		task.StatusCompleted {
		t.Fatalf(
			"expected task COMPLETED, got %s",
			stored.Status,
		)
	}
}
