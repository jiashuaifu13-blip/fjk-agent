package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/fjk-dev/fjk-agent/core/approval"
	"github.com/fjk-dev/fjk-agent/core/event"
	"github.com/fjk-dev/fjk-agent/core/permission"
	"github.com/fjk-dev/fjk-agent/core/task"
	"github.com/fjk-dev/fjk-agent/core/tool"
)

func TestApprovalTaskLifecycle(
	t *testing.T,
) {
	ctx := context.Background()

	bus := event.NewBus()

	taskManager :=
		task.NewManager(
			task.NewMemoryRepository(),
			bus,
		)

	created, err :=
		taskManager.Create(
			ctx,
			"session_approval_test",
			"Approval lifecycle test",
			"Verify WAITING_APPROVAL lifecycle",
		)

	if err != nil {
		t.Fatal(err)
	}

	if _, err :=
		taskManager.Transition(
			ctx,
			created.ID,
			task.StatusQueued,
		); err != nil {
		t.Fatal(err)
	}

	if _, err :=
		taskManager.Transition(
			ctx,
			created.ID,
			task.StatusRunning,
		); err != nil {
		t.Fatal(err)
	}

	registry :=
		tool.NewRegistry()

	if err := registry.Register(
		echoTool{
			risk: tool.RiskHigh,
		},
	); err != nil {
		t.Fatal(err)
	}

	approvalManager :=
		approval.NewManagerWithEvents(
			approval.NewMemoryRepository(),
			bus,
		)

	exec := New(
		registry,
		bus,
		permission.NewEngine(),
		approvalManager,
		taskManager,
	)

	arguments :=
		map[string]any{
			"text": "approved operation",
		}

	_, err =
		exec.ExecuteTool(
			ctx,
			tool.Request{
				TaskID: created.ID,

				SessionID: created.SessionID,

				ToolName: "test.echo",

				Arguments: arguments,
			},
		)

	var approvalRequired *ApprovalRequiredError

	if !errors.As(
		err,
		&approvalRequired,
	) {
		t.Fatalf(
			"expected ApprovalRequiredError, got %v",
			err,
		)
	}

	if approvalRequired.ApprovalID == "" {
		t.Fatal(
			"approval id is empty",
		)
	}

	waiting, err :=
		taskManager.Get(
			ctx,
			created.ID,
		)

	if err != nil {
		t.Fatal(err)
	}

	if waiting.Status !=
		task.StatusWaitingApproval {
		t.Fatalf(
			"expected WAITING_APPROVAL, got %s",
			waiting.Status,
		)
	}

	if err :=
		approvalManager.Approve(
			ctx,
			approvalRequired.ApprovalID,
		); err != nil {
		t.Fatal(err)
	}

	result, err :=
		exec.ExecuteTool(
			ctx,
			tool.Request{
				TaskID: created.ID,

				SessionID: created.SessionID,

				ToolName: "test.echo",

				Arguments: arguments,

				ApprovalID: approvalRequired.ApprovalID,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if result.ToolName != "test.echo" {
		t.Fatalf(
			"unexpected tool result: %s",
			result.ToolName,
		)
	}

	running, err :=
		taskManager.Get(
			ctx,
			created.ID,
		)

	if err != nil {
		t.Fatal(err)
	}

	if running.Status !=
		task.StatusRunning {
		t.Fatalf(
			"expected RUNNING, got %s",
			running.Status,
		)
	}

	storedApproval, err :=
		approvalManager.Get(
			ctx,
			approvalRequired.ApprovalID,
		)

	if err != nil {
		t.Fatal(err)
	}

	if storedApproval.Status !=
		approval.StatusConsumed {
		t.Fatalf(
			"expected CONSUMED, got %s",
			storedApproval.Status,
		)
	}
}
