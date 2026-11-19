package planner

import (
	"context"
	"testing"

	"github.com/fjk-dev/fjk-agent/core/model"
	"github.com/fjk-dev/fjk-agent/core/tool"
)

type plannerTestProvider struct {
	content string
}

func (p plannerTestProvider) Name() string {
	return "planner-test"
}

func (p plannerTestProvider) Generate(
	ctx context.Context,
	request model.Request,
) (*model.Response, error) {
	return &model.Response{
		Content: p.content,
	}, nil
}

func (p plannerTestProvider) Stream(
	ctx context.Context,
	request model.Request,
	handler model.StreamHandler,
) error {
	if handler != nil {
		return handler(
			ctx,
			p.content,
		)
	}

	return nil
}

type plannerTestTool struct{}

func (plannerTestTool) Definition() tool.Definition {
	return tool.Definition{
		Name:        "test.echo",
		Description: "Echo input",
		Risk:        tool.RiskLow,
	}
}

func (plannerTestTool) Execute(
	ctx context.Context,
	request tool.Request,
) (*tool.Result, error) {
	return &tool.Result{
		ToolName: "test.echo",
		Output:   request.Arguments,
	}, nil
}

func TestModelPlannerReply(
	t *testing.T,
) {
	provider :=
		plannerTestProvider{
			content: `{
"task_id":"task_test",
"steps":[
{
"id":"step_001",
"type":"REPLY",
"title":"Reply",
"content":"hello"
}
]
}`,
		}

	planner :=
		NewModelPlanner(
			provider,
			tool.NewRegistry(),
		)

	plan, err :=
		planner.Plan(
			context.Background(),
			Request{
				TaskID:    "task_test",
				SessionID: "session_test",
				Input:     "hello",
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(plan.Steps) != 1 {
		t.Fatalf(
			"expected 1 step, got %d",
			len(plan.Steps),
		)
	}

	if plan.Steps[0].Type !=
		StepTypeReply {
		t.Fatalf(
			"expected REPLY, got %s",
			plan.Steps[0].Type,
		)
	}
}

func TestModelPlannerTool(
	t *testing.T,
) {
	registry :=
		tool.NewRegistry()

	if err := registry.Register(
		plannerTestTool{},
	); err != nil {
		t.Fatal(err)
	}

	provider :=
		plannerTestProvider{
			content: `{
"task_id":"task_test",
"steps":[
{
"id":"step_001",
"type":"TOOL",
"title":"Echo",
"tool_name":"test.echo",
"arguments":{
"text":"hello"
}
}
]
}`,
		}

	planner :=
		NewModelPlanner(
			provider,
			registry,
		)

	plan, err :=
		planner.Plan(
			context.Background(),
			Request{
				TaskID: "task_test",
				Input:  "echo hello",
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if plan.Steps[0].ToolName !=
		"test.echo" {
		t.Fatalf(
			"unexpected tool: %s",
			plan.Steps[0].ToolName,
		)
	}
}

func TestModelPlannerRejectsUnknownTool(
	t *testing.T,
) {
	provider :=
		plannerTestProvider{
			content: `{
"task_id":"task_test",
"steps":[
{
"id":"step_001",
"type":"TOOL",
"tool_name":"danger.fake",
"arguments":{}
}
]
}`,
		}

	planner :=
		NewModelPlanner(
			provider,
			tool.NewRegistry(),
		)

	_, err :=
		planner.Plan(
			context.Background(),
			Request{
				TaskID: "task_test",
				Input:  "do something",
			},
		)

	if err == nil {
		t.Fatal(
			"expected unknown tool rejection",
		)
	}
}
