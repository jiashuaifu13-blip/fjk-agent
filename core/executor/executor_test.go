package executor

import (
	"context"
	"errors"
	"testing"

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

func TestExecuteLowRiskTool(t *testing.T) {
	ctx := context.Background()

	registry := tool.NewRegistry()
	bus := event.NewBus()
	permissions := permission.NewEngine()

	if err := registry.Register(
		echoTool{
			risk: tool.RiskLow,
		},
	); err != nil {
		t.Fatal(err)
	}

	exec := New(
		registry,
		bus,
		permissions,
	)

	result, err := exec.ExecuteTool(
		ctx,
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

func TestHighRiskToolRequiresApproval(
	t *testing.T,
) {
	ctx := context.Background()

	registry := tool.NewRegistry()
	bus := event.NewBus()
	permissions := permission.NewEngine()

	if err := registry.Register(
		echoTool{
			risk: tool.RiskHigh,
		},
	); err != nil {
		t.Fatal(err)
	}

	exec := New(
		registry,
		bus,
		permissions,
	)

	_, err := exec.ExecuteTool(
		ctx,
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

func TestHighRiskToolWithApproval(
	t *testing.T,
) {
	ctx := context.Background()

	registry := tool.NewRegistry()
	bus := event.NewBus()
	permissions := permission.NewEngine()

	if err := registry.Register(
		echoTool{
			risk: tool.RiskHigh,
		},
	); err != nil {
		t.Fatal(err)
	}

	exec := New(
		registry,
		bus,
		permissions,
	)

	_, err := exec.ExecuteTool(
		ctx,
		tool.Request{
			ToolName:     "test.echo",
			UserApproved: true,
		},
	)

	if err != nil {
		t.Fatal(err)
	}
}
