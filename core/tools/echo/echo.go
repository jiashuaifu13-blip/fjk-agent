package echo

import (
	"context"
	"fmt"

	"github.com/fjk-dev/fjk-agent/core/tool"
)

type Tool struct{}

func New() *Tool {
	return &Tool{}
}

func (t *Tool) Definition() tool.Definition {
	return tool.Definition{
		Name: "system.echo",

		Description: "Echo the provided text. Use only when the user explicitly asks to echo or repeat text.",

		Risk: tool.RiskLow,
	}
}

func (t *Tool) Execute(
	ctx context.Context,
	request tool.Request,
) (*tool.Result, error) {
	value, ok :=
		request.Arguments["text"]

	if !ok {
		return nil, fmt.Errorf(
			"text argument is required",
		)
	}

	text, ok :=
		value.(string)

	if !ok {
		return nil, fmt.Errorf(
			"text argument must be a string",
		)
	}

	return &tool.Result{
		ToolName: "system.echo",

		Output: text,
	}, nil
}
