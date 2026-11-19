package executor

import (
	"context"
	"testing"

	"github.com/fjk-dev/fjk-agent/core/event"
	"github.com/fjk-dev/fjk-agent/core/tool"
)

type echoTool struct{}

func (echoTool) Definition() tool.Definition {
	return tool.Definition{
		Name:        "test.echo",
		Description: "Echo arguments",
		Risk:        tool.RiskLow,
	}
}

func (echoTool) Execute(
	ctx context.Context,
	request tool.Request,
) (*tool.Result, error) {
	return &tool.Result{
		ToolName: "test.echo",
		Output:   request.Arguments,
	}, nil
}

func TestExecuteTool(t *testing.T) {
	ctx := context.Background()

	registry := tool.NewRegistry()
	bus := event.NewBus()

	if err := registry.Register(
		echoTool{},
	); err != nil {
		t.Fatal(err)
	}

	called := false
	resultReceived := false

	bus.Subscribe(
		event.TypeToolCalled,
		func(
			ctx context.Context,
			evt event.Event,
		) {
			called = true
		},
	)

	bus.Subscribe(
		event.TypeToolResult,
		func(
			ctx context.Context,
			evt event.Event,
		) {
			resultReceived = true
		},
	)

	exec := New(
		registry,
		bus,
	)

	result, err := exec.ExecuteTool(
		ctx,
		tool.Request{
			TaskID:    "task_test",
			SessionID: "session_test",
			ToolName:  "test.echo",
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
			"unexpected result tool: %s",
			result.ToolName,
		)
	}

	if !called {
		t.Fatal(
			"tool.called event missing",
		)
	}

	if !resultReceived {
		t.Fatal(
			"tool.result event missing",
		)
	}
}
