package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/fjk-dev/fjk-agent/core/approval"
	"github.com/fjk-dev/fjk-agent/core/event"
	"github.com/fjk-dev/fjk-agent/core/permission"
	"github.com/fjk-dev/fjk-agent/core/tool"
)

type echoTool struct {
	risk tool.RiskLevel
}

func (e echoTool) Definition() tool.Definition {
	return tool.Definition{
		Name:        "test.echo",
		Description: "Echo arguments",
		Risk:        e.risk,
	}
}

func (e echoTool) Execute(
	ctx context.Context,
	request tool.Request,
) (*tool.Result, error) {
	return &tool.Result{
		ToolName: "test.echo",
		Output:   request.Arguments,
	}, nil
}

func newTestExecutor(
	t *testing.T,
	risk tool.RiskLevel,
) (
	*Executor,
	*approval.Manager,
) {
	t.Helper()

	registry :=
		tool.NewRegistry()

	if err := registry.Register(
		echoTool{
			risk: risk,
		},
	); err != nil {
		t.Fatal(err)
	}

	approvalManager :=
		approval.NewManager(
			approval.NewMemoryRepository(),
		)

	exec := New(
		registry,
		event.NewBus(),
		permission.NewEngine(),
		approvalManager,
		nil,
	)

	return exec, approvalManager
}

func TestExecuteLowRiskTool(
	t *testing.T,
) {
	exec, _ :=
		newTestExecutor(
			t,
			tool.RiskLow,
		)

	result, err :=
		exec.ExecuteTool(
			context.Background(),
			tool.Request{
				ToolName: "test.echo",

				Arguments: map[string]any{
					"text": "hello",
				},
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if result.ToolName != "test.echo" {
		t.Fatalf(
			"unexpected tool: %s",
			result.ToolName,
		)
	}
}

func TestHighRiskRequiresApproval(
	t *testing.T,
) {
	exec, _ :=
		newTestExecutor(
			t,
			tool.RiskHigh,
		)

	_, err :=
		exec.ExecuteTool(
			context.Background(),
			tool.Request{
				ToolName: "test.echo",
			},
		)

	if !errors.Is(
		err,
		permission.ErrApprovalRequired,
	) {
		t.Fatalf(
			"expected approval required, got %v",
			err,
		)
	}
}

func TestHighRiskWithApproval(
	t *testing.T,
) {
	ctx :=
		context.Background()

	exec, approvals :=
		newTestExecutor(
			t,
			tool.RiskHigh,
		)

	arguments :=
		map[string]any{
			"text": "hello",
		}

	item, err :=
		approvals.Request(
			ctx,
			"task_test",
			"session_test",
			"test.echo",
			arguments,
		)

	if err != nil {
		t.Fatal(err)
	}

	if err :=
		approvals.Approve(
			ctx,
			item.ID,
		); err != nil {
		t.Fatal(err)
	}

	_, err =
		exec.ExecuteTool(
			ctx,
			tool.Request{
				TaskID: "task_test",

				SessionID: "session_test",

				ToolName: "test.echo",

				Arguments: arguments,

				ApprovalID: item.ID,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	stored, err :=
		approvals.Get(
			ctx,
			item.ID,
		)

	if err != nil {
		t.Fatal(err)
	}

	if stored.Status !=
		approval.StatusConsumed {
		t.Fatalf(
			"expected CONSUMED, got %s",
			stored.Status,
		)
	}
}

func TestApprovalScopeMismatch(
	t *testing.T,
) {
	ctx :=
		context.Background()

	exec, approvals :=
		newTestExecutor(
			t,
			tool.RiskHigh,
		)

	item, err :=
		approvals.Request(
			ctx,
			"task_test",
			"session_test",
			"test.echo",
			map[string]any{
				"text": "approved",
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if err :=
		approvals.Approve(
			ctx,
			item.ID,
		); err != nil {
		t.Fatal(err)
	}

	_, err =
		exec.ExecuteTool(
			ctx,
			tool.Request{
				TaskID: "task_test",

				SessionID: "session_test",

				ToolName: "test.echo",

				Arguments: map[string]any{
					"text": "modified",
				},

				ApprovalID: item.ID,
			},
		)

	if !errors.Is(
		err,
		approval.ErrScopeMismatch,
	) {
		t.Fatalf(
			"expected scope mismatch, got %v",
			err,
		)
	}
}
