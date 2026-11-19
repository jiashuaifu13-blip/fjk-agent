package tool

import (
	"context"
	"testing"
)

type testTool struct{}

func (testTool) Definition() Definition {
	return Definition{
		Name:        "test.echo",
		Description: "Test tool",
		Risk:        RiskLow,
	}
}

func (testTool) Execute(
	ctx context.Context,
	request Request,
) (*Result, error) {
	return &Result{
		ToolName: "test.echo",
		Output:   request.Arguments,
	}, nil
}

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	if err := registry.Register(
		testTool{},
	); err != nil {
		t.Fatal(err)
	}

	item, err := registry.Get(
		"test.echo",
	)

	if err != nil {
		t.Fatal(err)
	}

	if item.Definition().Name != "test.echo" {
		t.Fatal(
			"unexpected tool definition",
		)
	}

	if len(registry.Definitions()) != 1 {
		t.Fatal(
			"expected one registered tool",
		)
	}
}

func TestDuplicateTool(t *testing.T) {
	registry := NewRegistry()

	if err := registry.Register(
		testTool{},
	); err != nil {
		t.Fatal(err)
	}

	if err := registry.Register(
		testTool{},
	); err == nil {
		t.Fatal(
			"expected duplicate tool error",
		)
	}
}
